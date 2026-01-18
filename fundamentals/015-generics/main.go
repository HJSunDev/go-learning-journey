package main

import (
	"fmt"
)

// =============================================================================
// 本章核心问题：如何避免为每种类型写重复的代码？
// =============================================================================

// =============================================================================
// 第一步：痛点 —— 没有泛型时的代码重复
// =============================================================================

// 场景：订单系统需要找出切片中的最小值
// 比如：找出所有商品中最便宜的价格，或者找出最早的订单ID

// MinInt 找出 int 切片中的最小值
func MinInt(values []int) int {
	if len(values) == 0 {
		return 0
	}
	min := values[0]
	for _, v := range values[1:] {
		if v < min {
			min = v
		}
	}
	return min
}

// MinFloat64 找出 float64 切片中的最小值
// 逻辑和 MinInt 完全一样，只是类型不同
func MinFloat64(values []float64) float64 {
	if len(values) == 0 {
		return 0
	}
	min := values[0]
	for _, v := range values[1:] {
		if v < min {
			min = v
		}
	}
	return min
}

// MinString 找出 string 切片中的最小值（按字典序）
// 逻辑还是一样，又要重写一遍
func MinString(values []string) string {
	if len(values) == 0 {
		return ""
	}
	min := values[0]
	for _, v := range values[1:] {
		if v < min {
			min = v
		}
	}
	return min
}

// 问题总结：
// 1. 代码重复：三个函数的逻辑完全相同，只有类型不同
// 2. 维护困难：如果要修改逻辑（比如处理空切片的方式），要改三处
// 3. 扩展麻烦：新增类型（如 int64）就要再写一个函数

// =============================================================================
// 第二步：泛型函数 —— 用类型参数解决重复
// =============================================================================

// Min 是一个泛型函数，可以处理任意可比较的类型
//
// 语法拆解：
//   func Min[T cmp.Ordered](values []T) T
//   │    │  │  └── 类型约束：T 必须支持 < > 比较操作
//   │    │  └── 类型参数：T 是一个占位符，代表某种类型
//   │    └── 类型参数列表，用方括号 [] 包裹
//   └── func 关键字
//
// 由于 cmp.Ordered 需要导入 cmp 包，这里先用 comparable + 手动限制演示
// 后面会介绍 cmp.Ordered

// Ordered 是一个类型约束，定义了哪些类型可以使用 < > 比较
// 语法：interface 后面的 ~ 表示"底层类型是"
type Ordered interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 |
		~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 |
		~float32 | ~float64 |
		~string
}

// Min 找出切片中的最小值（泛型版本）
//
// 语法拆解：
//
//	[T Ordered] —— 类型参数列表
//	  T      —— 类型参数名，习惯用单个大写字母
//	  Ordered —— 类型约束，限制 T 必须是 Ordered 接口允许的类型
//
//	(values []T) —— 参数列表，values 是 T 类型的切片
//
//	T —— 返回值类型，也是 T
func Min[T Ordered](values []T) T {
	// 函数体内，T 就像一个具体类型一样使用
	if len(values) == 0 {
		// 返回 T 的零值
		var zero T
		return zero
	}

	min := values[0]
	for _, v := range values[1:] {
		// 因为 T 满足 Ordered 约束，所以可以用 < 比较
		if v < min {
			min = v
		}
	}
	return min
}

// =============================================================================
// 第三步：类型推断 —— 让调用更简洁
// =============================================================================

// 调用泛型函数时，Go 可以自动推断类型参数
// 完整写法：Min[int]([]int{3, 1, 2})
// 简化写法：Min([]int{3, 1, 2})  // Go 自动推断 T 是 int

// =============================================================================
// 第四步：类型约束详解
// =============================================================================

// any 约束：接受任意类型
// comparable 约束：接受可以用 == 和 != 比较的类型

// Contains 检查切片中是否包含某个元素
//
// comparable 是 Go 内置的约束，表示类型可以用 == 比较
// map 的键必须是 comparable，所以这个约束很常用
func Contains[T comparable](slice []T, target T) bool {
	for _, v := range slice {
		// 因为 T 满足 comparable，所以可以用 == 比较
		if v == target {
			return true
		}
	}
	return false
}

// Keys 返回 map 的所有键
//
// 这里有两个类型参数：K（键类型）和 V（值类型）
// K 必须是 comparable（map 的键必须可比较）
// V 是 any（值可以是任意类型）
func Keys[K comparable, V any](m map[K]V) []K {
	// make 创建一个切片，长度为 0，容量为 map 的长度
	keys := make([]K, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}

// =============================================================================
// 第五步：泛型类型 —— 创建可复用的数据结构
// =============================================================================

// 场景扩展：订单系统需要一个通用的"结果"类型
// 有时操作成功返回数据，有时失败返回错误信息

// Result 表示一个操作的结果，可能成功也可能失败
//
// 语法拆解：
//
//	type Result[T any] struct
//	│    │     │  └── 类型约束
//	│    │     └── 类型参数
//	│    └── 类型名
//	└── type 关键字
type Result[T any] struct {
	Value T      // 成功时的值
	Error string // 失败时的错误信息
	OK    bool   // 是否成功
}

// NewSuccess 创建一个成功的结果
func NewSuccess[T any](value T) Result[T] {
	return Result[T]{
		Value: value,
		OK:    true,
	}
}

// NewFailure 创建一个失败的结果
func NewFailure[T any](err string) Result[T] {
	return Result[T]{
		Error: err,
		OK:    false,
	}
}

// Unwrap 获取结果的值，如果失败则 panic
func (r Result[T]) Unwrap() T {
	if !r.OK {
		panic("called Unwrap on a failed Result: " + r.Error)
	}
	return r.Value
}

// UnwrapOr 获取结果的值，如果失败则返回默认值
func (r Result[T]) UnwrapOr(defaultValue T) T {
	if !r.OK {
		return defaultValue
	}
	return r.Value
}

// =============================================================================
// 第六步：泛型切片工具 —— 实用的集合操作
// =============================================================================

// Filter 过滤切片，返回满足条件的元素
//
// 与前面 Min 的区别：
//   - Min 的约束是 Ordered，只能处理可比较大小的基础类型（int、float64、string）
//   - Filter 的约束是 any，可以处理任意类型，包括自定义结构体（Order、User、Product 等）
//   - 没有泛型时，需要为每种类型写 FilterOrders、FilterUsers、FilterProducts...
//   - 有泛型后，一个 Filter 函数就能处理所有类型的切片
//
// 参数：
//   - slice: 要过滤的切片
//   - predicate: 判断函数，返回 true 的元素会被保留
//
// 语法拆解：
//
//	func(T) bool —— 这是一个函数类型
//	  T     —— 参数类型，和切片元素类型一致
//	  bool  —— 返回值类型
func Filter[T any](slice []T, predicate func(T) bool) []T {
	result := make([]T, 0)
	for _, v := range slice {
		if predicate(v) {
			result = append(result, v)
		}
	}
	return result
}

// Map 将切片中的每个元素转换为另一种类型
//
// 两个类型参数：T（输入类型）和 R（输出类型）
func Map[T any, R any](slice []T, transform func(T) R) []R {
	result := make([]R, len(slice))
	for i, v := range slice {
		result[i] = transform(v)
	}
	return result
}

// Reduce 将切片归约为单个值
//
// 参数：
//   - slice: 要归约的切片
//   - initial: 初始值
//   - reducer: 归约函数，接收累加器和当前元素，返回新的累加器
func Reduce[T any, R any](slice []T, initial R, reducer func(R, T) R) R {
	result := initial
	for _, v := range slice {
		result = reducer(result, v)
	}
	return result
}

// =============================================================================
// 第七步：实际应用 —— 订单处理场景
// =============================================================================

// Order 订单结构
type Order struct {
	ID       string
	Customer string
	Amount   float64
	Quantity int
}

// =============================================================================
// 主函数：演示所有概念
// =============================================================================

func main() {
	fmt.Println("===== 第一步：没有泛型时的问题 =====")
	fmt.Println()

	// 使用非泛型函数
	intPrices := []int{99, 45, 150, 30, 88}
	floatPrices := []float64{99.9, 45.5, 150.0, 30.3, 88.8}
	names := []string{"Bob", "Alice", "Charlie"}

	fmt.Printf("整数价格最小值: %d （使用 MinInt）\n", MinInt(intPrices))
	fmt.Printf("浮点价格最小值: %.1f （使用 MinFloat64）\n", MinFloat64(floatPrices))
	fmt.Printf("名字最小值: %s （使用 MinString）\n", MinString(names))
	fmt.Println()
	fmt.Println("问题：三个函数逻辑完全相同，只是类型不同，代码重复！")
	fmt.Println()

	fmt.Println("===== 第二步：泛型函数解决重复 =====")
	fmt.Println()

	// 使用泛型函数 Min
	// Go 自动推断类型，不需要显式写 Min[int]
	fmt.Printf("整数价格最小值: %d （使用 Min）\n", Min(intPrices))
	fmt.Printf("浮点价格最小值: %.1f （使用 Min）\n", Min(floatPrices))
	fmt.Printf("名字最小值: %s （使用 Min）\n", Min(names))
	fmt.Println()
	fmt.Println("一个函数处理所有类型！")
	fmt.Println()

	fmt.Println("===== 第三步：类型推断 =====")
	fmt.Println()

	// 显式指定类型参数（完整写法）
	result1 := Min[int](intPrices)
	// 类型推断（简化写法）
	result2 := Min(intPrices)
	fmt.Printf("显式写法 Min[int](...) = %d\n", result1)
	fmt.Printf("推断写法 Min(...) = %d\n", result2)
	fmt.Println()
	fmt.Println("两种写法结果相同，推荐使用类型推断（更简洁）")
	fmt.Println()

	fmt.Println("===== 第四步：comparable 约束 =====")
	fmt.Println()

	ids := []string{"order-003", "order-001", "order-002"}
	fmt.Printf("订单ID列表: %v\n", ids)
	fmt.Printf("是否包含 order-001: %v\n", Contains(ids, "order-001"))
	fmt.Printf("是否包含 order-999: %v\n", Contains(ids, "order-999"))
	fmt.Println()

	// Keys 函数示例
	orderMap := map[string]float64{
		"order-001": 99.9,
		"order-002": 150.0,
		"order-003": 45.5,
	}
	fmt.Printf("订单Map: %v\n", orderMap)
	fmt.Printf("所有订单ID: %v\n", Keys(orderMap))
	fmt.Println()

	fmt.Println("===== 第五步：泛型类型 Result =====")
	fmt.Println()

	// 模拟查找订单
	findOrder := func(id string) Result[Order] {
		orders := map[string]Order{
			"order-001": {ID: "order-001", Customer: "Alice", Amount: 99.9, Quantity: 2},
		}
		if order, exists := orders[id]; exists {
			return NewSuccess(order)
		}
		return NewFailure[Order]("订单不存在: " + id)
	}

	// 查找存在的订单
	result := findOrder("order-001")
	if result.OK {
		fmt.Printf("找到订单: %+v\n", result.Value)
	}

	// 查找不存在的订单
	result = findOrder("order-999")
	if !result.OK {
		fmt.Printf("查找失败: %s\n", result.Error)
	}

	// 使用 UnwrapOr 提供默认值
	defaultOrder := Order{ID: "default", Customer: "Guest", Amount: 0}
	order := findOrder("order-999").UnwrapOr(defaultOrder)
	fmt.Printf("使用默认值: %+v\n", order)
	fmt.Println()

	fmt.Println("===== 第六步：泛型集合操作 =====")
	fmt.Println()

	orders := []Order{
		{ID: "001", Customer: "Alice", Amount: 99.9, Quantity: 2},
		{ID: "002", Customer: "Bob", Amount: 150.0, Quantity: 1},
		{ID: "003", Customer: "Alice", Amount: 45.5, Quantity: 3},
		{ID: "004", Customer: "Charlie", Amount: 200.0, Quantity: 1},
	}
	fmt.Println("所有订单：")
	for _, o := range orders {
		fmt.Printf("  %+v\n", o)
	}
	fmt.Println()

	// Filter: 过滤 Alice 的订单
	aliceOrders := Filter(orders, func(o Order) bool {
		return o.Customer == "Alice"
	})
	fmt.Println("Alice 的订单（Filter）：")
	for _, o := range aliceOrders {
		fmt.Printf("  %+v\n", o)
	}
	fmt.Println()

	// Map: 提取所有订单金额
	amounts := Map(orders, func(o Order) float64 {
		return o.Amount
	})
	fmt.Printf("所有订单金额（Map）: %v\n", amounts)
	fmt.Println()

	// Reduce: 计算订单总金额
	total := Reduce(orders, 0.0, func(sum float64, o Order) float64 {
		return sum + o.Amount
	})
	fmt.Printf("订单总金额（Reduce）: %.2f\n", total)
	fmt.Println()

	// 链式操作：计算 Alice 的订单总金额
	aliceTotal := Reduce(
		Filter(orders, func(o Order) bool { return o.Customer == "Alice" }),
		0.0,
		func(sum float64, o Order) float64 { return sum + o.Amount },
	)
	fmt.Printf("Alice 的订单总金额: %.2f\n", aliceTotal)
	fmt.Println()

	fmt.Println("===== 总结 =====")
	fmt.Println(`
1. 泛型解决代码重复：一个函数/类型处理多种类型
2. 语法：func Name[T Constraint](params) ReturnType
3. 类型约束：限制泛型接受的类型
   - any: 任意类型
   - comparable: 可用 == 比较的类型
   - 自定义约束: interface { ~int | ~string | ... }
4. 类型推断：调用时通常不需要显式写类型参数
5. 泛型类型：type Name[T Constraint] struct { ... }`)
}
