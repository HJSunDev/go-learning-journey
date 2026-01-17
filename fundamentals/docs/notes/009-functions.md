# 009. Go 函数

> 掌握 Go 函数的定义与使用，理解参数传递机制、多返回值、可变参数，以及函数作为一等公民的高阶用法。

## 核心概念

函数是 Go 语言的核心构建块，用于将代码组织成可复用的逻辑单元。

Go 中的函数是**一等公民**（First-class Citizen），这意味着函数可以像普通值一样被赋值、传递和返回。

### 函数的核心特性

| 特性       | 说明                                         |
| ---------- | -------------------------------------------- |
| 多返回值   | 一个函数可以返回多个值，常用于返回结果和错误 |
| 命名返回值 | 可以在函数签名中为返回值命名，提高可读性     |
| 可变参数   | 使用 `...` 语法接收任意数量的参数          |
| 函数类型   | 函数有具体类型，可以作为变量、参数和返回值   |
| 匿名函数   | 支持没有名字的函数字面量                     |
| defer 机制 | 延迟执行函数调用，用于资源清理               |

---

## 1. 函数基础

### 1.1 函数声明语法

```go
func 函数名(参数列表) 返回值类型 {
    函数体
}

// 示例：无参无返回值
func greet() {
    fmt.Println("Hello, Go!")
}

// 示例：有参有返回值
func add(a int, b int) int {
    return a + b
}
```

### 1.2 可见性规则

Go 使用首字母大小写控制可见性，这个规则同样适用于函数：

```go
// 小写开头：包内私有，只能在同一个包内调用
func greet() { ... }

// 大写开头：公开导出，可以被其他包调用
func Greet() { ... }
```

| 命名      | 可见性   | 使用场景               |
| --------- | -------- | ---------------------- |
| `greet` | 包内私有 | 内部辅助函数           |
| `Greet` | 公开导出 | 提供给其他包使用的 API |

---

## 2. 参数传递

### 2.1 Go 是值传递语言

**核心原则**：Go 中所有的参数传递都是**值传递**。传递给函数的是参数的副本，而不是原始值本身。

```go
func swap(a, b int) {
    a, b = b, a  // 只交换了副本，原变量不受影响
}

x, y := 10, 20
swap(x, y)
fmt.Println(x, y)  // 输出: 10 20（未改变）
```

### 2.2 参数类型简写

当多个连续参数类型相同时，可以只在最后一个参数后声明类型：

```go
// 完整写法
func add(a int, b int) int { ... }

// 简写形式
func add(a, b int) int { ... }
```

### 2.3 指针参数

如果需要修改原变量，需要传递指针：

```go
func swapByPointer(a, b *int) {
    *a, *b = *b, *a  // 通过指针修改原值
}

x, y := 10, 20
swapByPointer(&x, &y)
fmt.Println(x, y)  // 输出: 20 10（已交换）
```

### 2.4 切片与映射参数

切片和映射本身是**引用类型**，传递时复制的是"头部结构"（指针、长度、容量），但共享底层数据：

```go
func modifySlice(s []int) {
    if len(s) > 0 {
        s[0] = 999  // 修改底层数组，原切片可见
    }
}

nums := []int{1, 2, 3}
modifySlice(nums)
fmt.Println(nums)  // 输出: [999 2 3]
```

### 2.5 参数传递总结

| 类型                       | 传递内容   | 能否修改原值       |
| -------------------------- | ---------- | ------------------ |
| 基本类型（int, string 等） | 值的副本   | 否                 |
| 指针                       | 地址的副本 | 是（通过解引用）   |
| 切片                       | 头部的副本 | 是（共享底层数组） |
| 映射                       | 引用的副本 | 是（共享底层数据） |

---

## 3. 返回值

### 3.1 单返回值

```go
func double(n int) int {
    return n * 2
}
```

### 3.2 多返回值

Go 支持函数返回多个值，这是区别于许多其他语言的重要特性：

```go
func divide(a, b int) (int, int) {
    return a / b, a % b  // 返回商和余数
}

quotient, remainder := divide(17, 5)
fmt.Println(quotient, remainder)  // 输出: 3 2
```

### 3.3 忽略部分返回值

使用空白标识符 `_` 忽略不需要的返回值：

```go
quotient, _ := divide(17, 5)  // 只需要商，忽略余数
```

### 3.4 错误处理惯用模式

Go 使用多返回值实现错误处理，这是最常见的模式：

```go
func divideWithError(a, b int) (int, error) {
    if b == 0 {
        return 0, fmt.Errorf("除数不能为零")
    }
    return a / b, nil
}

// 调用时检查错误
result, err := divideWithError(10, 0)
if err != nil {
    fmt.Println("错误:", err)
    return
}
fmt.Println("结果:", result)
```

### 3.5 命名返回值

可以在函数签名中为返回值命名，命名返回值会被自动初始化为零值：

```go
func rectangle(width, height int) (area int, perimeter int) {
    area = width * height
    perimeter = 2 * (width + height)
    return  // 裸返回：直接返回命名的返回值
}

// 推荐：显式返回更清晰
func rectangleExplicit(width, height int) (area int, perimeter int) {
    area = width * height
    perimeter = 2 * (width + height)
    return area, perimeter  // 显式返回
}
```

**最佳实践**：

- 命名返回值主要用于文档说明
- 复杂函数避免裸返回，显式返回更清晰
- 短小函数可以使用裸返回

---

## 4. 可变参数函数

### 4.1 基本语法

使用 `...T` 声明可变参数，在函数内部它是一个 `[]T` 切片：

```go
func sum(nums ...int) int {
    total := 0
    for _, n := range nums {
        total += n
    }
    return total
}

// 调用方式
sum()           // 0
sum(1)          // 1
sum(1, 2, 3)    // 6
sum(1, 2, 3, 4) // 10
```

### 4.2 传递切片

将切片传给可变参数函数时，需要使用 `...` 展开：

```go
numbers := []int{10, 20, 30}
result := sum(numbers...)  // 使用 ... 展开切片
```

### 4.3 固定参数 + 可变参数

可变参数必须是最后一个参数：

```go
func printf(format string, args ...interface{}) {
    fmt.Printf(format, args...)
}

printf("Name: %s, Age: %d\n", "Go", 15)
```

### 4.4 可变参数规则

| 规则     | 说明                                 |
| -------- | ------------------------------------ |
| 位置限制 | 可变参数必须是参数列表的最后一个     |
| 内部类型 | 在函数内部，可变参数是对应类型的切片 |
| 传递切片 | 使用 `slice...` 语法展开切片       |
| 空调用   | 可以不传任何参数，此时切片长度为 0   |

---

## 5. 函数类型与函数作为值

### 5.1 函数是一等公民

Go 中函数是一等公民，这意味着：

- 可以赋值给变量
- 可以作为参数传递
- 可以作为返回值
- 可以存储在数据结构中

### 5.2 定义函数类型

使用 `type` 为函数签名定义类型别名：

```go
// 定义一个函数类型
type Operation func(a, b int) int

// 符合该类型的函数
func add(a, b int) int      { return a + b }
func subtract(a, b int) int { return a - b }
func multiply(a, b int) int { return a * b }
```

### 5.3 函数赋值给变量

```go
var op Operation = add
fmt.Println(op(10, 5))  // 15

op = subtract
fmt.Println(op(10, 5))  // 5
```

### 5.4 函数存储在数据结构中

```go
operations := map[string]Operation{
    "add":      add,
    "subtract": subtract,
    "multiply": multiply,
}

for name, fn := range operations {
    fmt.Printf("%s(6, 3) = %d\n", name, fn(6, 3))
}
```

### 5.5 函数类型的零值

函数类型的零值是 `nil`，调用 nil 函数会导致 panic：

```go
var nilFunc Operation
fmt.Println(nilFunc == nil)  // true
// nilFunc(1, 2)  // panic: nil pointer dereference
```

---

## 6. 匿名函数

匿名函数是没有名字的函数，也称为**函数字面量**（Function Literal）。

### 6.1 匿名函数赋值给变量

```go
square := func(n int) int {
    return n * n
}

fmt.Println(square(5))  // 25
```

### 6.2 立即调用的匿名函数（IIFE）

```go
result := func(a, b int) int {
    return a + b
}(10, 20)  // 定义后立即调用

fmt.Println(result)  // 30
```

### 6.3 匿名函数的使用场景

| 场景      | 说明                          |
| --------- | ----------------------------- |
| 回调函数  | 作为参数传递给高阶函数        |
| goroutine | 作为并发执行的函数体          |
| defer     | 在 defer 中执行复杂的清理逻辑 |
| 闭包      | 捕获外部变量（下一章节详解）  |

---

## 7. 函数作为参数（高阶函数）

接收函数作为参数或返回函数的函数称为**高阶函数**。

### 7.1 Map 操作

对集合中的每个元素应用转换函数：

```go
func applyToAll(nums []int, fn func(int) int) []int {
    result := make([]int, len(nums))
    for i, n := range nums {
        result[i] = fn(n)
    }
    return result
}

nums := []int{1, 2, 3, 4, 5}
doubled := applyToAll(nums, func(n int) int { return n * 2 })
// doubled: [2 4 6 8 10]
```

### 7.2 Filter 操作

过滤集合，保留满足条件的元素：

```go
func filter(nums []int, predicate func(int) bool) []int {
    result := []int{}
    for _, n := range nums {
        if predicate(n) {
            result = append(result, n)
        }
    }
    return result
}

evens := filter(nums, func(n int) bool { return n%2 == 0 })
// evens: [2 4]
```

### 7.3 Reduce 操作

将集合归约为单个值：

```go
func reduce(nums []int, initial int, fn func(acc, curr int) int) int {
    result := initial
    for _, n := range nums {
        result = fn(result, n)
    }
    return result
}

sum := reduce(nums, 0, func(acc, curr int) int { return acc + curr })
// sum: 15
```

### 7.4 组合使用

```go
// 计算偶数的平方和
result := reduce(
    applyToAll(
        filter(nums, func(n int) bool { return n%2 == 0 }),
        func(n int) int { return n * n },
    ),
    0,
    func(acc, curr int) int { return acc + curr },
)
// result: 20 (2² + 4² = 4 + 16)
```

---

## 8. 函数作为返回值

### 8.1 工厂函数模式

根据参数生成定制化的函数：

```go
func makeMultiplier(factor int) func(int) int {
    return func(n int) int {
        return n * factor
    }
}

double := makeMultiplier(2)
triple := makeMultiplier(3)

fmt.Println(double(5))  // 10
fmt.Println(triple(5))  // 15
```

### 8.2 闭包与状态记忆（进阶）

闭包的核心特性是**捕获了外部变量**。即使外部函数已经返回，闭包依然持有对这些变量的引用。

```go
func makeCounter() func() int {
    // count 逃逸到堆上
    count := 0
  
    return func() int {
        // 闭包持有对 count 的引用
        count++
        return count
    }
}
```

#### 内存管理机制

1. **逃逸分析 (Escape Analysis)**:

   - 编译器发现 `count` 变量在 `makeCounter` 返回后仍被闭包使用。
   - 因此，`count` 不会分配在栈上（栈内存随函数结束销毁），而是**分配到堆（Heap）上**。
2. **生命周期与回收**:

   - **持有**: 只要你持有的闭包函数（如 `counter := makeCounter()`）还在使用，堆上的 `count` 就一直存在。
   - **回收**: 当闭包函数超出作用域或被置为 `nil`，且没有任何其他引用时，Go 的垃圾回收器（GC）会自动回收这块内存。
3. **安全性**:

   - **内存泄漏**: 在局部作用域使用（如示例中）是安全的，函数结束即回收。只有在将闭包保存在全局变量或长生命周期的 Map 中且不清理时，才可能导致泄漏。
   - **并发安全**: 这里的 `count` 没有加锁。如果多个 Goroutine 同时调用同一个 `counter` 实例，会发生数据竞争。并发场景需加锁保护。

### 8.3 格式化器工厂

```go
func makeFormatter(prefix, suffix string) func(string) string {
    return func(s string) string {
        return prefix + s + suffix
    }
}

wrapper := makeFormatter("[", "]")
htmlTag := makeFormatter("<p>", "</p>")

fmt.Println(wrapper("Hello"))    // [Hello]
fmt.Println(htmlTag("Content"))  // <p>Content</p>
```

### 8.3 函数作为返回值的应用场景

| 场景     | 说明                                 |
| -------- | ------------------------------------ |
| 工厂模式 | 根据参数生成定制函数                 |
| 延迟计算 | 返回的函数包含待执行的逻辑           |
| 闭包     | 捕获外部变量，保持状态（下一章详解） |
| 依赖注入 | 返回配置好依赖的函数                 |

---

## 9. 递归函数

递归是函数调用自身的编程技术，用于解决可以分解为相同子问题的问题。

### 9.1 递归的两个要素

1. **基准情况（Base Case）**：递归终止的条件
2. **递归情况（Recursive Case）**：缩小问题规模并调用自身

### 9.2 阶乘示例

```go
func factorial(n int) int {
    // 基准情况
    if n <= 1 {
        return 1
    }
    // 递归情况
    return n * factorial(n-1)
}

// 5! = 5 * 4 * 3 * 2 * 1 = 120
fmt.Println(factorial(5))  // 120
```

### 9.3 斐波那契数列

```go
func fibonacci(n int) int {
    if n <= 0 {
        return 0
    }
    if n == 1 {
        return 1
    }
    return fibonacci(n-1) + fibonacci(n-2)
}
```

### 9.4 递归注意事项

| 注意点   | 说明                                   |
| -------- | -------------------------------------- |
| 终止条件 | 必须有明确的基准情况，否则无限递归     |
| 栈溢出   | 深度递归可能导致栈空间耗尽             |
| 性能问题 | 某些问题（如斐波那契）存在大量重复计算 |
| 优化策略 | 使用记忆化（Memoization）或改为迭代    |

---

## 10. defer 延迟执行

`defer` 语句将函数调用推迟到外层函数返回之前执行。

### 10.1 基本用法

```go
func readFile(filename string) {
    fmt.Println("打开文件:", filename)
    defer fmt.Println("关闭文件:", filename)  // 延迟执行

    fmt.Println("读取文件内容...")
    // 即使这里发生错误提前返回，defer 也会执行
}
```

输出：

```
打开文件: config.yaml
读取文件内容...
关闭文件: config.yaml
```

### 10.2 执行顺序（LIFO）

多个 defer 按**后进先出**（栈结构）的顺序执行：

```go
func deferOrder() {
    fmt.Println("开始")
    defer fmt.Println("defer 1")
    defer fmt.Println("defer 2")
    defer fmt.Println("defer 3")
    fmt.Println("结束")
}
```

输出：

```
开始
结束
defer 3
defer 2
defer 1
```

### 10.3 参数求值时机

**defer 语句的参数在 defer 声明时就会求值**，而不是在执行时：

```go
func deferWithValue() {
    x := 10
    defer fmt.Println("defer 时 x =", x)  // x 在此时求值
    x = 20
    fmt.Println("当前 x =", x)
}
```

输出：

```
当前 x = 20
defer 时 x = 10
```

### 10.4 使用匿名函数获取最新值

如果需要在 defer 执行时获取变量的最新值，使用匿名函数：

```go
func deferWithClosure() {
    x := 10
    defer func() {
        fmt.Println("defer 闭包中 x =", x)  // 执行时才读取 x
    }()
    x = 20
    fmt.Println("当前 x =", x)
}
```

输出：

```
当前 x = 20
defer 闭包中 x = 20
```

### 10.5 defer + recover 处理 panic

`recover` 只能在 defer 函数中调用，用于捕获 panic：

```go
func safeDivide(a, b int) (result int) {
    defer func() {
        if r := recover(); r != nil {
            fmt.Println("捕获 panic:", r)
            result = 0  // 发生 panic 时返回默认值
        }
    }()

    if b == 0 {
        panic("除数不能为零")
    }
    return a / b
}

fmt.Println(safeDivide(10, 0))  // 输出: 0（捕获了 panic）
fmt.Println(safeDivide(10, 2))  // 输出: 5（正常执行）
```

### 10.6 defer 最佳实践

| 实践                 | 说明                     |
| -------------------- | ------------------------ |
| 资源获取后立即 defer | 确保资源一定会被释放     |
| 注意 LIFO 顺序       | 多个资源按逆序释放       |
| 了解参数求值时机     | 使用闭包获取最新值       |
| 配合 recover         | 处理 panic，实现优雅恢复 |

---

## 11. init 函数

`init` 是 Go 的特殊函数，在包加载时自动执行。

### 11.1 init 函数特点

```go
func init() {
    // 包初始化逻辑
    // 在 main() 之前自动执行
}
```

| 特点       | 说明                         |
| ---------- | ---------------------------- |
| 自动执行   | 无需调用，包加载时自动运行   |
| 无参无返回 | 签名必须是 `func init()`   |
| 可多次定义 | 每个文件可以有多个 init 函数 |
| 执行顺序   | 按文件名和定义顺序执行       |
| 不可调用   | 不能被显式调用               |

### 11.2 执行顺序

```
1. 导入的包的 init 函数
2. 当前包的包级变量初始化
3. 当前包的 init 函数
4. main 函数
```

### 11.3 常见用途

| 用途           | 示例               |
| -------------- | ------------------ |
| 初始化包级变量 | 复杂的初始化逻辑   |
| 注册驱动       | 数据库驱动注册     |
| 运行时检查     | 验证环境配置       |
| 配置验证       | 检查必需的环境变量 |

### 11.4 使用建议

- ❌ 避免在 init 中执行耗时操作
- ❌ 避免 init 函数间的相互依赖
- ✅ 优先使用显式初始化函数，便于测试
- ✅ 保持 init 函数简单，只做必要的初始化

---

## 12. 最佳实践总结

| 场景     | 推荐做法                                    |
| -------- | ------------------------------------------- |
| 函数命名 | 清晰表达功能，使用驼峰命名法                |
| 参数传递 | 理解值传递机制，需要修改时使用指针          |
| 返回值   | 使用 `(result, error)` 模式处理错误       |
| 可变参数 | 放在参数列表最后，传切片时使用 `...` 展开 |
| 函数类型 | 定义类型别名提高可读性                      |
| 高阶函数 | 用于通用算法（map、filter、reduce）         |
| defer    | 资源获取后立即 defer 释放                   |
| 递归     | 确保有明确的终止条件，注意性能              |
| init     | 保持简单，避免耗时操作                      |

---

## 13. 练习建议

1. **计算器函数**：实现 add、subtract、multiply、divide 四个函数，使用函数类型和 map 实现动态调用
2. **字符串处理器**：使用高阶函数实现 map、filter、reduce 操作处理字符串切片
3. **递归练习**：实现递归版本的二分查找
4. **defer 实践**：模拟文件操作，使用 defer 确保资源释放
5. **工厂函数**：创建一个日志格式化器工厂，生成不同级别（INFO、WARN、ERROR）的日志函数

---

## 参考资料

- [Go 语言规范 - Function declarations](https://go.dev/ref/spec#Function_declarations)
- [Go 语言规范 - Defer statements](https://go.dev/ref/spec#Defer_statements)
- [Effective Go - Functions](https://go.dev/doc/effective_go#functions)
