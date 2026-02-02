# 015. Swagger 接口文档集成

本章学习如何在 Go + Gin 项目中集成 Swagger 接口文档，实现 API 的自动化文档生成与可视化展示。

## 1. Swagger 概述

### 1.1 什么是 Swagger

Swagger 是一套围绕 **OpenAPI 规范**（OpenAPI Specification, OAS）的工具集。核心能力：

- **API 规范定义**：使用 JSON/YAML 描述 API 的端点、参数、响应
- **文档可视化**：Swagger UI 将规范渲染为可交互的网页
- **代码生成**：根据规范自动生成客户端 SDK、服务端桩代码

### 1.2 Go 生态的实现方案

Go 社区主流方案是 **swag + gin-swagger**：

```
源代码（Handler 注释）
        │
        ▼  swag init（命令行工具）
┌─────────────────────────────┐
│  internal/swagger/docs.go   │  生成的 Go 文件（包含规范数据）
│  internal/swagger/*.json    │  OpenAPI 2.0 规范 (JSON)
│  internal/swagger/*.yaml    │  OpenAPI 2.0 规范 (YAML)
└─────────────────────────────┘
        │
        ▼  gin-swagger（中间件）
┌─────────────────────────────┐
│  GET /swagger/*             │  Swagger UI 界面
└─────────────────────────────┘
```

### 1.3 与 NestJS 的对比

| 方面 | NestJS (@nestjs/swagger) | Go (swag + gin-swagger) |
|------|--------------------------|-------------------------|
| 标注方式 | TypeScript 装饰器 | 代码注释 |
| 类型推导 | 自动从 TS 类型推导 | 需在注释中显式声明 |
| 生成时机 | 运行时动态生成 | 编译前静态生成 |
| 文档位置 | 内存中动态构建 | 独立的 JSON/YAML 文件 |

## 2. 核心组件

### 2.1 swag CLI

命令行工具，扫描代码中的注释，生成 OpenAPI 规范文件。

```bash
# 安装
go install github.com/swaggo/swag/cmd/swag@latest

# 生成文档（在项目根目录执行）
swag init -g cmd/server/main.go
```

### 2.2 gin-swagger

Gin 中间件，提供 Swagger UI 的 HTTP 服务。

```go
import (
    swaggerFiles "github.com/swaggo/files"
    ginSwagger "github.com/swaggo/gin-swagger"
)

// 注册 Swagger 路由
r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
```

### 2.3 生成的文件结构

执行 `swag init` 后，会在项目中生成以下文件：

```
internal/swagger/
├── docs.go          # Go 代码，包含序列化的规范数据（package swagger）
├── swagger.json     # OpenAPI 2.0 规范 (JSON 格式)
└── swagger.yaml     # OpenAPI 2.0 规范 (YAML 格式)
```

**架构决策：为什么放在 `internal/swagger` 而不是 `docs/`？**

本项目的 `docs/` 目录存放的是学习笔记和架构文档（Human Readable），属于"脚手架资料"，在项目成熟后可能会被删除或归档。而 Swagger 生成的 `docs.go` 是**运行时代码**，服务启动依赖它（通过空导入注册规范）。

将运行时代码放在随时可能被删除的文件夹里是架构隐患。因此我们选择 `internal/swagger/`：

- **生命周期安全**：`internal/` 是应用核心代码，不会被误删
- **架构一致性**：Swagger 是内部实现的衍生品，放在 `internal` 下符合逻辑
- **职责分离**：`docs/` 保持纯净只放 Markdown，`internal/swagger/` 存放生成代码

## 3. 注释语法详解

swag 通过解析特定格式的注释来生成文档。注释分为两类：**全局注释**和**API 注释**。

### 3.1 全局注释（General API Info）

放在 `main.go` 的 `main` 函数或包声明上方，定义 API 的全局信息。

```go
// @title           Go API Template
// @version         1.0
// @description     一个基于整洁架构的 Go API 服务模板
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.example.com/support
// @contact.email  support@example.com

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8080
// @BasePath  /api/v1

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description 输入格式: Bearer {token}

func main() {
    // ...
}
```

常用全局注释：

| 注释 | 说明 | 示例 |
|------|------|------|
| `@title` | API 标题 | `@title My API` |
| `@version` | 版本号 | `@version 1.0` |
| `@description` | 描述信息 | `@description 这是一个示例 API` |
| `@host` | 服务地址 | `@host localhost:8080` |
| `@BasePath` | 基础路径 | `@BasePath /api/v1` |
| `@securityDefinitions.apikey` | 定义 API Key 认证 | 见上方示例 |

### 3.2 API 注释（API Operation）

放在每个 Handler 函数上方，定义单个 API 端点的详细信息。

```go
// SayHello 向指定用户发送问候
// @Summary      发送问候
// @Description  向指定用户发送问候消息，返回问候语和访问计数
// @Tags         greeter
// @Accept       json
// @Produce      json
// @Param        request body     dto.SayHelloRequest true "问候请求"
// @Success      200     {object} response.Response{data=v1.SayHelloResponse}
// @Failure      400     {object} response.Response
// @Failure      500     {object} response.Response
// @Router       /greeter/say-hello [post]
func handleSayHello(svc *service.GreeterService) gin.HandlerFunc {
    // ...
}
```

### 3.3 常用 API 注释详解

#### @Summary 和 @Description

```go
// @Summary      发送问候           // 简短摘要（显示在 API 列表中）
// @Description  详细描述...        // 完整描述（展开后显示）
```

#### @Tags

将 API 分组，便于在 Swagger UI 中组织展示。

```go
// @Tags greeter           // 单个标签
// @Tags greeter,hello     // 多个标签
```

#### @Accept 和 @Produce

定义请求和响应的 MIME 类型。

```go
// @Accept  json           // 接受 JSON 请求体
// @Produce json           // 返回 JSON 响应
```

常用值：`json`, `xml`, `plain`, `html`, `mpfd`（multipart/form-data）

#### @Param

定义请求参数，语法：`@Param 参数名 位置 类型 是否必填 "描述"`

```go
// 路径参数
// @Param id path int true "用户 ID"

// 查询参数
// @Param page query int false "页码" default(1)
// @Param size query int false "每页数量" default(10)

// 请求体
// @Param request body dto.SayHelloRequest true "请求体"

// Header
// @Param Authorization header string true "Bearer Token"
```

参数位置：`path`、`query`、`body`、`header`、`formData`

#### @Success 和 @Failure

定义响应，语法：`@Success 状态码 {类型} 数据类型 "描述"`

```go
// 简单响应
// @Success 200 {object} dto.UserResponse "成功"

// 嵌套泛型响应（使用统一响应结构）
// @Success 200 {object} response.Response{data=dto.UserResponse} "成功"

// 数组响应
// @Success 200 {array} dto.UserResponse "用户列表"

// 无响应体
// @Success 204 "删除成功"
```

#### @Router

定义路由路径和 HTTP 方法，语法：`@Router 路径 [方法]`

```go
// @Router /users/{id} [get]
// @Router /users [post]
// @Router /users/{id} [put]
// @Router /users/{id} [delete]
```

**注意**：路径相对于 `@BasePath`，路径参数使用 `{param}` 格式（不是 `:param`）

### 3.4 结构体注释

为结构体字段添加文档：

```go
// SayHelloRequest 问候请求
type SayHelloRequest struct {
    // 用户名称，长度 1-100 个字符
    Name string `json:"name" binding:"required,min=1,max=100" example:"World"`
}
```

结构体 Tag 说明：

| Tag | 说明 | 示例 |
|-----|------|------|
| `example` | 示例值，显示在文档中 | `example:"World"` |
| `enums` | 枚举值 | `enums:"active,inactive"` |
| `default` | 默认值 | `default:"active"` |
| `minimum` | 最小值（数字） | `minimum:"0"` |
| `maximum` | 最大值（数字） | `maximum:"100"` |
| `minLength` | 最小长度（字符串） | `minLength:"1"` |
| `maxLength` | 最大长度（字符串） | `maxLength:"100"` |

## 4. 项目实践

### 4.1 安装依赖

```bash
# 安装 swag CLI（开发环境）
go install github.com/swaggo/swag/cmd/swag@latest

# 安装 gin-swagger 依赖（项目依赖）
go get -u github.com/swaggo/gin-swagger
go get -u github.com/swaggo/files
```

### 4.2 项目结构

```
go-api-template/
├── cmd/server/
│   └── main.go              # 添加全局 Swagger 注释
├── docs/                    # 学习文档（与 Swagger 分离）
│   └── notes/
├── internal/
│   ├── swagger/             # Swagger 生成物（运行时代码）
│   │   ├── docs.go          # package swagger
│   │   ├── swagger.json
│   │   └── swagger.yaml
│   └── server/
│       ├── http.go          # 添加 API 注释
│       └── swagger.go       # Swagger UI 路由配置
└── Makefile                 # 添加 swagger 命令
```

### 4.3 配置全局注释

在 `cmd/server/main.go` 中添加：

```go
// @title           Go API Template
// @version         1.0
// @description     一个基于整洁架构的 Go API 服务模板

// @contact.name   开发者
// @contact.email  dev@example.com

// @host      localhost:8080
// @BasePath  /api/v1

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description 输入格式: Bearer {token}

func main() {
    // ...
}
```

### 4.4 添加 API 注释

在 Handler 函数上方添加注释，参见 `internal/server/http.go`。

### 4.5 注册 Swagger 路由

在 HTTP Server 初始化时注册 Swagger UI 路由：

```go
import (
    // 导入生成的 swagger 包（空导入，执行 init 函数注册规范）
    _ "go-api-template/internal/swagger"

    swaggerFiles "github.com/swaggo/files"
    ginSwagger "github.com/swaggo/gin-swagger"
)

// 注册 Swagger 路由（仅在非生产环境）
if cfg.App.Env != "production" {
    engine.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}
```

### 4.6 生成文档

```bash
# 生成 Swagger 文档到 internal/swagger 目录
swag init -g cmd/server/main.go -o internal/swagger --packageName swagger --parseDependency --parseInternal

# 或使用 Makefile（推荐）
make swagger
```

参数说明：
- `-g cmd/server/main.go`：指定包含全局注释的入口文件
- `-o internal/swagger`：输出目录
- `--packageName swagger`：生成的 Go 包名（默认是 docs）
- `--parseDependency`：解析外部依赖中的类型
- `--parseInternal`：解析 internal 包中的类型

### 4.7 访问 Swagger UI

启动服务后，访问：`http://localhost:8080/swagger/index.html`

## 5. 最佳实践

### 5.1 环境隔离

生产环境不应暴露 Swagger UI：

```go
// 仅在开发/测试环境启用 Swagger
if cfg.App.Env != "production" {
    engine.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}
```

### 5.2 与统一响应结构配合

项目使用了统一响应结构 `response.Response`，在 Swagger 中体现：

```go
// @Success 200 {object} response.Response{data=v1.SayHelloResponse} "成功"
// @Failure 400 {object} response.Response "请求参数错误"
```

### 5.3 保持注释与代码同步

- 将 `swag init` 命令加入 CI/CD 流程
- 代码审查时检查注释是否更新
- 在 Makefile 中将 swagger 生成作为构建前置步骤

### 5.4 DTO 与 Proto 的选择

在 Swagger 注释中引用类型时：

- **请求体**：优先使用 DTO（因为 DTO 包含 `binding` 和 `example` 等 tag）
- **响应体**：可以使用 Proto 生成的结构体（更贴近实际返回）

```go
// @Param   request body dto.SayHelloRequest true "请求参数"
// @Success 200 {object} response.Response{data=v1.SayHelloResponse}
```

## 6. 常见问题

### 6.1 swag init 报错找不到包

确保先执行 `go mod tidy`，且所有导入的包都存在。

### 6.2 Swagger UI 显示空白

检查是否导入了生成的 swagger 包：

```go
import _ "go-api-template/internal/swagger"
```

这个空导入会执行 `internal/swagger/docs.go` 中的 `init()` 函数，注册 Swagger 规范。

### 6.3 结构体字段不显示

确保字段是导出的（首字母大写），且有 `json` tag：

```go
type Request struct {
    Name string `json:"name"` // ✅ 显示
    age  int    `json:"age"`  // ❌ 不显示（未导出）
}
```

### 6.4 嵌套结构体显示不完整

使用泛型语法指定嵌套类型：

```go
// 错误：只显示 Response 结构，data 字段显示为 interface{}
// @Success 200 {object} response.Response

// 正确：明确指定 data 字段的类型
// @Success 200 {object} response.Response{data=dto.User}
```

## 7. 小结

| 组件 | 作用 |
|------|------|
| swag CLI | 解析注释，生成 OpenAPI 规范文件 |
| gin-swagger | Gin 中间件，提供 Swagger UI 界面 |
| internal/swagger/ | 存放生成的规范文件（运行时代码） |

开发流程：

1. 在 Handler 上方编写 Swagger 注释
2. 运行 `make swagger` 生成文档
3. 启动服务，访问 `/swagger/index.html` 查看

注意事项：

- 生产环境禁用 Swagger UI
- 保持注释与代码实现同步
- 善用 DTO 的 `example` tag 提供示例值
- Swagger 生成物放在 `internal/swagger/`，与学习文档 `docs/` 分离
