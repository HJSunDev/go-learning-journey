package main

import "fmt"

// =============================================================================
// 1. 定义结构体
// =============================================================================

// Player 表示游戏中的一个角色。
// 结构体将多个相关的数据字段组合成一个自定义类型。
// 这里 Player 有三个字段：名字、生命值、攻击力。
type Player struct {
	Name   string
	HP     int
	Attack int
}

// =============================================================================
// 2. 方法：值接收器 vs 指针接收器
// =============================================================================

// Status 返回玩家的状态描述。
// 使用值接收器 (p Player)：方法内部操作的是 Player 的副本。
// 适用于：只读操作，不需要修改结构体的场景。
func (p Player) Status() string {
	return fmt.Sprintf("%s - HP: %d, Attack: %d", p.Name, p.HP, p.Attack)
}

// TakeDamage 使玩家受到伤害。
// 使用指针接收器 (p *Player)：方法内部操作的是原始 Player。
// 适用于：需要修改结构体字段的场景。
func (p *Player) TakeDamage(damage int) {
	p.HP -= damage
	if p.HP < 0 {
		p.HP = 0
	}
	fmt.Printf("  %s 受到 %d 点伤害，剩余 HP: %d\n", p.Name, damage, p.HP)
}

// Heal 使玩家恢复生命值。
func (p *Player) Heal(amount int) {
	p.HP += amount
	fmt.Printf("  %s 恢复 %d 点生命，当前 HP: %d\n", p.Name, amount, p.HP)
}

// =============================================================================
// 2.5 方法声明 vs 函数类型：语法对比
// =============================================================================
//
// 方法声明（给结构体添加行为）:
//   func (p Player) Status() string { ... }
//        └───┬───┘ └──┬──┘ └┬┘└──┬─┘
//          接收器   方法名  参数  返回类型
//
// 函数类型（描述函数的"形状"）:
//   var fn func(Player) string
//          └──────┬──────────┘
//          类型，不是函数声明
//
// 判断技巧：看 func (...) 后面是否紧跟一个名字（方法名）

// MakePlayerFormatter 返回一个函数类型。
// 注意：这不是方法，而是一个普通函数，返回值类型是 func(Player) string。
func MakePlayerFormatter(prefix string) func(Player) string {
	// 返回一个闭包（匿名函数）
	return func(p Player) string {
		return prefix + p.Name
	}
}

// =============================================================================
// 3. 构造函数模式：使用工厂函数创建结构体
// =============================================================================

// NewPlayer 是 Player 的"构造函数"。
// Go 没有内置构造函数语法，惯例是用 New+类型名 的函数来模拟。
// 返回指针是惯例：避免复制，且方便后续调用指针接收器方法。
func NewPlayer(name string) *Player {
	return &Player{
		Name:   name,
		HP:     100, // 默认生命值
		Attack: 10,  // 默认攻击力
	}
}

// NewPlayerWithStats 是带自定义属性的构造函数。
// 展示如何提供多种构造函数以满足不同创建需求。
func NewPlayerWithStats(name string, hp, attack int) *Player {
	return &Player{
		Name:   name,
		HP:     hp,
		Attack: attack,
	}
}

// =============================================================================
// 4. 结构体组合（嵌入）：Go 的"继承"替代方案
// =============================================================================

// Skill 表示一个技能
type Skill struct {
	Name   string
	Damage int
}

// Hero 是一个更强大的角色类型。
// 通过嵌入 Player，Hero 自动获得 Player 的所有字段和方法。
// 这不是继承，而是组合——Hero "包含" 一个 Player。
type Hero struct {
	Player        // 嵌入 Player（匿名字段）
	Title  string // Hero 自己的字段
	Skill  Skill  // 组合另一个结构体（命名字段）
}

// UseSkill 是 Hero 特有的方法
func (h *Hero) UseSkill() {
	fmt.Printf("  %s 使用技能【%s】，造成 %d 点伤害！\n",
		h.Name, h.Skill.Name, h.Skill.Damage)
}

// NewHero 创建一个英雄
func NewHero(name, title string, skill Skill) *Hero {
	return &Hero{
		Player: Player{
			Name:   name,
			HP:     150,
			Attack: 20,
		},
		Title: title,
		Skill: skill,
	}
}

func main() {
	fmt.Println("===== 1. 结构体的创建方式 =====")

	// 方式1：使用字段名初始化（推荐）
	p1 := Player{Name: "战士", HP: 100, Attack: 15}
	fmt.Printf("字段名初始化: %+v\n", p1)

	// 方式2：零值初始化（所有字段为零值）
	var p2 Player
	fmt.Printf("零值初始化: %+v\n", p2)

	// 方式3：部分字段初始化（未指定的为零值）
	p3 := Player{Name: "法师"}
	fmt.Printf("部分初始化: %+v\n", p3)

	// 方式4：使用 new()，返回指针
	p4 := new(Player)
	p4.Name = "刺客"
	fmt.Printf("new() 创建: %+v\n", *p4)

	fmt.Println("\n===== 2. 字段访问与修改 =====")
	warrior := Player{Name: "剑客", HP: 80, Attack: 25}
	fmt.Printf("初始状态: %s\n", warrior.Status())

	// 直接修改字段
	warrior.HP = 90
	fmt.Printf("修改后: %s\n", warrior.Status())

	fmt.Println("\n===== 3. 方法调用：值接收器 vs 指针接收器 =====")
	player := Player{Name: "勇者", HP: 100, Attack: 20}
	fmt.Printf("初始状态: %s\n", player.Status())

	// 指针接收器方法会修改原始结构体
	player.TakeDamage(30)
	fmt.Printf("受伤后: %s\n", player.Status())

	player.Heal(20)
	fmt.Printf("治疗后: %s\n", player.Status())

	fmt.Println("\n===== 4. 构造函数模式 =====")
	// 使用构造函数创建（推荐方式）
	hero1 := NewPlayer("小白")
	fmt.Printf("默认构造: %s\n", hero1.Status())

	hero2 := NewPlayerWithStats("大侠", 200, 50)
	fmt.Printf("自定义构造: %s\n", hero2.Status())

	fmt.Println("\n===== 5. 匿名结构体：一次性使用 =====")
	// 当结构体只用一次时，无需单独定义类型
	config := struct {
		Host string
		Port int
	}{
		Host: "localhost",
		Port: 8080,
	}
	fmt.Printf("匿名结构体: Host=%s, Port=%d\n", config.Host, config.Port)

	// 常见场景：JSON 解析、测试数据、临时分组
	response := struct {
		Code    int
		Message string
		Data    []string
	}{
		Code:    200,
		Message: "success",
		Data:    []string{"item1", "item2"},
	}
	fmt.Printf("API 响应: Code=%d, Message=%s, Data=%v\n",
		response.Code, response.Message, response.Data)

	fmt.Println("\n===== 6. 结构体组合（嵌入） =====")
	fireball := Skill{Name: "火球术", Damage: 50}
	hero := NewHero("亚瑟", "圣骑士", fireball)

	// 嵌入的字段和方法被"提升"，可以直接访问
	fmt.Printf("英雄名字: %s（%s）\n", hero.Name, hero.Title)
	fmt.Printf("英雄状态: %s\n", hero.Status()) // 调用的是 Player.Status()

	// 也可以显式访问嵌入的结构体
	fmt.Printf("通过 Player 访问: %s\n", hero.Player.Status())

	// 调用 Hero 特有的方法
	hero.UseSkill()

	// 嵌入的方法同样可以修改内部状态
	hero.TakeDamage(40)
	fmt.Printf("受伤后: %s\n", hero.Status())

	fmt.Println("\n===== 7. 选择接收器类型的原则 =====")
	fmt.Println(`
┌──────────────────────────────────────────────────────────────┐
│  何时使用指针接收器 (*T)：                                     │
│    1. 方法需要修改接收器的字段                                 │
│    2. 结构体较大，避免复制开销                                 │
│    3. 保持一致性：如果某个方法需要指针，其他方法也用指针        │
├──────────────────────────────────────────────────────────────┤
│  何时使用值接收器 (T)：                                        │
│    1. 方法只读取数据，不修改                                   │
│    2. 结构体很小（如只有几个基本类型字段）                      │
│    3. 需要接收器不可变（方法不能意外修改原始值）               │
└──────────────────────────────────────────────────────────────┘`)
}
