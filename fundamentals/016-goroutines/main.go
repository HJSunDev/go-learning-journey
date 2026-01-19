package main

import (
	"fmt"
	"sync"
	"time"
)

// =============================================================================
// 016. Go 协程（Goroutine）：轻量级并发的基石
// =============================================================================
// 场景：订单系统 - 批量发送客户通知
// 当有多个客户需要通知时，如何高效地并发处理？
// =============================================================================

// Customer 表示一个客户
type Customer struct {
	ID    string
	Name  string
	Email string
}

// =============================================================================
// 第一部分：痛点 —— 顺序执行的问题
// =============================================================================

// sendNotificationSync 模拟同步发送通知（耗时操作）
// 每个通知需要 500ms 来模拟网络延迟
func sendNotificationSync(customer Customer) {
	fmt.Printf("[同步] 开始发送通知给 %s (%s)...\n", customer.Name, customer.Email)
	// 模拟网络延迟
	time.Sleep(500 * time.Millisecond)
	fmt.Printf("[同步] 通知已发送给 %s\n", customer.Name)
}

// demoSequential 演示顺序执行的问题
// 3个客户，每个500ms，总共需要1500ms
func demoSequential(customers []Customer) {
	fmt.Println("========== 顺序执行演示 ==========")
	start := time.Now()

	for _, customer := range customers {
		sendNotificationSync(customer)
	}

	elapsed := time.Since(start)
	fmt.Printf("顺序执行总耗时: %v\n\n", elapsed)
}

// =============================================================================
// 第二部分：go 关键字 —— 启动协程
// =============================================================================

// demoGoroutineBasic 演示最基础的协程用法
// 问题：主程序可能在协程执行完之前就退出了
func demoGoroutineBasic(customers []Customer) {
	fmt.Println("========== 基础协程演示（有问题的版本）==========")
	start := time.Now()

	for _, customer := range customers {
		// go 关键字启动一个协程
		// 这行代码会立即返回，不会等待 sendNotificationSync 执行完成
		go sendNotificationSync(customer)
	}

	elapsed := time.Since(start)
	// 这里几乎立即就打印了，因为 go 只是"安排"任务，不等待执行
	fmt.Printf("启动所有协程耗时: %v\n", elapsed)
	fmt.Println("注意：此时协程可能还没执行完！\n")

	// 临时方案：等待一段时间让协程有机会执行
	// 这是一个糟糕的解决方案，因为你不知道要等多久
	time.Sleep(600 * time.Millisecond)
}

// =============================================================================
// 第三部分：WaitGroup —— 等待所有协程完成
// =============================================================================

// sendNotificationWithWG 发送通知，完成后通知 WaitGroup
func sendNotificationWithWG(customer Customer, wg *sync.WaitGroup) {
	// defer 确保函数结束时一定会调用 Done()
	// 即使发生 panic 也会执行
	defer wg.Done()

	fmt.Printf("[并发] 开始发送通知给 %s (%s)...\n", customer.Name, customer.Email)
	time.Sleep(500 * time.Millisecond)
	fmt.Printf("[并发] 通知已发送给 %s\n", customer.Name)
}

// demoWaitGroup 演示使用 WaitGroup 等待所有协程完成
func demoWaitGroup(customers []Customer) {
	fmt.Println("========== WaitGroup 演示 ==========")
	start := time.Now()

	// sync.WaitGroup 是一个计数器
	// 用于等待一组协程全部完成
	var wg sync.WaitGroup

	for _, customer := range customers {
		// Add(1) 告诉 WaitGroup：有一个新的协程要启动
		// 必须在启动协程之前调用
		wg.Add(1)

		// 启动协程执行通知
		go sendNotificationWithWG(customer, &wg)
	}

	// Wait() 会阻塞，直到计数器变为 0
	// 每个协程完成时调用 Done()，计数器减 1
	wg.Wait()

	elapsed := time.Since(start)
	fmt.Printf("并发执行总耗时: %v\n", elapsed)
	fmt.Println("所有通知已发送完成！\n")
}

// =============================================================================
// 第四部分：闭包变量捕获（历史陷阱，Go 1.22 已修复）
// =============================================================================

// demoClosureTrap 演示闭包变量捕获的经典陷阱
// Go 1.22+ 已自动修复此问题（需要 go.mod 声明 go 1.22 或更高）
// 但在 Go 1.21 及之前版本中，这段代码会输出重复的名字
func demoClosureTrap(customers []Customer) {
	fmt.Println("========== 闭包捕获演示（Go 1.22+ 已自动修复）==========")
	var wg sync.WaitGroup

	for _, customer := range customers {
		wg.Add(1)

		// Go 1.21 及之前：所有协程共享同一个 customer 变量，结果全是最后一个
		// Go 1.22+：每次迭代创建新的 customer 变量，结果正确
		go func() {
			defer wg.Done()
			fmt.Printf("[闭包] 发送给: %s\n", customer.Name)
		}()
	}

	wg.Wait()
	fmt.Println("Go 1.22+：每次迭代创建新变量，结果正确\n")
}

// demoClosureFixed 演示传统修复方法（参数传递）
// 这种写法在所有 Go 版本中都正确，推荐用于需要兼容旧版本的代码
func demoClosureFixed(customers []Customer) {
	fmt.Println("========== 传统修复：参数传递（兼容所有版本）==========")
	var wg sync.WaitGroup

	for _, customer := range customers {
		wg.Add(1)

		// 将变量作为参数传入匿名函数
		// 参数传递是值拷贝，每个协程获得自己的副本
		// 这种写法在 Go 1.21 及之前是必须的，在 Go 1.22+ 中仍然推荐
		go func(c Customer) {
			defer wg.Done()
			fmt.Printf("[参数传递] 发送给: %s\n", c.Name)
		}(customer)
	}

	wg.Wait()
	fmt.Println()
}

// =============================================================================
// 第五部分：最佳实践示例
// =============================================================================

// NotificationService 通知服务
// 封装并发通知的逻辑
type NotificationService struct {
	// 最大并发数（这里仅作示意，实际控制需要用 channel，下一章讲解）
	MaxConcurrency int
}

// NotifyAll 向所有客户发送通知
// 返回成功发送的数量
func (s *NotificationService) NotifyAll(customers []Customer) int {
	var wg sync.WaitGroup
	successCount := 0

	// 用于安全地更新计数器的互斥锁
	var mu sync.Mutex

	for _, customer := range customers {
		wg.Add(1)

		// 将 customer 作为参数传入，避免闭包陷阱
		go func(c Customer) {
			defer wg.Done()

			// 模拟发送通知
			success := s.sendNotification(c)

			if success {
				// 使用互斥锁保护共享变量
				mu.Lock()
				successCount++
				mu.Unlock()
			}
		}(customer)
	}

	wg.Wait()
	return successCount
}

// sendNotification 发送单个通知
func (s *NotificationService) sendNotification(c Customer) bool {
	fmt.Printf("[服务] 正在通知 %s...\n", c.Name)
	time.Sleep(300 * time.Millisecond)
	fmt.Printf("[服务] %s 通知成功\n", c.Name)
	return true
}

// demoBestPractice 演示最佳实践
func demoBestPractice(customers []Customer) {
	fmt.Println("========== 最佳实践演示 ==========")
	start := time.Now()

	service := &NotificationService{MaxConcurrency: 10}
	count := service.NotifyAll(customers)

	elapsed := time.Since(start)
	fmt.Printf("成功通知 %d 位客户，耗时: %v\n\n", count, elapsed)
}

// =============================================================================
// 主函数
// =============================================================================

func main() {
	// 测试数据：3个客户
	customers := []Customer{
		{ID: "C001", Name: "张三", Email: "zhangsan@example.com"},
		{ID: "C002", Name: "李四", Email: "lisi@example.com"},
		{ID: "C003", Name: "王五", Email: "wangwu@example.com"},
	}

	// 1. 顺序执行 —— 展示痛点
	demoSequential(customers)

	// 2. 基础协程 —— go 关键字
	demoGoroutineBasic(customers)

	// 3. WaitGroup —— 正确等待协程完成
	demoWaitGroup(customers)

	// 4. 闭包变量捕获（Go 1.22+ 已自动修复）
	demoClosureTrap(customers)

	// 5. 传统修复方法（兼容旧版本）
	demoClosureFixed(customers)

	// 6. 最佳实践
	demoBestPractice(customers)
}
