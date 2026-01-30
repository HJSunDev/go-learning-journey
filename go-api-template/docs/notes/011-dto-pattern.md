# 011. DTO 模式与请求验证实践

## 1. 先看当前代码的问题

打开 `internal/server/http.go`，看看现在是怎么处理请求的：

```go
func handleSayHello(svc *service.GreeterService) gin.HandlerFunc {
    return func(c *gin.Context) {
        // 直接使用 Proto 生成的类型
        var req v1.SayHelloRequest

        if err := c.ShouldBindJSON(&req); err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": "name is required"})
            return
        }

        // 手动验证！因为 Proto 类型没有 binding tag
        if req.GetName() == "" {
            c.JSON(http.StatusBadRequest, gin.H{"error": "name is required"})
            return
        }
  
        // ...
    }
}
```

**问题在哪？**

看看 Proto 生成的 `SayHelloRequest` 长什么样（在 `api/helloworld/v1/greeter.pb.go`）：

```go
type SayHelloRequest struct {
    Name string `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
    // ... 其他 proto 内部字段
}
```

注意看 Tag：

- 有 `protobuf:"..."` - Proto 序列化用的
- 有 `json:"name,omitempty"` - JSON 序列化用的
- **没有 `binding:"required"`** - 所以 Validator 不会验证！

这就是为什么我们必须手动写 `if req.GetName() == ""`。

**能不能给 Proto 生成的代码加上 binding tag？**

不能。因为：

1. `.proto` 文件是手写的，但 `.pb.go` 文件是 `buf generate` 自动生成的
2. 每次运行 `buf generate` 都会重新生成 `.pb.go`，覆盖掉你手动加的任何修改

---

## 2. DTO 是什么

### 2.1 定义

**DTO = Data Transfer Object = 数据传输对象**

DTO 是一种**设计模式**，不是某个具体的库或框架。它的核心思想是：**定义专门用于跨层或跨系统传输数据的对象，与业务逻辑解耦**。

| 特性     | 说明                                                            |
| -------- | --------------------------------------------------------------- |
| 本质     | 一种设计模式，源自 Java 企业级开发，现广泛应用于各种语言        |
| 实现方式 | 在 Go 中就是普通的 struct，不需要引入任何依赖                   |
| 命名约定 | 通常以 `Request`、`Response` 结尾，如 `CreateUserRequest` |
| 存放位置 | 本项目放在 `internal/server/dto/` 目录                        |

在 Web 开发中，DTO 通常指：

- 接收 HTTP 请求的结构体（Request DTO）
- 返回 HTTP 响应的结构体（Response DTO）

### 2.2 我们有了proto以后为什么还需要 DTO

| 问题                        | DTO 怎么解决             |
| --------------------------- | ------------------------ |
| Proto 类型没有 binding tag  | DTO 可以自由添加任何 tag |
| HTTP 字段名可能和业务层不同 | DTO 可以做字段映射       |
| HTTP 请求需要验证           | DTO 定义验证规则         |
| 不想暴露内部结构给 HTTP     | DTO 只暴露该暴露的字段   |

---

## 3. DTO 和 Proto 的关系

### 3.1 它们分别是什么

| 类型                 | 定义位置                              | 用途                 | 谁使用           |
| -------------------- | ------------------------------------- | -------------------- | ---------------- |
| **Proto 类型** | `.proto` 文件定义，自动生成 Go 代码 | Service 层的输入输出 | Service 层、gRPC |
| **DTO**        | 手动编写的 Go struct                  | HTTP 层的输入输出    | HTTP Handler     |

### 3.2 它们怎么协作

```
┌─────────────────────────────────────────────────────────────
│  HTTP 请求   
│  {"name": "World"}  
└─────────────────────────────────────────────────────────────
                              │
                              ▼
┌─────────────────────────────────────────────────────────────
│  HTTP Handler (internal/server/http.go)   
│                   
│  1. 用 DTO 接收请求:   
│     var req dto.SayHelloRequest   
│     c.ShouldBindJSON(&req)  // 自动验证！  
│                             
│  2. DTO 转 Proto:         
│     protoReq := req.ToProto()   
│                                 
│  3. 调用 Service:                    
│     resp, err := svc.SayHello(ctx, protoReq)  
└─────────────────────────────────────────────────────────────
                              │
                              ▼
┌─────────────────────────────────────────────────────────────
│  Service 层 (internal/service/greeter.go)   
│                                             
│  func SayHello(ctx, req *v1.SayHelloRequest)  
│      → 使用 Proto 类型                           
│      → 调用业务层                                  
│      → 返回 Proto 类型                               
└─────────────────────────────────────────────────────────────
```

**关键点**：

- HTTP Handler 用 **DTO** 接收请求（因为需要验证）
- Service 层用 **Proto 类型**（因为这是 API 契约）
- DTO 提供 `ToProto()` 方法做转换

### 3.3 为什么 Service 层用 Proto 而不是 DTO？

1. **Schema First 原则**：Proto 文件是"契约"，Service 层遵守契约
2. **gRPC 兼容**：如果以后加 gRPC，Service 层不用改
3. **多入口复用**：HTTP 和 gRPC 都可以调用同一个 Service

```
HTTP 请求  ──》  DTO  ──》  Proto  ──》 Service
                             ↑
gRPC 请求  ──────────────────┘         (直接用 Proto)
```

---

## 4. 开发流程

### 4.1 先写 Proto 还是先写 DTO？

**推荐顺序：Proto → DTO**

| 步骤 | 做什么            | 文件                           |
| ---- | ----------------- | ------------------------------ |
| 1    | 写 Proto 定义接口 | `api/xxx/v1/xxx.proto`       |
| 2    | 生成代码          | `buf generate`               |
| 3    | 写 Service 层     | `internal/service/xxx.go`    |
| 4    | 写 DTO            | `internal/server/dto/xxx.go` |
| 5    | 写 HTTP Handler   | `internal/server/http.go`    |

**为什么这个顺序？**

- Proto 定义了"接口有哪些字段"
- DTO 只是给这些字段"加上验证规则"
- 所以先有 Proto，才知道 DTO 要有哪些字段

### 4.2 现阶段不用 gRPC，Proto 有用吗？

**有用，但主要用的是 Proto 生成的类型，不是 gRPC 功能。**

Proto 生成了两个文件：

- `greeter.pb.go` - 定义了 `SayHelloRequest`、`SayHelloResponse` 类型 ✅ **在用**
- `greeter_grpc.pb.go` - gRPC 服务代码 ❌ **暂时没用**

当前项目的 Service 层使用 Proto 生成的类型作为参数和返回值：

```go
// internal/service/greeter.go
func (s *GreeterService) SayHello(ctx context.Context, req *v1.SayHelloRequest) (*v1.SayHelloResponse, error) {
    // ...
}
```

这样设计的好处是：**以后想加 gRPC，Service 层一行都不用改**。

### 4.3 能不能完全不用 Proto，只用 DTO？

**可以，但会失去一些好处。**

如果只用 DTO：

```go
// 只用 DTO 的写法
func (s *GreeterService) SayHello(ctx context.Context, name string) (string, error) {
    // ...
}
```

| 对比项      | 只用 DTO          | Proto + DTO               |
| ----------- | ----------------- | ------------------------- |
| 复杂度      | 简单              | 稍复杂                    |
| 以后加 gRPC | 要重写 Service 层 | Service 层不用改          |
| 类型定义    | 分散在各处        | 集中在 Proto 文件         |
| 适合场景    | 小项目、学习      | 中大型项目、可能扩展 gRPC |

**本项目保留 Proto 的原因**：了解实践完整的架构模式。

---

## 5. 代码实现

### 5.1 创建 DTO 文件

创建 `internal/server/dto/greeter.go`：

```go
package dto

import (
    v1 "go-api-template/api/helloworld/v1"
)

// SayHelloRequest 是 POST /api/v1/greeter/say-hello 的请求体
type SayHelloRequest struct {
    // 添加 binding tag，Validator 会自动验证
    Name string `json:"name" binding:"required,min=1,max=100"`
}

// ToProto 将 DTO 转换为 Proto 类型
// HTTP Handler 用 DTO 接收请求后，调用这个方法转换，再传给 Service
func (r *SayHelloRequest) ToProto() *v1.SayHelloRequest {
    return &v1.SayHelloRequest{
        Name: r.Name,
    }
}
```

### 5.2 修改 HTTP Handler

修改 `internal/server/http.go`：

```go
import (
    // 新增
    "go-api-template/internal/server/dto"
)

func handleSayHello(svc *service.GreeterService) gin.HandlerFunc {
    return func(c *gin.Context) {
        // 改用 DTO 接收请求
        var req dto.SayHelloRequest

        // ShouldBindJSON 会自动验证 binding tag 的规则
        // 如果 name 为空，这里就会返回错误，不需要手动检查了！
        if err := c.ShouldBindJSON(&req); err != nil {
            c.JSON(http.StatusBadRequest, gin.H{
                "error": err.Error(),
            })
            return
        }

        // DTO 转 Proto，调用 Service
        resp, err := svc.SayHello(c.Request.Context(), req.ToProto())
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{
                "error": err.Error(),
            })
            return
        }

        c.JSON(http.StatusOK, gin.H{
            "message": resp.GetMessage(),
        })
    }
}
```

**对比改动前后**：

| 改动前                              | 改动后                               |
| ----------------------------------- | ------------------------------------ |
| `var req v1.SayHelloRequest`      | `var req dto.SayHelloRequest`      |
| 手动检查 `if req.GetName() == ""` | 删除，Validator 自动检查             |
| `svc.SayHello(ctx, &req)`         | `svc.SayHello(ctx, req.ToProto())` |

---

## 6. 验证测试

启动服务后测试：

```bash
# 测试 1: 不传 name（验证 required）
curl -X POST http://localhost:8080/api/v1/greeter/say-hello \
  -H "Content-Type: application/json" \
  -d '{}'

# 预期返回 400 错误，包含验证失败信息

# 测试 2: name 太长（验证 max=100）
curl -X POST http://localhost:8080/api/v1/greeter/say-hello \
  -H "Content-Type: application/json" \
  -d '{"name": "这是一个超级长的名字...（超过100个字符）..."}'

# 预期返回 400 错误

# 测试 3: 正常请求
curl -X POST http://localhost:8080/api/v1/greeter/say-hello \
  -H "Content-Type: application/json" \
  -d '{"name": "World"}'

# 预期返回 200，包含问候消息
```

---

## 7. 总结

### 核心概念

| 概念       | 说明                                                |
| ---------- | --------------------------------------------------- |
| DTO        | 专门用于 HTTP 请求/响应的结构体，可以加 binding tag |
| Proto 类型 | 自动生成的类型，用于 Service 层，不能加自定义 tag   |
| ToProto()  | DTO 提供的转换方法，把 DTO 转成 Proto 类型          |

### 为什么每个 DTO 都要写 ToProto()？

你可能会问：能不能写一个通用的 `ToProto()`，所有 DTO 都能用？

**答案是：Go 语言做不到。** Go 在编译时必须知道具体类型，无法写出"自动把任意 DTO 转成对应 Proto"的通用函数。

**结论：每个 DTO 都要写自己的 `ToProto()` 方法，这是 Go 社区的标准做法。**

```go
// 每个 DTO 都这样写，没有捷径
func (r *CreateUserRequest) ToProto() *v1.CreateUserRequest {
    return &v1.CreateUserRequest{
        Name:  r.Name,
        Email: r.Email,
        Age:   r.Age,
    }
}
```

**字段特别多怎么办？**

如果一个 DTO 有几十个字段，手写确实烦。这时可以用 `copier` 库简化：

```go
import "github.com/jinzhu/copier"

func (r *CreateUserRequest) ToProto() *v1.CreateUserRequest {
    var proto v1.CreateUserRequest
    copier.Copy(&proto, r)  // 自动复制所有同名字段
    return &proto
}
```

但注意：**还是要写 `ToProto()` 方法**，只是方法内部用 `copier` 代替手动赋值。方法本身省不掉。

### 数据流向

```
HTTP JSON → DTO (验证) → Proto 类型 → Service → Proto 类型 → HTTP JSON
```

### 文件职责

| 文件                           | 职责                                           |
| ------------------------------ | ---------------------------------------------- |
| `api/xxx/v1/xxx.proto`       | 定义接口契约                                   |
| `api/xxx/v1/xxx.pb.go`       | 自动生成的类型（Service 层用）                 |
| `internal/server/dto/xxx.go` | HTTP 请求的 DTO（Handler 用）                  |
| `internal/server/http.go`    | HTTP Handler，用 DTO 接收，转 Proto 调 Service |
| `internal/service/xxx.go`    | Service 层，用 Proto 类型                      |

### DTO 命名建议

| HTTP 方法         | DTO 命名                         |
| ----------------- | -------------------------------- |
| POST /users       | `CreateUserRequest`            |
| PUT /users/:id    | `UpdateUserRequest`            |
| GET /users (列表) | `ListUsersRequest`             |
| GET /users/:id    | 通常不需要 DTO，直接从 URL 取 id |
| DELETE /users/:id | 通常不需要 DTO                   |
