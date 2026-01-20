package main

import (
	"fmt"
	"sync"
	"time"
)

// =============================================================================
// 018. Go 互斥锁（Mutex）：并发安全的守护者
// =============================================================================
// 场景：订单系统 - 实时统计模块
// 多个客服协程同时处理订单，需要安全地更新共享的统计数据
// =============================================================================

// OrderStats 订单统计数据（共享资源）
// 多个协程会同时读写这个结构体
type OrderStats struct {
	TotalOrders   int     // 总订单数
	TotalAmount   float64 // 总金额
	SuccessOrders int     // 成功订单数
	FailedOrders  int     // 失败订单数
}

// Order 表示一个订单
type Order struct {
	ID     string
	Amount float64
}

// =============================================================================
// 第一部分：痛点 —— 没有锁时的数据竞争
// =============================================================================

// processOrderUnsafe 不安全的订单处理
// 直接修改共享变量，会导致数据竞争
func processOrderUnsafe(order Order, stats *OrderStats, wg *sync.WaitGroup) {
	defer wg.Done()

	// 模拟处理耗时
	time.Sleep(10 * time.Millisecond)

	// === 危险区域：多个协程同时执行以下代码 ===
	//
	// 问题解析：stats.TotalOrders++ 看起来是一行代码，
	// 但在 CPU 层面实际是三个操作：
	//   1. 读取 stats.TotalOrders 的当前值到寄存器
	//   2. 将寄存器中的值加 1
	//   3. 将新值写回 stats.TotalOrders
	//
	// 当两个协程同时执行时，可能发生：
	//   协程A 读取值 = 5
	//   协程B 读取值 = 5  （还没等 A 写回）
	//   协程A 写入值 = 6
	//   协程B 写入值 = 6  （覆盖了 A 的结果！）
	// 结果：两个订单只增加了 1，丢失了一次计数
	stats.TotalOrders++
	stats.TotalAmount += order.Amount
	stats.SuccessOrders++
}

// demoUnsafe 演示没有锁时的数据竞争问题
func demoUnsafe() {
	fmt.Println("========== 第一部分：数据竞争问题（不安全版本）==========")

	orders := generateOrders(1000)
	stats := &OrderStats{}
	var wg sync.WaitGroup

	start := time.Now()

	for _, order := range orders {
		wg.Add(1)
		go processOrderUnsafe(order, stats, &wg)
	}

	wg.Wait()
	elapsed := time.Since(start)

	fmt.Printf("期望订单数: %d\n", len(orders))
	fmt.Printf("实际统计数: %d\n", stats.TotalOrders)
	fmt.Printf("丢失订单数: %d\n", len(orders)-stats.TotalOrders)
	fmt.Printf("耗时: %v\n\n", elapsed)
}

// =============================================================================
// 第二部分：互斥锁基础 —— Lock/Unlock
// =============================================================================

// SafeOrderStats 线程安全的订单统计
// 将互斥锁和数据放在同一个结构体中，是 Go 的惯用模式
type SafeOrderStats struct {
	mu sync.Mutex // 互斥锁，保护下面的字段

	TotalOrders   int
	TotalAmount   float64
	SuccessOrders int
	FailedOrders  int
}

// processOrderSafe 安全的订单处理
// 使用互斥锁保护共享变量
func processOrderSafe(order Order, stats *SafeOrderStats, wg *sync.WaitGroup) {
	defer wg.Done()

	// 模拟处理耗时
	time.Sleep(10 * time.Millisecond)

	// Lock() 获取锁
	// 如果锁已被其他协程持有，当前协程会在这里阻塞等待
	// 直到锁被释放，当前协程才能继续执行
	stats.mu.Lock()

	// === 临界区：同一时刻只有一个协程能执行这里的代码 ===
	stats.TotalOrders++
	stats.TotalAmount += order.Amount
	stats.SuccessOrders++

	// Unlock() 释放锁
	// 其他等待的协程可以获取锁并继续执行
	stats.mu.Unlock()
}

// demoMutexBasic 演示互斥锁的基本用法
func demoMutexBasic() {
	fmt.Println("========== 第二部分：互斥锁基础（Lock/Unlock）==========")

	orders := generateOrders(1000)
	stats := &SafeOrderStats{}
	var wg sync.WaitGroup

	start := time.Now()

	for _, order := range orders {
		wg.Add(1)
		go processOrderSafe(order, stats, &wg)
	}

	wg.Wait()
	elapsed := time.Since(start)

	fmt.Printf("期望订单数: %d\n", len(orders))
	fmt.Printf("实际统计数: %d\n", stats.TotalOrders)
	fmt.Printf("丢失订单数: %d\n", len(orders)-stats.TotalOrders)
	fmt.Printf("耗时: %v\n\n", elapsed)
}

// =============================================================================
// 第三部分：defer 解锁模式 —— 推荐的安全写法
// =============================================================================

// processOrderDefer 使用 defer 确保解锁
// 即使函数中发生 panic，锁也会被释放
func processOrderDefer(order Order, stats *SafeOrderStats, wg *sync.WaitGroup) {
	defer wg.Done()

	time.Sleep(10 * time.Millisecond)

	stats.mu.Lock()
	// defer 会在函数返回时执行 Unlock
	// 无论是正常返回、提前 return、还是 panic，都会执行
	// 这避免了忘记解锁或异常时锁未释放的问题
	defer stats.mu.Unlock()

	// 临界区代码
	stats.TotalOrders++
	stats.TotalAmount += order.Amount
	stats.SuccessOrders++

	// 假设这里可能有多个 return 分支，或可能 panic
	// 使用 defer 确保无论如何都会解锁
}

// demoDefer 演示 defer 解锁模式
func demoDefer() {
	fmt.Println("========== 第三部分：defer 解锁模式 ==========")

	orders := generateOrders(1000)
	stats := &SafeOrderStats{}
	var wg sync.WaitGroup

	for _, order := range orders {
		wg.Add(1)
		go processOrderDefer(order, stats, &wg)
	}

	wg.Wait()
	fmt.Printf("统计结果: 总订单 %d, 成功 %d\n\n", stats.TotalOrders, stats.SuccessOrders)
}

// =============================================================================
// 第四部分：读写锁 RWMutex —— 读多写少场景优化
// =============================================================================
//
// 场景延续：订单统计系统现在要提供查询接口
//
// 问题：统计数据需要被频繁查询
// - 写操作：每处理完一个订单，更新一次统计（相对较少）
// - 读操作：前端页面、报表系统、告警系统都要查询统计数据（非常频繁）
//
// 使用普通 Mutex 的问题：
// - 用户A查询统计时，用户B也想查询，但必须等A查完
// - 可是查询不会修改数据，为什么不能同时查？
//
// RWMutex 解决方案：
// - 多个"读"可以同时进行（共享锁）
// - "写"必须独占（写的时候不能读，读的时候不能写）
// =============================================================================

// ----- 方案一：使用普通 Mutex（读也会互斥）-----

// StatsWithMutex 使用普通互斥锁的统计服务
type StatsWithMutex struct {
	mu          sync.Mutex
	TotalOrders int
	TotalAmount float64
}

// Query 查询统计（读操作）
// 问题：即使只是读取，也要等待其他读操作完成
func (s *StatsWithMutex) Query() (int, float64) {
	s.mu.Lock()
	defer s.mu.Unlock()

	// 模拟查询需要一点时间（比如聚合计算）
	time.Sleep(100 * time.Millisecond)

	return s.TotalOrders, s.TotalAmount
}

// ----- 方案二：使用 RWMutex（读可以并发）-----

// StatsWithRWMutex 使用读写锁的统计服务
type StatsWithRWMutex struct {
	mu          sync.RWMutex
	TotalOrders int
	TotalAmount float64
}

// Query 查询统计（读操作）
// 优势：多个查询可以同时进行
func (s *StatsWithRWMutex) Query() (int, float64) {
	// RLock = Read Lock，获取读锁
	// 多个协程可以同时持有读锁
	s.mu.RLock()
	defer s.mu.RUnlock()

	// 模拟查询需要一点时间
	time.Sleep(100 * time.Millisecond)

	return s.TotalOrders, s.TotalAmount
}

// Update 更新统计（写操作）
func (s *StatsWithRWMutex) Update(amount float64) {
	// Lock = Write Lock，获取写锁（独占）
	// 必须等所有读锁释放后才能获取
	s.mu.Lock()
	defer s.mu.Unlock()

	s.TotalOrders++
	s.TotalAmount += amount
}

// demoRWMutex 对比演示：Mutex vs RWMutex 在多读场景下的性能差异
func demoRWMutex() {
	fmt.Println("========== 第四部分：读写锁 RWMutex ==========")
	fmt.Println()

	// ----- 场景：3 个用户同时查询统计数据 -----

	// 测试 1：使用普通 Mutex
	fmt.Println("--- 测试 1：使用普通 Mutex（3 个查询必须排队）---")
	statsMutex := &StatsWithMutex{TotalOrders: 100, TotalAmount: 5000}
	var wg1 sync.WaitGroup

	start1 := time.Now()
	for i := 1; i <= 3; i++ {
		wg1.Add(1)
		go func(userID int) {
			defer wg1.Done()
			count, amount := statsMutex.Query()
			fmt.Printf("用户%d 查询完成: %d 订单, %.0f 元\n", userID, count, amount)
		}(i)
	}
	wg1.Wait()
	fmt.Printf("总耗时: %v\n", time.Since(start1))
	fmt.Println()

	// 测试 2：使用 RWMutex
	fmt.Println("--- 测试 2：使用 RWMutex（3 个查询可以同时进行）---")
	statsRW := &StatsWithRWMutex{TotalOrders: 100, TotalAmount: 5000}
	var wg2 sync.WaitGroup

	start2 := time.Now()
	for i := 1; i <= 3; i++ {
		wg2.Add(1)
		go func(userID int) {
			defer wg2.Done()
			count, amount := statsRW.Query()
			fmt.Printf("用户%d 查询完成: %d 订单, %.0f 元\n", userID, count, amount)
		}(i)
	}
	wg2.Wait()
	fmt.Printf("总耗时: %v\n", time.Since(start2))
	fmt.Println()

	// 结论
	fmt.Println("对比结论：")
	fmt.Println("- Mutex: 3 个查询排队执行，耗时约 300ms")
	fmt.Println("- RWMutex: 3 个查询同时执行，耗时约 100ms")
	fmt.Println("- 读多写少场景下，RWMutex 性能更好")
	fmt.Println()
}

// =============================================================================
// 第五部分：常见陷阱 —— 死锁演示
// =============================================================================

// demoDeadlock 演示死锁情况（已注释掉实际会卡住的代码）
func demoDeadlock() {
	fmt.Println("========== 第五部分：常见陷阱 ==========")

	// --- 陷阱 1：重复加锁（同一个协程对同一个锁加锁两次）---
	// var mu sync.Mutex
	// mu.Lock()
	// mu.Lock()  // 死锁！同一协程再次加锁会永远阻塞
	// mu.Unlock()
	fmt.Println("陷阱 1: 重复加锁 - 同一协程对同一锁连续 Lock() 会死锁")

	// --- 陷阱 2：忘记解锁 ---
	// func doSomething() {
	//     mu.Lock()
	//     if someCondition {
	//         return  // 忘记解锁！后续所有尝试加锁的操作都会永远阻塞
	//     }
	//     mu.Unlock()
	// }
	fmt.Println("陷阱 2: 忘记解锁 - 使用 defer 可以避免")

	// --- 陷阱 3：锁拷贝 ---
	// mu := sync.Mutex{}
	// mu2 := mu      // 错误！拷贝了锁
	// mu.Lock()
	// mu2.Unlock()   // 这是解锁另一个锁！原锁仍被持有
	fmt.Println("陷阱 3: 锁拷贝 - Mutex 是值类型，拷贝后是两个不同的锁")

	// --- 陷阱 4：交叉锁定（两个锁的获取顺序不一致）---
	//
	// 场景：系统中有两把锁，lockA 保护资源A，lockB 保护资源B
	// 某个操作需要同时访问资源A和资源B，所以要同时持有两把锁
	//
	// 问题代码示例：
	//
	// var lockA, lockB sync.Mutex
	//
	// // 协程 A 的代码：先锁 A，再锁 B
	// go func() {
	//     lockA.Lock()
	//     time.Sleep(1 * time.Millisecond)  // 模拟一些操作
	//     lockB.Lock()                       // 等待 lockB
	//     // ... 操作资源 ...
	//     lockB.Unlock()
	//     lockA.Unlock()
	// }()
	//
	// // 协程 B 的代码：先锁 B，再锁 A（顺序相反！）
	// go func() {
	//     lockB.Lock()
	//     time.Sleep(1 * time.Millisecond)  // 模拟一些操作
	//     lockA.Lock()                       // 等待 lockA
	//     // ... 操作资源 ...
	//     lockA.Unlock()
	//     lockB.Unlock()
	// }()
	//
	// 死锁过程：
	// 1. 协程A 获取了 lockA，然后想获取 lockB
	// 2. 协程B 获取了 lockB，然后想获取 lockA
	// 3. 协程A 等待 lockB 释放，但 lockB 被协程B 持有
	// 4. 协程B 等待 lockA 释放，但 lockA 被协程A 持有
	// 5. 双方互相等待，谁也无法继续，程序卡死
	//
	// 解决方案：所有协程都按相同顺序获取锁（比如都先锁A再锁B）
	fmt.Println("陷阱 4: 交叉锁定 - 两个协程以相反顺序获取多把锁会导致死锁")

	fmt.Println()
}

// =============================================================================
// 第六部分：最佳实践 —— 把锁封装在结构体里
// =============================================================================
//
// 核心思想：调用者不需要知道锁的存在，只管调用方法
//
// 对比：
//   不好的写法：调用者需要自己加锁
//     mu.Lock()
//     counter++
//     mu.Unlock()
//
//   好的写法：锁封装在方法里，调用者直接用
//     counter.Add(1)
//
// =============================================================================

// Counter 线程安全的计数器
// 这就是最佳实践的核心：把锁和数据放在一起，通过方法操作
type Counter struct {
	mu    sync.Mutex // 锁
	value int        // 被保护的数据
}

// Add 增加计数（调用者不需要关心锁）
func (c *Counter) Add(n int) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.value += n
}

// Value 获取当前值
func (c *Counter) Value() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.value
}

// demoBestPractice 演示最佳实践
func demoBestPractice() {
	fmt.Println("========== 第六部分：最佳实践 ==========")

	counter := &Counter{}
	var wg sync.WaitGroup

	// 启动 100 个协程，每个协程给计数器加 1
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			counter.Add(1) // 调用者不需要关心锁，直接用
		}()
	}

	wg.Wait()

	fmt.Printf("100 个协程并发执行 Add(1)\n")
	fmt.Printf("最终结果: %d（正确应该是 100）\n\n", counter.Value())
}

// =============================================================================
// 辅助函数
// =============================================================================

// generateOrders 生成测试订单
func generateOrders(count int) []Order {
	orders := make([]Order, count)
	for i := 0; i < count; i++ {
		orders[i] = Order{
			ID: fmt.Sprintf("ORD-%05d", i+1),
			// 设置订单金额：
			//  - i % 100 + 1 保证金额区间为 1 到 100，避免金额为 0
			//  - 再乘以 10 得到 10, 20, ..., 1000 之间的金额
			//  - float64 用于类型转换，将 int 转换为浮点数类型
			//    Go 中变量不能自动在 int/float64 之间转换，需要手动用 float64(x) 转为浮点数
			Amount: float64(i%100+1) * 10,
		}
	}
	return orders
}

// =============================================================================
// 主函数
// =============================================================================

func main() {
	// 1. 数据竞争问题（不安全版本）
	demoUnsafe()

	// 2. 互斥锁基础
	demoMutexBasic()

	// 3. defer 解锁模式
	demoDefer()

	// 4. 读写锁
	demoRWMutex()

	// 5. 常见陷阱
	demoDeadlock()

	// 6. 最佳实践
	demoBestPractice()
}
