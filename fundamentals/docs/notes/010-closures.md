# 010. Go 闭包 (Closures)：自带状态的函数

[返回索引](../README.md) | [查看代码](../../010-closures/main.go)

闭包（Closure）是 Go 语言中一个非常核心的概念。

**一句话定义**：闭包是一个**函数值**，它引用了函数体之外的变量。这个函数可以访问并修改这些被引用的变量。

## 1. 核心机制：变量捕获

在 Go 语言中，我们不能在函数内部定义命名函数（嵌套函数），但可以定义**匿名函数**（函数字面量）。

当这个匿名函数引用了外部作用域的变量时，Go 编译器会把这个变量和函数绑定在一起，形成一个**闭包**。

```go
func createCounter() func() int {
    count := 0 // 局部变量
    
    // 返回一个闭包
    // 这个匿名函数引用了外部的 count
    return func() int {
        count++ 
        return count
    }
}

func main() {
    c1 := createCounter()
    fmt.Println(c1()) // 输出 1
    fmt.Println(c1()) // 输出 2
    
    c2 := createCounter() // 创建一个新的闭包，拥有独立的 count
    fmt.Println(c2()) // 输出 1
}
```

### 为什么变量没有销毁？

正常情况下，`count` 是 `createCounter` 的局部变量，函数执行完应该就被销毁了。
但是在闭包中，**Go 编译器检测到 `count` 变量在函数返回后依然被引用**，因此会自动将 `count` 变量分配到**堆（Heap）**上，而不是栈上。这样即使 `createCounter` 函数结束了，`count` 依然存在，直到闭包不再被使用。

## 2. 只有匿名函数才能做闭包吗？

**是的，在 Go 语言中是这样。**

准确的表述是：**Go 语言中的闭包是通过函数字面量（Function Literals）实现的。**

*   **命名函数**（如 `func main()`, `func add()`）是定义在包级别的，它们没有“外部局部作用域”可以捕获，只能访问全局变量。
*   只有在函数内部定义的**匿名函数**，才能捕获该函数的局部变量，从而形成闭包。

> 注意：虽然你不能在函数内部写 `func myFunc() {}`，但你可以把匿名函数赋值给一个变量 `myFunc := func() {}`，但这本质上依然是匿名函数赋值。

## 3. 实战场景：函数适配与封装（延迟计算）

这是闭包最常见的工程用途之一。

**问题背景**：
你使用了一个第三方库，或者标准库的某个功能（如 `http.HandleFunc`，或者协程池任务），它要求传入的回调函数必须是无参数的 `func()`。
但你的业务逻辑函数 `sendEmail(content string, userID int)` 是需要参数的。

**解决方案**：
利用闭包，将参数“包裹”进一个无参函数中。

### 代码示例

```go
// 1. 假设这是第三方库提供的函数，只接收无参回调
func RunTask(task func()) {
    fmt.Println("Start task...")
    task() // 执行回调
}

// 2. 这是你的业务函数，需要具体参数
func sendEmail(msg string, userID int) {
    fmt.Printf("Email to %d: %s\n", userID, msg)
}

func main() {
    userID := 101
    msg := "Alert!"

    // 错误：直接传递签名不匹配
    // RunTask(sendEmail) // Compile Error

    // 正确：创建一个闭包来适配
    // 这个匿名函数符合 func() 的签名
    wrappedTask := func() {
        // 在这里，它捕获了外部的 msg 和 userID
        sendEmail(msg, userID)
    }

    // 将闭包传递给库函数
    RunTask(wrappedTask)
}
```

在这个例子中，`wrappedTask` 就是一个闭包。它把 `sendEmail` 需要的参数（状态）和调用逻辑（行为）打包在了一起，对外暴露成一个简单的 `func()` 接口。

## 4. 内存管理

*   **生命周期**：被闭包捕获的变量（如 `count`），其生命周期会延长到和闭包函数一致。
*   **垃圾回收 (GC)**：只要闭包函数还被变量引用，它捕获的变量就不会被回收。当闭包变量置为 `nil` 或离开作用域后，这些捕获的变量才会随之被 GC 回收。
