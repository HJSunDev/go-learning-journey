package main

import (
	"fmt"
	"strings"
)

func main() {
	fmt.Println("=== Go 函数：代码复用的基石 ===")
	fmt.Println()

	// 1. 函数基础
	demonstrateFunctionBasics()

	// 2. 参数传递
	demonstrateParameters()

	// 3. 返回值
	demonstrateReturnValues()

	// 4. 可变参数函数
	demonstrateVariadicFunctions()

	// 5. 函数类型与函数作为值
	demonstrateFunctionTypes()

	// 6. 匿名函数
	demonstrateAnonymousFunctions()

	// 7. 函数作为参数（高阶函数）
	demonstrateHigherOrderFunctions()

	// 8. 函数作为返回值
	demonstrateFunctionReturningFunction()

	// 9. 递归函数
	demonstrateRecursion()

	// 10. defer 延迟执行
	demonstrateDefer()

	// 11. init 函数说明
	demonstrateInitExplanation()
}

// ============================================================
// 1. 函数基础
// ============================================================

// greet 是一个简单的无参无返回值函数
// 函数名以小写开头，只能在同一个包内访问（未导出）
func greet() {
	fmt.Println("  Hello, Go!")
}

// Greet 以大写开头，可以被其他包访问（导出）
// 这是 Go 的可见性规则：大写=公开，小写=私有
func Greet() {
	fmt.Println("  Hello from exported function!")
}

func demonstrateFunctionBasics() {
	fmt.Println("--- 1. 函数基础 ---")

	// Go 函数的基本结构：
	// func 函数名(参数列表) 返回值类型 {
	//     函数体
	// }

	// 调用无参函数
	fmt.Println("调用 greet():")
	greet()

	// 函数命名规范
	fmt.Println("\n📌 函数命名规范:")
	fmt.Println("  - 小写开头 (greet): 包内私有，其他包无法调用")
	fmt.Println("  - 大写开头 (Greet): 公开导出，其他包可以调用")
	fmt.Println("  - 使用驼峰命名法 (calculateTotalPrice)")
	fmt.Println("  - 名称应清晰表达函数的作用")

	fmt.Println()
}

// ============================================================
// 2. 参数传递
// ============================================================

// add 接收两个整数参数并返回它们的和
func add(a int, b int) int {
	return a + b
}

// addShort 当多个参数类型相同时，可以简写
func addShort(a, b int) int {
	return a + b
}

// swap 演示值传递：函数内的修改不影响原变量
func swap(a, b int) {
	a, b = b, a
	fmt.Printf("  函数内交换后: a=%d, b=%d\n", a, b)
}

// swapByPointer 使用指针参数实现真正的交换
func swapByPointer(a, b *int) {
	*a, *b = *b, *a
}

// modifySlice 演示切片作为参数：可以修改底层数组
// 切片本身是值传递（复制切片头），但共享底层数组
func modifySlice(s []int) {
	if len(s) > 0 {
		s[0] = 999
	}
}

func demonstrateParameters() {
	fmt.Println("--- 2. 参数传递 ---")

	// 2.1 基本参数
	fmt.Println("基本参数:")
	result := add(3, 5)
	fmt.Printf("  add(3, 5) = %d\n", result)

	// 2.2 值传递机制
	fmt.Println("\n📌 Go 是值传递语言:")
	fmt.Println("  所有参数都是值的副本，函数内修改不影响原变量")

	x, y := 10, 20
	fmt.Printf("\n交换前: x=%d, y=%d\n", x, y)
	swap(x, y)
	fmt.Printf("交换后（原变量）: x=%d, y=%d （未改变）\n", x, y)

	// 2.3 指针参数
	fmt.Println("\n使用指针参数实现真正的交换:")
	fmt.Printf("交换前: x=%d, y=%d\n", x, y)
	swapByPointer(&x, &y)
	fmt.Printf("交换后: x=%d, y=%d （已交换）\n", x, y)

	// 2.4 切片参数
	fmt.Println("\n切片作为参数（共享底层数组）:")
	nums := []int{1, 2, 3}
	fmt.Printf("  修改前: %v\n", nums)
	modifySlice(nums)
	fmt.Printf("  修改后: %v （第一个元素被修改）\n", nums)

	fmt.Println("\n💡 参数传递总结:")
	fmt.Println("  - 基本类型：值的副本，修改无效")
	fmt.Println("  - 指针类型：地址的副本，可通过指针修改原值")
	fmt.Println("  - 切片/映射：头部的副本，共享底层数据")

	fmt.Println()
}

// ============================================================
// 3. 返回值
// ============================================================

// double 单返回值函数
func double(n int) int {
	return n * 2
}

// divide 多返回值函数：返回商和余数
// Go 的多返回值是一个强大的特性，常用于返回结果和错误
func divide(a, b int) (int, int) {
	return a / b, a % b
}

// divideWithError 返回结果和错误（Go 的惯用模式）
func divideWithError(a, b int) (int, error) {
	if b == 0 {
		return 0, fmt.Errorf("除数不能为零")
	}
	return a / b, nil
}

// rectangle 命名返回值：在函数签名中声明返回值变量名
// 命名返回值会被自动初始化为零值
func rectangle(width, height int) (area int, perimeter int) {
	area = width * height
	perimeter = 2 * (width + height)
	// 裸返回：直接 return，返回命名的返回值
	return
}

// rectangleExplicit 显式返回（推荐用于复杂函数）
func rectangleExplicit(width, height int) (area int, perimeter int) {
	area = width * height
	perimeter = 2 * (width + height)
	// 显式返回：更清晰，推荐使用
	return area, perimeter
}

func demonstrateReturnValues() {
	fmt.Println("--- 3. 返回值 ---")

	// 3.1 单返回值
	fmt.Println("单返回值:")
	fmt.Printf("  double(5) = %d\n", double(5))

	// 3.2 多返回值
	fmt.Println("\n多返回值:")
	quotient, remainder := divide(17, 5)
	fmt.Printf("  divide(17, 5) = 商: %d, 余数: %d\n", quotient, remainder)

	// 3.3 忽略部分返回值
	fmt.Println("\n忽略部分返回值（使用 _）:")
	q, _ := divide(17, 5)
	fmt.Printf("  只取商: %d\n", q)

	// 3.4 返回错误（Go 惯用模式）
	fmt.Println("\n返回错误（Go 惯用模式）:")
	if result, err := divideWithError(10, 0); err != nil {
		fmt.Printf("  错误: %v\n", err)
	} else {
		fmt.Printf("  结果: %d\n", result)
	}

	if result, err := divideWithError(10, 3); err != nil {
		fmt.Printf("  错误: %v\n", err)
	} else {
		fmt.Printf("  结果: %d\n", result)
	}

	// 3.5 命名返回值
	fmt.Println("\n命名返回值:")
	area, perimeter := rectangle(5, 3)
	fmt.Printf("  rectangle(5, 3) = 面积: %d, 周长: %d\n", area, perimeter)

	fmt.Println("\n💡 返回值最佳实践:")
	fmt.Println("  - 错误处理使用 (result, error) 模式")
	fmt.Println("  - 命名返回值用于文档说明")
	fmt.Println("  - 复杂函数避免裸返回，显式返回更清晰")

	fmt.Println()
}

// ============================================================
// 4. 可变参数函数
// ============================================================

// sum 接收任意数量的整数并返回总和
// ...int 表示可变参数，函数内部 nums 是 []int 类型
func sum(nums ...int) int {
	total := 0
	for _, n := range nums {
		total += n
	}
	return total
}

// printf 模拟 fmt.Printf 的签名：固定参数 + 可变参数
func printf(format string, args ...interface{}) {
	fmt.Printf("  [自定义] "+format, args...)
}

// joinStrings 连接多个字符串
func joinStrings(sep string, strs ...string) string {
	return strings.Join(strs, sep)
}

func demonstrateVariadicFunctions() {
	fmt.Println("--- 4. 可变参数函数 ---")

	// 4.1 基本用法
	fmt.Println("基本用法:")
	fmt.Printf("  sum() = %d\n", sum())
	fmt.Printf("  sum(1) = %d\n", sum(1))
	fmt.Printf("  sum(1, 2, 3) = %d\n", sum(1, 2, 3))
	fmt.Printf("  sum(1, 2, 3, 4, 5) = %d\n", sum(1, 2, 3, 4, 5))

	// 4.2 传递切片给可变参数函数
	fmt.Println("\n传递切片（使用 ... 展开）:")
	numbers := []int{10, 20, 30, 40}
	// 使用 slice... 语法将切片展开为可变参数
	fmt.Printf("  numbers = %v\n", numbers)
	fmt.Printf("  sum(numbers...) = %d\n", sum(numbers...))

	// 4.3 固定参数 + 可变参数
	fmt.Println("\n固定参数 + 可变参数:")
	printf("Hello, %s! You are %d years old.\n", "Go", 15)

	// 4.4 实用示例
	fmt.Println("\n实用示例 - 字符串连接:")
	result := joinStrings("-", "2024", "01", "15")
	fmt.Printf("  joinStrings(\"-\", \"2024\", \"01\", \"15\") = %s\n", result)

	fmt.Println("\n📌 可变参数规则:")
	fmt.Println("  - 可变参数必须是最后一个参数")
	fmt.Println("  - 函数内部可变参数是切片类型")
	fmt.Println("  - 传递切片时使用 slice... 展开")

	fmt.Println()
}

// ============================================================
// 5. 函数类型与函数作为值
// ============================================================

// Operation 定义一个函数类型
// 接收两个 int 参数，返回一个 int
type Operation func(a, b int) int

// 定义符合 Operation 类型的函数
func addOp(a, b int) int      { return a + b }
func subtractOp(a, b int) int { return a - b }
func multiplyOp(a, b int) int { return a * b }

func demonstrateFunctionTypes() {
	fmt.Println("--- 5. 函数类型与函数作为值 ---")

	// 5.1 函数是一等公民
	fmt.Println("📌 Go 中函数是一等公民（First-class citizen）:")
	fmt.Println("  - 可以赋值给变量")
	fmt.Println("  - 可以作为参数传递")
	fmt.Println("  - 可以作为返回值")
	fmt.Println("  - 可以存储在数据结构中")

	// 5.2 函数赋值给变量
	fmt.Println("\n函数赋值给变量:")
	var op Operation = addOp
	fmt.Printf("  op(10, 5) = %d (使用 addOp)\n", op(10, 5))

	op = subtractOp
	fmt.Printf("  op(10, 5) = %d (使用 subtractOp)\n", op(10, 5))

	// 5.3 函数存储在 map 中
	fmt.Println("\n函数存储在 map 中:")
	operations := map[string]Operation{
		"add":      addOp,
		"subtract": subtractOp,
		"multiply": multiplyOp,
	}

	for name, fn := range operations {
		fmt.Printf("  %s(6, 3) = %d\n", name, fn(6, 3))
	}

	// 5.4 函数类型的零值是 nil
	fmt.Println("\n函数类型的零值:")
	var nilFunc Operation
	fmt.Printf("  nilFunc == nil: %v\n", nilFunc == nil)
	fmt.Println("  调用 nil 函数会 panic，使用前需检查")

	fmt.Println()
}

// ============================================================
// 6. 匿名函数
// ============================================================

func demonstrateAnonymousFunctions() {
	fmt.Println("--- 6. 匿名函数 ---")

	// 匿名函数是没有名字的函数，也叫函数字面量（Function Literal）

	// 6.1 匿名函数赋值给变量
	fmt.Println("匿名函数赋值给变量:")
	square := func(n int) int {
		return n * n
	}
	fmt.Printf("  square(5) = %d\n", square(5))

	// 6.2 立即调用的匿名函数 (IIFE)
	fmt.Println("\n立即调用的匿名函数 (IIFE):")
	result := func(a, b int) int {
		return a + b
	}(10, 20) // 直接在定义后调用
	fmt.Printf("  立即计算 10 + 20 = %d\n", result)

	// 6.3 匿名函数作为 goroutine（预告）
	fmt.Println("\n匿名函数常见使用场景:")
	fmt.Println("  - 作为回调函数传递")
	fmt.Println("  - 作为 goroutine 的执行体")
	fmt.Println("  - 实现闭包（下一章节详解）")
	fmt.Println("  - 延迟执行 (defer)")

	// 6.4 defer 中使用匿名函数
	fmt.Println("\ndefer 中使用匿名函数:")
	func() {
		defer func() {
			fmt.Println("  [defer] 匿名函数在 defer 中执行")
		}()
		fmt.Println("  [normal] 普通语句先执行")
	}()

	fmt.Println()
}

// ============================================================
// 7. 函数作为参数（高阶函数）
// ============================================================

// applyToAll 对切片中的每个元素应用指定函数
// 接收一个函数作为参数，这就是高阶函数
func applyToAll(nums []int, fn func(int) int) []int {
	result := make([]int, len(nums))
	for i, n := range nums {
		result[i] = fn(n)
	}
	return result
}

// filter 过滤切片，保留满足条件的元素
func filter(nums []int, predicate func(int) bool) []int {
	result := []int{}
	for _, n := range nums {
		if predicate(n) {
			result = append(result, n)
		}
	}
	return result
}

// reduce 将切片归约为单个值
func reduce(nums []int, initial int, fn func(acc, curr int) int) int {
	result := initial
	for _, n := range nums {
		result = fn(result, n)
	}
	return result
}

func demonstrateHigherOrderFunctions() {
	fmt.Println("--- 7. 函数作为参数（高阶函数）---")

	// 高阶函数：接收函数作为参数或返回函数的函数

	nums := []int{1, 2, 3, 4, 5}
	fmt.Printf("原始数据: %v\n", nums)

	// 7.1 Map 操作：对每个元素应用转换
	fmt.Println("\nMap 操作 - applyToAll:")
	doubled := applyToAll(nums, func(n int) int {
		return n * 2
	})
	fmt.Printf("  每个元素 * 2: %v\n", doubled)

	squared := applyToAll(nums, func(n int) int {
		return n * n
	})
	fmt.Printf("  每个元素平方: %v\n", squared)

	// 7.2 Filter 操作：过滤元素
	fmt.Println("\nFilter 操作 - filter:")
	evens := filter(nums, func(n int) bool {
		return n%2 == 0
	})
	fmt.Printf("  偶数: %v\n", evens)

	greaterThan2 := filter(nums, func(n int) bool {
		return n > 2
	})
	fmt.Printf("  大于 2 的数: %v\n", greaterThan2)

	// 7.3 Reduce 操作：归约
	fmt.Println("\nReduce 操作 - reduce:")
	sum := reduce(nums, 0, func(acc, curr int) int {
		return acc + curr
	})
	fmt.Printf("  求和: %d\n", sum)

	product := reduce(nums, 1, func(acc, curr int) int {
		return acc * curr
	})
	fmt.Printf("  求积: %d\n", product)

	// 7.4 组合使用
	fmt.Println("\n组合使用（链式操作）:")
	// 先过滤出偶数，再将每个数平方，最后求和
	result := reduce(
		applyToAll(
			filter(nums, func(n int) bool { return n%2 == 0 }),
			func(n int) int { return n * n },
		),
		0,
		func(acc, curr int) int { return acc + curr },
	)
	fmt.Printf("  偶数的平方和: %d (2² + 4² = 4 + 16 = 20)\n", result)

	fmt.Println()
}

// ============================================================
// 8. 函数作为返回值
// ============================================================

// makeMultiplier 返回一个乘法函数
// 这是工厂函数模式
func makeMultiplier(factor int) func(int) int {
	return func(n int) int {
		return n * factor
	}
}

// makeCounter 返回一个计数器函数
// 这是一个典型的闭包：内部函数引用了外部函数的变量 count
func makeCounter() func() int {
	// count 变量发生了"逃逸分析" (Escape Analysis)：
	// 虽然它是在 makeCounter 中定义的局部变量，但因为被返回的闭包引用，
	// 编译器会将它分配到"堆"（Heap）上，而不是"栈"（Stack）上。
	count := 0

	return func() int {
		// 只要这个返回的函数还被持有（引用），堆上的 count 就会一直存在
		count++
		return count
	}
	// 当返回的函数不再被任何变量引用时（例如超出作用域）：demonstrateFunctionReturningFunction() 函数执行完毕
	// Go 的垃圾回收器（GC）会回收这个函数和它捕获的 count 变量。
}

// makeFormatter 返回一个格式化函数
func makeFormatter(prefix, suffix string) func(string) string {
	return func(s string) string {
		return prefix + s + suffix
	}
}

func demonstrateFunctionReturningFunction() {
	fmt.Println("--- 8. 函数作为返回值 ---")

	// 8.1 工厂函数模式
	fmt.Println("工厂函数模式 - makeMultiplier:")
	double := makeMultiplier(2)
	triple := makeMultiplier(3)

	fmt.Printf("  double(5) = %d\n", double(5))
	fmt.Printf("  triple(5) = %d\n", triple(5))

	// 8.2 计数器（闭包预告）
	fmt.Println("\n计数器函数:")
	counter := makeCounter()
	fmt.Printf("  第 1 次调用: %d\n", counter())
	fmt.Printf("  第 2 次调用: %d\n", counter())
	fmt.Printf("  第 3 次调用: %d\n", counter())
	fmt.Println("  (这是闭包的效果，下一章详解)")

	// 8.3 格式化器
	fmt.Println("\n格式化器函数 - makeFormatter:")
	wrapper := makeFormatter("[", "]")
	htmlTag := makeFormatter("<p>", "</p>")

	fmt.Printf("  wrapper(\"Hello\") = %s\n", wrapper("Hello"))
	fmt.Printf("  htmlTag(\"Content\") = %s\n", htmlTag("Content"))

	fmt.Println("\n💡 函数作为返回值的应用:")
	fmt.Println("  - 工厂模式：根据参数生成定制函数")
	fmt.Println("  - 延迟计算：返回的函数包含待执行的逻辑")
	fmt.Println("  - 闭包：捕获外部变量（下一章详解）")

	fmt.Println()
}

// ============================================================
// 9. 递归函数
// ============================================================

// factorial 计算阶乘：n! = n * (n-1) * ... * 1
func factorial(n int) int {
	// 递归终止条件（基准情况）
	if n <= 1 {
		return 1
	}
	// 递归调用
	return n * factorial(n-1)
}

// fibonacci 计算斐波那契数列
// F(n) = F(n-1) + F(n-2)，其中 F(0)=0, F(1)=1
func fibonacci(n int) int {
	if n <= 0 {
		return 0
	}
	if n == 1 {
		return 1
	}
	return fibonacci(n-1) + fibonacci(n-2)
}

// sumRecursive 递归求和
func sumRecursive(nums []int) int {
	// 终止条件：空切片
	if len(nums) == 0 {
		return 0
	}
	// 递归：第一个元素 + 剩余元素的和
	return nums[0] + sumRecursive(nums[1:])
}

func demonstrateRecursion() {
	fmt.Println("--- 9. 递归函数 ---")

	// 递归函数：函数调用自身
	// 必须有终止条件，否则会无限递归导致栈溢出

	fmt.Println("📌 递归的两个要素:")
	fmt.Println("  1. 基准情况（终止条件）")
	fmt.Println("  2. 递归情况（缩小问题规模）")

	// 9.1 阶乘
	fmt.Println("\n阶乘 factorial(n):")
	for i := 0; i <= 5; i++ {
		fmt.Printf("  %d! = %d\n", i, factorial(i))
	}

	// 9.2 斐波那契
	fmt.Println("\n斐波那契数列 fibonacci(n):")
	fmt.Print("  ")
	for i := 0; i <= 10; i++ {
		fmt.Printf("%d ", fibonacci(i))
	}
	fmt.Println()

	// 9.3 递归求和
	fmt.Println("\n递归求和 sumRecursive:")
	nums := []int{1, 2, 3, 4, 5}
	fmt.Printf("  sumRecursive(%v) = %d\n", nums, sumRecursive(nums))

	// 9.4 递归注意事项
	fmt.Println("\n⚠️ 递归注意事项:")
	fmt.Println("  - 必须有明确的终止条件")
	fmt.Println("  - 深度递归可能导致栈溢出")
	fmt.Println("  - 考虑使用尾递归优化或迭代替代")
	fmt.Println("  - 某些问题（如斐波那契）递归效率低，需用记忆化优化")

	fmt.Println()
}

// ============================================================
// 10. defer 延迟执行
// ============================================================

// readFile 模拟文件操作，演示 defer 用于资源清理
func readFile(filename string) {
	fmt.Printf("  打开文件: %s\n", filename)
	// defer 确保函数返回前执行清理操作
	defer fmt.Printf("  关闭文件: %s\n", filename)

	fmt.Printf("  读取文件内容...\n")
	// 即使这里发生错误提前返回，defer 也会执行
}

// deferOrder 演示多个 defer 的执行顺序: 开始 -> 结束 -> defer 3 -> defer 2 -> defer 1
func deferOrder() {
	fmt.Println("  开始")
	defer fmt.Println("  defer 1")
	defer fmt.Println("  defer 2")
	defer fmt.Println("  defer 3")
	fmt.Println("  结束")
}

// deferWithValue 演示 defer 参数的求值时机: 在 defer 行被执行时立即求值，
// 而不是等到函数返回时才取当时的变量值，因为 Go 的 defer 机制设计为能够捕获当下的参数状态（以便资源释放等场景下状态可预测）
func deferWithValue() {
	x := 10
	// defer 的参数在 defer 语句执行时就会求值
	defer fmt.Printf("  defer 时 x = %d\n", x)
	x = 20
	fmt.Printf("  当前 x = %d\n", x)
}

// deferWithClosure 演示 defer 与匿名函数
func deferWithClosure() {
	x := 10
	// 使用匿名函数可以获取函数返回时的值
	defer func() {
		fmt.Printf("  defer 闭包中 x = %d\n", x)
	}()
	x = 20
	fmt.Printf("  当前 x = %d\n", x)
}

// safeDivide 演示 defer + recover 处理 panic
func safeDivide(a, b int) (result int) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("  捕获 panic: %v\n", r)
			result = 0 // 发生 panic 时返回默认值
		}
	}()

	if b == 0 {
		panic("除数不能为零")
	}
	return a / b
}

func demonstrateDefer() {
	fmt.Println("--- 10. defer 延迟执行 ---")

	// defer 会将函数调用推迟到外层函数返回之前执行
	// 常用于：资源清理、解锁、关闭连接等

	// 10.1 基本用法
	fmt.Println("基本用法 - 资源清理:")
	readFile("config.yaml")

	// 10.2 执行顺序（LIFO - 后进先出）
	fmt.Println("\n执行顺序（LIFO - 栈结构）:")
	deferOrder()

	// 10.3 参数求值时机
	fmt.Println("\n参数求值时机:")
	fmt.Println("直接传参（defer 时求值）:")
	deferWithValue()

	fmt.Println("\n使用闭包（返回时求值）:")
	deferWithClosure()

	// 10.4 defer + recover 处理 panic
	fmt.Println("\ndefer + recover 处理 panic:")
	result := safeDivide(10, 0)
	fmt.Printf("  安全除法结果: %d\n", result)

	result = safeDivide(10, 2)
	fmt.Printf("  正常除法结果: %d\n", result)

	fmt.Println("\n💡 defer 最佳实践:")
	fmt.Println("  - 资源获取后立即 defer 释放")
	fmt.Println("  - 注意 LIFO 顺序")
	fmt.Println("  - 了解参数求值时机")
	fmt.Println("  - 配合 recover 处理 panic")

	fmt.Println()
}

// ============================================================
// 11. init 函数
// ============================================================

// init 函数在包加载时自动执行
// 特点：
// - 无参数、无返回值
// - 每个文件可以有多个 init 函数
// - 按文件名和定义顺序执行
// - 在 main 函数之前执行
// - 不能被显式调用

// 本文件的 init 函数
func init() {
	// 包初始化逻辑
	// 例如：配置加载、数据库连接、日志初始化等
	_ = "init 函数已执行"
}

func demonstrateInitExplanation() {
	fmt.Println("--- 11. init 函数 ---")

	fmt.Println("📌 init 函数特点:")
	fmt.Println("  - 自动执行，无需调用")
	fmt.Println("  - 无参数、无返回值")
	fmt.Println("  - 在 main() 之前执行")
	fmt.Println("  - 每个文件可以有多个 init")
	fmt.Println("  - 按依赖顺序执行（先执行被导入包的 init）")

	fmt.Println("\n📌 执行顺序:")
	fmt.Println("  1. 导入的包的 init 函数")
	fmt.Println("  2. 当前包的包级变量初始化")
	fmt.Println("  3. 当前包的 init 函数")
	fmt.Println("  4. main 函数")

	fmt.Println("\n📌 常见用途:")
	fmt.Println("  - 初始化包级变量")
	fmt.Println("  - 注册驱动（如数据库驱动）")
	fmt.Println("  - 运行时检查")
	fmt.Println("  - 配置验证")

	fmt.Println("\n📌 注意事项:")
	fmt.Println("  - 避免在 init 中执行耗时操作")
	fmt.Println("  - 避免 init 间的依赖关系")
	fmt.Println("  - 优先使用显式初始化函数")

	fmt.Println()
	fmt.Println("=== 函数章节演示完成 ===")
}
