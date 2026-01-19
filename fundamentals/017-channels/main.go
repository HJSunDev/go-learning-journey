package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

// =============================================================================
// 017. Go 通道（Channel）：协程间通信的桥梁
// =============================================================================
// 场景：订单系统 - 批量发送客户通知（续）
// 上一章解决了"如何并发执行"，本章解决"如何安全地交换数据和协调工作"
// =============================================================================

// Customer 表示一个客户
type Customer struct {
	ID    string
	Name  string
	Email string
}

// NotificationResult 表示单个通知的发送结果
type NotificationResult struct {
	CustomerID   string
	CustomerName string
	Success      bool
	Message      string
}

// =============================================================================
// 第一部分：痛点 —— 共享内存的问题
// =============================================================================

// demoSharedMemoryProblem 演示共享内存方式收集结果的问题
// 问题：代码复杂，需要手动管理锁，容易出错
func demoSharedMemoryProblem(customers []Customer) {
	fmt.Println("========== 共享内存方式（上一章的方法）==========")

	var wg sync.WaitGroup
	var mu sync.Mutex

	// 共享的结果切片
	results := make([]NotificationResult, 0, len(customers))

	for _, customer := range customers {
		wg.Add(1)
		go func(c Customer) {
			defer wg.Done()

			// 模拟发送通知
			time.Sleep(100 * time.Millisecond)
			result := NotificationResult{
				CustomerID:   c.ID,
				CustomerName: c.Name,
				Success:      true,
				Message:      "发送成功",
			}

			// 问题：每次操作共享变量都要加锁
			// 忘记加锁会导致数据竞争，程序行为不可预测
			mu.Lock()
			results = append(results, result)
			mu.Unlock()
		}(customer)
	}

	wg.Wait()

	fmt.Printf("共享内存方式收集到 %d 个结果\n", len(results))
	for _, r := range results {
		fmt.Printf("  - %s: %s\n", r.CustomerName, r.Message)
	}
	fmt.Println()
}

// =============================================================================
// 第二部分：Channel 基础 —— 创建、发送、接收
// =============================================================================

// demoChannelBasic 演示 Channel 的基本使用
func demoChannelBasic() {
	fmt.Println("========== Channel 基础演示 ==========")

	// 创建一个无缓冲的 string 类型通道
	// make(chan 元素类型) 创建通道
	ch := make(chan string)

	// 启动一个协程发送数据
	go func() {
		fmt.Println("[发送方] 准备发送消息...")
		// ch <- 值：向通道发送数据
		// 无缓冲通道：发送操作会阻塞，直到有接收方准备好接收
		ch <- "Hello, Channel!"
		fmt.Println("[发送方] 消息已发送")
	}()

	// 模拟接收方稍后准备好
	time.Sleep(200 * time.Millisecond)
	fmt.Println("[接收方] 准备接收消息...")

	// <- ch：从通道接收数据
	// 无缓冲通道：接收操作会阻塞，直到有发送方发送数据
	msg := <-ch
	fmt.Printf("[接收方] 收到消息: %s\n\n", msg)
}

// =============================================================================
// 第三部分：无缓冲通道 —— 同步通信
// =============================================================================

// demoUnbufferedChannel 演示无缓冲通道的同步特性
func demoUnbufferedChannel() {
	fmt.Println("========== 无缓冲通道（同步）==========")

	// 无缓冲通道：发送和接收必须同时就绪，否则阻塞
	// 类比：打电话 —— 双方必须同时在线才能通话
	ch := make(chan int)

	go func() {
		for i := 1; i <= 3; i++ {
			fmt.Printf("[发送] 准备发送 %d...\n", i)
			ch <- i
			// 只有当接收方取走数据后，这行才会执行
			fmt.Printf("[发送] %d 已被接收\n", i)
		}
	}()

	// 接收方慢一点，观察发送方的阻塞行为
	for i := 1; i <= 3; i++ {
		time.Sleep(300 * time.Millisecond)
		val := <-ch
		fmt.Printf("[接收] 收到 %d\n", val)
	}
	fmt.Println()
}

// =============================================================================
// 第四部分：有缓冲通道 —— 异步通信
// =============================================================================

// demoBufferedChannel 演示有缓冲通道的异步特性
func demoBufferedChannel() {
	fmt.Println("========== 有缓冲通道（异步）==========")

	// 有缓冲通道：缓冲区未满时发送不阻塞，缓冲区未空时接收不阻塞
	// make(chan 类型, 缓冲区大小)
	// 类比：快递柜 —— 柜子未满时可以直接放入，不用等收件人
	ch := make(chan int, 3)

	// 发送 3 个值，不会阻塞（因为缓冲区容量为 3）
	fmt.Println("[发送] 连续发送 3 个值...")
	ch <- 1
	ch <- 2
	ch <- 3
	fmt.Println("[发送] 全部发送完成，没有阻塞")

	// 查看当前缓冲区状态
	// len(ch)：缓冲区中的元素数量
	// cap(ch)：缓冲区容量
	fmt.Printf("[状态] 缓冲区: %d/%d\n", len(ch), cap(ch))

	// 依次接收
	fmt.Printf("[接收] %d\n", <-ch)
	fmt.Printf("[接收] %d\n", <-ch)
	fmt.Printf("[接收] %d\n", <-ch)
	fmt.Printf("[状态] 缓冲区: %d/%d\n\n", len(ch), cap(ch))
}

// =============================================================================
// 第五部分：Channel 收集协程结果（替代共享内存）
// =============================================================================

// demoChannelCollectResults 演示使用 Channel 收集协程的执行结果
// 对比第一部分的共享内存方式，这种方式更安全、更清晰
func demoChannelCollectResults(customers []Customer) {
	fmt.Println("========== Channel 收集结果 ==========")

	// 创建结果通道
	// 使用有缓冲通道，缓冲区大小等于客户数量
	// 发送方不需要等待接收方，可以立即返回
	resultCh := make(chan NotificationResult, len(customers))

	var wg sync.WaitGroup

	for _, customer := range customers {
		wg.Add(1)
		go func(c Customer) {
			defer wg.Done()

			// 模拟发送通知
			time.Sleep(100 * time.Millisecond)
			// rand 是 Go 标准库中的伪随机数生成器
			// Float32() 是 rand 包的函数，返回一个 [0.0, 1.0) 区间内的随机 float32 类型的数
			// float 表示浮点类型，float32 占用 4 字节 —— 通常用于数学概率计算
			// 因此下面一行表示以 70% 概率“发送成功”，30% 概率“发送失败”
			success := rand.Float32() > 0.3

			// 将结果发送到通道，无需加锁
			// Channel 内部已经处理好并发安全
			resultCh <- NotificationResult{
				CustomerID:   c.ID,
				CustomerName: c.Name,
				Success:      success,
				Message:      map[bool]string{true: "发送成功", false: "发送失败"}[success],
			}
		}(customer)
	}

	// 等待所有协程完成后关闭通道
	// 关闭通道会通知接收方：不会再有新数据了
	go func() {
		wg.Wait()
		close(resultCh)
	}()

	// 使用 for range 遍历通道
	// 重要：for range 遍历通道本身就是接收操作！
	// 它等价于：for { result, ok := <-resultCh; if !ok { break }; ... }
	// 每次循环都会从通道中"取走"一个值（不是像数组那样只读取）
	// 通道关闭且为空后，for range 会自动退出
	successCount := 0
	for result := range resultCh {
		if result.Success {
			successCount++
		}
		fmt.Printf("  - %s: %s\n", result.CustomerName, result.Message)
	}
	// 输出格式：成功数/总数，例如 "成功: 2/3" 表示 3 个客户中 2 个发送成功
	fmt.Printf("成功: %d/%d\n\n", successCount, len(customers))
}

// =============================================================================
// 第六部分：关闭 Channel 与检测关闭
// =============================================================================

// demoCloseChannel 演示关闭通道和检测通道关闭
func demoCloseChannel() {
	fmt.Println("========== 关闭 Channel ==========")

	ch := make(chan int, 5)

	// 发送一些数据后关闭
	ch <- 1
	ch <- 2
	ch <- 3
	// 关键理解：close(ch) 只是标记通道为"已关闭"，不会清空缓冲区！
	// 关闭后：
	//   - 不能再发送（会 panic）
	//   - 但缓冲区里的 1, 2, 3 还在，可以继续接收
	close(ch)
	fmt.Printf("关闭后，缓冲区仍有数据: %d 个\n", len(ch))

	// 方式一：使用 comma-ok 模式检测通道是否关闭
	// val, ok := <-ch
	// ok 为 true：成功收到值（即使通道已关闭，只要缓冲区有数据）
	// ok 为 false：通道已关闭 且 缓冲区为空（两个条件都满足才返回 false）
	fmt.Println("方式一：comma-ok 模式")
	for {
		val, ok := <-ch
		if !ok {
			// 只有当通道已关闭 且 缓冲区为空时，才会走到这里
			fmt.Println("  通道已关闭且缓冲区为空")
			break
		}
		// 即使通道已关闭，只要缓冲区有数据，ok 仍为 true
		fmt.Printf("  收到: %d (ok=%v)\n", val, ok)
	}

	// 方式二：for range 自动处理关闭
	ch2 := make(chan int, 3)
	ch2 <- 10
	ch2 <- 20
	close(ch2)

	fmt.Println("方式二：for range")
	for val := range ch2 {
		fmt.Printf("  收到: %d\n", val)
	}
	fmt.Println("  循环结束（通道已关闭）")
}

// =============================================================================
// 第七部分：select — 同时等待多个通道
// =============================================================================
// 核心价值：for range 在一个时刻只能等待一个通道，select 可以同时等待多个通道
//
// 用 for range？做不到同时监听！
//   for val := range ch1 { ... }  // 只能等 ch1，期间 ch2 来数据你收不到
//   for val := range ch2 { ... }  // 必须等 ch1 关闭后才能开始等 ch2
//
// 用 select？可以同时等待多个通道！
//   select {
//   case val := <-ch1: ...  // 同时等
//   case val := <-ch2: ...  // 同时等
//   case <-done: return     // 还能同时等退出信号
//   }
// =============================================================================

// demoSelectVsRange 演示 select 的不可替代性
func demoSelectVsRange() {
	fmt.Println("========== select vs for range ==========")
	fmt.Println("场景：同时监听「数据通道」和「退出信号」")
	fmt.Println("for range 做不到，select 可以")
	fmt.Println()

	dataCh := make(chan string)
	done := make(chan struct{})

	// 生产者：持续发送数据
	go func() {
		messages := []string{"消息1", "消息2", "消息3", "消息4", "消息5"}
		for _, msg := range messages {
			time.Sleep(100 * time.Millisecond)
			dataCh <- msg
		}
	}()

	// 250ms 后发送退出信号
	go func() {
		time.Sleep(250 * time.Millisecond)
		fmt.Println("  [信号] 发送退出命令")
		close(done)
	}()

	// 消费者：同时监听数据和退出信号
	// 如果用 for range dataCh { ... }，你就无法响应 done 信号！
	// 只有用 select 才能同时等待两个通道
	fmt.Println("  开始接收（select 同时等待 dataCh 和 done）:")
	for {
		select {
		case msg := <-dataCh:
			fmt.Printf("  [数据] 收到: %s\n", msg)
		case <-done:
			// 收到退出信号，立即停止，不管 dataCh 还有没有数据
			fmt.Println("  [退出] 收到退出信号，停止接收")
			fmt.Println()
			return
		}
	}
}

// demoSelect 演示 select 同时等待多个数据通道
func demoSelect() {
	fmt.Println("========== select 多路复用 ==========")
	fmt.Println("场景：两个数据源，谁先来处理谁")
	fmt.Println()

	ch1 := make(chan string)
	ch2 := make(chan string)

	// 协程1：500ms 后发送
	// time.Millisecond 是 time 包定义的常量，表示 1 毫秒
	// 500 * time.Millisecond = 500 毫秒 = 0.5 秒
	go func() {
		time.Sleep(500 * time.Millisecond)
		ch1 <- "来自通道1的消息"
	}()

	// 协程2：300ms 后发送（比协程1快）
	go func() {
		time.Sleep(300 * time.Millisecond)
		ch2 <- "来自通道2的消息"
	}()

	// select 同时等待两个通道，谁先就绪就处理谁
	// 如果用 for range，你必须先等完 ch1，再等 ch2，无法交替处理
	for i := 0; i < 2; i++ {
		select {
		// case msg := <-ch1 表示：等待从 ch1 接收数据
		case msg := <-ch1:
			fmt.Printf("  收到: %s\n", msg)
		// case msg := <-ch2 表示：等待从 ch2 接收数据
		case msg := <-ch2:
			fmt.Printf("  收到: %s\n", msg)
		}
		// 第一次循环：ch2 在 300ms 后就绪，先执行 ch2 的 case
		// 第二次循环：ch1 在 500ms 后就绪，执行 ch1 的 case
	}
	fmt.Println()
}

// demoSelectTimeout 演示 select 实现超时控制
func demoSelectTimeout() {
	fmt.Println("========== select 超时控制 ==========")

	ch := make(chan string)

	// 模拟一个需要很长时间的操作
	go func() {
		time.Sleep(2 * time.Second)
		ch <- "操作完成"
	}()

	// 使用 select + time.After 实现超时
	// time.After 返回一个通道，在指定时间后发送当前时间
	select {
	case result := <-ch:
		fmt.Printf("  收到结果: %s\n", result)
	case <-time.After(500 * time.Millisecond):
		fmt.Println("  操作超时！")
	}
	fmt.Println()
}

// demoSelectDefault 演示 select 的 default 分支（非阻塞操作）
func demoSelectDefault() {
	fmt.Println("========== select default（非阻塞）==========")

	ch := make(chan int, 1)

	// 尝试发送，如果通道满了就跳过
	ch <- 1

	select {
	case ch <- 2:
		fmt.Println("  成功发送 2")
	default:
		// 通道已满，立即执行 default
		fmt.Println("  通道已满，跳过发送")
	}

	// 尝试接收，如果通道空了就跳过
	select {
	case val := <-ch:
		fmt.Printf("  收到: %d\n", val)
	default:
		fmt.Println("  通道为空，没有数据")
	}

	select {
	case val := <-ch:
		fmt.Printf("  收到: %d\n", val)
	default:
		fmt.Println("  通道为空，没有数据")
	}
	fmt.Println()
}

// =============================================================================
// 第八部分：死锁演示与避免
// =============================================================================

// 注意：以下死锁示例被注释掉，取消注释运行会导致程序崩溃

// demoDeadlockExamples 解释常见的死锁场景（不实际执行）
func demoDeadlockExamples() {
	fmt.Println("========== 死锁场景说明（不执行）==========")

	fmt.Println(`死锁场景 1：无缓冲通道，同一协程中发送后接收
--------------------------------------------
ch := make(chan int)
ch <- 1    // 阻塞！等待接收方
val := <-ch // 永远执行不到
原因：发送操作阻塞，等待接收方；但接收代码在发送之后，永远不会执行

死锁场景 2：没有发送方，接收方永远等待
--------------------------------------------
ch := make(chan int)
val := <-ch  // 阻塞！永远等待
原因：没有协程向通道发送数据

死锁场景 3：没有接收方，发送方永远等待（无缓冲通道）
--------------------------------------------
ch := make(chan int)
ch <- 1  // 阻塞！永远等待
原因：没有协程从通道接收数据

死锁场景 4：循环等待
--------------------------------------------
ch1 := make(chan int)
ch2 := make(chan int)
go func() {
    <-ch1    // 等待 ch1
    ch2 <- 1 // 然后发送 ch2
}()
<-ch2    // 等待 ch2
ch1 <- 1 // 然后发送 ch1
原因：主协程等待 ch2，协程等待 ch1，互相等待`)

	fmt.Println("避免死锁的要点：")
	fmt.Println("  1. 确保每个发送都有对应的接收")
	fmt.Println("  2. 使用有缓冲通道可以缓解部分问题")
	fmt.Println("  3. 使用 select + default 实现非阻塞操作")
	fmt.Println("  4. 使用 select + timeout 防止永久等待")
	fmt.Println()
}

// =============================================================================
// 第九部分：实战模式 —— 工作池（Worker Pool）
// =============================================================================
// 关键问题：多个 worker 同时 range 同一个 tasks 通道，会不会数据竞争？
// 答案：不会！
//
// Channel 的接收操作是原子的（内部有锁保护）：
//   - 当多个协程同时从一个通道接收时，每条数据只会被一个协程收到
//   - Go 运行时保证"一个数据只发给一个接收方"
//   - 不需要额外加锁
//
// 类比：取号机。3 个窗口同时等号，每次只有一个窗口能拿到号，不会重复。
// =============================================================================

// Task 表示一个待处理的任务
type Task struct {
	ID       int
	Customer Customer
}

// worker 是工作池中的一个工作者
// 从 tasks 通道获取任务，将结果发送到 results 通道
func worker(id int, tasks <-chan Task, results chan<- NotificationResult, wg *sync.WaitGroup) {
	defer wg.Done()

	// 持续从 tasks 通道获取任务
	// 多个 worker 同时 range 同一个 tasks 通道是安全的：
	//   - Channel 内部保证每个任务只会被一个 worker 获取
	//   - 不会出现两个 worker 拿到同一个任务的情况
	// 当 tasks 通道关闭且为空时，for range 退出
	for task := range tasks {
		fmt.Printf("[Worker %d] 处理任务 %d: 通知 %s\n", id, task.ID, task.Customer.Name)

		// 模拟处理时间
		time.Sleep(200 * time.Millisecond)

		// 发送结果
		results <- NotificationResult{
			CustomerID:   task.Customer.ID,
			CustomerName: task.Customer.Name,
			Success:      true,
			Message:      fmt.Sprintf("由 Worker %d 处理完成", id),
		}
	}
	fmt.Printf("[Worker %d] 退出\n", id)
}

// demoWorkerPool 演示工作池模式
// 使用固定数量的工作者处理大量任务，控制并发度
func demoWorkerPool(customers []Customer) {
	fmt.Println("========== 工作池模式 ==========")

	// 创建任务通道和结果通道
	tasks := make(chan Task, len(customers))
	results := make(chan NotificationResult, len(customers))

	// 启动 3 个工作者（限制并发度为 3）
	var wg sync.WaitGroup
	workerCount := 3
	for i := 1; i <= workerCount; i++ {
		wg.Add(1)
		go worker(i, tasks, results, &wg)
	}

	// 发送所有任务
	for i, customer := range customers {
		tasks <- Task{ID: i + 1, Customer: customer}
	}
	// 关闭任务通道，通知工作者没有更多任务了
	close(tasks)

	// 等待所有工作者完成后关闭结果通道
	go func() {
		wg.Wait()
		close(results)
	}()

	// 收集结果
	fmt.Println("\n处理结果:")
	for result := range results {
		fmt.Printf("  - %s: %s\n", result.CustomerName, result.Message)
	}
	fmt.Println()
}

// =============================================================================
// 第十部分：使用 Channel 实现优雅退出
// =============================================================================
// 【重要】select 不是 switch！
//
// switch 思维（错误）：立即检查条件，满足就执行，不满足跳过，不阻塞
// select 实际行为：阻塞等待通道，哪个通道先有数据，就执行哪个 case
//
// for { select { } } 的执行流程：
//   1. 进入 select，阻塞等待（协程休眠，不占 CPU，不是疯狂转圈）
//   2. 某个通道有数据了 → 被唤醒，执行对应 case
//   3. 执行完后，回到 for，再次进入 select，阻塞等待
//   4. 重复...直到 return 或 break
//
// 本例时间线：
//   0ms:   启动协程，进入 select，阻塞等待（done 是空的，ticker.C 也没数据）
//   100ms: ticker.C 有数据！执行任务 #1，回到 for，继续阻塞等待...
//   200ms: ticker.C 有数据！执行任务 #2，回到 for，继续阻塞等待...
//   300ms: ticker.C 有数据！执行任务 #3，回到 for，继续阻塞等待...
//   350ms: 主协程 close(done)，done 就绪！执行 case <-done，return 退出
// =============================================================================

// demoGracefulShutdown 演示使用 Channel 信号控制协程退出
func demoGracefulShutdown() {
	fmt.Println("========== 优雅退出 ==========")

	// done 通道用于发送退出信号
	// struct{} 是空结构体，不占用内存，只用于传递"信号"，不传递数据
	done := make(chan struct{})

	// 启动一个持续运行的工作协程
	go func() {
		// time.NewTicker 创建一个定时器
		// 它会每隔 100ms 自动往 ticker.C 这个通道里发送一个时间值
		// ticker.C 的类型是 <-chan time.Time（只读通道）
		ticker := time.NewTicker(100 * time.Millisecond)
		defer ticker.Stop() // 退出时释放定时器资源

		count := 0

		// 这是一个死循环，但【不是疯狂转圈】！
		// 因为 select 会【阻塞】，协程在等待期间是休眠的，不占 CPU
		for {
			// select 在这里【阻塞等待】，同时监听两个通道：
			//   - done：退出信号（目前是空的，没人往里发数据）
			//   - ticker.C：定时信号（每 100ms 自动发一次）
			// 哪个通道先有数据，就执行哪个 case
			select {
			case <-done:
				// 【问】这个 case 第一次执行了吗？
				// 【答】没有！因为 done 通道是空的，没有数据
				//      select 不会执行这个 case，而是继续阻塞等待
				//      直到 350ms 时主协程 close(done)，这个 case 才会执行
				fmt.Println("[工作协程] 收到退出信号，正在清理...")
				time.Sleep(50 * time.Millisecond)
				fmt.Println("[工作协程] 清理完成，退出")
				return
			// <-ticker.C 表达式的含义是“每隔指定时间，从 ticker.C 这个通道接收一个当前时间点信号”。
			// C 是 ticker 内部维护的只读时间通道，定时器每到达设定的间隔就把对应的时间发送到此通道，需要通过接收操作（<-ticker.C）来获得下一个“定时唤醒”信号。
			case <-ticker.C:
				// 定时执行的任务
				count++
				fmt.Printf("[工作协程] 执行任务 #%d\n", count)
			}
		}
	}()

	// 主协程等待一段时间后发送退出信号
	time.Sleep(350 * time.Millisecond)
	fmt.Println("[主协程] 发送退出信号")

	// 【关键】close(done) 的行为：
	// close() 不只是"释放资源"，它会让所有正在等待 <-done 的操作【立即就绪】！
	// 就绪后返回通道元素类型的零值（对于 struct{}，零值是 struct{}{}）
	//
	// 对比：
	//   done <- struct{}{}  → 只有一个接收者能收到
	//   close(done)         → 所有等待的接收者都能收到（广播效果）
	//
	// 所以 close(done) 就是"向所有监听者发送退出信号"
	close(done)

	// 等待工作协程完成清理
	time.Sleep(100 * time.Millisecond)
	fmt.Println()
}

// =============================================================================
// 主函数
// =============================================================================

func main() {
	// 测试数据
	customers := []Customer{
		{ID: "C001", Name: "张三", Email: "zhangsan@example.com"},
		{ID: "C002", Name: "李四", Email: "lisi@example.com"},
		{ID: "C003", Name: "王五", Email: "wangwu@example.com"},
		{ID: "C004", Name: "赵六", Email: "zhaoliu@example.com"},
		{ID: "C005", Name: "钱七", Email: "qianqi@example.com"},
	}

	// 1. 共享内存方式的问题（对比）
	demoSharedMemoryProblem(customers[:3])

	// 2. Channel 基础
	demoChannelBasic()

	// 3. 无缓冲通道
	demoUnbufferedChannel()

	// 4. 有缓冲通道
	demoBufferedChannel()

	// 5. 使用 Channel 收集结果
	demoChannelCollectResults(customers[:3])

	// 6. 关闭 Channel
	demoCloseChannel()

	// 7. select vs for range（展示 select 的不可替代性）
	demoSelectVsRange()

	// 8. select 多路复用
	demoSelect()

	// 9. select 超时控制
	demoSelectTimeout()

	// 10. select default 非阻塞
	demoSelectDefault()

	// 11. 死锁说明
	demoDeadlockExamples()

	// 12. 工作池模式
	demoWorkerPool(customers)

	// 13. 优雅退出
	demoGracefulShutdown()
}
