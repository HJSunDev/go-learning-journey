package main

import "fmt"

func main() {
	fmt.Println("========== Go 类型系统核心演示 ==========")

	// 各类型的值展示
	demoTypeValues()

	// 值类型 vs 引用类型
	demoValueVsReference()

	// 数组 vs 切片
	demoArrayVsSlice()

	// 结构体是值类型，想共享用指针
	demoStructAndPointer()

}

// 展示各类型的值
func demoTypeValues() {
	fmt.Println("各类型的值展示")
	fmt.Println("─────────────────────────────")

	// 基本类型
	var b bool = true
	var i int = 42
	var f float64 = 3.14159
	var s string = "Hello"
	var r rune = '中'

	fmt.Printf("bool:    %v\n", b)
	fmt.Printf("int:     %d\n", i)
	fmt.Printf("float64: %.5f\n", f)
	fmt.Printf("string:  %s\n", s)
	fmt.Printf("rune:    %c (码点: %d)\n", r, r)

	// 整数细分
	var i8 int8 = 127
	var u8 uint8 = 255
	fmt.Printf("int8:    %d (范围: -128~127)\n", i8)
	fmt.Printf("uint8:   %d (范围: 0~255)\n", u8)

	// 零值示例
	var zeroInt int
	var zeroString string
	var zeroSlice []int
	fmt.Printf("零值 int:    %d\n", zeroInt)
	fmt.Printf("零值 string: %q\n", zeroString)
	fmt.Printf("零值 slice:  %v\n", zeroSlice)

	// 引用类型：映射、通道、指针
	m := map[string]int{"Alice": 100, "Bob": 90}
	ch := make(chan int, 1)
	x := 42
	ptr := &x

	fmt.Printf("映射 map:    %v\n", m)
	fmt.Printf("通道 chan:   %T (已创建)\n", ch)
	fmt.Printf("指针 *int:   地址=%p, 值=%d\n", ptr, *ptr)
}

// 值类型复制数据，引用类型共享数据
func demoValueVsReference() {
	fmt.Println("值类型 vs 引用类型")
	fmt.Println("─────────────────────────────")

	// 值类型：复制后互不影响
	a := 10
	b := a
	b = 99
	fmt.Printf("值类型: a=%d, b=%d (互不影响)\n", a, b)

	// 引用类型：复制后共享数据
	s1 := []int{1, 2, 3}
	s2 := s1
	s2[0] = 99
	fmt.Printf("引用类型: s1=%v, s2=%v (共享数据)\n\n", s1, s2)
}

// 数组是值类型，切片是引用类型
func demoArrayVsSlice() {
	fmt.Println("数组 vs 切片")
	fmt.Println("─────────────────────────────")

	// 数组 [3]int：值类型
	arr1 := [3]int{1, 2, 3}
	arr2 := arr1
	arr2[0] = 99
	fmt.Printf("数组: arr1=%v, arr2=%v (独立副本)\n", arr1, arr2)

	// 切片 []int：引用类型
	slc1 := []int{1, 2, 3}
	slc2 := slc1
	slc2[0] = 99
	fmt.Printf("切片: slc1=%v, slc2=%v (共享底层)\n\n", slc1, slc2)
}

// 结构体默认值复制，用指针实现共享
func demoStructAndPointer() {
	fmt.Println("结构体与指针")
	fmt.Println("─────────────────────────────")

	type User struct{ Name string }

	// 结构体：值复制
	u1 := User{Name: "Alice"}
	u2 := u1
	u2.Name = "Bob"
	fmt.Printf("结构体: u1=%s, u2=%s (独立副本)\n", u1.Name, u2.Name)

	// 指针：共享
	u3 := User{Name: "Alice"}
	u4 := &u3
	u4.Name = "Carol"
	fmt.Printf("指针:   u3=%s, u4=%s (共享)\n", u3.Name, u4.Name)
}
