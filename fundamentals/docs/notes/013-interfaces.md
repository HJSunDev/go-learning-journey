# 013. Go 接口（Interface）：定义行为的契约

[返回索引](../README.md) | [查看代码](../../013-interfaces/main.go)

## 本章要解决的问题

**如何让一个函数能接受多种不同的类型？**

---

## 第一步：场景设定

假设你在开发支付系统，有三种支付方式：

```go
type WechatPay struct{ UserID string }
type Alipay struct{ Account string }
type BankCard struct{ CardNo string }
```

每种支付方式都有 `Pay` 方法：

```go
func (w WechatPay) Pay(amount float64) { fmt.Printf("微信支付 %.2f 元\n", amount) }
func (a Alipay) Pay(amount float64)    { fmt.Printf("支付宝支付 %.2f 元\n", amount) }
func (b BankCard) Pay(amount float64)  { fmt.Printf("银行卡支付 %.2f 元\n", amount) }
```

直接使用没问题：

```go
wechat := WechatPay{UserID: "wx_001"}
wechat.Pay(100)  // 微信支付 100.00 元
```

**问题来了**：如果想写一个函数，能处理任意支付方式，参数类型写什么？

```go
func Checkout(??? , amount float64) {
    ???.Pay(amount)
}
```

写 `WechatPay`？那就只能处理微信。写 `Alipay`？那就只能处理支付宝。

**这就是接口要解决的问题。**

---

## 第二步：定义接口

接口定义了"必须有什么方法"：

```go
// Payer 接口：任何想成为 Payer 的类型，必须有 Pay(float64) 方法
type Payer interface {
    Pay(amount float64)
}
```

因为 `WechatPay`、`Alipay`、`BankCard` 都有 `Pay(float64)` 方法，所以它们**自动满足** `Payer` 接口。

Go 不需要写 `implements Payer`，只要方法匹配就行。这叫**隐式实现**。

---

## 第三步：接口类型的变量

**关键理解**：接口类型的变量，可以存储任何满足该接口的值。

```go
var p Payer          // 声明一个 Payer 接口类型的变量

p = wechat           // 合法！WechatPay 有 Pay 方法，满足 Payer 接口
p.Pay(100)           // 调用的是 WechatPay 的 Pay 方法
// 输出：微信支付 100.00 元

p = alipay           // 合法！Alipay 也有 Pay 方法
p.Pay(200)           // 调用的是 Alipay 的 Pay 方法
// 输出：支付宝支付 200.00 元

p = bank             // 合法！BankCard 也有 Pay 方法
p.Pay(300)           // 调用的是 BankCard 的 Pay 方法
// 输出：银行卡支付 300.00 元
```

同一个变量 `p`，存不同的值，调用 `Pay` 方法时表现不同。**这就是多态**。

---

## 第四步：接口作为函数参数

现在可以解决最初的问题了：

```go
// 参数类型是 Payer 接口
// 任何满足 Payer 接口的类型都可以传进来
func Checkout(p Payer, amount float64) {
    p.Pay(amount)
}
```

使用：

```go
Checkout(wechat, 100)   // 传入 WechatPay，输出：微信支付 100.00 元
Checkout(alipay, 200)   // 传入 Alipay，输出：支付宝支付 200.00 元
Checkout(bank, 300)     // 传入 BankCard，输出：银行卡支付 300.00 元
```

**一个函数处理所有支付方式**。新增支付方式时，只要它有 `Pay` 方法，就能直接用，`Checkout` 函数不用改。

---

## 第五步：接口变量内部是什么

接口变量内部存储两样东西：

| 名称 | 说明 |
|------|------|
| 动态类型 | 实际存的是什么类型 |
| 动态值 | 实际存的值是什么 |

```go
var p Payer
// p 内部：(nil, nil)

p = WechatPay{UserID: "wx_001"}
// p 内部：(WechatPay, WechatPay{UserID: "wx_001"})

p = Alipay{Account: "alice"}
// p 内部：(Alipay, Alipay{Account: "alice"})
```

调用 `p.Pay(100)` 时，Go 查看 `p` 内部的动态类型，找到对应的 `Pay` 方法执行。

---

## 第六步：类型断言 —— 从接口取出具体类型

接口变量只能调用接口定义的方法。如果想访问具体类型的字段，需要**类型断言**。

### 语法

```go
value, ok := 接口变量.(具体类型)
```

- 如果接口变量里存的确实是这个类型：`ok` 为 `true`，`value` 是取出的值
- 如果不是这个类型：`ok` 为 `false`，`value` 是零值

**注意**：`.(Type)` 不是函数调用，是类型断言的特殊语法。

### 示例

```go
var p Payer = WechatPay{UserID: "wx_001"}

// 尝试断言为 WechatPay
w, ok := p.(WechatPay)
// w 的类型是 WechatPay
// ok 的类型是 bool

if ok {
    fmt.Println("是微信支付，用户ID:", w.UserID)  // 可以访问 UserID 字段了
}

// 尝试断言为 Alipay
a, ok := p.(Alipay)
if ok {
    fmt.Println("是支付宝")
} else {
    fmt.Println("不是支付宝")  // 会执行这个
}
```

### 常见写法：if 语句中使用

Go 允许在 if 条件中同时赋值和判断：

```go
if w, ok := p.(WechatPay); ok {
    fmt.Println("用户ID:", w.UserID)
}
```

语法拆解：
- 分号前面 `w, ok := p.(WechatPay)` 是赋值语句
- 分号后面 `ok` 是条件判断
- `w` 变量只在 if 块内有效

---

## 第七步：type switch —— 判断多种类型

当需要判断多种类型时，用 type switch 更简洁：

```go
switch v := p.(type) {
case WechatPay:
    fmt.Println("微信支付，用户:", v.UserID)
case Alipay:
    fmt.Println("支付宝，账户:", v.Account)
case BankCard:
    fmt.Println("银行卡，卡号:", v.CardNo)
default:
    fmt.Println("未知类型")
}
```

**注意**：`.(type)` 只能在 switch 语句中使用，这是特殊语法。

---

## 第八步：空接口 any

空接口没有任何方法要求，所以**所有类型都满足它**：

```go
var x any       // any 是 interface{} 的别名

x = 42          // 可以存 int
x = "hello"     // 可以存 string
x = wechat      // 可以存 WechatPay
```

从 `any` 取出具体值也需要类型断言：

```go
if num, ok := x.(int); ok {
    fmt.Println("是整数:", num)
}
```

### 使用场景

```go
// 1. JSON 解析未知结构
var data map[string]any
json.Unmarshal(jsonBytes, &data)

// 2. 可变参数
fmt.Printf(format string, args ...any)
```

### 代价

使用 `any` 意味着失去类型安全，取值时必须类型断言。

---

## 总结

| 步骤 | 内容 |
|------|------|
| 定义接口 | `type Payer interface { Pay(float64) }` |
| 隐式实现 | 类型有接口要求的方法，就自动满足该接口 |
| 接口变量 | 可以存储任何满足该接口的值 |
| 多态 | 同一个接口变量，存不同值，调用方法时表现不同 |
| 函数参数 | 用接口类型，可以接受任何满足该接口的值 |
| 类型断言 | `v, ok := i.(Type)` 从接口取出具体类型 |
| type switch | `switch v := i.(type)` 判断多种类型 |
| 空接口 any | 可存任何类型，但失去类型安全 |
