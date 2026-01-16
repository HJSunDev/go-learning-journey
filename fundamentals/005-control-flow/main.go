package main

import "fmt"

func main() {
	fmt.Println("=== Go 控制流演示 ===")
	fmt.Println()

	// 1. if-else 条件语句
	demonstrateIfElse()

	// 2. switch 语句
	demonstrateSwitch()

	// 3. for 循环
	demonstrateForLoop()

	// 4. 循环控制：break 和 continue
	demonstrateLoopControl()

	// 5. defer 延迟执行
	demonstrateDefer()
}

// demonstrateIfElse 演示 if-else 条件语句的各种形式
func demonstrateIfElse() {
	fmt.Println("--- 1. if-else 条件语句 ---")

	// 1.1 基本 if 语句
	age := 20
	if age >= 18 {
		fmt.Println("已成年")
	}

	// 1.2 if-else
	score := 75
	if score >= 60 {
		fmt.Println("考试通过")
	} else {
		fmt.Println("考试未通过")
	}

	// 1.3 if-else if-else 链
	grade := 85
	if grade >= 90 {
		fmt.Println("等级: A")
	} else if grade >= 80 {
		fmt.Println("等级: B")
	} else if grade >= 70 {
		fmt.Println("等级: C")
	} else if grade >= 60 {
		fmt.Println("等级: D")
	} else {
		fmt.Println("等级: F")
	}

	// 1.4 带初始化语句的 if（Go 特有，非常实用）
	// 变量 n 的作用域仅限于 if-else 块内部
	if n := 10; n%2 == 0 {
		fmt.Printf("%d 是偶数\n", n)
	} else {
		fmt.Printf("%d 是奇数\n", n)
	}
	// n 在这里不可访问

	// 1.5 实际应用：错误检查（Go 的惯用模式）
	value, err := divide(10, 2)
	if err != nil {
		fmt.Println("错误:", err)
	} else {
		fmt.Println("结果:", value)
	}

	// 更常见的写法：初始化 + 错误检查
	if result, err := divide(10, 0); err != nil {
		fmt.Println("错误:", err)
	} else {
		fmt.Println("结果:", result)
	}

	fmt.Println()
}

// divide 演示返回值和错误处理
func divide(a, b int) (int, error) {
	if b == 0 {
		return 0, fmt.Errorf("除数不能为零")
	}
	return a / b, nil
}

// demonstrateSwitch 演示 switch 语句的各种形式
func demonstrateSwitch() {
	fmt.Println("--- 2. switch 语句 ---")

	// 2.1 基本 switch（不需要 break！）
	day := 3
	switch day {
	case 1:
		fmt.Println("星期一")
	case 2:
		fmt.Println("星期二")
	case 3:
		fmt.Println("星期三")
	case 4:
		fmt.Println("星期四")
	case 5:
		fmt.Println("星期五")
	case 6, 7: // 多值匹配
		fmt.Println("周末")
	default:
		fmt.Println("无效的日期")
	}

	// 2.2 带初始化语句的 switch
	switch num := 15; {
	case num < 0:
		fmt.Println("负数")
	case num == 0:
		fmt.Println("零")
	case num > 0:
		fmt.Println("正数")
	}

	// 2.3 无表达式 switch（相当于 if-else if 链）
	hour := 14
	switch {
	case hour < 12:
		fmt.Println("上午好")
	case hour < 18:
		fmt.Println("下午好")
	default:
		fmt.Println("晚上好")
	}

	// 2.4 fallthrough：强制执行下一个 case
	level := 1
	fmt.Print("你的权限: ")
	switch level {
	case 1:
		fmt.Print("读取 ")
		fallthrough
	case 2:
		fmt.Print("写入 ")
		fallthrough
	case 3:
		fmt.Print("执行")
	}
	fmt.Println()

	// 2.5 类型 switch（检查接口的实际类型）
	checkType(42)
	checkType("hello")
	checkType(3.14)
	checkType([]int{1, 2, 3})

	fmt.Println()
}

// checkType 演示类型 switch
func checkType(x interface{}) {
	switch v := x.(type) {
	case int:
		fmt.Printf("%v 是 int 类型\n", v)
	case string:
		fmt.Printf("%v 是 string 类型\n", v)
	case float64:
		fmt.Printf("%v 是 float64 类型\n", v)
	default:
		fmt.Printf("%v 是未知类型: %T\n", v, v)
	}
}

// demonstrateForLoop 演示 for 循环的各种形式
func demonstrateForLoop() {
	fmt.Println("--- 3. for 循环 ---")

	// 3.1 标准三段式 for 循环
	fmt.Print("三段式: ")
	for i := 0; i < 5; i++ {
		fmt.Print(i, " ")
	}
	fmt.Println()

	// 3.2 while 风格（只有条件）
	fmt.Print("while风格: ")
	count := 0
	for count < 5 {
		fmt.Print(count, " ")
		count++
	}
	fmt.Println()

	// 3.3 无限循环（通常配合 break 使用）
	fmt.Print("无限循环+break: ")
	n := 0
	for {
		if n >= 5 {
			break
		}
		fmt.Print(n, " ")
		n++
	}
	fmt.Println()

	// 3.4 for range 遍历切片
	nums := []int{10, 20, 30, 40, 50}
	fmt.Print("遍历切片(索引+值): ")
	for i, v := range nums {
		fmt.Printf("[%d]=%d ", i, v)
	}
	fmt.Println()

	// 只需要值，用 _ 忽略索引
	fmt.Print("只要值: ")
	for _, v := range nums {
		fmt.Print(v, " ")
	}
	fmt.Println()

	// 只需要索引
	fmt.Print("只要索引: ")
	for i := range nums {
		fmt.Print(i, " ")
	}
	fmt.Println()

	// 3.5 for range 遍历字符串（按 rune）
	str := "Go语言"
	fmt.Println("遍历字符串:")
	for i, r := range str {
		fmt.Printf("  索引 %d: %c (Unicode: %U)\n", i, r, r)
	}

	// 3.6 for range 遍历 map
	ages := map[string]int{
		"Alice": 25,
		"Bob":   30,
		"Carol": 28,
	}
	fmt.Println("遍历 map:")
	for name, age := range ages {
		fmt.Printf("  %s: %d 岁\n", name, age)
	}

	fmt.Println()
}

// demonstrateLoopControl 演示 break、continue 和标签
func demonstrateLoopControl() {
	fmt.Println("--- 4. 循环控制 ---")

	// 4.1 break：跳出循环
	fmt.Print("break 示例: ")
	for i := 0; i < 10; i++ {
		if i == 5 {
			break
		}
		fmt.Print(i, " ")
	}
	fmt.Println()

	// 4.2 continue：跳过当前迭代
	fmt.Print("continue 示例(跳过偶数): ")
	for i := 0; i < 10; i++ {
		if i%2 == 0 {
			continue
		}
		fmt.Print(i, " ")
	}
	fmt.Println()

	// 4.3 带标签的 break（跳出外层循环）
	fmt.Println("带标签的 break:")
outer:
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			if i == 1 && j == 1 {
				fmt.Println("  在 i=1, j=1 处跳出外层循环")
				break outer
			}
			fmt.Printf("  i=%d, j=%d\n", i, j)
		}
	}

	// 4.4 带标签的 continue（继续外层循环的下一次迭代）
	fmt.Println("带标签的 continue:")
nextRow:
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			if j == 1 {
				fmt.Printf("  跳过 i=%d 的剩余部分\n", i)
				continue nextRow
			}
			fmt.Printf("  i=%d, j=%d\n", i, j)
		}
	}

	fmt.Println()
}

// demonstrateDefer 演示 defer 延迟执行
func demonstrateDefer() {
	fmt.Println("--- 5. defer 延迟执行 ---")

	// 5.1 基本用法：函数返回前执行
	fmt.Println("调用 deferBasic():")
	deferBasic()

	// 5.2 多个 defer：后进先出（LIFO）
	fmt.Println("\n调用 deferStack():")
	deferStack()

	// 5.3 defer 与参数求值
	fmt.Println("\n调用 deferEvaluation():")
	deferEvaluation()

	// 5.4 实际应用：资源清理（模拟）
	fmt.Println("\n调用 processFile():")
	processFile("config.json")

	fmt.Println()
}

func deferBasic() {
	defer fmt.Println("  3. 这是 defer，最后执行")
	fmt.Println("  1. 函数开始")
	fmt.Println("  2. 函数中间")
	// defer 语句在函数返回前执行
}

func deferStack() {
	// 多个 defer 按 LIFO（后进先出）顺序执行
	defer fmt.Println("  第一个 defer")
	defer fmt.Println("  第二个 defer")
	defer fmt.Println("  第三个 defer")
	fmt.Println("  函数主体")
	// 输出顺序: 函数主体 → 第三个 → 第二个 → 第一个
}

func deferEvaluation() {
	x := 10
	// defer 的参数在声明时求值，而不是执行时
	defer fmt.Printf("  defer 中的 x = %d\n", x)
	x = 20
	fmt.Printf("  函数中的 x = %d\n", x)
	// 输出: 函数中的 x = 20, defer 中的 x = 10
}

func processFile(filename string) {
	// 模拟打开文件
	fmt.Printf("  打开文件: %s\n", filename)

	// 立即注册清理操作（关闭文件）
	// 无论函数如何返回（正常返回或 panic），defer 都会执行
	defer fmt.Printf("  关闭文件: %s\n", filename)

	// 模拟处理文件
	fmt.Printf("  处理文件: %s\n", filename)

	// 即使这里发生错误返回，defer 仍会执行
	// return
}
