# 012. Go 结构体 (Structs)：构建自定义类型的基石

[返回索引](../README.md) | [查看代码](../../012-structs/main.go)

结构体是 Go 语言中将多个相关数据组合成一个整体的核心机制。它是构建复杂程序的基础。

**一句话定义**：结构体是一种自定义类型，将多个不同类型的字段组合在一起，形成一个有意义的数据单元。

## 1. 结构体是什么

假设你在开发一个游戏，需要表示一个玩家角色。玩家有名字、生命值、攻击力。你可以用三个独立变量：

```go
name := "勇者"
hp := 100
attack := 20
```

问题很明显：这三个变量是分散的，传递时需要逐个传递，容易出错。

**结构体解决这个问题**——将相关数据组合成一个整体：

```go
type Player struct {
    Name   string
    HP     int
    Attack int
}

// 现在只需要一个变量就能表示完整的玩家
player := Player{Name: "勇者", HP: 100, Attack: 20}
```

## 2. 定义与创建

### 定义结构体

```go
type Player struct {
    Name   string
    HP     int
    Attack int
}
```

- `type` 关键字声明一个新类型
- `Player` 是类型名（首字母大写表示可导出）
- `struct { ... }` 包含字段列表
- 每个字段有名称和类型

### 创建结构体实例

Go 提供多种创建方式：

```go
// 方式1：使用字段名初始化（推荐）
p1 := Player{Name: "战士", HP: 100, Attack: 15}

// 方式2：零值初始化
var p2 Player  // Name="", HP=0, Attack=0

// 方式3：部分字段初始化
p3 := Player{Name: "法师"}  // HP=0, Attack=0

// 方式4：使用 new()，返回指针
p4 := new(Player)  // *Player，所有字段为零值
p4.Name = "刺客"
```

**推荐使用方式1**：字段名明确，代码可读性高，字段顺序变化也不影响。

### 字段访问

使用点号 `.` 访问和修改字段：

```go
player := Player{Name: "勇者", HP: 100, Attack: 20}

fmt.Println(player.Name)   // 读取：勇者
player.HP = 80             // 修改
```

## 3. 方法：给结构体赋予行为

光有数据不够，我们需要结构体能"做事"。Go 通过**方法**实现这一点。

### 定义方法

方法就是带有**接收器**的函数：

```go
// 在 func 和函数名之间添加接收器 (p Player)
func (p Player) Status() string {
    return fmt.Sprintf("%s - HP: %d, Attack: %d", p.Name, p.HP, p.Attack)
}
```

- `(p Player)` 是接收器，表示这个方法属于 `Player` 类型
- `p` 是接收器变量名，在方法内部代表调用该方法的结构体实例

### 方法声明 vs 匿名函数：三种形态放一起对比

下面三种写法都以 `func(...)` 开头，但含义完全不同：

```go
// 1、方法声明：给 Player 添加行为（第一个括号是“接收器”）
func (p Player) Greet(prefix string) string { ... }
//   └────────┘ └───┬──┘
//    接收器        方法名（这里出现“自定义名字”，就一定是方法）

// 2、匿名函数：一个函数值（第一个括号是“参数”）
fn := func(p Player, prefix string) string { ... }
//        └───────────────┘
//               参数（匿名函数没有名字）

// 3、匿名函数：返回函数类型（长得最像方法！）
factory := func(p Player, prefix string) func() string { ... }
//             └───────────────┘ └───────────┘
//                    参数          返回类型（这里的 func 是“类型”，不是名字）
```

**核心差别（最简规则）**：

- **方法**：`func (接收器) 方法名(...) ...` —— 第一个括号后面紧跟的是你起的**名字**（如 `Greet`）。
- **匿名函数**：`func(参数...) 返回类型 { ... }` —— 第一个括号后面不会出现你起的名字，只会进入**返回类型**或直接进入 `{`。

### 调用方法

```go
player := Player{Name: "勇者", HP: 100, Attack: 20}
fmt.Println(player.Status())  // 勇者 - HP: 100, Attack: 20
```

### 值接收器 vs 指针接收器

这是结构体方法的核心概念。

**值接收器**：方法收到的是结构体的**副本**

```go
func (p Player) Status() string {
    return fmt.Sprintf("%s - HP: %d", p.Name, p.HP)
}
```

**指针接收器**：方法收到的是结构体的**地址**

```go
func (p *Player) TakeDamage(damage int) {
    p.HP -= damage  // 直接修改原始结构体
}
```

### 对比示例

```go
// 值接收器：不能修改原始结构体
func (p Player) HealWrong(amount int) {
    p.HP += amount  // 修改的是副本，原始数据不变
}

// 指针接收器：可以修改原始结构体
func (p *Player) Heal(amount int) {
    p.HP += amount  // 修改原始数据
}

func main() {
    player := Player{Name: "勇者", HP: 50}
  
    player.HealWrong(20)
    fmt.Println(player.HP)  // 50（没变！）
  
    player.Heal(20)
    fmt.Println(player.HP)  // 70（成功修改）
}
```

### 如何选择？

| 场景               | 接收器类型      | 原因                             |
| ------------------ | --------------- | -------------------------------- |
| 需要修改结构体字段 | 指针 `*T`     | 值接收器无法修改原始数据         |
| 结构体较大         | 指针 `*T`     | 避免复制开销                     |
| 只读取不修改       | 值 `T` 或指针 | 都可以，小结构体用值更清晰       |
| 保持一致性         | 指针 `*T`     | 如果某个方法用指针，其他也用指针 |

**实践建议**：大多数情况使用指针接收器。只有确定不需要修改且结构体很小时才用值接收器。

## 4. 构造函数模式

Go 没有内置的构造函数语法（不像 Java/Python 的 `__init__`）。惯例是用**工厂函数**模拟。

### 基本构造函数

```go
// 惯例命名：New + 类型名
func NewPlayer(name string) *Player {
    return &Player{
        Name:   name,
        HP:     100,  // 默认值
        Attack: 10,   // 默认值
    }
}

// 使用
player := NewPlayer("勇者")
```

### 为什么返回指针？

1. **避免复制**：返回值类型会复制整个结构体，指针只复制地址
2. **调用指针方法更方便**：返回的指针可以直接调用指针接收器方法
3. **表示"创建了一个对象"的语义**：符合面向对象的习惯

### 带参数的构造函数

```go
func NewPlayerWithStats(name string, hp, attack int) *Player {
    return &Player{
        Name:   name,
        HP:     hp,
        Attack: attack,
    }
}

// 使用
tank := NewPlayerWithStats("坦克", 300, 5)
```

### 带验证的构造函数

```go
func NewPlayer(name string) (*Player, error) {
    if name == "" {
        return nil, fmt.Errorf("player name cannot be empty")
    }
    return &Player{
        Name:   name,
        HP:     100,
        Attack: 10,
    }, nil
}
```

## 5. 匿名结构体

当结构体只使用一次时，无需单独定义类型。

### 语法

```go
config := struct {
    Host string
    Port int
}{
    Host: "localhost",
    Port: 8080,
}
```

### 典型场景

**场景1：临时配置**

```go
config := struct {
    Debug   bool
    Timeout int
}{Debug: true, Timeout: 30}
```

**场景2：API 响应解析**

```go
response := struct {
    Code    int      `json:"code"`
    Message string   `json:"message"`
    Data    []string `json:"data"`
}{}

json.Unmarshal(jsonData, &response)
```

**场景3：测试数据**

```go
tests := []struct {
    input    string
    expected int
}{
    {"hello", 5},
    {"world", 5},
    {"go", 2},
}

for _, tc := range tests {
    result := len(tc.input)
    if result != tc.expected {
        t.Errorf("len(%s) = %d, want %d", tc.input, result, tc.expected)
    }
}
```

## 6. 结构体组合（嵌入）

Go 没有继承，但通过**组合**实现代码复用。

### 基本组合

一个结构体包含另一个结构体作为字段：

```go
type Skill struct {
    Name   string
    Damage int
}

type Hero struct {
    Player Player  // 命名字段：通过 hero.Player 访问
    Title  string
    Skill  Skill
}

hero := Hero{
    Player: Player{Name: "亚瑟", HP: 150, Attack: 20},
    Title:  "圣骑士",
    Skill:  Skill{Name: "火球术", Damage: 50},
}

fmt.Println(hero.Player.Name)  // 亚瑟
fmt.Println(hero.Skill.Name)   // 火球术
```

### 嵌入（匿名字段）

将结构体作为**匿名字段**嵌入，字段和方法会被"提升"：

```go
type Hero struct {
    Player        // 匿名字段：嵌入 Player
    Title  string
    Skill  Skill
}

hero := Hero{
    Player: Player{Name: "亚瑟", HP: 150, Attack: 20},
    Title:  "圣骑士",
    Skill:  Skill{Name: "火球术", Damage: 50},
}

// 字段提升：可以直接访问
fmt.Println(hero.Name)    // 亚瑟（等同于 hero.Player.Name）
fmt.Println(hero.HP)      // 150

// 方法提升：可以直接调用
fmt.Println(hero.Status())  // 调用 Player.Status()

// 仍然可以显式访问
fmt.Println(hero.Player.Name)
```

### 嵌入 vs 继承

| 特性     | 嵌入（Go）       | 继承（Java/Python） |
| -------- | ---------------- | ------------------- |
| 关系     | "包含"           | "是一个"            |
| 类型兼容 | Hero 不是 Player | 子类是父类          |
| 访问方式 | 字段/方法提升    | 直接继承            |
| 灵活性   | 可嵌入多个       | 单继承限制          |

### 添加 Hero 自己的方法

```go
func (h *Hero) UseSkill() {
    fmt.Printf("%s 使用【%s】，造成 %d 点伤害！\n",
        h.Name, h.Skill.Name, h.Skill.Damage)
}
```

### 方法覆盖

如果 Hero 定义了与 Player 同名的方法，调用时优先使用 Hero 的：

```go
func (h *Hero) Status() string {
    return fmt.Sprintf("[%s] %s - HP: %d", h.Title, h.Name, h.HP)
}

// hero.Status() 调用 Hero.Status()
// hero.Player.Status() 调用 Player.Status()
```

## 7. 完整示例回顾

```go
// 定义基础结构体
type Player struct {
    Name   string
    HP     int
    Attack int
}

// 添加方法
func (p Player) Status() string {
    return fmt.Sprintf("%s - HP: %d", p.Name, p.HP)
}

func (p *Player) TakeDamage(damage int) {
    p.HP -= damage
}

// 构造函数
func NewPlayer(name string) *Player {
    return &Player{Name: name, HP: 100, Attack: 10}
}

// 通过组合扩展
type Hero struct {
    Player
    Title string
}

func NewHero(name, title string) *Hero {
    return &Hero{
        Player: Player{Name: name, HP: 150, Attack: 20},
        Title:  title,
    }
}
```

## 8. 总结

| 概念       | 要点                                        |
| ---------- | ------------------------------------------- |
| 结构体定义 | `type Name struct { fields }`             |
| 创建实例   | 推荐使用字段名：`Player{Name: "勇者"}`    |
| 方法       | 函数 + 接收器：`func (p Player) Method()` |
| 值接收器   | `(p T)`：操作副本，不能修改原始值         |
| 指针接收器 | `(p *T)`：操作原始值，可以修改            |
| 构造函数   | 工厂函数 `NewXxx()`，通常返回指针         |
| 匿名结构体 | 一次性使用，无需定义类型                    |
| 组合/嵌入  | 匿名字段实现字段和方法提升                  |
