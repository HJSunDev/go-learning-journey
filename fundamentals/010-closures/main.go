package main

import (
	"fmt"
	"time"
)

// 1. 基础闭包
// adder 返回一个闭包。
// 这里的关键是：返回的函数引用了 adder 函数内部的 sum 变量。
// 即使 adder 执行结束了，sum 变量依然会被保留，因为它被返回的函数（闭包）捕获了。
func adder() func(int) int {
	sum := 0
	// 这是一个匿名函数（函数字面量）。
	// 它引用了外部的 sum。
	return func(x int) int {
		sum += x
		return sum
	}
}

// 2. 模拟第三方库的回调场景
// 假设这是一个我们无法修改的第三方库函数。
// 它要求传入的 handler 必须是 func() 类型（无参数，无返回值）。
func RunTask(handler func()) {
	fmt.Println("System: Starting task...")
	// 模拟耗时操作
	time.Sleep(100 * time.Millisecond)
	// 执行回调
	handler()
	fmt.Println("System: Task finished.")
}

// 这是一个具体的业务函数，它有特定的参数。
// 它的签名 func(string, int) 与 RunTask 要求的 func() 不匹配。
func sendEmail(content string, userID int) {
	fmt.Printf("  -> Email sent to user %d: %s\n", userID, content)
}

func main() {
	fmt.Println("--- 1. 基础闭包：变量捕获 ---")
	// pos 是一个闭包函数，它拥有自己独立的 sum 变量
	pos := adder()
	for i := 0; i < 3; i++ {
		fmt.Printf("pos(%d) = %d\n", i, pos(i))
	}

	fmt.Println("\n--- 2. 实战场景：封装函数以适配接口 ---")
	
	user := 10086
	msg := "Welcome to Go!"

	// 直接传递 sendEmail 是不行的，因为签名不匹配：
	// RunTask(sendEmail) // 编译错误：cannot use sendEmail (value of type func(string, int)) as func() value

	// 解决方案：使用闭包进行“包装”。
	// 我们定义一个匿名函数，它符合 func() 的签名。
	// 在这个匿名函数内部，我们要调用 sendEmail，并“捕获”外部的 user 和 msg 变量。
	wrappedTask := func() {
		// 这里捕获了 main 函数中的 msg 和 user 变量
		sendEmail(msg, user)
	}

	// 现在可以将 wrappedTask 传递给 RunTask 了
	RunTask(wrappedTask)

	// 也可以写成更紧凑的形式，直接传入匿名函数：
	fmt.Println("\n--- 3. 实战场景：直接传递闭包 ---")
	RunTask(func() {
		fmt.Printf("  -> Quick logging: User %d is active\n", user)
	})
}
