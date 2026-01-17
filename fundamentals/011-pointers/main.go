package main

import "fmt"

// =============================================================================
// 1. 失败的交换函数（值传递）
// =============================================================================

// swapWrong 尝试交换两个整数的值，但注定失败。
// Go 语言中，函数参数是"值传递"：a 和 b 是原始变量的副本。
// 在函数内部交换的只是副本，原始变量完全不受影响。
func swapWrong(a, b int) {
	a, b = b, a
	fmt.Printf("  函数内部: a=%d, b=%d (交换成功，但这是副本)\n", a, b)
}

// =============================================================================
// 2. 成功的交换函数（指针传递）
// =============================================================================

// swapRight 通过指针交换两个整数的值。
// 参数类型是 *int，表示"指向 int 的指针"。
// 通过 *p 可以访问指针指向的实际内存位置的值。
func swapRight(a, b *int) {
	// *a 表示"a 这个指针所指向的值"
	// *b 表示"b 这个指针所指向的值"
	*a, *b = *b, *a
	fmt.Printf("  函数内部: *a=%d, *b=%d (通过指针修改了原始变量)\n", *a, *b)
}

// =============================================================================
// 3. 实战场景：修改结构体
// =============================================================================

// User 表示一个用户
type User struct {
	Name  string
	Level int
}

// upgradeWrong 尝试给用户升级，但因为是值传递，修改无效。
// u 是原始 User 的一份完整拷贝。
func upgradeWrong(u User) {
	u.Level++
	fmt.Printf("  函数内部: Level=%d\n", u.Level)
}

// upgradeRight 通过指针修改用户等级。
// u 是指向原始 User 的指针，通过它可以直接修改原始数据。
func upgradeRight(u *User) {
	// Go 语法糖：u.Level 等同于 (*u).Level
	// 编译器会自动处理，无需手动解引用
	u.Level++
	fmt.Printf("  函数内部: Level=%d\n", u.Level)
}

// =============================================================================
// 4. 演示：引用类型（slice）不需要指针也能修改
// =============================================================================

// doubleSlice 将切片中的每个元素翻倍。
// 虽然参数是 []int 而不是 *[]int，但修改依然有效。
// 原因：slice 内部是一个包含指针的结构体，传递的是这个结构体的副本，
// 但副本中的指针仍然指向同一块底层数组。
func doubleSlice(nums []int) {
	for i := range nums {
		nums[i] *= 2
	}
}

func main() {
	fmt.Println("===== 1. 指针基础：& 和 * =====")
	x := 42
	// &x 获取变量 x 的内存地址
	p := &x
	fmt.Printf("变量 x 的值: %d\n", x)
	fmt.Printf("变量 x 的地址 (&x): %p\n", p)
	fmt.Printf("指针 p 指向的值 (*p): %d\n", *p)

	// 通过指针修改原始变量
	*p = 100
	fmt.Printf("通过 *p = 100 修改后，x 的值: %d\n", x)

	fmt.Println("\n===== 2. 值传递的问题：失败的交换 =====")
	a, b := 1, 2
	fmt.Printf("交换前: a=%d, b=%d\n", a, b)
	swapWrong(a, b)
	fmt.Printf("交换后: a=%d, b=%d (没有变化！)\n", a, b)

	fmt.Println("\n===== 3. 指针传递的解决：成功的交换 =====")
	a, b = 1, 2
	fmt.Printf("交换前: a=%d, b=%d\n", a, b)
	// &a 和 &b 传递的是变量的地址
	swapRight(&a, &b)
	fmt.Printf("交换后: a=%d, b=%d (成功交换！)\n", a, b)

	fmt.Println("\n===== 4. 实战场景：修改结构体 =====")
	user := User{Name: "Alice", Level: 1}
	fmt.Printf("升级前: %+v\n", user)

	fmt.Println("--- 使用值传递 (upgradeWrong) ---")
	upgradeWrong(user)
	fmt.Printf("升级后: %+v (没有变化！)\n", user)

	fmt.Println("--- 使用指针传递 (upgradeRight) ---")
	upgradeRight(&user)
	fmt.Printf("升级后: %+v (成功升级！)\n", user)

	fmt.Println("\n===== 5. 引用类型（slice）的特殊性 =====")
	nums := []int{1, 2, 3}
	fmt.Printf("翻倍前: %v\n", nums)
	// 不需要传递 &nums，因为 slice 内部已经包含指向底层数组的指针
	doubleSlice(nums)
	fmt.Printf("翻倍后: %v (修改成功，尽管没用指针)\n", nums)

	fmt.Println("\n===== 6. nil 指针 =====")
	var nilPtr *int
	fmt.Printf("未初始化的指针值: %v\n", nilPtr)
	fmt.Printf("nilPtr == nil ? %t\n", nilPtr == nil)

	// 使用 nil 指针前必须检查，否则会 panic
	if nilPtr != nil {
		fmt.Println(*nilPtr)
	} else {
		fmt.Println("指针为 nil，不能解引用")
	}
}
