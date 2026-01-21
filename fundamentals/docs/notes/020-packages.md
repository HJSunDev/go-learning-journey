# 020. Go 包（Package）：代码组织的基本单元

## 从一个问题开始

假设你写了一个计算器程序，代码越来越多：

```go
// main.go - 所有代码都在这一个文件里
package main

func add(a, b int) int { return a + b }
func subtract(a, b int) int { return a - b }
func multiply(a, b int) int { return a * b }
func divide(a, b int) int { return a / b }
func formatResult(result int) string { ... }
func logError(msg string) { ... }
// ... 还有 50 个函数

func main() { ... }
```

问题来了：

- 文件太长，找函数要滚动半天
- 想在另一个项目复用 `add` 函数？只能复制粘贴
- 团队协作时，大家都改这一个文件，冲突不断

**包（Package）就是用来解决这些问题的。**

---

## 1. 先理解"模块"（Module）

在讲"包"之前，必须先理解"模块"，因为**你必须先创建模块，才能创建和使用包**。

### 1.1 模块是什么？

**模块 = 一个项目**。它由一个 `go.mod` 文件定义。

当你开始一个新的 Go 项目时，第一步就是创建模块：

```bash
mkdir myapp          # 创建项目目录
cd myapp             # 进入目录
go mod init myapp    # 创建模块
```

执行 `go mod init myapp` 后，会生成一个 `go.mod` 文件：

```
myapp/
└── go.mod
```

打开 `go.mod`，内容很简单：

```
module myapp

go 1.21
```

**逐行解释：**

| 行               | 含义                                               |
| ---------------- | -------------------------------------------------- |
| `module myapp` | 模块的名字叫 `myapp`，后面导入包时会用到这个名字 |
| `go 1.21`      | 这个项目需要 Go 1.21 或更高版本                    |

### 1.2 模块名的选择

模块名可以随便起吗？

- **本地项目**：可以用简单的名字，如 `myapp`、`calculator`
- **要发布到 GitHub**：应该用完整路径，如 `github.com/yourname/myapp`

后者的好处是全球唯一，别人可以直接 `go get github.com/yourname/myapp` 下载你的代码。

---

## 2. 包（Package）是什么？

**包 = 一个目录里的所有 .go 文件**。

来看一个具体的例子：

```
myapp/                    <-- 模块根目录
├── go.mod                <-- 模块定义文件
├── main.go               <-- main 包（程序入口）
└── greetings/            <-- greetings 包（一个目录 = 一个包）
    └── greetings.go
```

这个结构里有 **2 个包**：

1. `main` 包：由 `main.go` 组成，是程序的入口
2. `greetings` 包：由 `greetings/` 目录下的所有 `.go` 文件组成

### 2.1 包的规则

1. **同一目录下的所有 .go 文件必须声明相同的包名**
2. **包名通常与目录名相同**（不是强制，但这是惯例）
3. **程序入口必须是 `main` 包的 `main` 函数**

---

## 3. 动手：创建和使用一个包

### 3.1 创建 greetings 包

在 `greetings/greetings.go` 中：

```go
// 第一行必须声明包名
// 这个文件属于 greetings 包
package greetings

// Hello 返回问候语
// 函数名首字母大写 = 导出（其他包可以调用）
func Hello(name string) string {
    return "你好, " + name + "!"
}

// goodbye 是私有函数
// 函数名首字母小写 = 未导出（只能在本包内使用）
func goodbye(name string) string {
    return "再见, " + name + "!"
}
```

**关键点：大小写决定可见性**

| 命名方式                | 可见性 | 谁能访问                      |
| ----------------------- | ------ | ----------------------------- |
| `Hello`（大写开头）   | 导出   | 任何包都能调用                |
| `goodbye`（小写开头） | 未导出 | 只有 `greetings` 包内部能用 |

### 3.2 使用 greetings 包

在 `main.go` 中：

```go
package main

import (
    "fmt"

    // 导入自定义包
    // 格式：模块名/包目录名
    "myapp/greetings"
)

func main() {
    // 调用包中的函数
    // 语法：包名.函数名(参数)
    message := greetings.Hello("世界")
    fmt.Println(message)  // 输出：你好, 世界!

    // 下面这行会编译失败（私有函数不能被外部调用）：
    // greetings.goodbye("世界")
}
```

**导入路径拆解：**

```
"myapp/greetings"
   │       │
   │       └── 包所在的目录名
   └── go.mod 中定义的模块名
```

### 3.3 运行程序

```bash
cd myapp
go run .
```

输出：

```
你好, 世界!
```

---

## 4. 导出规则详解

Go 不用 `public`、`private` 这些关键字，而是用**首字母大小写**来控制可见性。

这个规则适用于：

- 函数：`Hello` vs `hello`
- 变量：`Name` vs `name`
- 常量：`Pi` vs `pi`
- 类型：`User` vs `user`
- 结构体字段：`User.Name` vs `User.name`

```go
package greetings

// 导出的常量（其他包可以用 greetings.MaxLength）
const MaxLength = 100

// 未导出的变量（其他包无法访问）
var defaultPrefix = "Hello, "

// 导出的类型
type Greeter struct {
    Name    string  // 导出的字段
    counter int     // 未导出的字段
}
```

---

## 5. 第三方包

Go 社区有大量现成的包可以使用。

### 5.1 在哪找？

官方包索引：**https://pkg.go.dev**

这里可以搜索、查看文档、了解使用方法。

### 5.2 如何安装？

使用 `go get` 命令：

```bash
go get github.com/fatih/color
```

这会：

1. 下载代码到本地
2. 在 `go.mod` 中添加依赖记录

### 5.3 如何使用？

```go
import "github.com/fatih/color"

func main() {
    color.Red("这是红色文字")
    color.Green("这是绿色文字")
}
```

**包路径 = 域名/用户名/仓库名**，这就是为什么你在教程里看到很多 GitHub 地址。

### 5.4 直接依赖 vs 间接依赖

执行 `go get` 后，`go.mod` 会新增依赖记录：

```
require (
    github.com/fatih/color v1.18.0 // indirect
    github.com/mattn/go-colorable v0.1.13 // indirect
    ...
)
```

注意 `// indirect` 这个注释，它表示"间接依赖"——你的代码没有直接 import 它，是其他包依赖了它。

**问题来了**：如果你在代码中直接 `import "github.com/fatih/color"`，但 go.mod 里它被标记为 `// indirect`，Go 会报错：

```
github.com/fatih/color should be direct
```

**原因**：你明明直接用了这个包，却被标记成"间接"，这是矛盾的。

**解决方案**：运行 `go mod tidy`

```bash
go mod tidy
```

执行后，go.mod 会自动修正：

```
require github.com/fatih/color v1.18.0

require (
    github.com/mattn/go-colorable v0.1.13 // indirect
    github.com/mattn/go-isatty v0.0.20 // indirect
    golang.org/x/sys v0.25.0 // indirect
)
```

现在 `color` 变成了直接依赖（没有 `// indirect`），而它依赖的其他包保持间接依赖。

### 5.5 go mod tidy 是什么？

`go mod tidy` 是依赖清理命令，它会：

1. **添加缺失的依赖**：代码里 import 了但 go.mod 没记录的
2. **移除多余的依赖**：go.mod 记录了但代码没用到的
3. **修正依赖标记**：把直接依赖和间接依赖正确分类

**建议**：每次添加或删除 import 后，都运行一次 `go mod tidy`。

### 5.6 go.sum 是什么？

当执行 `go get` 或 `go mod tidy` 后，会自动生成一个 `go.sum` 文件：

```
github.com/fatih/color v1.18.0 h1:S8gINlzdQ840/4pfAwic/ZE0djQEH3wM94VfqLTZcOM=
github.com/fatih/color v1.18.0/go.mod h1:4FelSpRwEGDpQ12mAdzqdOukCy4u8WUtOY6lkT/6HfU=
...
```

**这是什么？**

`go.sum` 记录了每个依赖包的**校验和（checksum）**，就像文件的"指纹"。

**为什么需要它？**

想象这个场景：

1. 你今天下载了 `color v1.18.0`，代码正常运行
2. 一个月后，有人恶意篡改了 GitHub 上 `color v1.18.0` 的代码
3. 你的同事 clone 项目后执行 `go mod download`，下载到的是被篡改的代码

`go.sum` 就是为了防止这种情况。Go 会：

1. 下载依赖时，计算代码的校验和
2. 与 `go.sum` 中记录的校验和对比
3. 如果不匹配，拒绝使用，报错提示

**前因后果：**

| 时机                      | 发生了什么                         |
| ------------------------- | ---------------------------------- |
| `go get`                | 下载依赖，计算校验和，写入 go.sum  |
| `go mod tidy`           | 同上，同时清理不需要的记录         |
| `go build` / `go run` | 验证依赖的校验和是否与 go.sum 匹配 |

**要提交到 Git 吗？**

**是的**。`go.mod` 和 `go.sum` 都应该提交到版本控制，确保团队所有人使用完全相同的依赖版本。

---

## 6. 本章项目结构

```
020-packages/
├── go.mod               # 模块定义，记录依赖版本
├── go.sum               # 依赖校验和，确保安全
├── main.go              # 主程序
└── greetings/           # greetings 包
    └── greetings.go
```

---

## 7. go.work：多模块协作

### 7.1 遇到的问题

本章的 `020-packages` 目录有自己的 `go.mod`，是一个独立模块。但它位于 `fundamentals` 目录下，而 `fundamentals` 也有自己的 `go.mod`。

当你在 `020-packages` 目录执行 `go run .` 时，Go 会向上查找 `go.mod`，找到 `fundamentals/go.mod`，然后报错：

```
main module (github.com/xxx/fundamentals) does not contain package ...
```

**原因**：Go 认为你在 `fundamentals` 模块里，但 `fundamentals` 模块不包含 `020-packages` 这个包（因为 `020-packages` 有自己的 `go.mod`，是独立模块）。

### 7.2 go.work 是什么？

`go.work` 是 Go 1.18 引入的**工作区文件**，用于同时开发多个模块。

它告诉 Go："这些目录下的模块是一起工作的，请把它们当作一个整体来处理。"

### 7.3 如何创建？

在项目根目录执行：

```bash
go work init ./fundamentals ./fundamentals/020-packages
```

这会生成 `go.work` 文件：

```
go 1.21

use (
    ./fundamentals
    ./fundamentals/020-packages
)
```

**逐行解释：**

| 行 | 含义 |
|---|------|
| `go 1.21` | 工作区需要的 Go 版本 |
| `use (...)` | 列出工作区包含的所有模块目录 |

### 7.4 如何添加新模块？

如果后续创建了新模块（比如 `021-testing`），需要手动添加到 `go.work`：

```bash
go work use ./fundamentals/021-testing
```

或者直接编辑 `go.work` 文件，添加一行。

**注意**：`go work use` 是手动命令，不会自动执行。当你创建新的独立模块时，必须手动把它加入工作区。

### 7.5 本项目的情况

本项目的 `go.work` 是这样的：

```
go 1.25.5

use (
    ./fundamentals
    ./fundamentals/020-packages
)
```

当我创建 `020-packages/go.mod` 后，**手动**执行了 `go work use ./fundamentals/020-packages`，把新模块加入工作区。这不是自动发生的。

### 7.6 什么时候需要 go.work？

| 场景 | 是否需要 go.work |
|------|-----------------|
| 单模块项目（只有一个 go.mod） | 不需要 |
| 多模块项目（多个独立的 go.mod） | 需要 |
| 同时开发主项目和本地依赖库 | 需要 |

**普通项目通常不需要 go.work**。本项目使用它是因为学习目录结构特殊——每个章节可能是独立模块。

### 7.7 要提交到 Git 吗？

**看情况**：

- **个人学习项目**：可以提交，方便自己使用
- **团队项目**：通常**不提交**，加入 `.gitignore`，因为每个人的本地开发环境可能不同

---

## 核心概念总结

| 概念                     | 是什么                           | 文件                 |
| ------------------------ | -------------------------------- | -------------------- |
| **模块（Module）** | 一个项目，包含一个或多个包       | `go.mod`           |
| **包（Package）**  | 一个目录里的所有 .go 文件        | `package xxx` 声明 |
| **导出**           | 首字母大写的标识符，其他包可访问 | `func Hello()`     |
| **未导出**         | 首字母小写的标识符，仅本包可访问 | `func hello()`     |

---

## 运行示例

```bash
cd fundamentals/020-packages
go run .
```

输出：

```
你好, 世界!
```
