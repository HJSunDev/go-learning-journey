# 016. Go 协程（Goroutine）：轻量级并发的基石

[返回索引](../README.md) | [查看代码](../../016-goroutines/main.go)

## 本章要解决的问题

**当有多个独立任务需要执行时，如何避免顺序执行导致的时间浪费？**

---

## 第一步：场景设定

你在开发订单系统的通知模块。每当有订单状态变更，需要通知相关客户。

现在有 3 个客户需要通知，每个通知需要 500ms（模拟网络延迟）：

```go
customers := []Customer{
    {ID: "C001", Name: "张三", Email: "zhangsan@example.com"},
    {ID: "C002", Name: "李四", Email: "lisi@example.com"},
    {ID: "C003", Name: "王五", Email: "wangwu@example.com"},
}
```

---

## 第二步：没有协程时的痛点

按顺序逐个发送通知：

```go
func sendNotification(customer Customer) {
    fmt.Printf("开始发送通知给 %s...\n", customer.Name)
    time.Sleep(500 * time.Millisecond)  // 模拟网络延迟
    fmt.Printf("通知已发送给 %s\n", customer.Name)
}

// 顺序执行
for _, customer := range customers {
    sendNotification(customer)
}
```

**执行结果**：

```
开始发送通知给 张三...
通知已发送给 张三
开始发送通知给 李四...
通知已发送给 李四
开始发送通知给 王五...
通知已发送给 王五
总耗时: 1.5秒
```

**问题**：3 个独立的任务，却必须等上一个完成才能开始下一个。如果有 100 个客户，就需要 50 秒！

---

## 第三步：协程是什么

### 本质

协程是 **用户态（运行时 Runtime）管理的，自带栈和寄存器状态的，可暂停/恢复的函数执行流**。

用更直白的话说：

| 概念        | 说明                                                |
| ----------- | --------------------------------------------------- |
| 用户态管理  | Go 运行时自己调度，不需要操作系统介入，切换成本极低 |
| 自带栈      | 每个协程有独立的执行栈，初始仅 2KB，按需增长        |
| 可暂停/恢复 | 遇到 I/O 等待时自动暂停，让出 CPU 给其他协程        |

### 协程 vs 线程

| 对比项   | 线程（Thread）       | 协程（Goroutine）     |
| -------- | -------------------- | --------------------- |
| 管理者   | 操作系统内核         | Go 运行时             |
| 创建成本 | 高（约 1MB 栈空间）  | 极低（约 2KB 栈空间） |
| 切换成本 | 高（需要内核态切换） | 极低（用户态切换）    |
| 数量上限 | 通常几千个           | 可轻松创建几十万个    |
| 调度方式 | 抢占式               | 协作式 + 抢占式       |

### 为什么选择协程

传统多线程的问题：

1. 线程创建/销毁开销大
2. 线程切换需要陷入内核态，成本高
3. 线程数量受限，无法为每个请求创建一个线程

协程的优势：

1. 创建开销极小，可以为每个任务创建一个协程
2. 切换在用户态完成，不需要内核介入
3. Go 运行时会自动将协程分配到多个 CPU 核心上执行

---

## 第四步：go 关键字 —— 启动协程

### 语法

```go
go 函数调用
```

`go` 关键字后面跟一个函数调用，该函数会在新的协程中执行。

### 示例

```go
for _, customer := range customers {
    go sendNotification(customer)
}
```

### 语法拆解

```go
go sendNotification(customer)
│  └── 要在协程中执行的函数调用
└── go 关键字：启动一个新协程
```

**关键点**：

1. `go` 语句会立即返回，不会等待函数执行完成
2. 函数在独立的协程中异步执行
3. 主协程（main 函数）和新协程并发运行

### 问题：主程序提前退出

```go
func main() {
    for _, customer := range customers {
        go sendNotification(customer)
    }
    fmt.Println("主程序结束")
    // 主程序在这里结束，所有协程被强制终止！
}
```

**执行结果**：

```
主程序结束
```

通知根本没发出去！因为 `go` 只是"安排"任务，主程序不会等待协程完成。当 `main` 函数返回时，所有协程都被终止。

---

## 第五步：WaitGroup —— 等待所有协程完成

### sync.WaitGroup 是什么

`sync.WaitGroup` 是一个计数器，用于等待一组协程全部完成。

| 方法       | 作用                              |
| ---------- | --------------------------------- |
| `Add(n)` | 计数器加 n，表示有 n 个协程要等待 |
| `Done()` | 计数器减 1，表示一个协程已完成    |
| `Wait()` | 阻塞，直到计数器变为 0            |

### 使用流程

```
1. 创建 WaitGroup
2. 启动协程前调用 Add(1)
3. 协程结束时调用 Done()
4. 主程序调用 Wait() 等待所有协程
```

### 完整示例

```go
func sendNotificationWithWG(customer Customer, wg *sync.WaitGroup) {
    // defer 确保函数结束时一定会调用 Done()
    defer wg.Done()
  
    fmt.Printf("开始发送通知给 %s...\n", customer.Name)
    time.Sleep(500 * time.Millisecond)
    fmt.Printf("通知已发送给 %s\n", customer.Name)
}

func main() {
    customers := []Customer{...}
  
    // 创建 WaitGroup
    var wg sync.WaitGroup
  
    for _, customer := range customers {
        // 启动协程前，计数器加 1
        wg.Add(1)
  
        // 启动协程，传入 wg 的指针
        go sendNotificationWithWG(customer, &wg)
    }
  
    // Wait() 会阻塞，直到计数器变为 0
    // 每个协程完成时调用 Done()，计数器减 1
    wg.Wait()
  
    fmt.Println("所有通知发送完成")
}
```

### 语法逐行拆解

```go
var wg sync.WaitGroup
```

- `sync` 是 Go 标准库的同步包
- `WaitGroup` 是该包中的类型
- `var wg sync.WaitGroup` 声明一个 WaitGroup 变量，零值可用

```go
wg.Add(1)
```

- 告诉 WaitGroup：即将启动一个新协程
- 必须在 `go` 语句之前调用，否则可能出现竞态条件

```go
go sendNotificationWithWG(customer, &wg)
```

- `&wg` 是 WaitGroup 的指针
- 必须传指针，否则协程内操作的是副本，对原始 wg 无效

```go
defer wg.Done()
```

- `defer` 确保函数结束时执行 `Done()`
- 即使函数发生 panic，`defer` 的代码也会执行
- `Done()` 让计数器减 1

```go
wg.Wait()
```

- 阻塞当前协程，直到计数器变为 0
- 所有协程都调用 `Done()` 后，`Wait()` 返回

### 执行结果

```
开始发送通知给 张三...
开始发送通知给 李四...
开始发送通知给 王五...
通知已发送给 张三
通知已发送给 李四
通知已发送给 王五
总耗时: 500ms
```

三个通知并发执行，总耗时从 1.5 秒降到了 500 毫秒！

---

## 第六步：闭包变量捕获（经典陷阱与 Go 1.22 修复）

### 背景：Go 1.22 的重大变化

Go 1.22 改变了 `for` 循环变量的语义：

| Go 版本 | 循环变量行为 |
|---------|-------------|
| Go 1.21 及之前 | 循环变量在循环开始前声明一次，所有迭代**共享同一个变量** |
| Go 1.22+ | 每次迭代都会**创建新的变量实例** |

**触发条件**：go.mod 文件中声明 `go 1.22` 或更高版本。

### 经典陷阱代码

```go
for _, customer := range customers {
    wg.Add(1)
  
    go func() {
        defer wg.Done()
        fmt.Printf("发送给: %s\n", customer.Name)
    }()
}
```

**Go 1.21 及之前的执行结果**：

```
发送给: 王五
发送给: 王五
发送给: 王五
```

所有协程都打印了"王五"！

**Go 1.22+ 的执行结果**：

```
发送给: 张三
发送给: 李四
发送给: 王五
```

结果正确，因为每次迭代的 `customer` 是独立的变量。

### 问题分析（Go 1.21 及之前）

```go
for _, customer := range customers {
    //     └── customer 在循环开始前声明一次，所有迭代共享
  
    go func() {
        // 闭包捕获的是变量的引用，不是值的副本
        fmt.Printf("发送给: %s\n", customer.Name)
    }()
    // 协程启动快，但执行有延迟
    // 当协程真正执行时，循环可能已经结束
}
```

时间线（Go 1.21 及之前）：

1. 循环第 1 次：customer = 张三，启动协程（还没执行）
2. 循环第 2 次：customer = 李四，启动协程（还没执行）
3. 循环第 3 次：customer = 王五，启动协程（还没执行）
4. 循环结束，此时 customer = 王五
5. 三个协程开始执行，都读取到 customer = 王五

### 传统修复方法：参数传递

```go
for _, customer := range customers {
    wg.Add(1)
  
    // 将 customer 作为参数传入匿名函数
    go func(c Customer) {
        defer wg.Done()
        fmt.Printf("发送给: %s\n", c.Name)
    }(customer)  // 这里传入当前的 customer 值
}
```

**为什么有效**：

- 函数参数是值传递，`customer` 的值被复制给 `c`
- 每个协程有自己独立的 `c` 副本

### 为什么仍然推荐显式传参

虽然 Go 1.22+ 已自动修复此问题，但仍建议使用参数传递的方式：

1. **代码兼容性**：代码可能需要在旧版本 Go 上运行
2. **意图清晰**：明确表达"每个协程使用独立的值"
3. **代码可读性**：其他开发者（尤其是熟悉旧版本的）更容易理解

---

## 第七步：最佳实践

### 1. WaitGroup 的 Add 必须在 go 之前

```go
// 正确 ✓
wg.Add(1)
go func() {
    defer wg.Done()
    // ...
}()

// 错误 ✗ —— 可能在 Add 之前就调用了 Wait
go func() {
    wg.Add(1)  // 错误位置
    defer wg.Done()
    // ...
}()
```

### 2. 始终使用 defer wg.Done()

```go
go func() {
    defer wg.Done()  // 始终放在第一行
  
    // 即使这里发生 panic，Done() 也会被调用
    // 避免 WaitGroup 永远等待
    doSomething()
}()
```

### 3. 避免闭包陷阱

```go
// 推荐：将循环变量作为参数传入
for _, item := range items {
    wg.Add(1)
    go func(i Item) {
        defer wg.Done()
        process(i)
    }(item)
}
```

### 4. 保护共享变量

当多个协程需要修改同一个变量时，必须使用互斥锁：

```go
var (
    count int
    mu    sync.Mutex  // 互斥锁
)

for _, item := range items {
    wg.Add(1)
    go func(i Item) {
        defer wg.Done()
      
        result := process(i)
      
        // 修改共享变量前加锁
        mu.Lock()
        count += result
        mu.Unlock()
    }(item)
}
```

### sync.Mutex 语法说明

```go
var mu sync.Mutex
```

- `Mutex` 是互斥锁（Mutual Exclusion）的缩写
- 零值可用，不需要初始化

```go
mu.Lock()
```

- 获取锁。如果锁已被其他协程持有，当前协程会阻塞等待

```go
mu.Unlock()
```

- 释放锁。其他等待的协程可以继续

### 5. 协程生命周期管理

协程泄漏是常见问题：启动的协程永远不会结束，持续消耗资源。

**常见原因**：

- 协程在等待一个永远不会到来的事件
- 无限循环没有退出条件

**下一章会讲解如何使用 Channel 来控制协程的退出。**

---

## 总结

| 概念           | 语法/方法                       | 说明                         |
| -------------- | ------------------------------- | ---------------------------- |
| 启动协程       | `go 函数调用`                 | 在新协程中异步执行函数       |
| 创建 WaitGroup | `var wg sync.WaitGroup`       | 用于等待一组协程完成         |
| 增加计数       | `wg.Add(n)`                   | 在启动协程前调用             |
| 标记完成       | `wg.Done()`                   | 协程结束时调用，建议用 defer |
| 等待完成       | `wg.Wait()`                   | 阻塞直到计数器为 0           |
| 互斥锁         | `sync.Mutex`                  | 保护共享变量                 |
| 加锁/解锁      | `mu.Lock()` / `mu.Unlock()` | 确保同一时刻只有一个协程访问 |

### 协程的核心特点

1. **轻量**：初始栈仅 2KB，可创建数十万个
2. **快速切换**：用户态调度，无需内核介入
3. **简单语法**：`go` 关键字即可启动
4. **自动调度**：Go 运行时自动分配到多个 CPU 核心
