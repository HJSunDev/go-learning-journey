package main

import (
	"fmt"
	"maps"
	"slices"
)

func main() {
	fmt.Println("=== Go 映射（Map）演示 ===")
	fmt.Println()

	// 1. Map 创建
	demonstrateMapCreation()

	// 2. Map 基本操作
	demonstrateMapOperations()

	// 3. 键存在性检查（comma-ok 模式）
	demonstrateCommaOkPattern()

	// 4. Map 遍历
	demonstrateMapIteration()

	// 5. Map 键类型要求
	demonstrateMapKeyTypes()

	// 6. maps 标准库函数（Go 1.21+）
	demonstrateMapsPackage()

	// 7. Map 常见陷阱
	demonstrateMapPitfalls()

}

// demonstrateMapCreation 演示 Map 的创建方式
func demonstrateMapCreation() {
	fmt.Println("--- 1. Map 创建 ---")

	// 1.1 字面量创建
	// map[KeyType]ValueType{key1: value1, key2: value2}
	scores := map[string]int{
		"Alice": 95,
		"Bob":   87,
		"Carol": 92,
	}
	fmt.Println("字面量创建:", scores)

	// 1.2 空 map 字面量
	emptyMap := map[string]int{}
	fmt.Println("空 map 字面量:", emptyMap, "len:", len(emptyMap))

	// 1.3 make 创建
	// make(map[KeyType]ValueType) - 无初始容量
	// make(map[KeyType]ValueType, capacity) - 指定初始容量
	userAges := make(map[string]int)
	fmt.Println("make 创建:", userAges)

	// 1.4 make 预分配容量（性能优化）
	// 已知元素数量时，预分配可减少内存重新分配
	largeMap := make(map[string]int, 1000)
	fmt.Println("预分配容量:", len(largeMap), "(容量无法直接获取)")

	// 1.5 nil map（零值）
	var nilMap map[string]int
	fmt.Printf("nil map: %v, nil=%v\n", nilMap, nilMap == nil)
	// ⚠️ nil map 可以读取，但不能写入（会 panic）

	fmt.Println()
}

// demonstrateMapOperations 演示 Map 的基本操作
func demonstrateMapOperations() {
	fmt.Println("--- 2. Map 基本操作 ---")

	// 创建一个 map
	inventory := map[string]int{
		"apple":  50,
		"banana": 30,
		"orange": 45,
	}

	// 2.1 读取元素
	fmt.Println("读取 apple 数量:", inventory["apple"])

	// 2.2 读取不存在的键 → 返回值类型的零值
	fmt.Println("读取不存在的键:", inventory["grape"]) // 0（int 的零值）

	// 2.3 添加/修改元素
	inventory["grape"] = 20  // 添加新键
	inventory["apple"] = 100 // 修改已有键
	fmt.Println("添加/修改后:", inventory)

	// 2.4 删除元素：使用 delete() 内置函数
	// delete(map, key)
	delete(inventory, "banana")
	fmt.Println("删除 banana 后:", inventory)

	// 2.5 删除不存在的键 → 不报错，静默忽略
	delete(inventory, "mango") // 什么都不会发生
	fmt.Println("删除不存在的键后:", inventory)

	// 2.6 获取 map 长度
	fmt.Println("元素数量:", len(inventory))

	// 2.7 clear() 清空 map（Go 1.21+）
	temp := map[string]int{"a": 1, "b": 2}
	fmt.Println("clear 前:", temp)
	clear(temp)
	fmt.Println("clear 后:", temp)

	fmt.Println()
}

// demonstrateCommaOkPattern 演示键存在性检查（comma-ok 模式）
func demonstrateCommaOkPattern() {
	fmt.Println("--- 3. 键存在性检查（comma-ok 模式） ---")

	userScores := map[string]int{
		"Alice": 95,
		"Bob":   0, // 注意：Bob 的分数是 0
	}

	// 3.1 问题：无法区分"键不存在"和"值为零值"
	fmt.Println("问题演示:")
	fmt.Println("  Alice 分数:", userScores["Alice"]) // 95
	fmt.Println("  Bob 分数:", userScores["Bob"])     // 0（真实值）
	fmt.Println("  Eve 分数:", userScores["Eve"])     // 0（键不存在）

	// 3.2 解决方案：comma-ok 模式
	// value, ok := map[key]
	// ok 为 true 表示键存在，false 表示不存在
	fmt.Println("\ncomma-ok 模式:")

	// 检查 Alice
	if score, ok := userScores["Alice"]; ok {
		fmt.Printf("  Alice 存在，分数: %d\n", score)
	}

	// 检查 Bob（分数为 0）
	if score, ok := userScores["Bob"]; ok {
		fmt.Printf("  Bob 存在，分数: %d\n", score)
	}

	// 检查 Eve（不存在）
	if score, ok := userScores["Eve"]; ok {
		fmt.Printf("  Eve 存在，分数: %d\n", score)
	} else {
		fmt.Println("  Eve 不存在")
	}

	// 3.3 只检查存在性，忽略值（使用 _ 空白标识符）
	fmt.Println("\n只检查存在性（使用 _ 忽略值）:")
	if _, exists := userScores["Alice"]; exists {
		fmt.Println("  Alice 在名单中")
	}
	if _, exists := userScores["Eve"]; !exists {
		fmt.Println("  Eve 不在名单中")
	}

	fmt.Println()
}

// demonstrateMapIteration 演示 Map 遍历
func demonstrateMapIteration() {
	fmt.Println("--- 4. Map 遍历 ---")

	capitals := map[string]string{
		"China":  "Beijing",
		"Japan":  "Tokyo",
		"France": "Paris",
		"USA":    "Washington",
	}

	// 4.1 遍历键值对
	fmt.Println("遍历键值对:")
	for country, capital := range capitals {
		fmt.Printf("  %s → %s\n", country, capital)
	}

	// 4.2 只遍历键
	fmt.Println("\n只遍历键:")
	for country := range capitals {
		fmt.Printf("  %s\n", country)
	}

	// 4.3 只遍历值（使用 _ 忽略键）
	fmt.Println("\n只遍历值:")
	for _, capital := range capitals {
		fmt.Printf("  %s\n", capital)
	}

	// 4.4 ⚠️ 遍历顺序是随机的！
	fmt.Println("\n⚠️ 多次遍历，顺序可能不同:")
	for i := 0; i < 3; i++ {
		keys := []string{}
		for k := range capitals {
			keys = append(keys, k)
		}
		fmt.Printf("  第%d次: %v\n", i+1, keys)
	}

	// 4.5 如需有序遍历，先对键排序
	fmt.Println("\n有序遍历（先排序键）:")
	sortedKeys := make([]string, 0, len(capitals))
	for k := range capitals {
		sortedKeys = append(sortedKeys, k)
	}
	slices.Sort(sortedKeys)
	for _, k := range sortedKeys {
		fmt.Printf("  %s → %s\n", k, capitals[k])
	}

	fmt.Println()
}

// demonstrateMapKeyTypes 演示 Map 键类型要求
func demonstrateMapKeyTypes() {
	fmt.Println("--- 5. Map 键类型要求 ---")

	// 5.1 ✅ 可作为键的类型：所有可比较类型
	// 基本类型：int, string, float64, bool 等
	intKeyMap := map[int]string{1: "one", 2: "two"}
	fmt.Println("int 键:", intKeyMap)

	// 5.2 ✅ 数组可以作为键（长度固定，可比较）
	type Point [2]int
	pointMap := map[Point]string{
		{0, 0}: "origin",
		{1, 2}: "point A",
	}
	fmt.Println("数组键:", pointMap)

	// 5.3 ✅ 结构体可以作为键（如果所有字段都可比较）
	type Coordinate struct {
		X, Y int
	}
	coordMap := map[Coordinate]string{
		{0, 0}: "起点",
		{5, 5}: "终点",
	}
	fmt.Println("结构体键:", coordMap)

	// 5.4 ❌ 切片、map、函数不能作为键
	// 以下代码会编译错误：
	// invalidMap := map[[]int]string{}  // ❌ 切片不可比较
	// invalidMap := map[map[string]int]string{}  // ❌ map 不可比较

	fmt.Println("\n✅ 可作为键的类型:")
	fmt.Println("  - 基本类型: int, string, float64, bool, ...")
	fmt.Println("  - 数组: [n]T（长度固定）")
	fmt.Println("  - 结构体: struct（所有字段可比较）")
	fmt.Println("  - 指针: *T")
	fmt.Println("  - 接口: interface{}（如果动态值可比较）")

	fmt.Println("\n❌ 不能作为键的类型:")
	fmt.Println("  - 切片: []T")
	fmt.Println("  - 映射: map[K]V")
	fmt.Println("  - 函数: func()")

	fmt.Println()
}

// demonstrateMapsPackage 演示 maps 标准库（Go 1.21+）
func demonstrateMapsPackage() {
	fmt.Println("--- 6. maps 标准库（Go 1.21+） ---")

	// 6.1 maps.Equal - 比较两个 map 是否相等
	m1 := map[string]int{"a": 1, "b": 2, "c": 3}
	m2 := map[string]int{"a": 1, "b": 2, "c": 3}
	m3 := map[string]int{"a": 1, "b": 2, "c": 4}

	fmt.Println("maps.Equal 比较:")
	fmt.Printf("  m1 == m2: %v\n", maps.Equal(m1, m2)) // true
	fmt.Printf("  m1 == m3: %v\n", maps.Equal(m1, m3)) // false

	// 6.2 maps.Clone - 深拷贝 map
	original := map[string]int{"x": 10, "y": 20}
	cloned := maps.Clone(original)
	cloned["x"] = 999

	fmt.Println("\nmaps.Clone 深拷贝:")
	fmt.Println("  original:", original) // {x:10, y:20}
	fmt.Println("  cloned:", cloned)     // {x:999, y:20}

	// 6.3 maps.Copy - 将源 map 复制到目标 map
	dest := map[string]int{"a": 1}
	src := map[string]int{"b": 2, "c": 3}
	maps.Copy(dest, src)

	fmt.Println("\nmaps.Copy 复制:")
	fmt.Println("  dest after copy:", dest) // {a:1, b:2, c:3}

	// 6.4 maps.DeleteFunc - 按条件删除元素
	numbers := map[string]int{"one": 1, "two": 2, "three": 3, "four": 4}
	// 删除所有偶数值
	maps.DeleteFunc(numbers, func(k string, v int) bool {
		return v%2 == 0
	})

	fmt.Println("\nmaps.DeleteFunc 条件删除（删除偶数）:")
	fmt.Println("  result:", numbers) // {one:1, three:3}

	fmt.Println()
}

// demonstrateMapPitfalls 演示 Map 常见陷阱
func demonstrateMapPitfalls() {
	fmt.Println("--- 7. Map 常见陷阱 ---")

	// 7.1 陷阱：nil map 写入会 panic
	fmt.Println("陷阱1: nil map 写入会 panic")
	var nilMap map[string]int
	fmt.Printf("  nil map 读取: %d（返回零值）\n", nilMap["key"])
	// nilMap["key"] = 1  // ❌ panic: assignment to entry in nil map
	fmt.Println("  ⚠️ 写入 nil map 会 panic！")
	fmt.Println("  ✅ 解决方案：先用 make 初始化")

	// 正确做法
	nilMap = make(map[string]int)
	nilMap["key"] = 1 // ✅ 正常

	// 7.2 陷阱：map 不能用 == 比较
	fmt.Println("\n陷阱2: map 不能用 == 比较")
	// m1 := map[string]int{"a": 1}
	// m2 := map[string]int{"a": 1}
	// m1 == m2  // ❌ 编译错误：map 只能与 nil 比较
	fmt.Println("  ❌ m1 == m2 会编译错误")
	fmt.Println("  ✅ 使用 maps.Equal(m1, m2)")

	// 7.3 陷阱：遍历顺序不固定
	fmt.Println("\n陷阱3: 遍历顺序不固定")
	fmt.Println("  Go 故意使 map 遍历顺序随机化")
	fmt.Println("  ✅ 如需有序，先收集键并排序")

	// 7.4 陷阱：并发读写会 panic
	fmt.Println("\n陷阱4: 并发读写不安全")
	fmt.Println("  多个 goroutine 同时读写同一个 map 会 panic")
	fmt.Println("  ✅ 解决方案：")
	fmt.Println("     - 使用 sync.RWMutex 保护")
	fmt.Println("     - 使用 sync.Map（适合读多写少场景）")

	// 7.5 陷阱：不能对 map 元素取地址
	fmt.Println("\n陷阱5: 不能对 map 元素取地址")
	m := map[string]int{"a": 1}
	// ptr := &m["a"]  // ❌ 编译错误
	_ = m
	fmt.Println("  ❌ &m[\"a\"] 会编译错误")
	fmt.Println("  原因：map 内部可能重新分配内存，地址会失效")

	fmt.Println()
}
