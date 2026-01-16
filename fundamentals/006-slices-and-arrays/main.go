package main

import (
	"fmt"
	"slices"
)

func main() {
	fmt.Println("=== Go 切片与数组演示 ===")
	fmt.Println()

	// 1. 切片基础
	demonstrateSliceBasics()

	// 2. 切片操作符
	demonstrateSliceOperator()

	// 3. append 追加与扩容
	demonstrateAppend()

	// 4. slices 标准库函数
	demonstrateSlicesPackage()

	// 5. 二维切片
	demonstrate2DSlice()

	// 6. 切片陷阱
	demonstrateSlicePitfalls()

	// 7. 数组
	demonstrateArray()

	// 8. 数组的真实使用场景
	demonstrateArrayUseCases()
}

// demonstrateSliceBasics 演示切片的创建和基本操作
func demonstrateSliceBasics() {
	fmt.Println("--- 1. 切片基础 ---")

	// 1.1 字面量创建
	nums := []int{10, 20, 30, 40, 50}
	fmt.Println("字面量创建:", nums)

	// 1.2 make 创建：指定长度
	// make([]T, length) - 创建指定长度的切片，元素为零值
	s1 := make([]int, 5)
	fmt.Println("make([]int, 5):", s1) // [0 0 0 0 0]

	// 1.3 make 创建：指定长度和容量
	// make([]T, length, capacity) - 预分配容量，避免频繁扩容
	s2 := make([]int, 0, 10)
	fmt.Printf("make([]int, 0, 10): %v, len=%d, cap=%d\n", s2, len(s2), cap(s2))

	// 1.4 len() 和 cap()
	fmt.Printf("nums: len=%d, cap=%d\n", len(nums), cap(nums))

	// 1.5 索引访问
	fmt.Println("第一个元素 nums[0]:", nums[0])
	fmt.Println("最后一个元素 nums[len(nums)-1]:", nums[len(nums)-1])

	// 1.6 修改元素
	nums[0] = 100
	fmt.Println("修改后:", nums)

	fmt.Println()
}

// demonstrateSliceOperator 演示切片操作符
func demonstrateSliceOperator() {
	fmt.Println("--- 2. 切片操作符 [start:end] ---")

	nums := []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
	fmt.Println("原切片:", nums)

	// 2.1 基本切片操作
	// [start:end] - 从 start 到 end-1（左闭右开）
	fmt.Println("nums[2:5]:", nums[2:5]) // [2 3 4]

	// 2.2 省略 start：从头开始
	fmt.Println("nums[:3]:", nums[:3]) // [0 1 2]

	// 2.3 省略 end：到末尾
	fmt.Println("nums[7:]:", nums[7:]) // [7 8 9]

	// 2.4 省略 start 和 end：完整复制（但共享底层数组）
	fmt.Println("nums[:]:", nums[:]) // [0 1 2 3 4 5 6 7 8 9]

	// 2.5 负数索引？Go 不支持！
	// nums[-1] // ❌ 编译错误

	// 2.6 完整切片表达式 [start:end:max]
	// max 限制新切片的容量为 max-start
	sub := nums[2:5:7]
	fmt.Printf("nums[2:5:7]: %v, len=%d, cap=%d\n", sub, len(sub), cap(sub))
	// len = 5-2 = 3, cap = 7-2 = 5

	fmt.Println()
}

// demonstrateAppend 演示 append 函数和扩容机制
func demonstrateAppend() {
	fmt.Println("--- 3. append 追加与扩容 ---")

	// 3.1 基本追加
	s := []int{1, 2, 3}
	s = append(s, 4)
	fmt.Println("追加一个:", s) // [1 2 3 4]

	// 3.2 追加多个元素
	s = append(s, 5, 6, 7)
	fmt.Println("追加多个:", s) // [1 2 3 4 5 6 7]

	// 3.3 追加另一个切片（使用 ... 展开）
	extra := []int{8, 9, 10}
	s = append(s, extra...)
	fmt.Println("追加切片:", s) // [1 2 3 4 5 6 7 8 9 10]

	// 3.4 观察扩容
	fmt.Println("\n扩容演示:")
	demo := make([]int, 0)
	for i := 1; i <= 10; i++ {
		oldCap := cap(demo)
		demo = append(demo, i)
		newCap := cap(demo)
		if newCap != oldCap {
			fmt.Printf("  元素数: %d, 容量: %d -> %d\n", len(demo), oldCap, newCap)
		}
	}

	// 3.5 预分配容量的最佳实践
	// 已知大小时，使用 make 预分配可避免多次扩容
	fmt.Println("\n预分配容量（推荐）:")
	size := 1000
	efficient := make([]int, 0, size)
	for i := 0; i < size; i++ {
		efficient = append(efficient, i)
	}
	fmt.Printf("  预分配: len=%d, cap=%d\n", len(efficient), cap(efficient))

	fmt.Println()
}

// demonstrateSlicesPackage 演示 slices 标准库（Go 1.21+）
func demonstrateSlicesPackage() {
	fmt.Println("--- 4. slices 标准库函数 ---")

	// 4.1 slices.Equal 比较切片
	a := []int{1, 2, 3}
	b := []int{1, 2, 3}
	c := []int{1, 2, 4}

	fmt.Printf("slices.Equal(%v, %v): %v\n", a, b, slices.Equal(a, b)) // true
	fmt.Printf("slices.Equal(%v, %v): %v\n", a, c, slices.Equal(a, c)) // false

	// 注意：切片不能用 == 直接比较（只能与 nil 比较）
	// a == b // ❌ 编译错误

	// 4.2 slices.Sort 排序
	unsorted := []int{3, 1, 4, 1, 5, 9, 2, 6}
	slices.Sort(unsorted)
	fmt.Println("slices.Sort:", unsorted) // [1 1 2 3 4 5 6 9]

	// 4.3 slices.Contains 检查是否包含
	nums := []int{10, 20, 30, 40, 50}
	fmt.Printf("slices.Contains(%v, 30): %v\n", nums, slices.Contains(nums, 30)) // true
	fmt.Printf("slices.Contains(%v, 99): %v\n", nums, slices.Contains(nums, 99)) // false

	// 4.4 slices.Index 查找索引
	fmt.Printf("slices.Index(%v, 30): %v\n", nums, slices.Index(nums, 30)) // 2
	fmt.Printf("slices.Index(%v, 99): %v\n", nums, slices.Index(nums, 99)) // -1（未找到）

	// 4.5 slices.Reverse 反转
	toReverse := []int{1, 2, 3, 4, 5}
	slices.Reverse(toReverse)
	fmt.Println("slices.Reverse:", toReverse) // [5 4 3 2 1]

	// 4.6 slices.Clone 深拷贝
	original := []int{1, 2, 3}
	cloned := slices.Clone(original)
	cloned[0] = 999
	fmt.Println("original:", original) // [1 2 3] - 不受影响
	fmt.Println("cloned:", cloned)     // [999 2 3]

	// 4.7 slices.Max / slices.Min
	values := []int{5, 2, 8, 1, 9, 3}
	fmt.Printf("slices.Max(%v): %v\n", values, slices.Max(values)) // 9
	fmt.Printf("slices.Min(%v): %v\n", values, slices.Min(values)) // 1

	fmt.Println()
}

// demonstrate2DSlice 演示二维切片
func demonstrate2DSlice() {
	fmt.Println("--- 5. 二维切片 ---")

	// 5.1 创建规则二维切片
	rows, cols := 3, 4

	// 方法一：逐行分配
	matrix := make([][]int, rows)
	for i := range matrix {
		matrix[i] = make([]int, cols)
	}

	// 填充数据
	for i := 0; i < rows; i++ {
		for j := 0; j < cols; j++ {
			matrix[i][j] = i*cols + j + 1
		}
	}

	fmt.Println("3x4 矩阵:")
	for _, row := range matrix {
		fmt.Println(" ", row)
	}

	// 5.2 访问元素
	fmt.Printf("matrix[1][2] = %d\n", matrix[1][2]) // 第2行第3列

	// 5.3 字面量创建二维切片
	grid := [][]int{
		{1, 2, 3},
		{4, 5, 6},
		{7, 8, 9},
	}
	fmt.Println("\n字面量创建:")
	for _, row := range grid {
		fmt.Println(" ", row)
	}

	// 5.4 不规则二维切片（锯齿数组）
	// Go 的二维切片本质是"切片的切片"，每行可以有不同长度
	jagged := [][]int{
		{1},
		{2, 3},
		{4, 5, 6},
		{7, 8, 9, 10},
	}
	fmt.Println("\n不规则切片（锯齿数组）:")
	for i, row := range jagged {
		fmt.Printf("  第%d行（长度%d）: %v\n", i, len(row), row)
	}

	fmt.Println()
}

// demonstrateSlicePitfalls 演示切片的常见陷阱
func demonstrateSlicePitfalls() {
	fmt.Println("--- 6. 切片陷阱 ---")

	// 6.1 陷阱：子切片共享底层数组
	fmt.Println("陷阱1: 子切片共享底层数组")
	original := []int{1, 2, 3, 4, 5}
	sub := original[1:4] // [2, 3, 4]

	fmt.Println("  修改前:")
	fmt.Println("    original:", original)
	fmt.Println("    sub:", sub)

	sub[0] = 999 // 修改 sub 也会影响 original

	fmt.Println("  修改 sub[0] = 999 后:")
	fmt.Println("    original:", original) // [1 999 3 4 5]
	fmt.Println("    sub:", sub)           // [999 3 4]

	// 解决方案：使用 copy 或 slices.Clone
	fmt.Println("\n  解决方案：使用 slices.Clone")
	original2 := []int{1, 2, 3, 4, 5}
	safeCopy := slices.Clone(original2[1:4])
	safeCopy[0] = 999
	fmt.Println("    original2:", original2) // [1 2 3 4 5] - 不受影响
	fmt.Println("    safeCopy:", safeCopy)   // [999 3 4]

	// 6.2 陷阱：nil 切片 vs 空切片
	fmt.Println("\n陷阱2: nil 切片 vs 空切片")
	var nilSlice []int          // nil 切片
	emptySlice := []int{}       // 空切片（非 nil）
	makeEmpty := make([]int, 0) // 也是空切片（非 nil）

	fmt.Printf("  nilSlice:   %v, len=%d, nil=%v\n", nilSlice, len(nilSlice), nilSlice == nil)
	fmt.Printf("  emptySlice: %v, len=%d, nil=%v\n", emptySlice, len(emptySlice), emptySlice == nil)
	fmt.Printf("  makeEmpty:  %v, len=%d, nil=%v\n", makeEmpty, len(makeEmpty), makeEmpty == nil)

	// 好消息：nil 切片可以正常使用 append、len、cap
	nilSlice = append(nilSlice, 1, 2, 3)
	fmt.Println("  append 后 nilSlice:", nilSlice)

	// 6.3 陷阱：append 可能返回新切片
	fmt.Println("\n陷阱3: append 可能返回新切片")
	s := make([]int, 3, 3) // len=3, cap=3（已满）
	s[0], s[1], s[2] = 1, 2, 3

	fmt.Printf("  append 前: %v, cap=%d\n", s, cap(s))

	// append 返回新切片，原切片不变（如果发生扩容）
	s2 := append(s, 4)
	fmt.Printf("  append 后 s:  %v, cap=%d\n", s, cap(s))
	fmt.Printf("  append 后 s2: %v, cap=%d\n", s2, cap(s2))

	s[0] = 999 // 修改 s 不会影响 s2（因为已扩容）
	fmt.Println("  修改 s[0]=999 后:")
	fmt.Printf("    s:  %v\n", s)  // [999 2 3]
	fmt.Printf("    s2: %v\n", s2) // [1 2 3 4]

	fmt.Println()
}

// demonstrateArray 演示数组基础
func demonstrateArray() {
	fmt.Println("--- 7. 数组 ---")

	// 7.1 数组声明
	var arr1 [5]int                     // 零值初始化
	arr2 := [5]int{1, 2, 3, 4, 5}       // 字面量
	arr3 := [...]int{1, 2, 3}           // 编译器推断长度
	arr4 := [5]int{0: 10, 2: 30, 4: 50} // 指定索引初始化

	fmt.Println("零值数组:", arr1)
	fmt.Println("字面量:", arr2)
	fmt.Println("[...]推断长度:", arr3, "长度:", len(arr3))
	fmt.Println("指定索引:", arr4)

	// 7.2 数组是值类型
	fmt.Println("\n数组是值类型（赋值会复制）:")
	a := [3]int{1, 2, 3}
	b := a // 复制整个数组
	b[0] = 999
	fmt.Println("  a:", a) // [1 2 3] - 不受影响
	fmt.Println("  b:", b) // [999 2 3]

	// 7.3 数组长度是类型的一部分
	// [3]int 和 [5]int 是不同类型！
	var x [3]int
	var y [5]int
	fmt.Printf("\n[3]int 类型: %T\n", x)
	fmt.Printf("[5]int 类型: %T\n", y)
	// x = y // ❌ 编译错误：类型不匹配

	// 7.4 数组可以用 == 比较（切片不行）
	arr5 := [3]int{1, 2, 3}
	arr6 := [3]int{1, 2, 3}
	arr7 := [3]int{1, 2, 4}
	fmt.Println("\n数组比较:")
	fmt.Printf("  %v == %v: %v\n", arr5, arr6, arr5 == arr6) // true
	fmt.Printf("  %v == %v: %v\n", arr5, arr7, arr5 == arr7) // false

	// 7.5 数组遍历
	fmt.Println("\n数组遍历:")
	for i, v := range arr2 {
		fmt.Printf("  索引 %d: %d\n", i, v)
	}

	fmt.Println()
}

// demonstrateArrayUseCases 演示数组的真实使用场景
func demonstrateArrayUseCases() {
	fmt.Println("--- 8. 数组的真实使用场景 ---")

	// 场景1: RGB 颜色值（固定3个分量）
	fmt.Println("场景1: RGB 颜色")
	type RGB [3]uint8
	red := RGB{255, 0, 0}
	green := RGB{0, 255, 0}
	blue := RGB{0, 0, 255}
	fmt.Printf("  红: RGB%v, 绿: RGB%v, 蓝: RGB%v\n", red, green, blue)

	// 场景2: 坐标点（固定维度）
	fmt.Println("\n场景2: 坐标点")
	type Point2D [2]float64
	type Point3D [3]float64

	p2 := Point2D{3.5, 4.5}
	p3 := Point3D{1.0, 2.0, 3.0}
	fmt.Printf("  2D点: %v, 3D点: %v\n", p2, p3)

	// 场景3: 密码学哈希值（SHA-256 是 32 字节）
	fmt.Println("\n场景3: 密码学哈希")
	type SHA256Hash [32]byte
	// 实际使用中由 crypto/sha256 包生成
	var hash SHA256Hash
	hash[0] = 0xab
	hash[1] = 0xcd
	fmt.Printf("  SHA256 哈希（部分）: %x...\n", hash[:4])

	// 场景4: 固定大小的缓冲区
	fmt.Println("\n场景4: 固定大小缓冲区")
	// 在栈上分配，避免堆分配开销
	var buffer [4096]byte
	copy(buffer[:], "Hello, World!")
	fmt.Printf("  缓冲区内容: %s\n", buffer[:13])
	fmt.Printf("  缓冲区大小: %d bytes\n", len(buffer))

	// 场景5: IPv4 地址（固定4字节）
	fmt.Println("\n场景5: IPv4 地址")
	type IPv4 [4]byte
	localhost := IPv4{127, 0, 0, 1}
	gateway := IPv4{192, 168, 1, 1}
	fmt.Printf("  本机: %d.%d.%d.%d\n", localhost[0], localhost[1], localhost[2], localhost[3])
	fmt.Printf("  网关: %d.%d.%d.%d\n", gateway[0], gateway[1], gateway[2], gateway[3])

	// 场景6: 星期（固定7天）
	fmt.Println("\n场景6: 一周的日程")
	weekSchedule := [7]string{
		"周一: 开会",
		"周二: 编码",
		"周三: 编码",
		"周四: 代码评审",
		"周五: 部署",
		"周六: 休息",
		"周日: 休息",
	}
	for i, schedule := range weekSchedule {
		if i < 3 { // 只打印前3天
			fmt.Printf("  %s\n", schedule)
		}
	}
	fmt.Println("  ...")

	fmt.Println("\n💡 数组使用原则:")
	fmt.Println("  1. 大小在编译时已知且固定不变")
	fmt.Println("  2. 需要值语义（赋值即复制）")
	fmt.Println("  3. 需要作为 map 的键（切片不能作为键）")
	fmt.Println("  4. 性能敏感场景（避免堆分配）")
	fmt.Println()
}
