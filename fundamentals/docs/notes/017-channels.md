# 017. Go 通道（Channel）：协程间通信的桥梁

[返回索引](../README.md) | [查看代码](../../017-channels/main.go)

## 本章要解决的问题

**当多个协程需要交换数据或协调工作时，如何安全、优雅地实现？**

---

## 第一步：场景回顾与新问题

上一章我们学习了协程，可以并发执行多个任务。继续订单通知系统的场景：

```go
customers := []Customer{
    {ID: "C001", Name: "张三", Email: "zhangsan@example.com"},
    {ID: "C002", Name: "李四", Email: "lisi@example.com"},
    {ID: "C003", Name: "王五", Email: "wangwu@example.com"},
}
```

上一章解决了"让通知并发发送"。现在有新需求：

1. **收集结果**：每个通知发送后，需要知道成功还是失败
2. **限制并发**：不能同时发太多请求，否则会压垮邮件服务器
3. **随时取消**：用户点击"取消"时，正在进行的任务能够停下来

这些问题的核心是：**协程之间如何安全地交换数据和协调工作？**

---

## 第二步：共享内存的痛点

上一章我们用 `sync.Mutex` 保护共享变量来收集结果：

```go
var wg sync.WaitGroup
var mu sync.Mutex

// 多个协程要往这个切片里追加结果
results := make([]NotificationResult, 0)

for _, customer := range customers {
    wg.Add(1)
    go func(c Customer) {
        defer wg.Done()
  
        result := sendNotification(c)
  
	// 问题：每次操作共享变量都要加锁
	// 忘记加锁会导致数据竞争，程序行为不可预测
        mu.Lock()
        results = append(results, result)
        mu.Unlock()
    }(customer)
}

wg.Wait()
```

### 这种方式的问题

| 问题     | 说明                                                  |
| -------- | ----------------------------------------------------- |
| 容易出错 | 忘记加锁会导致数据竞争（data race），程序行为不可预测 |
| 代码复杂 | 每个访问共享变量的地方都要加锁、解锁                  |
| 难以扩展 | 当协程数量增多、共享变量变复杂时，锁的管理变得困难    |

### Go 的并发哲学

Go 语言的设计者提出了一个著名的原则：

> **不要通过共享内存来通信，而是通过通信来共享内存。**
> — Go 谚语

这句话的意思是：与其让多个协程直接操作同一块内存（共享内存），不如让协程之间通过"发消息"的方式传递数据（通信）。

**Channel 就是实现这个理念的工具。**

---

## 第三步：Channel 是什么

### 本质

Channel（通道）是 Go 语言提供的一种数据类型，用于协程之间安全地传递数据。

可以把 Channel 想象成一个**管道**：

```
协程 A ──发送数据──> [ Channel ] ──接收数据──> 协程 B
```

- 一个协程往管道里放数据（发送）
- 另一个协程从管道里取数据（接收）
- Channel 内部保证并发安全，无需手动加锁

### 类比

| 类比物 | 对应关系                           |
| ------ | ---------------------------------- |
| 传送带 | Channel 是传送带，数据是货物       |
| 快递柜 | 发件人放入，收件人取走             |
| 电话   | 一方说话（发送），另一方听（接收） |

---

## 第四步：Channel 基础语法

### 创建 Channel

```go
ch := make(chan string)
```

**语法拆解**：

```go
ch := make(chan string)
│     │    │    └── string：通道传输的数据类型
│     │    └── chan：表示这是一个通道类型
│     └── make：内置函数，用于创建 slice、map、channel
└── ch：变量名，指向创建的通道
```

Channel 是引用类型，必须用 `make` 创建，零值是 `nil`（不可用）。

### 发送数据

```go
ch <- "Hello"
```

**语法拆解**：

```go
ch <- "Hello"
│  │   └── 要发送的值
│  └── <- 发送操作符（箭头指向通道，表示数据流入通道）
└── ch：目标通道
```

### 接收数据

```go
msg := <-ch
```

**语法拆解**：

```go
msg := <-ch
│      │  └── 源通道
│      └── <- 接收操作符（箭头从通道出来，表示数据流出通道）
└── msg：接收到的值
```

### 完整示例

```go
func main() {
    // 创建一个传输 string 的通道
    ch := make(chan string)
  
    // 启动一个协程发送数据
    go func() {
        ch <- "Hello, Channel!"  // 发送
    }()
  
    // 主协程接收数据
    msg := <-ch  // 接收
    fmt.Println(msg)  // 输出: Hello, Channel!
}
```

---

## 第五步：无缓冲通道 — 同步通信

### 什么是无缓冲通道

```go
ch := make(chan int)  // 没有第二个参数，创建的是无缓冲通道
```

无缓冲通道的特点：**发送和接收必须同时就绪，否则阻塞**。

- 发送方执行 `ch <- 值` 时，会阻塞等待，直到有接收方执行 `<-ch`
- 接收方执行 `<-ch` 时，会阻塞等待，直到有发送方执行 `ch <- 值`

### 类比：打电话

无缓冲通道就像打电话：

- 双方必须同时在线才能通话
- 打电话的人（发送方）要等接听的人（接收方）接起电话
- 接电话的人要等有人打电话进来

### 示例

```go
ch := make(chan int)

go func() {
    fmt.Println("准备发送...")
    ch <- 42
    // 只有当接收方取走数据后，这行才会执行
    fmt.Println("发送完成")
}()

time.Sleep(1 * time.Second)  // 让发送方先阻塞一会
fmt.Println("准备接收...")
val := <-ch
fmt.Println("收到:", val)
```

**输出顺序**：

```
准备发送...
（1秒后）
准备接收...
收到: 42
发送完成
```

发送方在 `ch <- 42` 处阻塞了 1 秒，直到接收方准备好。

### 同步的意义

无缓冲通道提供了一种**同步机制**：确保发送和接收在某个时间点"握手"。这可以用来协调协程的执行顺序。

---

## 第六步：有缓冲通道 — 异步通信

### 什么是有缓冲通道

```go
ch := make(chan int, 3)  // 第二个参数指定缓冲区大小
```

有缓冲通道的特点：

- 缓冲区未满时，发送不阻塞
- 缓冲区未空时，接收不阻塞
- 缓冲区满了，发送阻塞
- 缓冲区空了，接收阻塞

### 类比：快递柜

有缓冲通道就像快递柜：

- 快递员（发送方）可以把包裹放进柜子，不用等收件人
- 柜子满了，快递员才需要等待
- 收件人（接收方）可以随时来取
- 柜子空了，收件人才需要等待

### 示例

```go
ch := make(chan int, 3)  // 缓冲区容量为 3

// 连续发送 3 个值，不会阻塞
ch <- 1
ch <- 2
ch <- 3
fmt.Println("全部发送完成")

// 查看缓冲区状态
fmt.Printf("缓冲区: %d/%d\n", len(ch), cap(ch))
// 输出: 缓冲区: 3/3

// 依次接收
fmt.Println(<-ch)  // 1
fmt.Println(<-ch)  // 2
fmt.Println(<-ch)  // 3
```

### len 和 cap

```go
len(ch)  // 当前缓冲区中的元素数量
cap(ch)  // 缓冲区的容量
```

### 接收后数据会从通道中移除

**重要概念**：每次从通道接收数据后，该数据就从通道中**移除**了。

```go
ch := make(chan int, 3)
ch <- 1
ch <- 2
ch <- 3
fmt.Printf("接收前: %d/%d\n", len(ch), cap(ch))  // 3/3

val := <-ch  // 接收一个值
fmt.Println("收到:", val)  // 1
fmt.Printf("接收后: %d/%d\n", len(ch), cap(ch))  // 2/3  ← 少了一个
```

通道不是数组，不能随机访问。它是一个**队列**（先进先出），接收操作会把数据"取走"。

### 无缓冲 vs 有缓冲

| 特性     | 无缓冲通道         | 有缓冲通道             |
| -------- | ------------------ | ---------------------- |
| 创建     | `make(chan T)`   | `make(chan T, size)` |
| 发送阻塞 | 立即阻塞，等待接收 | 缓冲区满了才阻塞       |
| 接收阻塞 | 立即阻塞，等待发送 | 缓冲区空了才阻塞       |
| 同步性   | 强同步（握手）     | 弱同步（解耦）         |
| 使用场景 | 需要严格同步       | 允许发送方先行         |

---

## 第七步：关闭 Channel 与遍历

### 关闭 Channel

```go
close(ch)
```

**语法拆解**：

```go
close(ch)
│     └── 要关闭的通道
└── close：内置函数，关闭通道
```

### 关键理解：关闭 ≠ 清空

**关闭通道不会清空缓冲区！**\

```go
ch := make(chan int, 5)  // 有缓冲通道

ch <- 1  // 缓冲区: [1]
ch <- 2  // 缓冲区: [1, 2]
ch <- 3  // 缓冲区: [1, 2, 3]
close(ch)  // 关闭通道，但缓冲区里的数据还在！

fmt.Println(len(ch))  // 输出: 3 ← 数据还在
```

**关闭后的状态**：

| 操作     | 结果                                       |
| -------- | ------------------------------------------ |
| 发送数据 | panic（不能再发送）                        |
| 接收数据 | 先返回缓冲区中的数据，缓冲区空了才返回零值 |
| len(ch)  | 返回缓冲区中剩余的元素数量                 |

### 关键理解：close() 会让所有等待者立即就绪

**`close()` 不只是"释放资源"，它会让所有正在等待 `<-ch` 的操作立即就绪！**

```go
done := make(chan struct{})

go func() {
    <-done  // 阻塞等待...直到 close(done) 被调用
    fmt.Println("收到信号，退出")
}()

close(done)  // 这会让上面的 <-done 立即就绪，返回 struct{}{} 零值
```

**`close(done)` vs `done <- struct{}{}`**：

| 操作                   | 效果                               |
| ---------------------- | ---------------------------------- |
| `done <- struct{}{}` | 只有**一个**接收者能收到     |
| `close(done)`        | **所有**等待的接收者都能收到 |

这就是为什么优雅退出通常用 `close(done)` 而不是 `done <- struct{}{}`：它可以同时通知多个协程退出。

### 检测通道关闭：comma-ok 模式

```go
val, ok := <-ch
```

**语法拆解**：

```go
val, ok := <-ch
│    │     └── 源通道
│    └── ok：布尔值
│          - true：成功收到值（即使通道已关闭，只要缓冲区有数据）
│          - false：通道已关闭 且 缓冲区为空（两个条件都满足）
└── val：收到的值（如果 ok 为 false，val 是零值）
```

**示例**：

```go
ch := make(chan int, 2)
ch <- 1
ch <- 2
close(ch)  // 关闭，但缓冲区里还有 1 和 2

for {
    val, ok := <-ch
    if !ok {
        // 只有当通道已关闭 且 缓冲区为空时，才会走到这里
        fmt.Println("通道已关闭且为空")
        break
    }
    // 即使通道已关闭，只要缓冲区有数据，ok 仍为 true
    fmt.Println("收到:", val, "ok:", ok)
}
```

**输出**：

```
收到: 1 ok: true    ← 通道已关闭，但缓冲区有数据，ok 仍为 true
收到: 2 ok: true    ← 同上
通道已关闭且为空     ← 缓冲区空了，ok 才变成 false
```

### for range 遍历 Channel — 本质是接收操作

**关键理解**：`for range` 遍历通道**本身就是接收操作**，不是像数组那样的读取。

```go
for val := range ch {
    fmt.Println("收到:", val)
}
```

**这段代码完全等价于**：

```go
for {
    val, ok := <-ch      // ← 这里就是接收操作！
    if !ok {             // 通道关闭且为空
        break
    }
    fmt.Println("收到:", val)
}
```

**语法拆解**：

```go
for val := range ch {
│   │      │     └── 要遍历的通道
│   │      └── range：在这里的作用是"持续接收，直到通道关闭"
│   └── val：每次迭代接收到的值（数据被取走了！）
└── for：循环
```

**为什么这样设计**：

- `for range` 对切片是遍历元素
- `for range` 对通道是**持续接收**，直到通道关闭
- 每次循环都是一次 `<-ch` 接收操作，数据从通道中移除

### 关闭的注意事项

| 规则                         | 说明                           |
| ---------------------------- | ------------------------------ |
| 只有发送方应该关闭通道       | 接收方不知道是否还有数据要发送 |
| 关闭已关闭的通道会 panic     | 不要重复关闭                   |
| 向已关闭的通道发送会 panic   | 关闭后不能再发送               |
| 从已关闭的通道接收不会 panic | 返回缓冲区中的值，或零值       |

---

## 第八步：用 Channel 收集协程结果

现在用 Channel 来解决"收集多个协程执行结果"的问题。

### 完整示例

```go
func demoChannelCollectResults(customers []Customer) {
    // 创建结果通道，缓冲区大小 = 客户数量
    resultCh := make(chan NotificationResult, len(customers))

    var wg sync.WaitGroup

    for _, customer := range customers {
        wg.Add(1)
        go func(c Customer) {
            defer wg.Done()
          
            // 模拟发送通知
            time.Sleep(100 * time.Millisecond)
            success := rand.Float32() > 0.3
          
            // 将结果发送到通道（无需加锁！）
            resultCh <- NotificationResult{
                CustomerID:   c.ID,
                CustomerName: c.Name,
                Success:      success,
                Message:      map[bool]string{true: "发送成功", false: "发送失败"}[success],
            }
        }(customer)
    }

    // 另起一个协程：等所有工作协程完成后，关闭通道
    go func() {
        wg.Wait()
        close(resultCh)  // 关闭通道，通知接收方没有更多数据了
    }()

    // 遍历通道收集结果
    // for range 会持续接收，直到 resultCh 被关闭
    successCount := 0
    for result := range resultCh {  // ← 这里就是接收操作！
        if result.Success {
            successCount++
        }
        fmt.Printf("  - %s: %s\n", result.CustomerName, result.Message)
    }
    fmt.Printf("成功: %d/%d\n", successCount, len(customers))
}
```

### 逐行解析关键部分

```go
resultCh := make(chan NotificationResult, len(customers))
```

创建有缓冲通道。缓冲区大小设为 `len(customers)`（比如 3），这样 3 个协程可以同时发送结果，不会阻塞。

```go
resultCh <- NotificationResult{...}
```

每个协程把自己的结果**发送**到通道。无需加锁，通道内部保证并发安全。

```go
go func() {
    wg.Wait()
    close(resultCh)
}()
```

另起一个协程，等待所有工作协程完成后关闭通道。关闭通道会让 `for range` 知道"没有更多数据了，可以结束循环"。

```go
for result := range resultCh {
    // result 是从通道接收到的值
}
```

**这就是接收操作**！`for range` 会不断执行 `<-resultCh`，直到通道关闭且为空。

### 输出示例

```
  - 张三: 发送成功
  - 李四: 发送失败
  - 王五: 发送成功
成功: 2/3
```

`2/3` 表示 3 个客户中有 2 个发送成功。

---

## 第九步：select — 同时等待多个通道

### 痛点：for range 同一时间只能等一个通道

假设你有一个数据通道和一个退出信号通道：

```go
dataCh := make(chan string)
done := make(chan struct{})
```

如果你用 `for range`：

```go
for msg := range dataCh {
    fmt.Println(msg)
}
// 问题：你根本没法响应 done 信号！
// 只有等 dataCh 关闭后，循环才会结束
```

**问题**：`for range` 同一时间只能等**一个**通道。期间如果 `done` 收到退出信号，你收不到。

### select 的核心价值

`select` 可以**同时**等待多个通道，任何一个就绪就响应：

```go
for {
    select {
    case msg := <-dataCh:
        fmt.Println(msg)
    case <-done:
        fmt.Println("收到退出信号，停止")
        return
    }
}
```

**这是 `for range` 做不到的！**

### 完整示例：同时监听数据和退出信号

```go
dataCh := make(chan string)
done := make(chan struct{})

// 生产者：持续发送数据
go func() {
    for i := 1; i <= 5; i++ {
        time.Sleep(100 * time.Millisecond)
        dataCh <- fmt.Sprintf("消息%d", i)
    }
}()

// 250ms 后发送退出信号
go func() {
    time.Sleep(250 * time.Millisecond)
    close(done)
}()

// 消费者：同时监听数据和退出信号
for {
    select {
    case msg := <-dataCh:
        fmt.Println("收到:", msg)
    case <-done:
        fmt.Println("收到退出信号，停止接收")
        return
    }
}
```

**输出**：

```
收到: 消息1
收到: 消息2
收到退出信号，停止接收
```

在 250ms 时收到退出信号，立即停止，不管 dataCh 还有没有数据。

### for range vs select 对比

| 能力                             | for range | select |
| -------------------------------- | --------- | ------ |
| 等待一个通道                     | ✓        | ✓     |
| 同时等待多个通道                 | ✗        | ✓     |
| 同时等待数据和退出信号           | ✗        | ✓     |
| 实现超时控制                     | ✗        | ✓     |
| 非阻塞尝试（有就收，没有就跳过） | ✗        | ✓     |

### select 语法

```go
select {
case val := <-ch1:
    // ch1 可以接收时执行
case val := <-ch2:
    // ch2 可以接收时执行
case ch3 <- data:
    // ch3 可以发送时执行
}
```

**特点**：

- `select` 会阻塞，直到其中一个 case 就绪
- 如果多个 case 同时就绪，随机选择一个执行
- 每次 `select` 只执行一个 case

### select 和 switch 的区别

| 对比项    | switch         | select                           |
| --------- | -------------- | -------------------------------- |
| 作用      | 根据值选择分支 | 根据通道就绪状态选择分支         |
| case 内容 | 普通表达式     | 通道操作（发送或接收）           |
| 执行时机  | 立即判断       | 阻塞等待，哪个通道先就绪就选哪个 |
| 多个匹配  | 执行第一个匹配 | 如果多个同时就绪，随机选一个     |

---

## 第十步：select 实现超时控制

### 痛点：等待超时怎么办？

调用外部服务时，如果对方很慢或卡住了，你不想永远等下去。需要一个"超时"机制。

### time.After 是什么

```go
time.After(500 * time.Millisecond)
```

**语法拆解**：

```go
time.After(500 * time.Millisecond)
│          │         └── Millisecond 是 time 包定义的常量，表示 1 毫秒
│          └── 500 * time.Millisecond = 500 毫秒
└── time.After：返回一个通道，在指定时间后往这个通道发送一个值
```

`time.After` 返回一个 `chan time.Time` 类型的通道。500ms 后，这个通道会收到一个值（当前时间）。

### 结合 select 实现超时

```go
ch := make(chan string)

// 模拟一个很慢的操作（2秒）
go func() {
    time.Sleep(2 * time.Second)
    ch <- "操作完成"
}()

// 只愿意等 500ms
select {
case result := <-ch:
    fmt.Println("收到结果:", result)
case <-time.After(500 * time.Millisecond):
    fmt.Println("超时！不等了")
}
```

**执行流程**：

1. `select` 同时等待两个通道：`ch` 和 `time.After` 返回的通道
2. 500ms 后，`time.After` 的通道先就绪
3. 执行超时分支，程序继续，不会卡住 2 秒

**输出**：

```
超时！不等了
```

---

## 第十一步：select 的 default 分支 — 非阻塞操作

### 痛点：不想等待怎么办？

有时候你想"试一下"，有数据就收，没数据就算了，不要阻塞。

### default 分支

```go
select {
case val := <-ch:
    fmt.Println("收到:", val)
default:
    fmt.Println("没有数据，不等待")
}
```

**执行逻辑**：

- 如果 `ch` 可以接收 → 执行 case 分支
- 如果 `ch` 不能接收（没有数据）→ **立即**执行 default 分支，不阻塞

### 非阻塞发送

```go
ch := make(chan int, 1)
ch <- 1  // 缓冲区已满

select {
case ch <- 2:
    fmt.Println("发送成功")
default:
    fmt.Println("通道已满，跳过")  // ← 会执行这里
}
```

### 非阻塞接收

```go
ch := make(chan int, 1)
// 通道为空

select {
case val := <-ch:
    fmt.Println("收到:", val)
default:
    fmt.Println("通道为空，没有数据")  // ← 会执行这里
}
```

---

## 第十二步：死锁 — 原因与避免

### 什么是死锁

**死锁（Deadlock）**是指程序中的协程互相等待，导致所有协程都无法继续执行的状态。

Go 运行时能够检测到死锁，并报错：

```
fatal error: all goroutines are asleep - deadlock!
```

### 常见死锁场景

**场景 1：无缓冲通道，同一协程发送后接收**

```go
ch := make(chan int)
ch <- 1      // 阻塞！等待接收方
val := <-ch  // 永远执行不到
```

发送操作阻塞，等待接收方；但接收代码在发送之后，永远不会执行。

**场景 2：没有发送方，接收方永远等待**

```go
ch := make(chan int)
val := <-ch  // 阻塞！没有人发送数据
```

**场景 3：没有接收方，发送方永远等待**

```go
ch := make(chan int)
ch <- 1  // 阻塞！没有人接收数据
```

**场景 4：循环等待**

```go
ch1 := make(chan int)
ch2 := make(chan int)

go func() {
    <-ch1    // 等待 ch1 有数据
    ch2 <- 1 // 然后发送到 ch2
}()

<-ch2    // 主协程等待 ch2 有数据
ch1 <- 1 // 然后发送到 ch1
```

主协程在等待 ch2，协程在等待 ch1，互相等待，死锁。

### 如何避免死锁

| 方法                  | 说明                         |
| --------------------- | ---------------------------- |
| 确保发送和接收配对    | 每个发送都应该有对应的接收   |
| 使用有缓冲通道        | 可以暂存数据，缓解同步压力   |
| 使用 select + default | 实现非阻塞操作，避免永久等待 |
| 使用 select + timeout | 设置超时，防止无限等待       |
| 正确关闭通道          | 使用 for range 自动处理关闭  |

---

## 第十三步：实战模式 — 工作池

### 问题场景

有 100 个客户需要发送通知，但邮件服务器最多允许同时 3 个连接。如何限制并发数？

### 工作池模式

工作池（Worker Pool）是一种常见的并发模式：

```
任务队列 tasks: [任务1, 任务2, 任务3, 任务4, 任务5]
                          ↓
                    ┌─────┴─────┐
                    │  Channel  │  ← 内部有锁，保证每个任务只被一个 Worker 拿到
                    └─────┬─────┘
           ┌──────────────┼──────────────┐
           ↓              ↓              ↓
       Worker 1       Worker 2       Worker 3
       拿到任务1       拿到任务2       拿到任务3
           ↓              ↓              ↓
       results ←──────────┴──────────────┘
```

- 固定数量的 Worker 协程
- Worker 从任务通道获取任务
- Worker 将结果发送到结果通道

### 关键问题：多个 Worker 同时接收，会不会数据竞争？

**不会！** 这是 Channel 的核心设计。

**Channel 的接收操作是原子的**（内部有锁保护）：

- 当多个协程同时从一个通道接收时，**每条数据只会被一个协程收到**
- Go 运行时保证"一个数据只发给一个接收方"
- 不需要额外加锁

**类比**：取号机。3 个窗口同时等号，每次只有一个窗口能拿到号，不会出现两个窗口拿到同一个号。

这正是 Channel 比共享内存更安全的原因：你不需要自己加锁，Channel 内部已经处理好了。

### 实现代码

```go
// Task 表示一个任务
type Task struct {
    ID       int
    Customer Customer
}

// worker 是工作池中的一个工作者
func worker(id int, tasks <-chan Task, results chan<- NotificationResult, wg *sync.WaitGroup) {
    defer wg.Done()
  
    // 持续从 tasks 通道获取任务
    for task := range tasks {
        fmt.Printf("[Worker %d] 处理任务 %d\n", id, task.ID)
      
        // 处理任务...
        result := NotificationResult{
            CustomerID:   task.Customer.ID,
            CustomerName: task.Customer.Name,
            Success:      true,
        }
      
        // 发送结果
        results <- result
    }
}
```

**语法拆解**：

```go
func worker(id int, tasks <-chan Task, results chan<- NotificationResult, wg *sync.WaitGroup)
                     │      │           │       │
                     │      │           │       └── chan<- 只写通道（只能发送）
                     │      │           └── chan<- 只写通道
                     │      └── <-chan 只读通道（只能接收）
                     └── 通道的元素类型
```

**单向通道**：

| 类型         | 说明                   |
| ------------ | ---------------------- |
| `chan T`   | 双向通道，可发送可接收 |
| `<-chan T` | 只读通道，只能接收     |
| `chan<- T` | 只写通道，只能发送     |

单向通道用于限制函数对通道的操作权限，提高代码安全性。

### 启动工作池

```go
func main() {
    tasks := make(chan Task, 100)
    results := make(chan NotificationResult, 100)
  
    var wg sync.WaitGroup
  
    // 启动 3 个 Worker（限制并发度为 3）
    for i := 1; i <= 3; i++ {
        wg.Add(1)
        go worker(i, tasks, results, &wg)
    }
  
    // 发送所有任务
    for i, customer := range customers {
        tasks <- Task{ID: i + 1, Customer: customer}
    }
    close(tasks)  // 关闭任务通道，通知 Worker 没有更多任务
  
    // 等待所有 Worker 完成后关闭结果通道
    go func() {
        wg.Wait()
        close(results)
    }()
  
    // 收集结果
    for result := range results {
        fmt.Printf("%s: %v\n", result.CustomerName, result.Success)
    }
}
```

### 工作池的优点

| 优点     | 说明                                  |
| -------- | ------------------------------------- |
| 限制并发 | 固定数量的 Worker 限制了最大并发数    |
| 复用资源 | Worker 可以复用，避免频繁创建销毁协程 |
| 解耦     | 任务生产者和消费者通过通道解耦        |

---

## 第十四步：优雅退出

### 问题场景

有一个持续运行的工作协程，用户想要随时能停止它。

### 使用 Channel 发送退出信号

```go
// done 通道用于发送退出信号
done := make(chan struct{})

go func() {
    for {
        select {
        case <-done:
            // 收到退出信号
            fmt.Println("正在退出...")
            return
        default:
            // 正常工作
            doWork()
        }
    }
}()

// 需要停止时，关闭 done 通道
close(done)
```

**为什么用 `chan struct{}`**：

```go
done := make(chan struct{})
```

- `struct{}` 是空结构体，不占用内存
- 用于只传递信号、不传递数据的场景
- 关闭通道时，所有等待接收的协程都会收到通知

---

## 最佳实践

### 1. 通道所有权

**原则：发送方关闭通道，接收方不要关闭。**

为什么？因为只有发送方知道"什么时候不会再发送了"。接收方不知道还有没有数据要来。

```go
func main() {
    // producer 创建通道并返回
    ch := producer()
    
    // 接收方使用这个通道
    for val := range ch {
        fmt.Println(val)
    }
}

// producer：创建通道 + 发送数据 + 关闭通道
func producer() chan int {
    ch := make(chan int)
    go func() {
        for i := 0; i < 5; i++ {
            ch <- i
        }
        close(ch)  // 发送方关闭，因为只有发送方知道"发完了"
    }()
    return ch
}
```

**执行流程**：
1. `main` 调用 `producer()`，得到通道 `ch`
2. `producer` 内部启动协程，往 `ch` 发送 0, 1, 2, 3, 4
3. 发送完后，`producer` 内部的协程执行 `close(ch)`
4. `main` 中的 `for range ch` 收到所有数据后，因为通道关闭，自动退出循环

**总结**：
- 发送方关闭（发送方知道"发完了"）
- 接收方不关闭（接收方不知道"还有没有"）

### 2. 防止 goroutine 泄漏

确保每个协程都有退出条件：

```go
// 错误：协程可能永远阻塞
go func() {
    val := <-ch  // 如果没人发送，永远阻塞
    process(val)
}()

// 正确：使用 select 和 done 通道
go func() {
    select {
    case val := <-ch:
        process(val)
    case <-done:
        return
    }
}()
```

### 3. 有缓冲 vs 无缓冲的选择

| 场景                   | 推荐                    |
| ---------------------- | ----------------------- |
| 需要严格同步           | 无缓冲                  |
| 发送方不需要等待接收方 | 有缓冲                  |
| 收集固定数量的结果     | 有缓冲（容量=结果数）   |
| 限流、控制并发         | 有缓冲（容量=并发限制） |

### 4. 使用 for range 遍历

比 comma-ok 模式更简洁，自动处理通道关闭：

```go
// 推荐
for val := range ch {
    process(val)
}

// 不推荐
for {
    val, ok := <-ch
    if !ok {
        break
    }
    process(val)
}
```

---

## 总结

| 概念           | 语法/方法                   | 说明                   |
| -------------- | --------------------------- | ---------------------- |
| 创建通道       | `make(chan T)`            | 无缓冲通道             |
| 创建有缓冲通道 | `make(chan T, size)`      | 指定缓冲区大小         |
| 发送           | `ch <- 值`                | 箭头指向通道           |
| 接收           | `值 := <-ch`              | 箭头从通道出来         |
| 关闭           | `close(ch)`               | 只有发送方应该关闭     |
| 检测关闭       | `val, ok := <-ch`         | ok 为 false 表示已关闭 |
| 遍历           | `for val := range ch`     | 自动处理关闭           |
| 多路复用       | `select { case ... }`     | 同时监听多个通道       |
| 超时           | `case <-time.After(d)`    | 配合 select 使用       |
| 非阻塞         | `select { ... default: }` | default 分支实现非阻塞 |
| 只读通道       | `<-chan T`                | 只能接收               |
| 只写通道       | `chan<- T`                | 只能发送               |

### Channel 的核心特点

1. **并发安全**：内部处理好同步，无需手动加锁
2. **阻塞语义**：发送和接收可能阻塞，用于协程同步
3. **类型安全**：通道只能传输指定类型的数据
4. **可关闭**：关闭通道通知接收方没有更多数据

### Go 并发哲学

> **通过通信来共享内存，而不是通过共享内存来通信。**

Channel 是 Go 实现这一理念的核心工具。
