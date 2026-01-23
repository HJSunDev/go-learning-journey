# 006. 领域层与模拟数据

本章记录阶段三的完整操作步骤：实现领域层业务逻辑，使用内存存储作为模拟数据层。

---

## 1. 阶段目标

- 在 `internal/biz/` 定义领域实体和 Repository 接口
- 在 `internal/data/` 使用内存 Map 实现 Repository
- 在 `internal/service/` 实现 proto 定义的服务接口
- 在 `internal/server/` 配置 HTTP 路由
- 体验**依赖倒置**的威力——业务逻辑不关心数据存储在哪里

---

## 2. 核心概念

### 2.1 依赖倒置原则 (DIP)

传统分层架构中，依赖方向是：`Service → Repository → Database`

依赖倒置后，依赖方向变为：

```
┌─────────────────────────────────────────────────────────────┐
│                        biz (领域层)                          │
│  - 定义 GreeterRepo 接口                                     │
│  - 实现 GreeterUsecase（业务逻辑）                            │
│  - 不知道数据存在哪里                                         │
└─────────────────────────────────────────────────────────────┘
                             ▲
                             │ 实现接口
                             │
┌─────────────────────────────────────────────────────────────┐
│                        data (数据层)                         │
│  - 实现 biz.GreeterRepo 接口                                 │
│  - 当前：内存 Map                                            │
│  - 未来：可替换为 MySQL、PostgreSQL、MongoDB 等               │
└─────────────────────────────────────────────────────────────┘
```

**关键理解**：`biz` 层定义 `GreeterRepo` 接口，`data` 层实现该接口。当需要更换存储时，只需修改 `data` 层，`biz` 层代码一行都不用改。

### 2.2 四层架构

```
┌─────────────┐
│   cmd/      │  入口：初始化所有依赖，启动服务
└──────┬──────┘
       │
       ▼
┌─────────────┐
│   server/   │  传输层：配置 HTTP 路由，处理请求/响应
└──────┬──────┘
       │
       ▼
┌─────────────┐
│  service/   │  应用层：实现 proto 接口，调用业务用例
└──────┬──────┘
       │
       ▼
┌─────────────┐      ┌─────────────┐
│    biz/     │◄─────│    data/    │
│   (定义)    │ 实现 │   (实现)     │
└─────────────┘      └─────────────┘
```

---

## 3. 代码实现

### 3.1 领域层 (internal/biz/)

#### biz.go - 包声明

```go
// Package biz 是领域层，包含核心业务逻辑和领域模型。
// 这是整洁架构的核心，不依赖任何外部层。
package biz
```

#### greeter.go - 领域实体与业务用例

```go
package biz

import (
    "context"
    "fmt"
    "time"
)

// Greeter 是领域实体，表示一条问候记录
type Greeter struct {
    ID        int64     // 唯一标识
    Name      string    // 被问候者名称
    Message   string    // 问候消息
    CreatedAt time.Time // 创建时间
}

// GreeterRepo 定义了问候数据的存储接口
// 这是依赖倒置的关键：接口定义在领域层，实现在数据层
type GreeterRepo interface {
    Save(ctx context.Context, g *Greeter) (*Greeter, error)
    GetByName(ctx context.Context, name string) (*Greeter, error)
    Count(ctx context.Context) (int64, error)
}

// GreeterUsecase 是问候业务用例，包含核心业务逻辑
type GreeterUsecase struct {
    repo GreeterRepo
}

// NewGreeterUsecase 创建 GreeterUsecase 实例
func NewGreeterUsecase(repo GreeterRepo) *GreeterUsecase {
    return &GreeterUsecase{repo: repo}
}

// SayHello 执行问候业务逻辑
func (uc *GreeterUsecase) SayHello(ctx context.Context, name string) (*Greeter, error) {
    count, err := uc.repo.Count(ctx)
    if err != nil {
        return nil, fmt.Errorf("failed to get count: %w", err)
    }

    message := fmt.Sprintf("Hello, %s! You are visitor #%d.", name, count+1)

    greeter := &Greeter{
        Name:      name,
        Message:   message,
        CreatedAt: time.Now(),
    }

    saved, err := uc.repo.Save(ctx, greeter)
    if err != nil {
        return nil, fmt.Errorf("failed to save greeter: %w", err)
    }

    return saved, nil
}
```

**代码解析**：

| 组件 | 职责 |
|------|------|
| `Greeter` 结构体 | 领域实体，表示核心业务对象 |
| `GreeterRepo` 接口 | 数据访问抽象，定义需要的数据操作 |
| `GreeterUsecase` | 业务用例，实现核心业务逻辑 |
| `NewGreeterUsecase` | 构造函数，通过依赖注入接收 Repo |

### 3.2 数据层 (internal/data/)

#### data.go - 数据层初始化

```go
package data

import "sync"

// Data 是数据层的核心结构
// 阶段三：使用内存存储
// 阶段五：将替换为真实数据库连接
type Data struct {
    greeterStore *sync.Map
    idCounter    int64
    mu           sync.Mutex
}

// NewData 创建并初始化 Data 实例
func NewData() (*Data, error) {
    return &Data{
        greeterStore: &sync.Map{},
        idCounter:    0,
    }, nil
}

// NextID 生成下一个自增 ID（并发安全）
func (d *Data) NextID() int64 {
    d.mu.Lock()
    defer d.mu.Unlock()
    d.idCounter++
    return d.idCounter
}
```

#### greeter.go - Repository 实现

```go
package data

import (
    "context"
    "sync/atomic"

    "go-api-template/internal/biz"
)

// greeterRepo 实现 biz.GreeterRepo 接口
type greeterRepo struct {
    data *Data
}

// NewGreeterRepo 创建 GreeterRepo 实例
// 返回接口类型，隐藏实现细节
func NewGreeterRepo(data *Data) biz.GreeterRepo {
    return &greeterRepo{data: data}
}

func (r *greeterRepo) Save(ctx context.Context, g *biz.Greeter) (*biz.Greeter, error) {
    g.ID = r.data.NextID()
    r.data.greeterStore.Store(g.ID, g)
    r.data.greeterStore.Store("name:"+g.Name, g)
    return g, nil
}

func (r *greeterRepo) GetByName(ctx context.Context, name string) (*biz.Greeter, error) {
    value, ok := r.data.greeterStore.Load("name:" + name)
    if !ok {
        return nil, nil
    }
    return value.(*biz.Greeter), nil
}

func (r *greeterRepo) Count(ctx context.Context) (int64, error) {
    return atomic.LoadInt64(&r.data.idCounter), nil
}
```

**代码解析**：

| 技术选择 | 说明 |
|----------|------|
| `sync.Map` | Go 内置的并发安全 Map，适合读多写少场景 |
| `sync.Mutex` | 保护 ID 计数器的并发访问 |
| 返回接口类型 | `NewGreeterRepo` 返回 `biz.GreeterRepo`，隐藏实现细节 |

### 3.3 应用层 (internal/service/)

```go
package service

import (
    "context"

    v1 "go-api-template/api/helloworld/v1"
    "go-api-template/internal/biz"
)

// GreeterService 实现 proto 定义的 GreeterServiceServer 接口
type GreeterService struct {
    v1.UnimplementedGreeterServiceServer
    uc *biz.GreeterUsecase
}

func NewGreeterService(uc *biz.GreeterUsecase) *GreeterService {
    return &GreeterService{uc: uc}
}

// SayHello 实现 GreeterServiceServer.SayHello 方法
func (s *GreeterService) SayHello(ctx context.Context, req *v1.SayHelloRequest) (*v1.SayHelloResponse, error) {
    greeter, err := s.uc.SayHello(ctx, req.GetName())
    if err != nil {
        return nil, err
    }

    return &v1.SayHelloResponse{
        Message: greeter.Message,
    }, nil
}
```

**代码解析**：

| 组件 | 职责 |
|------|------|
| 嵌入 `UnimplementedGreeterServiceServer` | 保持向前兼容，新增 RPC 方法时不会破坏编译 |
| 依赖 `*biz.GreeterUsecase` | 调用领域层的业务逻辑 |
| DTO 转换 | 将领域对象 `Greeter` 转换为 proto 响应 |

### 3.4 传输层 (internal/server/)

```go
package server

import (
    "net/http"

    "github.com/gin-gonic/gin"

    v1 "go-api-template/api/helloworld/v1"
    "go-api-template/internal/service"
)

// NewHTTPServer 创建并配置 HTTP 服务器
func NewHTTPServer(greeterSvc *service.GreeterService) *gin.Engine {
    engine := gin.Default()

    engine.GET("/health", func(c *gin.Context) {
        c.JSON(http.StatusOK, gin.H{"status": "ok"})
    })

    engine.GET("/", func(c *gin.Context) {
        c.JSON(http.StatusOK, gin.H{
            "name":    "go-api-template",
            "version": "0.1.0",
        })
    })

    registerGreeterRoutes(engine, greeterSvc)

    return engine
}

func registerGreeterRoutes(engine *gin.Engine, svc *service.GreeterService) {
    v1Group := engine.Group("/api/v1")
    {
        v1Group.POST("/greeter/say-hello", handleSayHello(svc))
        v1Group.GET("/greeter/say-hello/:name", handleSayHelloByPath(svc))
    }
}
```

### 3.5 入口文件 (cmd/server/main.go)

```go
package main

import (
    "log"

    "go-api-template/internal/biz"
    "go-api-template/internal/data"
    "go-api-template/internal/server"
    "go-api-template/internal/service"
)

func main() {
    // 手动依赖注入（阶段四将使用 Wire 自动化）
    // 依赖组装顺序：Data -> Repo -> Usecase -> Service -> Server

    // 1. 初始化数据层
    dataLayer, err := data.NewData()
    if err != nil {
        log.Fatalf("Failed to create data layer: %v", err)
    }

    // 2. 创建 Repository
    greeterRepo := data.NewGreeterRepo(dataLayer)

    // 3. 创建业务用例
    greeterUsecase := biz.NewGreeterUsecase(greeterRepo)

    // 4. 创建服务
    greeterService := service.NewGreeterService(greeterUsecase)

    // 5. 创建 HTTP 服务器
    httpServer := server.NewHTTPServer(greeterService)

    // 启动服务
    log.Println("Starting server on :8080")
    if err := httpServer.Run(":8080"); err != nil {
        log.Fatalf("Failed to start server: %v", err)
    }
}
```

**依赖组装顺序**：

```
Data → GreeterRepo → GreeterUsecase → GreeterService → HTTPServer
```

---

## 4. 验证

### 4.1 启动服务

```powershell
go run ./cmd/server
```

### 4.2 测试端点

```powershell
# 健康检查
$response = Invoke-RestMethod -Uri "http://localhost:8080/health" -Method Get
$response | ConvertTo-Json
# {"status": "ok"}

# GET 方式问候
$response = Invoke-RestMethod -Uri "http://localhost:8080/api/v1/greeter/say-hello/World" -Method Get
$response | ConvertTo-Json
# {"message": "Hello, World! You are visitor #1."}

# POST 方式问候
$response = Invoke-RestMethod -Uri "http://localhost:8080/api/v1/greeter/say-hello" -Method Post -ContentType "application/json" -Body '{"name":"Alice"}'
$response | ConvertTo-Json
# {"message": "Hello, Alice! You are visitor #2."}
```

---

## 5. 阶段三完成后的目录结构

```
go-api-template/
├── api/
│   ├── buf.yaml
│   └── helloworld/v1/
│       ├── greeter.proto
│       ├── greeter.pb.go
│       └── greeter_grpc.pb.go
├── cmd/
│   └── server/
│       └── main.go                ← 手动依赖注入
├── internal/
│   ├── biz/
│   │   ├── biz.go                 ← 包声明
│   │   └── greeter.go             ← 领域实体 + 业务用例
│   ├── data/
│   │   ├── data.go                ← 数据层初始化
│   │   └── greeter.go             ← Repository 实现（内存）
│   ├── server/
│   │   ├── server.go              ← 包声明
│   │   └── http.go                ← HTTP 路由配置
│   └── service/
│       ├── service.go             ← 包声明
│       └── greeter.go             ← API 实现
├── go.mod
├── go.sum
└── Makefile
```

---

## 6. 依赖倒置的威力

阶段三的核心收获是理解**依赖倒置**带来的好处：

| 好处 | 说明 |
|------|------|
| 可测试性 | 可以用 Mock 实现替换真实 Repository，方便单元测试 |
| 可替换性 | 切换存储（内存 → MySQL）只需修改 data 层 |
| 关注点分离 | 业务逻辑不关心技术细节（如何存储、如何传输） |
| 团队协作 | 不同团队可以并行开发不同层 |

---

## 7. 下一步：阶段四

阶段四将引入 **Google Wire** 自动化依赖注入：

1. 安装 Wire CLI
2. 编写 `wire.go` 声明依赖关系
3. 运行 `wire` 生成 `wire_gen.go`
4. 删除 `main.go` 中的手动组装代码

目标：让依赖注入代码自动生成，减少样板代码。
