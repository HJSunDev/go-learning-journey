package main

import "fmt"

// =============================================================================
// 本章核心问题：如何让一个函数能接受多种不同的类型？
// =============================================================================

// =============================================================================
// 第一步：定义三种支付方式
// =============================================================================

type WechatPay struct{ UserID string }
type Alipay struct{ Account string }
type BankCard struct{ CardNo string }

// 每种支付方式都有 Pay 方法
func (w WechatPay) Pay(amount float64) {
	fmt.Printf("微信支付: 用户 %s 支付 %.2f 元\n", w.UserID, amount)
}
func (a Alipay) Pay(amount float64) {
	fmt.Printf("支付宝支付: 账户 %s 支付 %.2f 元\n", a.Account, amount)
}
func (b BankCard) Pay(amount float64) {
	fmt.Printf("银行卡支付: 卡号 %s 支付 %.2f 元\n", b.CardNo, amount)
}

// =============================================================================
// 第二步：定义接口
// =============================================================================

// Payer 是一个接口，它规定：任何想成为 Payer 的类型，必须有 Pay(float64) 方法
type Payer interface {
	Pay(amount float64)
}

// 因为 WechatPay、Alipay、BankCard 都有 Pay(float64) 方法
// 所以它们都自动满足 Payer 接口（Go 不需要显式声明 implements）

// =============================================================================
// 第三步：用接口类型作为函数参数
// =============================================================================

// Checkout 的参数类型是 Payer 接口
// 这意味着：任何满足 Payer 接口的类型，都可以传进来
func Checkout(p Payer, amount float64) {
	// 不管 p 实际是 WechatPay、Alipay 还是 BankCard
	// 都可以调用 Pay 方法，因为接口保证了它一定有这个方法
	p.Pay(amount)
}

func main() {
	fmt.Println("===== 第一步：没有接口时的问题 =====")
	fmt.Println()

	// 创建三种支付方式
	wechat := WechatPay{UserID: "wx_001"}
	alipay := Alipay{Account: "alice@example.com"}
	bank := BankCard{CardNo: "6222****1234"}

	// 直接调用各自的 Pay 方法，没问题
	wechat.Pay(100)
	alipay.Pay(200)
	bank.Pay(300)

	fmt.Println()
	fmt.Println("问题来了：如果我想写一个函数，能处理任意支付方式怎么办？")
	fmt.Println("函数参数写什么类型？WechatPay？Alipay？BankCard？都不行，因为它们是不同的类型")
	fmt.Println()

	fmt.Println("===== 第二步：接口解决问题 =====")
	fmt.Println()

	// 关键理解：接口类型的变量，可以存储任何实现了该接口的值
	//
	// 语法：var 变量名 接口类型 = 具体类型的值
	var p Payer // 声明一个 Payer 接口类型的变量
	p = wechat  // 把 WechatPay 类型的值赋给它（合法，因为 WechatPay 有 Pay 方法）
	fmt.Println("p 现在存的是 WechatPay：")
	p.Pay(100) // 调用的是 WechatPay 的 Pay 方法

	p = alipay // 把 Alipay 类型的值赋给它（合法，因为 Alipay 有 Pay 方法）
	fmt.Println("p 现在存的是 Alipay：")
	p.Pay(200) // 调用的是 Alipay 的 Pay 方法

	p = bank // 把 BankCard 类型的值赋给它
	fmt.Println("p 现在存的是 BankCard：")
	p.Pay(300) // 调用的是 BankCard 的 Pay 方法

	fmt.Println()
	fmt.Println("===== 第三步：接口作为函数参数 =====")
	fmt.Println()

	fmt.Println("现在 Checkout 函数的参数是 Payer 接口类型")
	fmt.Println("所以任何满足 Payer 接口的类型都能传进去：")
	fmt.Println()

	Checkout(wechat, 100) // 传入 WechatPay，合法
	Checkout(alipay, 200) // 传入 Alipay，合法
	Checkout(bank, 300)   // 传入 BankCard，合法

	fmt.Println()
	fmt.Println("核心价值：一个 Checkout 函数处理所有支付方式！")
	fmt.Println("新增支付方式时，只要它有 Pay 方法，就能直接用，Checkout 不用改")

	fmt.Println()
	fmt.Println("===== 第四步：接口变量内部是什么？ =====")
	fmt.Println()

	// 接口变量内部存储两样东西：
	// 1. 具体类型（动态类型）：实际存的是什么类型
	// 2. 具体值（动态值）：实际存的值是什么

	var p2 Payer
	fmt.Printf("空接口变量 p2: 值=%v, 类型=%T, 是否为nil=%v\n", p2, p2, p2 == nil)

	p2 = wechat
	fmt.Printf("赋值后 p2: 值=%v, 类型=%T\n", p2, p2)

	p2 = alipay
	fmt.Printf("重新赋值 p2: 值=%v, 类型=%T\n", p2, p2)

	fmt.Println()
	fmt.Println("===== 第五步：类型断言 —— 从接口取出具体类型 =====")
	fmt.Println()

	// 有时候需要知道接口变量里存的到底是什么类型
	// 比如想访问 WechatPay 的 UserID 字段，但接口变量只能调用接口定义的方法

	var p3 Payer = wechat

	// 语法解释：value, ok := 接口变量.(具体类型)
	//
	// 这叫"类型断言"，作用是：尝试把接口变量里的值取出来，转成指定的具体类型
	// - 如果接口变量里存的确实是这个类型，ok 为 true，value 就是取出的值
	// - 如果不是这个类型，ok 为 false，value 是零值
	//
	// 注意：.(Type) 不是函数调用，是类型断言的特殊语法

	w, ok := p3.(WechatPay)
	// w 的类型是 WechatPay
	// ok 的类型是 bool

	fmt.Printf("尝试断言为 WechatPay: ok=%v, w=%v\n", ok, w)
	if ok {
		// 断言成功，可以访问 WechatPay 的字段了
		fmt.Printf("  成功！UserID = %s\n", w.UserID)
	}

	// 再试试断言为 Alipay
	a, ok := p3.(Alipay)
	fmt.Printf("尝试断言为 Alipay: ok=%v\n", ok)
	if ok {
		fmt.Printf("  成功！Account = %s\n", a.Account)
	} else {
		fmt.Println("  失败！p3 里存的不是 Alipay")
	}

	fmt.Println()
	fmt.Println("===== 常见写法：if 语句中使用类型断言 =====")
	fmt.Println()

	// Go 允许在 if 条件中同时进行赋值和判断
	// 语法：if value, ok := 接口变量.(类型); ok { ... }
	//
	// 分号前面是赋值语句，分号后面是条件判断
	// 这样 value 变量只在 if 块内有效

	p3 = alipay
	if ali, ok := p3.(Alipay); ok {
		fmt.Printf("p3 是 Alipay，账户: %s\n", ali.Account)
	}

	if wx, ok := p3.(WechatPay); ok {
		fmt.Printf("p3 是 WechatPay，用户: %s\n", wx.UserID)
	} else {
		fmt.Println("p3 不是 WechatPay")
	}

	fmt.Println()
	fmt.Println("===== type switch：判断多种类型 =====")
	fmt.Println()

	// 当需要判断多种类型时，用 type switch 更简洁
	// 语法：switch v := 接口变量.(type) { case 类型1: ... case 类型2: ... }
	//
	// 注意：.(type) 只能在 switch 语句中使用，这是特殊语法

	payers := []Payer{wechat, alipay, bank}
	for _, payer := range payers {
		switch v := payer.(type) {
		case WechatPay:
			fmt.Printf("微信支付，用户ID: %s\n", v.UserID)
		case Alipay:
			fmt.Printf("支付宝，账户: %s\n", v.Account)
		case BankCard:
			fmt.Printf("银行卡，卡号: %s\n", v.CardNo)
		}
	}

	fmt.Println()
	fmt.Println("===== 空接口 any：能存任何类型 =====")
	fmt.Println()

	// any 是 interface{} 的别名
	// 空接口没有任何方法要求，所以所有类型都满足它

	var x any // 声明一个空接口变量
	x = 42    // 可以存 int
	fmt.Printf("x = %v (类型: %T)\n", x, x)

	x = "hello" // 可以存 string
	fmt.Printf("x = %v (类型: %T)\n", x, x)

	x = wechat // 可以存 WechatPay
	fmt.Printf("x = %v (类型: %T)\n", x, x)

	// 从 any 取出具体值也需要类型断言
	if num, ok := x.(int); ok {
		fmt.Printf("x 是 int: %d\n", num)
	} else {
		fmt.Println("x 不是 int")
	}

	fmt.Println()
	fmt.Println("===== 总结 =====")
	fmt.Println(`
1. 接口是方法签名的集合，定义"必须有什么方法"
2. 类型只要有接口要求的方法，就自动满足该接口（隐式实现）
3. 接口类型的变量可以存储任何满足该接口的值
4. 调用接口变量的方法时，实际执行的是内部具体类型的方法（多态）
5. 类型断言 v, ok := i.(Type) 用于从接口取出具体类型
6. 空接口 any 可以存储任何类型的值
`)
}
