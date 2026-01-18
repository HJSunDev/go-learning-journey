# 014. Go 枚举：用 const 和 iota 构建类型安全的常量集合

[返回索引](../README.md) | [查看代码](../../014-enums/main.go)

## 本章要解决的问题

**如何优雅地表示一组固定的、有限的选项？**

---

## 第一步：场景设定 —— 订单状态管理

假设你在开发电商系统，订单有五种状态：

| 状态 | 含义 |
|------|------|
| 待支付 | 用户下单但未付款 |
| 已支付 | 用户已付款，等待发货 |
| 已发货 | 商品已寄出 |
| 已完成 | 用户确认收货 |
| 已取消 | 订单被取消 |

最直接的做法是用整数表示：

```go
// 0=待支付, 1=已支付, 2=已发货, 3=已完成, 4=已取消
func ProcessOrder(status int) {
    if status == 1 {
        fmt.Println("订单已支付，准备发货")
    }
}
```

**问题来了**：

1. **魔法数字**：`status == 1` 是什么意思？看代码的人不知道
2. **类型不安全**：`ProcessOrder(100)` 不会报错，但 100 不是有效状态
3. **难以维护**：如果状态编号改了，要到处找和改

**这就是"枚举"要解决的问题。**

---

## 第二步：用 const 定义常量

Go 使用 `const` 关键字定义常量：

```go
const (
    StatusPending   = 0  // 待支付
    StatusPaid      = 1  // 已支付
    StatusShipped   = 2  // 已发货
    StatusCompleted = 3  // 已完成
    StatusCancelled = 4  // 已取消
)
```

### 语法拆解

```go
const (
    常量名1 = 值1
    常量名2 = 值2
)
```

- `const` 关键字声明常量
- 圆括号 `()` 将多个常量分组
- 常量一旦赋值，不可修改

### 使用效果

```go
func ProcessOrder(status int) {
    if status == StatusPaid {  // 可读性好多了
        fmt.Println("订单已支付，准备发货")
    }
}
```

**改进**：代码可读性提升了。

**遗留问题**：`status` 的类型还是 `int`，任何整数都能传进来。

---

## 第三步：用 iota 简化常量定义

手动写 0, 1, 2, 3, 4 太麻烦，Go 提供了 `iota` 自动生成递增数值。

```go
const (
    StatusPending   = iota  // 0
    StatusPaid              // 1
    StatusShipped           // 2
    StatusCompleted         // 3
    StatusCancelled         // 4
)
```

### iota 工作原理

| 规则 | 说明 |
|------|------|
| 初始值 | `iota` 在每个 `const` 块开始时重置为 0 |
| 自动递增 | 每定义一个常量，`iota` 自动加 1 |
| 表达式继承 | 后续常量不写 `= iota`，自动继承上一行的表达式 |

### 语法拆解

```go
const (
    StatusPending   = iota  // iota = 0，StatusPending = 0
    StatusPaid              // iota = 1，继承上一行的 "= iota"，所以 StatusPaid = 1
    StatusShipped           // iota = 2，StatusShipped = 2
)
```

- 第一行必须写 `= iota`
- 后续行省略 `= iota`，Go 自动继承

---

## 第四步：自定义类型 —— 类型安全

这是 Go 枚举的核心技巧：**定义新类型**。

```go
// 定义一个新类型 OrderStatus，底层是 int
type OrderStatus int

const (
    Pending   OrderStatus = iota  // 类型是 OrderStatus，值为 0
    Paid                          // 类型是 OrderStatus，值为 1
    Shipped                       // 类型是 OrderStatus，值为 2
    Completed                     // 类型是 OrderStatus，值为 3
    Cancelled                     // 类型是 OrderStatus，值为 4
)
```

### 语法拆解

```go
type OrderStatus int
```

- `type` 关键字定义新类型
- `OrderStatus` 是新类型的名字
- `int` 是底层类型（OrderStatus 本质上还是个整数）

```go
Pending OrderStatus = iota
```

- `Pending` 是常量名
- `OrderStatus` 是常量的类型
- `= iota` 是常量的值

### 类型安全的效果

```go
func ProcessOrder(status OrderStatus) {
    // 只接受 OrderStatus 类型的参数
}

ProcessOrder(Paid)     // 合法
ProcessOrder(Shipped)  // 合法
ProcessOrder(1)        // 编译错误！1 是 int，不是 OrderStatus
```

**关键理解**：虽然 `OrderStatus` 底层是 `int`，但 Go 认为它们是不同的类型。必须显式转换：

```go
ProcessOrder(OrderStatus(1))  // 合法，显式转换为 OrderStatus
```

---

## 第五步：给枚举类型添加方法

Go 允许给自定义类型添加方法，这是实现"枚举行为"的关键。

### 5.1 实现 String() 方法

`fmt` 包在打印时会自动调用 `String()` 方法（如果类型实现了的话）。

```go
func (s OrderStatus) String() string {
    names := []string{"待支付", "已支付", "已发货", "已完成", "已取消"}
    
    // 边界检查：确保状态值在有效范围内
    if s < 0 || int(s) >= len(names) {
        return "未知状态"
    }
    
    return names[s]
}
```

#### 语法拆解

```go
func (s OrderStatus) String() string {
```

- `(s OrderStatus)` 是接收器，表示这个方法属于 `OrderStatus` 类型
- `String()` 是方法名
- `string` 是返回值类型

```go
names := []string{"待支付", "已支付", "已发货", ...}
```

- 定义一个字符串切片，索引和枚举值对应
- `names[0]` = "待支付"，对应 `Pending` (值为 0)
- `names[1]` = "已支付"，对应 `Paid` (值为 1)

```go
return names[s]
```

- `s` 是 `OrderStatus` 类型，底层是 `int`
- 直接用 `s` 作为索引访问切片

#### 使用效果

```go
fmt.Println(Paid)           // 输出：已支付
fmt.Printf("状态：%s\n", Shipped)  // 输出：状态：已发货
```

### 5.2 添加业务方法

```go
// IsValid 检查状态值是否有效
func (s OrderStatus) IsValid() bool {
    return s >= Pending && s <= Cancelled
}

// CanCancel 判断当前状态下订单是否可以取消
func (s OrderStatus) CanCancel() bool {
    return s == Pending || s == Paid
}

// CanShip 判断当前状态下订单是否可以发货
func (s OrderStatus) CanShip() bool {
    return s == Paid
}
```

#### 使用效果

```go
status := Paid
if status.CanCancel() {
    fmt.Println("订单可以取消")
}

if status.CanShip() {
    fmt.Println("订单可以发货")
}
```

**核心价值**：业务逻辑封装在枚举类型内部，调用方不需要知道具体规则。

---

## 第六步：iota 进阶用法

### 6.1 从 1 开始编号

```go
type Priority int

const (
    Low    Priority = iota + 1  // 0 + 1 = 1
    Medium                      // 1 + 1 = 2
    High                        // 2 + 1 = 3
)
```

`iota + 1` 这个表达式会被后续常量继承，所以每个值都是 `iota + 1`。

### 6.2 位运算生成标志位

```go
type FilePermission int

const (
    Read  FilePermission = 1 << iota  // 1 << 0 = 1  (二进制: 001)
    Write                             // 1 << 1 = 2  (二进制: 010)
    Exec                              // 1 << 2 = 4  (二进制: 100)
)
```

#### 位运算语法

`1 << iota` 表示将 1 左移 iota 位：

| iota | 1 << iota | 二进制 |
|------|-----------|--------|
| 0 | 1 | 001 |
| 1 | 2 | 010 |
| 2 | 4 | 100 |

#### 组合权限

```go
// 用位或运算组合多个权限
readWrite := Read | Write  // 1 | 2 = 3 (二进制: 011)

// 用位与运算检查是否有某权限
if readWrite & Read != 0 {
    fmt.Println("有读权限")
}
```

### 6.3 跳过某个值

```go
const (
    Monday = iota + 1  // 1
    Tuesday            // 2
    Wednesday          // 3
    _                  // 4，用下划线跳过
    Friday             // 5
)
```

下划线 `_` 是空白标识符，用于丢弃不需要的值。

---

## 最佳实践

### 1. 始终使用自定义类型

```go
// 好：有类型安全
type Color int
const (
    Red Color = iota
    Green
    Blue
)

// 差：没有类型安全
const (
    Red   = 0
    Green = 1
    Blue  = 2
)
```

### 2. 实现 String() 方法

让枚举值可以输出可读的名称，方便调试和日志。

### 3. 添加验证方法

```go
func (c Color) IsValid() bool {
    return c >= Red && c <= Blue
}
```

### 4. 零值要有意义或明确标记

```go
const (
    Unknown OrderStatus = iota  // 0 是未知状态
    Pending
    Paid
)
```

或者从 1 开始，让 0 表示"未设置"。

---

## 总结

| 步骤 | 内容 |
|------|------|
| 痛点 | 魔法数字不可读、没有类型安全 |
| const | 用常量替代魔法数字，提高可读性 |
| iota | 自动生成递增数值，简化定义 |
| 自定义类型 | `type X int` 增加类型安全 |
| 添加方法 | `String()`、业务方法封装逻辑 |
| iota 进阶 | `iota+1` 偏移、`1<<iota` 位运算 |
