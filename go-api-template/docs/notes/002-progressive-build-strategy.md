# 002. 渐进式构建策略

本章介绍从零构建 Go API 服务的渐进式策略，并详细记录阶段一的完整操作步骤。

---

## 1. 五阶段构建路线图

构建一个生产级 API 服务，采用渐进式策略，每完成一个阶段都确保项目可运行。

| 阶段 | 名称 | 目标 | 关键产物 |
|------|------|------|----------|
| 一 | 骨架与传输层 | 搭建目录结构，启动 HTTP 服务 | `main.go`, Gin Server |
| 二 | API 定义 | 使用 Protobuf 定义接口契约 | `.proto` 文件, 生成的 Go 代码 |
| 三 | 领域层与模拟数据 | 实现业务逻辑，内存存储 | `biz/`, `service/`, `data/` |
| 四 | 依赖注入 | 引入 Wire，自动组装依赖 | `wire.go`, `wire_gen.go` |
| 五 | 持久化 | 接入真实数据库 | Ent Schema, 数据库迁移 |

**核心原则**：每个阶段完成后，服务都能正常启动和响应请求。

---

## 2. 阶段一：骨架与传输层

### 2.1 目标

- 创建符合 Go Standard Layout 的目录结构
- 初始化 Go 模块
- 使用 Gin 启动一个监听 8080 端口的 HTTP 服务
- 不涉及任何业务逻辑

### 2.2 前置条件

- Go 环境已安装 (推荐 1.21+)
- 了解基本的终端操作

### 2.3 操作步骤

#### 步骤 1：创建项目目录

在你的工作区创建项目根目录：

```powershell
# 创建项目目录并进入
mkdir go-api-template
cd go-api-template
```

#### 步骤 2：创建目录结构

```powershell
# 创建所有必需的目录
New-Item -ItemType Directory -Force -Path "api/helloworld/v1", "cmd/server", "configs", "internal/biz", "internal/data", "internal/service", "internal/server", "ent", "third_party", "docs/notes"
```

目录说明：

| 目录 | 用途 |
|------|------|
| `api/` | Protobuf 定义文件 |
| `cmd/server/` | 程序入口 |
| `configs/` | 配置文件 |
| `internal/biz/` | 领域层（核心业务逻辑） |
| `internal/data/` | 数据层（仓储实现） |
| `internal/service/` | 应用层（API 实现） |
| `internal/server/` | 传输层（HTTP/gRPC 配置） |
| `ent/` | ORM 生成代码 |
| `third_party/` | 第三方 proto 文件 |
| `docs/` | 项目文档 |

#### 步骤 3：初始化 Go 模块

```powershell
go mod init go-api-template
```

执行后会生成 `go.mod` 文件：

```
module go-api-template

go 1.21
```

#### 步骤 4：添加 Gin 依赖

```powershell
go get github.com/gin-gonic/gin
```

#### 步骤 5：创建程序入口

创建 `cmd/server/main.go` 文件：

```powershell
New-Item -ItemType File -Force -Path "cmd/server/main.go"
```

写入以下内容：

```go
package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	engine := gin.Default()

	// 健康检查端点，用于验证服务是否正常运行
	engine.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
		})
	})

	// 根路径，返回服务基本信息
	engine.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"name":    "go-api-template",
			"version": "0.1.0",
			"message": "Welcome to Go API Template",
		})
	})

	log.Println("Starting server on :8080")
	if err := engine.Run(":8080"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
```

**代码解析**：

| 代码 | 说明 |
|------|------|
| `gin.Default()` | 创建带有 Logger 和 Recovery 中间件的 Gin 引擎 |
| `engine.GET(path, handler)` | 注册 GET 路由 |
| `c.JSON(code, obj)` | 返回 JSON 响应 |
| `gin.H{}` | `map[string]interface{}` 的快捷写法 |
| `engine.Run(":8080")` | 启动 HTTP 服务，监听 8080 端口 |

#### 步骤 6：创建 Makefile

创建 `Makefile` 文件用于统一管理构建命令：

```powershell
New-Item -ItemType File -Force -Path "Makefile"
```

写入以下内容：

```makefile
.PHONY: all build run clean proto wire help

# 默认目标
all: help

# 构建可执行文件
build:
	go build -o bin/server ./cmd/server

# 运行服务
run:
	go run ./cmd/server

# 清理构建产物
clean:
	rm -rf bin/

# 生成 Proto 代码 (阶段二实现)
proto:
	@echo "Proto generation not configured yet. See Phase 2."

# 生成 Wire 依赖注入代码 (阶段四实现)
wire:
	@echo "Wire generation not configured yet. See Phase 4."

# 整理依赖
tidy:
	go mod tidy

# 运行测试
test:
	go test -v ./...

# 显示帮助信息
help:
	@echo "Available targets:"
	@echo "  build  - Build the server binary"
	@echo "  run    - Run the server"
	@echo "  clean  - Remove build artifacts"
	@echo "  proto  - Generate code from proto files (Phase 2)"
	@echo "  wire   - Generate dependency injection code (Phase 4)"
	@echo "  tidy   - Run go mod tidy"
	@echo "  test   - Run tests"
```

#### 步骤 7：创建占位文件

为空目录创建 `.gitkeep` 文件，确保 Git 能跟踪这些目录：

```powershell
New-Item -ItemType File -Force -Path "api/helloworld/v1/.gitkeep", "configs/.gitkeep", "ent/.gitkeep", "third_party/.gitkeep"
```

#### 步骤 8：添加到 Go Workspace（如适用）

如果你的项目在 Go Workspace 中（存在 `go.work` 文件），需要将模块添加进去：

```powershell
go work use .
```

### 2.4 验证

#### 启动服务

```powershell
go run ./cmd/server
```

预期输出：

```
[GIN-debug] GET    /health                   --> main.main.func1 (3 handlers)
[GIN-debug] GET    /                         --> main.main.func2 (3 handlers)
2026/01/21 19:55:00 Starting server on :8080
[GIN-debug] Listening and serving HTTP on :8080
```

#### 停止服务

服务启动后会持续运行，占用当前终端。停止服务的方法：

| 方式 | 操作 | 说明 |
|------|------|------|
| 快捷键 | `Ctrl + C` | 在运行服务的终端窗口按下，发送中断信号 |
| 关闭终端 | 直接关闭终端窗口 | 终端进程结束，服务随之停止 |

按下 `Ctrl + C` 后，终端会显示类似以下内容，表示服务已停止：

```
^C
```

之后终端恢复可输入状态，你可以继续执行其他命令。

#### 测试端点

打开新的终端窗口，执行：

```powershell
# 测试健康检查
Invoke-RestMethod -Uri "http://localhost:8080/health" -Method Get

# 预期输出：
# status
# ------
# ok

# 测试根路径
Invoke-RestMethod -Uri "http://localhost:8080/" -Method Get

# 预期输出：
# name            version message
# ----            ------- -------
# go-api-template 0.1.0   Welcome to Go API Template
```

或使用 curl（如已安装）：

```bash
curl http://localhost:8080/health
# {"status":"ok"}

curl http://localhost:8080/
# {"message":"Welcome to Go API Template","name":"go-api-template","version":"0.1.0"}
```

### 2.5 阶段一完成后的目录结构

```
go-api-template/
├── api/
│   └── helloworld/v1/.gitkeep
├── cmd/
│   └── server/
│       └── main.go              ← Gin HTTP 服务入口
├── configs/.gitkeep
├── docs/
│   ├── notes/
│   │   └── 001-progressive-build-strategy.md
│   └── README.md
├── ent/.gitkeep
├── internal/
│   ├── biz/                     ← 待实现
│   ├── data/                    ← 待实现
│   ├── server/                  ← 待实现
│   └── service/                 ← 待实现
├── third_party/.gitkeep
├── go.mod
├── go.sum
├── LEARNING.md
└── Makefile
```

### 2.6 本阶段小结

| 完成项 | 说明 |
|--------|------|
| 目录结构 | 符合 Go Standard Layout + Kratos 分层 |
| go.mod | 模块初始化完成 |
| Gin 依赖 | HTTP 框架已引入 |
| main.go | 可启动的 HTTP 服务 |
| Makefile | 统一的构建命令入口 |
| 验证 | `/health` 和 `/` 端点正常响应 |

---

## 3. 下一步：阶段二

阶段二将引入 **Protobuf** 和 **Buf** 工具链：

1. 安装 Buf CLI
2. 编写 `api/helloworld/v1/greeter.proto`
3. 配置 `buf.yaml` 和 `buf.gen.yaml`
4. 生成 Go 代码

目标：理解 "Schema First" —— 先定义接口契约，再编写业务代码。
