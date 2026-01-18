# 015. Go 泛型（Generics）：一次编写，多种类型

[返回索引](../README.md) | [查看代码](../../015-generics/main.go)

## 本章要解决的问题

**如何避免为每种类型写重复的代码？**

---

## 第一步：场景设定

假设你在开发订单系统，需要实现一个功能：**找出切片中的最小值**。

比如：找出所有商品中最便宜的价格，或者找出字典序最小的订单ID。

```go
intPrices := []int{99, 45, 150, 30, 88}
floatPrices := []float64{99.9, 45.5, 150.0, 30.3}
names := []string{"Bob", "Alice", "Charlie"}
```

你需要从这些切片中分别找出最小值。

---

## 第二步：没有泛型时的痛点

没有泛型，你只能为每种类型写一个函数：

```go
// 处理 int 类型
func MinInt(values []int) int {
    if len(values) == 0 {
        return 0
    }
    min := values[0]
    for _, v := range values[1:] {
        if v < min {
            min = v
        }
    }
    return min
}

// 处理 float64 类型 —— 逻辑完全相同！
func MinFloat64(values []float64) float64 {
    if len(values) == 0 {
        return 0
    }
    min := values[0]
    for _, v := range values[1:] {
        if v < min {
            min = v
        }
    }
    return min
}

// 处理 string 类型 —— 又要重写一遍！
func MinString(values []string) string {
    if len(values) == 0 {
        return ""
    }
    min := values[0]
    for _, v := range values[1:] {
        if v < min {
            min = v
        }
    }
    return min
}
```

**问题**：

1. **代码重复**：三个函数的逻辑完全相同，只有类型不同
2. **维护困难**：如果要修改逻辑（比如处理空切片的方式），要改多处
3. **扩展麻烦**：新增类型（如 `int64`）就要再写一个函数

---

## 第三步：泛型函数 —— 一个函数处理所有类型

泛型允许你定义一个函数，用**类型参数**代替具体类型：

```go
func Min[T Ordered](values []T) T {
    if len(values) == 0 {
        var zero T
        return zero
    }
    min := values[0]
    for _, v := range values[1:] {
        if v < min {
            min = v
        }
    }
    return min
}
```

### 语法逐行拆解

```go
func Min[T Ordered](values []T) T
│    │  │  │       │          └── 返回值类型是 T
│    │  │  │       └── 参数 values 是 T 类型的切片
│    │  │  └── 类型约束：T 必须满足 Ordered 接口
│    │  └── 类型参数：T 是占位符，代表某种类型
│    └── 函数名
└── func 关键字
```

**关键概念**：

| 术语 | 说明 |
|------|------|
| 类型参数 | `[T ...]` 中的 `T`，是一个占位符，代表调用时传入的具体类型 |
| 类型约束 | `Ordered`，限制 `T` 必须是什么类型（必须支持 `<` 比较） |
| 类型参数列表 | 用方括号 `[]` 包裹，放在函数名和参数列表之间 |

### 使用泛型函数

```go
// 类型推断：Go 自动推断 T 是 int
result := Min([]int{99, 45, 150, 30, 88})
// result = 30

// 显式指定类型（效果相同，通常不需要）
result := Min[int]([]int{99, 45, 150, 30, 88})

// 处理 float64
minPrice := Min([]float64{99.9, 45.5, 150.0})
// minPrice = 45.5

// 处理 string（按字典序）
minName := Min([]string{"Bob", "Alice", "Charlie"})
// minName = "Alice"
```

**一个函数，处理所有类型**。

---

## 第四步：类型约束 —— 限制泛型接受的类型

### 为什么需要约束

如果 `T` 可以是任意类型，那函数体内就不能用 `<` 比较，因为不是所有类型都支持比较。

```go
// 如果没有约束，这行代码无法编译
if v < min {  // 错误：不能确定 T 类型支持 < 操作
```

约束告诉编译器：`T` 只能是支持 `<` 操作的类型。

### 内置约束

Go 提供两个内置约束：

| 约束 | 说明 |
|------|------|
| `any` | 任意类型（等同于 `interface{}`） |
| `comparable` | 可以用 `==` 和 `!=` 比较的类型 |

```go
// Contains 检查切片中是否包含某个元素
func Contains[T comparable](slice []T, target T) bool {
    for _, v := range slice {
        if v == target {  // comparable 保证可以用 ==
            return true
        }
    }
    return false
}
```

### 自定义约束

用接口定义哪些类型可以使用：

```go
// Ordered 约束：允许使用 < > 比较的类型
type Ordered interface {
    ~int | ~int8 | ~int16 | ~int32 | ~int64 |
    ~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 |
    ~float32 | ~float64 |
    ~string
}
```

### 语法拆解

```go
type Ordered interface {
    ~int | ~string | ~float64
}
```

| 符号 | 含义 |
|------|------|
| `\|` | "或"，表示类型可以是左边或右边 |
| `~` | "底层类型是"，表示包括基于该类型定义的新类型 |

**`~` 的作用**：

```go
type Price float64  // Price 的底层类型是 float64

// 没有 ~：只接受 float64，不接受 Price
type Constraint1 interface { float64 }

// 有 ~：接受 float64，也接受 Price
type Constraint2 interface { ~float64 }
```

### 标准库约束

Go 1.21+ 提供 `cmp.Ordered` 约束，可以直接使用：

```go
import "cmp"

func Min[T cmp.Ordered](values []T) T {
    // ...
}
```

---

## 第五步：多个类型参数

一个泛型可以有多个类型参数：

```go
// Keys 返回 map 的所有键
// K 是键类型，V 是值类型
func Keys[K comparable, V any](m map[K]V) []K {
    keys := make([]K, 0, len(m))
    for k := range m {
        keys = append(keys, k)
    }
    return keys
}
```

### 语法拆解

```go
func Keys[K comparable, V any](m map[K]V) []K
        │  │           │  └── V 可以是任意类型
        │  │           └── 多个类型参数用逗号分隔
        │  └── K 必须是可比较的（map 的键必须可比较）
        └── 两个类型参数：K 和 V
```

使用：

```go
orderMap := map[string]float64{
    "order-001": 99.9,
    "order-002": 150.0,
}
allKeys := Keys(orderMap)
// allKeys = ["order-001", "order-002"]
```

---

## 第六步：泛型类型 —— 可复用的数据结构

除了泛型函数，还可以定义泛型类型（结构体）。

### 场景扩展

订单系统中，操作可能成功也可能失败。定义一个通用的"结果"类型：

```go
// Result 表示操作的结果
type Result[T any] struct {
    Value T      // 成功时的值
    Error string // 失败时的错误信息
    OK    bool   // 是否成功
}
```

### 语法拆解

```go
type Result[T any] struct {
│    │     │  └── 类型约束
│    │     └── 类型参数
│    └── 类型名
└── type 关键字
```

### 为泛型类型添加方法

```go
// NewSuccess 创建成功的结果
func NewSuccess[T any](value T) Result[T] {
    return Result[T]{
        Value: value,
        OK:    true,
    }
}

// NewFailure 创建失败的结果
func NewFailure[T any](err string) Result[T] {
    return Result[T]{
        Error: err,
        OK:    false,
    }
}

// UnwrapOr 获取值，如果失败则返回默认值
func (r Result[T]) UnwrapOr(defaultValue T) T {
    if !r.OK {
        return defaultValue
    }
    return r.Value
}
```

### 使用泛型类型

```go
// 模拟查找订单
func findOrder(id string) Result[Order] {
    if order, exists := orders[id]; exists {
        return NewSuccess(order)
    }
    return NewFailure[Order]("订单不存在")
}

// 使用
result := findOrder("order-001")
if result.OK {
    fmt.Println(result.Value)
}

// 使用默认值
order := findOrder("order-999").UnwrapOr(defaultOrder)
```

---

## 第七步：实用的泛型集合操作

### Filter —— 过滤

```go
func Filter[T any](slice []T, predicate func(T) bool) []T {
    result := make([]T, 0)
    for _, v := range slice {
        if predicate(v) {
            result = append(result, v)
        }
    }
    return result
}
```

**`predicate func(T) bool` 语法拆解**：
- 这是一个函数类型的参数
- `func(T)` —— 接受一个 `T` 类型的参数
- `bool` —— 返回布尔值

使用：

```go
orders := []Order{
    {Customer: "Alice", Amount: 99.9},
    {Customer: "Bob", Amount: 150.0},
    {Customer: "Alice", Amount: 45.5},
}

// 过滤 Alice 的订单
aliceOrders := Filter(orders, func(o Order) bool {
    return o.Customer == "Alice"
})
```

### Map —— 转换

```go
func Map[T any, R any](slice []T, transform func(T) R) []R {
    result := make([]R, len(slice))
    for i, v := range slice {
        result[i] = transform(v)
    }
    return result
}
```

使用：

```go
// 提取所有订单金额
amounts := Map(orders, func(o Order) float64 {
    return o.Amount
})
// amounts = [99.9, 150.0, 45.5]
```

### Reduce —— 归约

```go
func Reduce[T any, R any](slice []T, initial R, reducer func(R, T) R) R {
    result := initial
    for _, v := range slice {
        result = reducer(result, v)
    }
    return result
}
```

使用：

```go
// 计算订单总金额
total := Reduce(orders, 0.0, func(sum float64, o Order) float64 {
    return sum + o.Amount
})
// total = 295.4
```

### 链式组合

```go
// 计算 Alice 的订单总金额
aliceTotal := Reduce(
    Filter(orders, func(o Order) bool { return o.Customer == "Alice" }),
    0.0,
    func(sum float64, o Order) float64 { return sum + o.Amount },
)
```

---

## 第八步：最佳实践

### 何时使用泛型

| 场景 | 是否使用泛型 |
|------|-------------|
| 通用数据结构（栈、队列、链表） | ✅ 使用 |
| 集合操作（Filter、Map、Reduce） | ✅ 使用 |
| 不同类型有相同的操作逻辑 | ✅ 使用 |
| 只处理特定类型 | ❌ 不需要 |
| 可以用接口解决 | ⚠️ 优先考虑接口 |

### 泛型 vs 接口

```go
// 接口：定义行为契约，不同类型有不同实现
type Payer interface {
    Pay(amount float64)
}

// 泛型：相同逻辑，不同类型
func Min[T Ordered](values []T) T
```

**选择原则**：
- 如果不同类型有**不同的行为实现** → 用接口
- 如果不同类型有**相同的操作逻辑** → 用泛型

### 命名约定

| 约定 | 说明 |
|------|------|
| `T` | 通用类型（Type） |
| `K`, `V` | 键值对（Key, Value） |
| `E` | 元素（Element） |
| `R` | 结果（Result） |

---

## 总结

| 概念 | 语法 | 说明 |
|------|------|------|
| 泛型函数 | `func Name[T Constraint](...)` | 用类型参数处理多种类型 |
| 类型参数 | `[T ...]` | 占位符，代表调用时的具体类型 |
| 类型约束 | `any`, `comparable`, 自定义接口 | 限制类型参数可以是什么 |
| 类型推断 | `Min(slice)` | 编译器自动推断类型参数 |
| 泛型类型 | `type Name[T Constraint] struct` | 可复用的数据结构 |
| `~` 符号 | `~int` | 包括底层类型是 int 的所有类型 |
| `\|` 符号 | `int \| string` | 类型可以是 int 或 string |
