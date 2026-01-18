package main

import "fmt"

// =============================================================================
// 本章核心问题：如何优雅地表示一组固定的、有限的选项？
// =============================================================================

// =============================================================================
// 第一步：痛点 —— 没有枚举时的问题
// =============================================================================

// 假设我们用整数表示订单状态
// 0=待支付, 1=已支付, 2=已发货, 3=已完成, 4=已取消

// ProcessOrderBad 处理订单状态（有问题的写法）
func ProcessOrderBad(status int) {
	// 问题1：魔法数字，代码可读性差
	// 看到 status == 1，不知道 1 代表什么
	if status == 1 {
		fmt.Println("订单已支付，准备发货")
	}

	// 问题2：没有类型约束，任何 int 都能传进来
	// 比如传入 100，代码不会报错，但逻辑上是错误的
	if status == 100 {
		fmt.Println("这是个无效状态，但编译器不会报错")
	}
}

// =============================================================================
// 第二步：用 const 定义常量 —— 解决魔法数字问题
// =============================================================================

// 把状态值定义成常量，用有意义的名字
const (
	StatusPending   = 0 // 待支付
	StatusPaid      = 1 // 已支付
	StatusShipped   = 2 // 已发货
	StatusCompleted = 3 // 已完成
	StatusCancelled = 4 // 已取消
)

// ProcessOrderBetter 使用常量的版本（好一点，但还不够）
func ProcessOrderBetter(status int) {
	// 现在代码可读性好多了
	if status == StatusPaid {
		fmt.Println("订单已支付，准备发货")
	}

	// 但问题2还没解决：status 的类型是 int，任何 int 都能传进来
	// 比如调用 ProcessOrderBetter(100) 不会报错
}

// =============================================================================
// 第三步：用 iota 简化常量定义
// =============================================================================

// iota 是 Go 的常量生成器，在 const 块中自动递增
// 第一个常量 iota = 0，第二个 iota = 1，依此类推

const (
	StatusPendingV2   = iota // iota = 0
	StatusPaidV2             // iota = 1，不用写 = iota，自动继承上一行的表达式
	StatusShippedV2          // iota = 2
	StatusCompletedV2        // iota = 3
	StatusCancelledV2        // iota = 4
)

// =============================================================================
// 第四步：自定义类型 —— 增加类型安全
// =============================================================================

// 定义一个新类型 OrderStatus，底层是 int
// 语法：type 新类型名 底层类型
type OrderStatus int

// 用 OrderStatus 类型定义常量
const (
	// iota 在每个 const 块开始时重置为 0
	Pending   OrderStatus = iota // Pending 的类型是 OrderStatus，值为 0
	Paid                         // 自动继承 OrderStatus 类型和 iota 表达式，值为 1
	Shipped                      // 值为 2
	Completed                    // 值为 3
	Cancelled                    // 值为 4
)

// ProcessOrder 使用 OrderStatus 类型作为参数
// 现在只有 OrderStatus 类型的值才能传进来
func ProcessOrder(status OrderStatus) {
	switch status {
	case Pending:
		fmt.Println("订单待支付")
	case Paid:
		fmt.Println("订单已支付，准备发货")
	case Shipped:
		fmt.Println("订单已发货")
	case Completed:
		fmt.Println("订单已完成")
	case Cancelled:
		fmt.Println("订单已取消")
	default:
		fmt.Println("未知状态")
	}
}

// =============================================================================
// 第五步：给枚举类型添加方法 —— 实现 String() 方法
// =============================================================================

// String 方法让 OrderStatus 可以输出可读的字符串
// 这是 Go 的 fmt.Stringer 接口，fmt 包会自动调用它
func (s OrderStatus) String() string {
	// 定义一个状态名称的映射表
	names := []string{"待支付", "已支付", "已发货", "已完成", "已取消"}

	// 边界检查：确保状态值在有效范围内
	if s < 0 || int(s) >= len(names) {
		return "未知状态"
	}

	return names[s]
}

// IsValid 方法检查状态值是否有效
func (s OrderStatus) IsValid() bool {
	return s >= Pending && s <= Cancelled
}

// CanCancel 方法：判断当前状态下订单是否可以取消
// 业务规则：只有待支付和已支付状态可以取消
func (s OrderStatus) CanCancel() bool {
	return s == Pending || s == Paid
}

// CanShip 方法：判断当前状态下订单是否可以发货
// 业务规则：只有已支付状态可以发货
func (s OrderStatus) CanShip() bool {
	return s == Paid
}

// =============================================================================
// 第六步：iota 进阶用法
// =============================================================================

// 示例1：从 1 开始编号
type Priority int

const (
	Low    Priority = iota + 1 // 0 + 1 = 1
	Medium                     // 1 + 1 = 2
	High                       // 2 + 1 = 3
)

// 示例2：跳过某个值
type FilePermission int

const (
	Read  FilePermission = 1 << iota // 1 << 0 = 1  (二进制: 001)
	Write                            // 1 << 1 = 2  (二进制: 010)
	Exec                             // 1 << 2 = 4  (二进制: 100)
)

// 示例3：跳过某个常量
type Weekday int

const (
	Monday    Weekday = iota + 1 // 1
	Tuesday                      // 2
	Wednesday                    // 3
	Thursday                     // 4
	Friday                       // 5
	_                            // 6，用下划线跳过（表示不需要这个值）
	Sunday                       // 7
)

// =============================================================================
// 主函数：演示所有概念
// =============================================================================

func main() {
	fmt.Println("===== 第一步：没有枚举时的问题 =====")
	fmt.Println()

	fmt.Println("使用魔法数字的问题：")
	fmt.Println("  status == 1 是什么意思？看代码的人不知道")
	fmt.Println("  ProcessOrderBad(100) 不会报错，但 100 不是有效状态")
	fmt.Println()

	fmt.Println("===== 第二步：用 const 定义常量 =====")
	fmt.Println()

	fmt.Printf("StatusPending = %d\n", StatusPending)
	fmt.Printf("StatusPaid = %d\n", StatusPaid)
	fmt.Printf("StatusShipped = %d\n", StatusShipped)
	fmt.Println()
	fmt.Println("现在代码可读性好多了，但类型还是 int，没有类型安全")
	fmt.Println()

	fmt.Println("===== 第三步：用 iota 简化常量定义 =====")
	fmt.Println()

	fmt.Println("iota 是常量生成器，在 const 块中自动递增：")
	fmt.Printf("StatusPendingV2 = %d\n", StatusPendingV2)
	fmt.Printf("StatusPaidV2 = %d\n", StatusPaidV2)
	fmt.Printf("StatusShippedV2 = %d\n", StatusShippedV2)
	fmt.Printf("StatusCompletedV2 = %d\n", StatusCompletedV2)
	fmt.Printf("StatusCancelledV2 = %d\n", StatusCancelledV2)
	fmt.Println()

	fmt.Println("===== 第四步：自定义类型 —— 类型安全 =====")
	fmt.Println()

	// 使用自定义类型 OrderStatus
	var myStatus OrderStatus = Paid
	fmt.Printf("myStatus = %d\n", myStatus)

	// 类型安全：只有 OrderStatus 类型才能传给 ProcessOrder
	ProcessOrder(Pending)
	ProcessOrder(Paid)
	ProcessOrder(Shipped)
	fmt.Println()

	// 下面这行代码如果取消注释，会编译报错：
	// ProcessOrder(1)  // 错误：cannot use 1 (untyped int constant) as OrderStatus
	fmt.Println("类型安全：ProcessOrder(1) 会编译报错，必须传 OrderStatus 类型")
	fmt.Println()

	fmt.Println("===== 第五步：给枚举类型添加方法 =====")
	fmt.Println()

	// String() 方法让枚举值可以输出可读的名称
	fmt.Println("实现 String() 方法后，fmt 会自动调用它：")
	for status := Pending; status <= Cancelled; status++ {
		// %v 和 %s 会调用 String() 方法
		fmt.Printf("  状态 %d 的名称是：%s\n", int(status), status)
	}
	fmt.Println()

	// IsValid() 方法检查状态是否有效
	fmt.Println("IsValid() 方法检查状态是否有效：")
	fmt.Printf("  Paid.IsValid() = %v\n", Paid.IsValid())
	fmt.Printf("  OrderStatus(100).IsValid() = %v\n", OrderStatus(100).IsValid())
	fmt.Println()

	// 业务方法演示
	fmt.Println("业务方法演示：")
	orders := []OrderStatus{Pending, Paid, Shipped, Completed, Cancelled}
	for _, status := range orders {
		fmt.Printf("  订单状态[%s]: 可取消=%v, 可发货=%v\n",
			status, status.CanCancel(), status.CanShip())
	}
	fmt.Println()

	fmt.Println("===== 第六步：iota 进阶用法 =====")
	fmt.Println()

	fmt.Println("从 1 开始编号（iota + 1）：")
	fmt.Printf("  Low = %d, Medium = %d, High = %d\n", Low, Medium, High)
	fmt.Println()

	fmt.Println("位运算生成权限标志（1 << iota）：")
	fmt.Printf("  Read = %d (二进制: %03b)\n", Read, Read)
	fmt.Printf("  Write = %d (二进制: %03b)\n", Write, Write)
	fmt.Printf("  Exec = %d (二进制: %03b)\n", Exec, Exec)
	// 位运算组合权限
	readWrite := Read | Write
	fmt.Printf("  Read | Write = %d (二进制: %03b)\n", readWrite, readWrite)
	fmt.Println()

	fmt.Println("===== 总结 =====")
	fmt.Println(`
1. Go 没有 enum 关键字，使用 const + iota 实现枚举
2. iota 在 const 块中自动递增，第一个值为 0
3. 使用自定义类型（如 type OrderStatus int）增加类型安全
4. 给枚举类型添加方法（如 String()）提高可读性和封装业务逻辑
5. iota 支持表达式（如 iota+1, 1<<iota）实现灵活编号`)
}
