# 002. Go 工程化：项目结构与工作区模式 (Workspaces)

## 1. 核心问题与概念

### 1.1 我们要解决什么问题？

当一个 Go 项目变得复杂，你可能需要在**同一个仓库**下管理多个相对独立的代码库（例如：基础练习 `fundamentals`、Web 服务 `web-service`、工具库 `pkg`）。

传统做法是创建多个 Git 仓库，但这带来了麻烦：
*   本地联调困难（需要先 `git push`，再 `go get` 拉取更新）。
*   版本管理分散。

**Go Workspaces** 正是为了解决这个问题而设计的：它允许你在**一个项目根目录下**管理多个独立的 Go 模块，并且在本地开发时让它们无缝协作。

### 1.2 核心概念：`go.mod` 与 `go.work`

在深入之前，必须先理解 Go 的两个核心配置文件：

#### `go.mod` - 模块定义文件 (Module Definition)

*   **本质**：一个文件夹的"身份证"。它声明了这个文件夹是一个独立的 **Go Module (模块)**。
*   **作用**：
    1.  定义模块的**唯一路径/名称** (Module Path)，例如 `github.com/HJSunDev/go-learning-journey/fundamentals`。
    2.  记录该模块依赖的第三方包及其版本。
*   **生成方式**：`go mod init <module-path>`
*   **位置**：放在模块的根目录。每个模块有且只有一个 `go.mod`。

#### `go.work` - 工作区定义文件 (Workspace Definition)

*   **本质**：一个项目的"调度中心"。它告诉 Go 编译器："我这里有多个模块，请把它们当作一个整体来编译和运行。"
*   **作用**：
    1.  列出当前工作区包含哪些本地模块。
    2.  让这些模块之间可以**直接相互引用**，无需发布到网上或使用 `replace` 指令。
*   **生成方式**：`go work init <module-path-1> <module-path-2> ...`
*   **位置**：放在项目根目录（所有模块的共同父目录）。一个工作区只有一个 `go.work`。
*   **版本要求**：Go 1.18+

#### 它们的关系

```
go-journey/          <-- go.work 在这里 (工作区层级)
├── fundamentals/    <-- go.mod 在这里 (模块层级)
│   └── ...
└── web-service/     <-- 未来也会有自己的 go.mod (另一个模块)
    └── ...
```

`go.work` 是"管理者"，`go.mod` 是"被管理的单元"。

---

## 2. 项目结构设计

这是我们当前项目的完整目录规划：

```text
go-journey/
├── .git/               # Git 版本控制目录
├── .gitignore          # Git 忽略规则
├── go.work             # 【核心】工作区定义，管理下方的多个 module
├── README.md           # 项目说明 (GitHub 门面)
│
├── fundamentals/       # 【阶段一】Go 基础语法练习，独立 module
│   ├── go.mod          # 该模块的身份证
│   ├── docs/           # 该阶段的学习文档
│   ├── variables/      # 知识点：变量
│   │   └── main.go
│   └── concurrency/    # 知识点：并发
│       └── main.go
│
└── web-service/        # 【阶段二】未来的 Web 框架项目，独立 module
    ├── go.mod
    └── ...
```

### 2.1 为什么每个知识点一个子目录？

Go 语言规定：**一个目录下只能有一个 `main` 函数**（属于 `package main`）。

如果你把所有练习代码都放在 `fundamentals/` 根目录，就会因为有多个 `main` 函数而编译失败。所以我们为每个知识点创建独立的子目录（如 `variables/`、`concurrency/`），让它们各自拥有独立的 `main.go`，互不冲突。

---

## 3. 核心命令详解

以下是初始化项目结构时用到的核心 Go 命令。

### 3.1 `go mod init` - 初始化模块

```bash
# 进入 fundamentals 目录
cd fundamentals

# 初始化模块，指定模块路径
go mod init github.com/HJSunDev/go-learning-journey/fundamentals
```

**命令拆解**：
*   `go mod init`：初始化一个新的 Go 模块。
*   `github.com/HJSunDev/go-learning-journey/fundamentals`：模块的唯一标识符 (Module Path)。
    *   **为什么要加 `/fundamentals`？** 因为我们的 Git 仓库是 `go-learning-journey`，而 `fundamentals` 只是其中一个子模块。如果将来有人要引用这个模块，完整路径能明确指向仓库中的特定目录。

**执行结果**：在 `fundamentals/` 目录下生成 `go.mod` 文件。

### 3.2 `go work init` - 初始化工作区

```bash
# 回到项目根目录 (go-journey/)
cd ..

# 初始化工作区，并将 fundamentals 模块纳入管理
go work init ./fundamentals
```

**命令拆解**：
*   `go work init`：在当前目录创建 `go.work` 文件。
*   `./fundamentals`：指定要纳入工作区的模块路径（相对路径）。

**执行结果**：在 `go-journey/` 根目录生成 `go.work` 文件，内容类似：

```
go 1.25

use ./fundamentals
```

### 3.3 `go work use` - 添加新模块到工作区

当你将来创建 `web-service` 模块后，需要将其也加入工作区：

```bash
# 先初始化 web-service 模块
cd web-service
go mod init github.com/HJSunDev/go-learning-journey/web-service
cd ..

# 将新模块加入工作区
go work use ./web-service
```

执行后，`go.work` 会更新为：

```
go 1.25

use (
    ./fundamentals
    ./web-service
)
```

---

## 4. 最佳实践与注意事项

### ✅ 推荐做法

*   **模块路径使用 GitHub 仓库地址**：即使是本地练习项目，也使用 `github.com/username/repo/module` 格式，这是 Go 社区的标准规范，方便未来发布。
*   **`go.work` 文件纳入 `.gitignore`（可选）**：有些团队认为 `go.work` 是本地开发配置，不应提交。但对于学习项目，提交它可以帮助你在其他电脑上快速还原环境。

### ❌ 避免做法

*   **不要在模块根目录直接写 `main.go`**：如果 `fundamentals/` 下将来有多个练习，直接在根目录写会导致 `main` 函数冲突。始终为每个练习创建子目录。
*   **不要忘记 `go.mod` 中的模块路径**：如果你只写 `go mod init fundamentals`（不带 GitHub 前缀），虽然能运行，但不符合规范，将来发布或被他人引用时会出问题。

---

## 5. 行动导向 (Action Guide)

以下是从零初始化本项目结构的完整步骤清单。

### Step 1: 创建项目根目录

```bash
mkdir go-journey
cd go-journey
```

### Step 2: 初始化 Git 仓库并关联远程

```bash
git init
git remote add origin https://github.com/HJSunDev/go-learning-journey.git
git branch -M main
```

### Step 3: 创建 fundamentals 子模块

```bash
mkdir fundamentals
cd fundamentals

# 初始化 Go 模块，路径与 GitHub 仓库结构对应
go mod init github.com/HJSunDev/go-learning-journey/fundamentals

cd ..
```

### Step 4: 初始化工作区

```bash
# 在项目根目录执行
go work init ./fundamentals
```

### Step 5: 验证

检查生成的文件：
*   `go-journey/go.work` 存在，内容包含 `use ./fundamentals`。
*   `go-journey/fundamentals/go.mod` 存在，内容包含 `module github.com/.../fundamentals`。

运行 `go work sync` 确保工作区同步正常（无报错即成功）。
