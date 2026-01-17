package main

import (
	"fmt"
	"strings"
	"unicode/utf8"
)

func main() {
	fmt.Println("=== Go Range 迭代器演示 ===")
	fmt.Println()

	// 1. Range 基础语法
	demonstrateRangeBasics()

	// 2. 切片遍历：传统 for vs range
	demonstrateSliceIteration()

	// 3. 数组遍历
	demonstrateArrayIteration()

	// 4. 映射遍历
	demonstrateMapIteration()

	// 5. 字符串遍历
	demonstrateStringIteration()

	// 6. strings 标准库函数
	demonstrateStringsPackage()

	// 7. Range 常见陷阱
	demonstrateRangePitfalls()

	// 8. 通道遍历（预告）
	demonstrateChannelRange()
}

// demonstrateRangeBasics 演示 range 的基础语法
func demonstrateRangeBasics() {
	fmt.Println("--- 1. Range 基础语法 ---")

	// range 是 Go 的迭代器，可用于遍历多种数据结构
	// 根据数据类型不同，range 返回的值也不同：
	//
	// | 数据类型 | 第一个返回值 | 第二个返回值 |
	// |---------|-------------|-------------|
	// | 切片    | 索引 (int)   | 元素值       |
	// | 数组    | 索引 (int)   | 元素值       |
	// | 映射    | 键           | 值           |
	// | 字符串  | 字节索引(int)| rune 值      |
	// | 通道    | 元素值       | 无           |

	nums := []int{10, 20, 30}

	// 1.1 完整形式：获取索引和值
	fmt.Println("完整形式（索引 + 值）:")
	for index, value := range nums {
		fmt.Printf("  索引: %d, 值: %d\n", index, value)
	}

	// 1.2 只获取索引（省略第二个变量）
	fmt.Println("\n只获取索引:")
	for index := range nums {
		fmt.Printf("  索引: %d\n", index)
	}

	// 1.3 只获取值（使用 _ 忽略索引）
	fmt.Println("\n只获取值（使用 _ 忽略索引）:")
	for _, value := range nums {
		fmt.Printf("  值: %d\n", value)
	}

	// 1.4 只遍历，不需要任何值
	fmt.Println("\n只遍历（不使用返回值）:")
	count := 0
	for range nums {
		count++
	}
	fmt.Printf("  遍历了 %d 次\n", count)

	fmt.Println()
}

// demonstrateSliceIteration 演示切片遍历
func demonstrateSliceIteration() {
	fmt.Println("--- 2. 切片遍历：传统 for vs range ---")

	fruits := []string{"苹果", "香蕉", "橙子", "葡萄", "西瓜"}

	// 2.1 传统 for 循环（C 风格）
	fmt.Println("传统 for 循环:")
	for i := 0; i < len(fruits); i++ {
		fmt.Printf("  [%d] %s\n", i, fruits[i])
	}

	// 2.2 for range 循环
	fmt.Println("\nfor range 循环:")
	for i, fruit := range fruits {
		fmt.Printf("  [%d] %s\n", i, fruit)
	}

	// 2.3 对比分析
	fmt.Println("\n📊 对比分析:")
	fmt.Println("  ┌────────────────┬─────────────────────────────────────┐")
	fmt.Println("  │ 方式           │ 适用场景                             │")
	fmt.Println("  ├────────────────┼─────────────────────────────────────┤")
	fmt.Println("  │ 传统 for       │ 需要控制步长、倒序、或复杂索引操作    │")
	fmt.Println("  │ for range      │ 顺序遍历所有元素（推荐，更简洁安全）  │")
	fmt.Println("  └────────────────┴─────────────────────────────────────┘")

	// 2.4 传统 for 的独特能力：步长控制
	fmt.Println("\n传统 for 的独特能力:")
	fmt.Print("  每隔一个元素: ")
	for i := 0; i < len(fruits); i += 2 {
		fmt.Printf("%s ", fruits[i])
	}
	fmt.Println()

	// 2.5 传统 for 的独特能力：倒序遍历
	fmt.Print("  倒序遍历: ")
	for i := len(fruits) - 1; i >= 0; i-- {
		fmt.Printf("%s ", fruits[i])
	}
	fmt.Println()

	// 2.6 range 的优势：更安全
	fmt.Println("\n✅ range 的优势:")
	fmt.Println("  - 不会出现索引越界")
	fmt.Println("  - 代码更简洁易读")
	fmt.Println("  - 自动处理空切片")

	// 空切片测试：演示 range 可以安全处理 nil 切片
	var empty []int
	fmt.Print("  遍历空切片: ")
	for _, v := range empty {
		// 循环体不会执行，也不会 panic
		fmt.Print(v)
	}
	fmt.Println("（无输出，安全通过）")

	fmt.Println()
}

// demonstrateArrayIteration 演示数组遍历
func demonstrateArrayIteration() {
	fmt.Println("--- 3. 数组遍历 ---")

	// 数组的 range 遍历方式与切片完全相同
	weekdays := [5]string{"周一", "周二", "周三", "周四", "周五"}

	fmt.Println("遍历数组:")
	for i, day := range weekdays {
		fmt.Printf("  [%d] %s\n", i, day)
	}

	// 数组作为参数传递时的区别
	fmt.Println("\n💡 数组 vs 切片的遍历区别:")
	fmt.Println("  - 数组 range 是对【数组副本】进行遍历")
	fmt.Println("  - 切片 range 是对【底层数组】进行遍历")

	// 3.1 数组：修改原数组不影响后续遍历（因为遍历的是副本）
	fmt.Println("\n测试1: 数组（遍历副本）")
	arr := [3]int{1, 2, 3}
	fmt.Printf("  初始数组: %v\n", arr)
	for i, v := range arr {
		if i == 0 {
			// 在遍历第1个元素时，修改原数组的第2个元素
			arr[1] = 100
			fmt.Println("  -> i=0 时修改 arr[1] = 100")
		}
		// 观察 i=1 时，v 是旧值(2)还是新值(100)？
		fmt.Printf("  遍历 i=%d, v=%d\n", i, v)
	}
	fmt.Println("  结论: v 保持旧值 2，说明 range 遍历的是数组开始时的副本")

	// 3.2 切片：修改原切片会影响后续遍历（因为共享底层数组）
	fmt.Println("\n测试2: 切片（遍历底层数组）")
	sli := []int{1, 2, 3}
	fmt.Printf("  初始切片: %v\n", sli)
	for i, v := range sli {
		if i == 0 {
			sli[1] = 100
			fmt.Println("  -> i=0 时修改 sli[1] = 100")
		}
		// 观察 i=1 时，v 是旧值(2)还是新值(100)？
		fmt.Printf("  遍历 i=%d, v=%d\n", i, v)
	}
	fmt.Println("  结论: v 变成新值 100，说明 range 实时反映底层数组的变化")

	fmt.Println()
}

// demonstrateMapIteration 演示映射遍历
func demonstrateMapIteration() {
	fmt.Println("--- 4. 映射遍历 ---")

	scores := map[string]int{
		"Alice":   95,
		"Bob":     87,
		"Charlie": 92,
		"Diana":   88,
	}

	// 4.1 遍历键值对
	fmt.Println("遍历键值对:")
	for name, score := range scores {
		fmt.Printf("  %s: %d 分\n", name, score)
	}

	// 4.2 只遍历键
	fmt.Println("\n只遍历键:")
	for name := range scores {
		fmt.Printf("  学生: %s\n", name)
	}

	// 4.3 只遍历值
	fmt.Println("\n只遍历值:")
	total := 0
	for _, score := range scores {
		total += score
	}
	fmt.Printf("  总分: %d, 平均分: %.1f\n", total, float64(total)/float64(len(scores)))

	// 4.4 ⚠️ 重要：Map 遍历顺序是随机的
	fmt.Println("\n⚠️ 多次遍历，顺序不同:")
	for i := 0; i < 3; i++ {
		names := []string{}
		for name := range scores {
			names = append(names, name)
		}
		fmt.Printf("  第 %d 次: %v\n", i+1, names)
	}
	fmt.Println("  Go 故意让 map 遍历顺序随机化，以避免程序依赖特定顺序")

	fmt.Println()
}

// demonstrateStringIteration 演示字符串遍历
func demonstrateStringIteration() {
	fmt.Println("--- 5. 字符串遍历 ---")

	// Go 字符串是 UTF-8 编码的字节序列
	// range 遍历字符串时，自动按 Unicode 码点（rune）解码

	text := "Hello, 世界!"
	fmt.Printf("字符串: %q\n", text)
	fmt.Printf("字节长度 len(): %d\n", len(text))
	fmt.Printf("字符数量 utf8.RuneCountInString(): %d\n", utf8.RuneCountInString(text))

	// 5.1 使用 range 遍历（按 rune）
	fmt.Println("\n使用 range 遍历（按 rune/Unicode 码点）:")
	for i, r := range text {
		fmt.Printf("  字节索引: %2d, 字符: %c, Unicode: U+%04X\n", i, r, r)
	}

	// 5.2 使用传统 for 遍历（按字节）
	fmt.Println("\n使用传统 for 遍历（按字节）:")
	for i := 0; i < len(text); i++ {
		fmt.Printf("  索引: %2d, 字节: 0x%02X\n", i, text[i])
	}

	// 5.3 关键区别说明
	fmt.Println("\n📊 两种遍历的关键区别:")
	fmt.Println("  ┌─────────────┬──────────────────────────────────────┐")
	fmt.Println("  │ 遍历方式    │ 说明                                  │")
	fmt.Println("  ├─────────────┼──────────────────────────────────────┤")
	fmt.Println("  │ range       │ 按 rune 遍历，自动处理多字节 UTF-8    │")
	fmt.Println("  │ 传统 for    │ 按字节遍历，中文等字符会被拆开        │")
	fmt.Println("  └─────────────┴──────────────────────────────────────┘")

	// 5.4 处理中文字符串
	chinese := "中国"
	fmt.Printf("\n中文字符串 %q:\n", chinese)
	fmt.Printf("  字节长度: %d（每个中文占 3 字节）\n", len(chinese))
	fmt.Printf("  字符数量: %d\n", utf8.RuneCountInString(chinese))

	fmt.Print("  range 遍历: ")
	for _, r := range chinese {
		fmt.Printf("%c ", r)
	}
	fmt.Println()

	// 5.5 字符串转切片遍历
	fmt.Println("\n将字符串转换为 rune 切片:")
	runes := []rune(chinese)
	fmt.Printf("  []rune: %v\n", runes)
	fmt.Printf("  长度: %d\n", len(runes))

	// 转换为字节切片
	bytes := []byte(chinese)
	fmt.Printf("  []byte: %v\n", bytes)
	fmt.Printf("  长度: %d\n", len(bytes))

	fmt.Println()
}

// demonstrateStringsPackage 演示 strings 标准库函数
func demonstrateStringsPackage() {
	fmt.Println("--- 6. strings 标准库函数 ---")

	s := "  Hello, Go World!  "
	fmt.Printf("原字符串: %q\n\n", s)

	// 6.1 修剪（Trim）
	fmt.Println("📌 修剪函数:")
	fmt.Printf("  TrimSpace:      %q\n", strings.TrimSpace(s))
	fmt.Printf("  Trim(s, \" !\"):  %q\n", strings.Trim(s, " !"))
	fmt.Printf("  TrimLeft:       %q\n", strings.TrimLeft(s, " "))
	fmt.Printf("  TrimRight:      %q\n", strings.TrimRight(s, " !"))
	fmt.Printf("  TrimPrefix:     %q\n", strings.TrimPrefix(strings.TrimSpace(s), "Hello"))
	fmt.Printf("  TrimSuffix:     %q\n", strings.TrimSuffix(strings.TrimSpace(s), "!"))

	// 6.2 查找
	text := "Go is awesome. Go is fast."
	fmt.Println("\n📌 查找函数:")
	fmt.Printf("  原字符串: %q\n", text)
	fmt.Printf("  Contains(\"awesome\"): %v\n", strings.Contains(text, "awesome"))
	fmt.Printf("  HasPrefix(\"Go\"):     %v\n", strings.HasPrefix(text, "Go"))
	fmt.Printf("  HasSuffix(\".\"):      %v\n", strings.HasSuffix(text, "."))
	fmt.Printf("  Index(\"is\"):         %d（首次出现位置）\n", strings.Index(text, "is"))
	fmt.Printf("  LastIndex(\"is\"):     %d（最后出现位置）\n", strings.LastIndex(text, "is"))
	fmt.Printf("  Count(\"Go\"):         %d（出现次数）\n", strings.Count(text, "Go"))

	// 6.3 转换
	fmt.Println("\n📌 转换函数:")
	sample := "Hello, World"
	fmt.Printf("  原字符串: %q\n", sample)
	fmt.Printf("  ToUpper:   %q\n", strings.ToUpper(sample))
	fmt.Printf("  ToLower:   %q\n", strings.ToLower(sample))
	fmt.Printf("  ToTitle:   %q\n", strings.ToTitle(sample))

	// 6.4 替换
	fmt.Println("\n📌 替换函数:")
	fmt.Printf("  原字符串: %q\n", text)
	fmt.Printf("  Replace(Go, Rust, 1):  %q\n", strings.Replace(text, "Go", "Rust", 1))
	fmt.Printf("  Replace(Go, Rust, -1): %q\n", strings.Replace(text, "Go", "Rust", -1))
	fmt.Printf("  ReplaceAll(Go, Rust):  %q\n", strings.ReplaceAll(text, "Go", "Rust"))

	// 6.5 分割与连接
	fmt.Println("\n📌 分割与连接函数:")
	csv := "apple,banana,orange,grape"
	fmt.Printf("  原字符串: %q\n", csv)

	// Split
	parts := strings.Split(csv, ",")
	fmt.Printf("  Split(\",\"):  %v\n", parts)

	// SplitN
	partsN := strings.SplitN(csv, ",", 2)
	fmt.Printf("  SplitN(\",\", 2): %v\n", partsN)

	// Fields（按空白分割）
	sentence := "  hello   world  go  "
	fields := strings.Fields(sentence)
	fmt.Printf("  Fields(%q): %v\n", sentence, fields)

	// Join
	joined := strings.Join(parts, " | ")
	fmt.Printf("  Join(\" | \"): %q\n", joined)

	// 6.6 重复与填充
	fmt.Println("\n📌 重复函数:")
	fmt.Printf("  Repeat(\"Go\", 3):  %q\n", strings.Repeat("Go ", 3))
	fmt.Printf("  Repeat(\"-\", 20): %q\n", strings.Repeat("-", 20))

	// 6.7 比较
	fmt.Println("\n📌 比较函数:")
	fmt.Printf("  EqualFold(\"GO\", \"go\"): %v（忽略大小写）\n", strings.EqualFold("GO", "go"))
	fmt.Printf("  Compare(\"a\", \"b\"):     %d（-1 表示 a < b）\n", strings.Compare("a", "b"))
	fmt.Printf("  Compare(\"b\", \"a\"):     %d（1 表示 b > a）\n", strings.Compare("b", "a"))
	fmt.Printf("  Compare(\"a\", \"a\"):     %d（0 表示相等）\n", strings.Compare("a", "a"))

	// 6.8 Builder（高效字符串拼接）
	fmt.Println("\n📌 strings.Builder（高效拼接）:")
	var builder strings.Builder
	for i := 0; i < 5; i++ {
		builder.WriteString("Go")
		builder.WriteByte(' ')
	}
	result := builder.String()
	fmt.Printf("  构建结果: %q\n", result)
	fmt.Println("  ✅ Builder 比 + 拼接更高效，避免频繁内存分配")

	fmt.Println()
}

// demonstrateRangePitfalls 演示 range 的常见陷阱
func demonstrateRangePitfalls() {
	fmt.Println("--- 7. Range 常见陷阱 ---")

	// 7.1 陷阱：range 返回的是值的副本
	fmt.Println("陷阱1: range 返回值的副本，修改无效")
	nums := []int{1, 2, 3, 4, 5}
	fmt.Println("  原切片:", nums)

	// 错误方式：修改 v 不会影响原切片
	for _, v := range nums {
		v *= 2 // 这只修改了副本
		_ = v
	}
	fmt.Println("  错误修改后:", nums) // 仍然是 [1 2 3 4 5]

	// 正确方式：使用索引修改
	for i := range nums {
		nums[i] *= 2
	}
	fmt.Println("  正确修改后:", nums) // [2 4 6 8 10]

	// 7.2 陷阱：在 range 中修改切片
	fmt.Println("\n陷阱2: 在 range 中修改切片")
	data := []int{1, 2, 3}
	fmt.Println("  原切片:", data)

	// range 在开始时确定遍历范围，追加的元素不会被遍历
	for i, v := range data {
		if i == 0 {
			data = append(data, 100, 200)
		}
		fmt.Printf("    遍历: 索引=%d, 值=%d\n", i, v)
	}
	fmt.Println("  遍历后切片:", data)
	fmt.Println("  ⚠️ 追加的 100, 200 没有被遍历到")

	// 7.3 陷阱：range map 时删除/添加元素
	fmt.Println("\n陷阱3: range map 时的增删操作")
	fmt.Println("  ✅ 可以安全删除当前遍历的键")
	fmt.Println("  ⚠️ 新添加的键可能被遍历，也可能不被遍历（不确定）")
	fmt.Println("  💡 建议：遍历时避免修改 map，先收集操作再统一执行")

	m := map[string]int{"a": 1, "b": 2, "c": 3}
	fmt.Println("  原 map:", m)
	for k := range m {
		if k == "a" {
			delete(m, k) // 安全：删除当前遍历的键
		}
	}
	fmt.Println("  删除 'a' 后:", m)

	fmt.Println()
}

// demonstrateChannelRange 演示通道遍历（预告）
func demonstrateChannelRange() {
	fmt.Println("--- 8. 通道遍历（预告） ---")

	// 通道（Channel）是 Go 并发编程的核心
	// range 可以持续接收通道的值，直到通道关闭

	// 创建一个缓冲通道
	ch := make(chan int, 3)

	// 发送数据
	ch <- 10
	ch <- 20
	ch <- 30
	close(ch) // 关闭通道后，range 才会结束

	// 使用 range 遍历通道
	fmt.Println("使用 range 遍历通道:")
	for value := range ch {
		fmt.Printf("  接收: %d\n", value)
	}

	// 说明
	fmt.Println("\n💡 通道 range 的特点:")
	fmt.Println("  - range 会阻塞等待通道数据")
	fmt.Println("  - 通道关闭后，range 自动结束")
	fmt.Println("  - 只返回一个值（通道元素），没有索引")
	fmt.Println("  - 详细内容将在并发编程章节介绍")

	fmt.Println()
}
