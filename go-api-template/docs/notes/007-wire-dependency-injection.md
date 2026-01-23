# 007. Wire 依赖注入实战

本章记录阶段四的完整操作步骤：使用 Google Wire 实现编译期依赖注入，替代手动依赖组装。

---

## 1. 核心问题与概念

### 1.1 解决什么问题

在阶段三中，我们的 `main.go` 包含大量手动依赖组装代码：

```go
dataLayer, _ := data.NewData()
greeterRepo := data.NewGreeterRepo(dataLayer)
greeterUsecase := biz.NewGreeterUsecase(greeterRepo)
greeterService := service.NewGreeterService(greeterUsecase)
httpServer := server.NewHTTPServer(greeterService)
```

随着项目规模增长，这种手动组装会带来以下问题：

| 问题     | 描述                                                |
| -------- | --------------------------------------------------- |
| 代码冗长 | 每新增一个服务都需要在 main.go 中添加多行初始化代码 |
| 顺序依赖 | 必须按正确顺序调用构造函数，否则编译失败            |
| 容易出错 | 参数传错、遗漏依赖等错误只在编译时发现              |
| 难以维护 | 修改依赖关系需要在多处同步更新                      |

### 1.2 什么是依赖注入 (DI)

依赖注入是一种设计模式：组件不自己创建依赖，而是由外部传入。

```go
// ❌ 传统方式：组件自己创建依赖
type GreeterUsecase struct {
    repo *MySQLGreeterRepo  // 硬编码具体实现
}
func NewGreeterUsecase() *GreeterUsecase {
    return &GreeterUsecase{
        repo: &MySQLGreeterRepo{},  // 内部创建
    }
}

// ✅ 依赖注入：依赖由外部传入
type GreeterUsecase struct {
    repo GreeterRepo  // 接口类型
}
func NewGreeterUsecase(repo GreeterRepo) *GreeterUsecase {
    return &GreeterUsecase{repo: repo}  // 外部传入
}
```

### 1.3 Wire 是什么

Wire 是 Google 开发的 Go 依赖注入工具，特点：

| 特性       | 说明                                     |
| ---------- | ---------------------------------------- |
| 编译期生成 | 在编译前生成依赖组装代码，不是运行时反射 |
| 类型安全   | 依赖关系错误在编译时发现，不会运行时崩溃 |
| 零开销     | 生成的是普通 Go 代码，无运行时性能损耗   |
| 代码可读   | 生成的 `wire_gen.go` 是可读的 Go 代码  |

### 1.4 核心概念

| 概念        | 定义                       | 示例                                     |
| ----------- | -------------------------- | ---------------------------------------- |
| Provider    | 构造函数，返回某个依赖     | `func NewData() (*Data, error)`        |
| ProviderSet | 一组相关 Provider 的集合   | `wire.NewSet(NewData, NewGreeterRepo)` |
| Injector    | 入口函数，声明最终需要什么 | `func wireApp() (*gin.Engine, error)`  |

---

## 2. 核心用法

### 2.1 Provider（提供者）

**Provider 就是构造函数**，不是数据结构、DTO 或实体。Wire 通过分析构造函数的签名来理解依赖关系。

```go
// Provider：一个构造函数
// - 参数：该组件需要的依赖
// - 返回值：该组件提供的类型
func NewGreeterRepo(data *Data) biz.GreeterRepo {
    return &greeterRepo{data: data}
}
```

Wire 分析这个函数签名：

- **输入参数** `*Data` → 这是它需要的依赖
- **返回值** `biz.GreeterRepo` → 这是它提供的类型

**常见误解澄清**：

| 放入 ProviderSet 的是 | 不是 |
|----------------------|------|
| `NewGreeterRepo`（构造函数） | `greeterRepo`（结构体类型） |
| `NewGreeterUsecase`（构造函数） | `GreeterUsecase`（结构体类型） |
| `NewData`（构造函数） | `Data`（结构体类型） |

```go
// ✅ 正确：放入构造函数
var ProviderSet = wire.NewSet(NewGreeterRepo)

// ❌ 错误：不能放结构体类型
var ProviderSet = wire.NewSet(greeterRepo{})  // 编译错误
```

### 2.2 ProviderSet（提供者集合）

将相关的 Provider（构造函数）组织成集合，便于管理：

```go
// internal/data/data.go
var ProviderSet = wire.NewSet(
    NewData,        // 构造函数，提供 *Data
    NewGreeterRepo, // 构造函数，提供 biz.GreeterRepo
)
```

**ProviderSet 里应该放什么？**

| 应该放 | 示例 | 说明 |
|--------|------|------|
| 构造函数 | `NewData`, `NewGreeterRepo` | 返回某个依赖实例的函数 |
| 其他 ProviderSet | `GreeterProviderSet` | 用于聚合多个模块 |

**不应该放什么？**

| 不应该放 | 示例 | 原因 |
|----------|------|------|
| 结构体类型 | `Data{}`, `greeterRepo{}` | Wire 需要知道如何创建，不是类型本身 |
| 接口类型 | `biz.GreeterRepo` | 接口是抽象，Wire 需要具体的构造函数 |
| 实例变量 | `myRepo` | Wire 在编译期生成代码，不接受运行时变量 |

**如何判断某个函数是否是 Provider？**

满足以下条件的函数就是 Provider：
1. 返回一个或多个值（第一个是提供的类型，可选第二个是 error）
2. 参数是它需要的依赖（可以为空）

```go
// ✅ 有效的 Provider 签名
func NewData() (*Data, error)              // 无依赖，可能失败
func NewGreeterRepo(d *Data) biz.GreeterRepo  // 有依赖，不会失败
func NewHTTPServer(svc *GreeterService) *gin.Engine  // 有依赖，不会失败
```

**实际项目中每层放什么？**

以本项目为例，每层的构造函数和职责：

| 层 | 构造函数 | 它创建什么 | 职责 |
|---|---------|-----------|------|
| data | `NewData` | 数据库连接/内存存储 | 管理数据源 |
| data | `NewGreeterRepo` | 问候数据访问器 | 实现 CRUD 操作 |
| biz | `NewGreeterUsecase` | 问候业务处理器 | 执行业务逻辑 |
| service | `NewGreeterService` | 问候 API 服务 | 处理请求/响应转换 |
| server | `NewHTTPServer` | HTTP 服务器 | 路由注册、中间件 |

**规律：每层放的是"创建该层核心对象的函数"**

| 层 | 核心对象类型 | 构造函数命名规范 |
|---|-------------|-----------------|
| data | Repository（数据访问器） | `NewXxxRepo` |
| biz | Usecase（业务处理器） | `NewXxxUsecase` |
| service | Service（API 服务） | `NewXxxService` |
| server | Server（协议服务器） | `NewHTTPServer` / `NewGRPCServer` |

**扩展示例：如果新增 User 模块**

| 层 | 新增的构造函数 | 放入哪个 ProviderSet |
|---|---------------|---------------------|
| data | `NewUserRepo` | `data.UserProviderSet` |
| biz | `NewUserUsecase` | `biz.UserProviderSet` |
| service | `NewUserService` | `service.UserProviderSet` |

### 2.3 Injector（注入器）

声明最终需要的依赖，Wire 会自动生成组装代码：

```go
//go:build wireinject

func wireApp() (*gin.Engine, error) {
    wire.Build(
        data.ProviderSet,
        biz.ProviderSet,
        service.ProviderSet,
        server.ProviderSet,
    )
    return nil, nil  // 占位，Wire 会替换
}
```

---

## 3. 深度原理与机制

### 3.1 Wire 工作流程

```
┌─────────────────────────────────────────────────────────────┐
│  1. 开发者编写 wire.go                                       │
│     - 定义 Injector 函数（声明需要什么）                      │
│     - 引用所有 ProviderSet                                   │
└──────────────────────────┬──────────────────────────────────┘
                           │
                           ▼
┌─────────────────────────────────────────────────────────────┐
│  2. 运行 wire ./cmd/server/                                  │
│     - 分析所有 Provider 的函数签名                            │
│     - 构建依赖图                                             │
│     - 检测循环依赖                                           │
└──────────────────────────┬──────────────────────────────────┘
                           │
                           ▼
┌─────────────────────────────────────────────────────────────┐
│  3. 生成 wire_gen.go                                         │
│     - 按正确顺序调用构造函数                                  │
│     - 传递依赖参数                                           │
│     - 处理错误                                               │
└──────────────────────────┬──────────────────────────────────┘
                           │
                           ▼
┌─────────────────────────────────────────────────────────────┐
│  4. 编译时                                                   │
│     - wire.go 被忽略（wireinject 构建标签）                   │
│     - wire_gen.go 被编译（!wireinject 构建标签）              │
└─────────────────────────────────────────────────────────────┘
```

### 3.2 构建标签的作用

Wire 使用 Go 的构建标签来控制哪个文件被编译：

```go
// wire.go
//go:build wireinject     // 只在 Wire 运行时读取
// +build wireinject      // 兼容旧版 Go

// wire_gen.go
//go:build !wireinject    // 正常编译时使用
// +build !wireinject     // 排除 Wire 运行时
```

### 3.3 依赖图分析

Wire 分析每个 Provider 的输入输出，构建依赖图：

```
wireApp() 需要 *gin.Engine
    │
    └── NewHTTPServer 需要 *service.GreeterService
            │
            └── NewGreeterService 需要 *biz.GreeterUsecase
                    │
                    └── NewGreeterUsecase 需要 biz.GreeterRepo
                            │
                            └── NewGreeterRepo 需要 *data.Data
                                    │
                                    └── NewData 不需要输入
```

然后反向遍历，生成正确的调用顺序：

```go
dataData, err := data.NewData()
greeterRepo := data.NewGreeterRepo(dataData)
greeterUsecase := biz.NewGreeterUsecase(greeterRepo)
greeterService := service.NewGreeterService(greeterUsecase)
engine := server.NewHTTPServer(greeterService)
```

---

## 4. 最佳实践与坑

### ✅ 推荐做法

| 实践                         | 说明                                                     |
| ---------------------------- | -------------------------------------------------------- |
| 模块级 ProviderSet + 层级聚合 | 每个模块定义自己的 ProviderSet，层级主文件只做聚合       |
| 返回接口类型                 | `NewGreeterRepo` 返回 `biz.GreeterRepo` 而非具体类型 |
| 错误优先                     | Provider 如有初始化错误，返回 error                      |

### 模块级 ProviderSet 模式（扩展性最佳）

当业务模块增多时，采用**模块级 ProviderSet + 层级聚合**模式：

```
internal/biz/
├── biz.go           ← 只做聚合
├── greeter.go       ← 定义 GreeterProviderSet
├── user.go          ← 定义 UserProviderSet（未来）
└── order.go         ← 定义 OrderProviderSet（未来）
```

**模块文件 (greeter.go)**：

```go
package biz

import "github.com/google/wire"

// GreeterProviderSet 是 Greeter 模块的依赖提供者集合
var GreeterProviderSet = wire.NewSet(NewGreeterUsecase)

// ... 业务代码
```

**聚合文件 (biz.go)**：

```go
package biz

import "github.com/google/wire"

// ProviderSet 聚合 biz 层所有模块的 ProviderSet
var ProviderSet = wire.NewSet(
    GreeterProviderSet,
    // UserProviderSet,    // 未来：用户模块
    // OrderProviderSet,   // 未来：订单模块
)
```

**优势**：
- 新增模块时，只需在新文件定义 ProviderSet，然后添加一行到聚合文件
- 符合开闭原则：扩展不修改已有模块代码
- 各模块独立管理自己的依赖

### ❌ 避免做法

| 反模式                     | 问题                                 |
| -------------------------- | ------------------------------------ |
| 修改 wire_gen.go           | 这是自动生成的，重新运行 wire 会覆盖 |
| 循环依赖                   | A 依赖 B，B 依赖 A，Wire 会报错      |
| 忘记添加到 ProviderSet     | 新增 Provider 后必须加入 ProviderSet |
| 多个 Provider 返回相同类型 | 会导致歧义，Wire 无法决定用哪个      |

---

## 5. 阶段四完成后的目录结构

```
go-api-template/
├── cmd/
│   └── server/
│       ├── main.go           ← 简化后的入口
│       ├── wire.go           ← Wire 注入器声明
│       └── wire_gen.go       ← Wire 自动生成的代码
├── internal/
│   ├── biz/
│   │   ├── biz.go            ← 包含 ProviderSet
│   │   └── greeter.go
│   ├── data/
│   │   ├── data.go           ← 包含 ProviderSet
│   │   └── greeter.go
│   ├── server/
│   │   ├── server.go         ← 包含 ProviderSet
│   │   └── http.go
│   └── service/
│       ├── service.go        ← 包含 ProviderSet
│       └── greeter.go
├── go.mod
├── go.sum
└── Makefile                  ← 包含 wire 命令
```

---

## 6. 行动导向 (Action Guide)

### Step 1: 安装 Wire CLI

**这一步在干什么**：Wire CLI 是代码生成工具，用于分析 `wire.go` 并生成 `wire_gen.go`。

```powershell
go install github.com/google/wire/cmd/wire@latest
```

验证安装：

```powershell
wire help
```

### Step 2: 添加 Wire 库依赖

**这一步在干什么**：Wire 库提供 `wire.NewSet` 和 `wire.Build` 等 API，用于声明依赖关系。

```powershell
cd go-api-template
go get github.com/google/wire
```

### Step 3: 为每层创建 ProviderSet

**这一步在干什么**：采用模块级 ProviderSet + 层级聚合模式，便于未来扩展。

**biz 层**：

```go
// internal/biz/greeter.go - 模块文件
package biz

import "github.com/google/wire"

// GreeterProviderSet 是 Greeter 模块的依赖提供者集合
var GreeterProviderSet = wire.NewSet(NewGreeterUsecase)

// ... 业务代码
```

```go
// internal/biz/biz.go - 聚合文件
package biz

import "github.com/google/wire"

// ProviderSet 聚合 biz 层所有模块的 ProviderSet
var ProviderSet = wire.NewSet(
    GreeterProviderSet,
    // UserProviderSet,  // 未来扩展
)
```

**data 层**：

```go
// internal/data/greeter.go - 模块文件
var GreeterProviderSet = wire.NewSet(NewGreeterRepo)
```

```go
// internal/data/data.go - 聚合文件
var ProviderSet = wire.NewSet(
    NewData,             // 基础设施
    GreeterProviderSet,  // Greeter 模块
)
```

**service 层**：

```go
// internal/service/greeter.go
var GreeterProviderSet = wire.NewSet(NewGreeterService)

// internal/service/service.go
var ProviderSet = wire.NewSet(GreeterProviderSet)
```

**server 层**：

```go
// internal/server/server.go
// server 层按协议划分（HTTP/gRPC），不按业务模块划分
var ProviderSet = wire.NewSet(
    NewHTTPServer,
    // NewGRPCServer, // 未来扩展
)
```

### Step 4: 创建 Wire 注入器

**这一步在干什么**：创建 `wire.go` 文件，声明需要 Wire 生成什么（最终需要的依赖类型）。

**cmd/server/wire.go**：

```go
//go:build wireinject
// +build wireinject

package main

import (
    "github.com/gin-gonic/gin"
    "github.com/google/wire"

    "go-api-template/internal/biz"
    "go-api-template/internal/data"
    "go-api-template/internal/server"
    "go-api-template/internal/service"
)

// wireApp 是 Wire 的注入器函数
func wireApp() (*gin.Engine, error) {
    wire.Build(
        data.ProviderSet,
        biz.ProviderSet,
        service.ProviderSet,
        server.ProviderSet,
    )
    return nil, nil
}
```

### Step 5: 运行 Wire 生成代码

**这一步在干什么**：Wire 分析 `wire.go`，自动生成依赖组装代码到 `wire_gen.go`。

```powershell
wire ./cmd/server/
```

成功输出：

```
wire: go-api-template/cmd/server: wrote E:\Dev\go-journey\go-api-template\cmd\server\wire_gen.go
```

查看生成的代码：

```go
// wire_gen.go (自动生成，不要手动修改)
func wireApp() (*gin.Engine, error) {
    dataData, err := data.NewData()
    if err != nil {
        return nil, err
    }
    greeterRepo := data.NewGreeterRepo(dataData)
    greeterUsecase := biz.NewGreeterUsecase(greeterRepo)
    greeterService := service.NewGreeterService(greeterUsecase)
    engine := server.NewHTTPServer(greeterService)
    return engine, nil
}
```

### Step 6: 简化 main.go

**这一步在干什么**：使用 Wire 生成的 `wireApp()` 替代手动依赖组装代码。

**cmd/server/main.go**：

```go
package main

import "log"

func main() {
    // 使用 Wire 生成的函数初始化所有依赖
    httpServer, err := wireApp()
    if err != nil {
        log.Fatalf("Failed to initialize application: %v", err)
    }

    log.Println("Starting server on :8080")
    if err := httpServer.Run(":8080"); err != nil {
        log.Fatalf("Failed to start server: %v", err)
    }
}
```

### Step 7: 更新 Makefile

**这一步在干什么**：添加 `wire` 命令，方便后续修改依赖后重新生成。

```makefile
# 生成 Wire 依赖注入代码
wire:
	wire ./cmd/server/
```

### Step 8: 验证

**这一步在干什么**：确保使用 Wire 后服务仍能正常运行。

```powershell
# 构建
go build -o bin/server ./cmd/server

# 运行
go run ./cmd/server

# 测试
$response = Invoke-RestMethod -Uri "http://localhost:8080/api/v1/greeter/say-hello/Wire" -Method Get
$response | ConvertTo-Json
# {"message": "Hello, Wire! You are visitor #1."}
```

---

## 7. 阶段四小结

| 完成项      | 说明                                  |
| ----------- | ------------------------------------- |
| Wire CLI    | 已安装                                |
| Wire 库     | 已添加到项目依赖                      |
| ProviderSet | 四层（data/biz/service/server）各一个 |
| wire.go     | 声明注入器函数                        |
| wire_gen.go | Wire 自动生成的依赖组装代码           |
| main.go     | 简化为调用 wireApp()                  |
| Makefile    | 添加 wire 命令                        |

**代码对比**：

| 阶段三 (手动)          | 阶段四 (Wire)                        |
| ---------------------- | ------------------------------------ |
| 15+ 行依赖组装         | 1 行 wireApp() 调用                  |
| 手动维护顺序           | Wire 自动排序                        |
| 新增依赖需修改 main.go | 只需更新 ProviderSet 并重新运行 wire |

---

## 8. 下一步：阶段五

阶段五将引入 **Ent ORM**，替换内存存储为真实数据库：

1. 安装 Ent CLI
2. 定义 Schema（数据模型）
3. 生成 Ent 代码
4. 修改 `internal/data` 实现真实数据库操作
5. 配置数据库连接

**关键收益**：`biz` 和 `service` 层的代码一行都不需要改——这就是依赖倒置的威力。
