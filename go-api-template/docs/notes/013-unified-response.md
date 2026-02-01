# 013. 统一响应格式

## 1. 核心问题

在 012 章节中，我们实现了统一的错误响应格式，但成功响应和错误响应的结构仍然不一致：

**成功响应**（012 章节之前）：

```json
{"message": "Hello, World!"}
```

**错误响应**（012 章节实现）：

```json
{
    "code": "INVALID_PARAMS",
    "message": "请求参数验证失败",
    "http_code": 400,
    "details": [...]
}
```

**问题**：前端需要用不同的逻辑处理成功和失败的响应，代码复杂且容易出错。

---

## 2. 设计目标

让所有 API 响应都有一致的顶层结构：

```json
{
    "code": "状态码",
    "message": "状态消息",
    "http_code": 200,
    "data": { ... },      // 成功时：业务数据
    "details": [ ... ]    // 失败时：错误详情
}
```

前端处理逻辑变得简单：

```javascript
// 前端统一处理
const response = await fetch('/api/v1/greeter/say-hello', { ... });
const result = await response.json();

if (result.code === 'SUCCESS') {
    // 使用 result.data
    console.log(result.data.message);
} else {
    // 显示错误
    console.error(result.message);
    if (result.details) {
        // 显示字段错误
    }
}
```

---

## 3. 架构设计

### 3.1 三层职责分离

为了实现真正的解耦和语义清晰，我们将系统分为三个部分：

| 组件               | 职责          | 位置                          | 说明                                                           |
| ------------------ | ------------- | ----------------------------- | -------------------------------------------------------------- |
| **Reason**   | 业务状态字典  | `internal/pkg/reason/`      | 纯字典，定义业务状态原因（如 `SUCCESS`, `INVALID_PARAMS`） |
| **AppError** | 内部错误传递  | `internal/pkg/apperrors/`   | 控制流，Service 层抛出的异常，Handler 层捕获                   |
| **Response** | HTTP 响应格式 | `internal/server/response/` | 表现层，面向客户端的 JSON 结构                                 |

### 3.2 Response 与 AppError 的关系

**核心原则**：Response **使用** AppError

```
┌─────────────────────────────────────────────────────────────────
│  Service 层                             
│    - 不知道 HTTP 是什么                    
│    - 返回业务数据，或抛出 AppError            
└──────────────────────────┬──────────────────────────────────────
                           │
                           ▼
┌─────────────────────────────────────────────────────────────────
│  Handler 层 (HTTP Controller)                     
│    - 调用 Service，获取结果或捕获 AppError            
│    - 使用 Response 包装输出                            
│      - 成功：response.SuccessJSON(c, data)               
│      - 失败：response.ErrorJSON(c, appErr)                 
└─────────────────────────────────────────────────────────────────
```

**为什么需要两者？**

| 组件               | 存在于                   | 作用                                                |
| ------------------ | ------------------------ | --------------------------------------------------- |
| **AppError** | Service 层 → Handler 层 | Go 代码内部的错误传递，支持错误链、Unwrap、日志记录 |
| **Response** | Handler 层 → 客户端     | HTTP 响应的 JSON 格式，面向前端                     |

如果 Service 层直接返回 `Response`，就违反了**整洁架构**的依赖规则（业务层不应依赖传输层）。

### 3.3 文件结构

```
internal/
├── pkg/
│   ├── reason/
│   │   └── reason.go     # 业务状态原因定义
│   └── apperrors/
│       └── errors.go     # AppError（引用 reason）
└── server/
    └── response/
        └── response.go   # Response（引用 reason 和 apperrors）
```

---

## 4. 代码实现

### 4.1 定义业务原因 (Reason)

创建 `internal/pkg/reason/reason.go`：

```go
package reason

type Reason string

const (
    // 成功
    Success Reason = "SUCCESS"

    // 客户端问题
    InvalidParams Reason = "INVALID_PARAMS"
    Unauthorized  Reason = "UNAUTHORIZED"
    NotFound      Reason = "NOT_FOUND"

    // 服务端问题
    InternalError      Reason = "INTERNAL_ERROR"
    ServiceUnavailable Reason = "SERVICE_UNAVAILABLE"
)

func (r Reason) HTTPStatus() int { ... }
```

### 4.2 AppError 引用 Reason

修改 `internal/pkg/apperrors/errors.go`：

```go
package apperrors

import "go-api-template/internal/pkg/reason"

type AppError struct {
    Code     reason.Reason `json:"code"`
    Message  string        `json:"message"`
    HTTPCode int           `json:"http_code"`
    Details  []FieldError  `json:"details,omitempty"`
    Cause    error         `json:"-"` // 不暴露给客户端
}

// 快捷构造方法
func NotFound(message string) *AppError {
    return New(reason.NotFound, message)
}
```

### 4.3 Response 引用 Reason 和 AppError

修改 `internal/server/response/response.go`：

```go
package response

import (
    "go-api-template/internal/pkg/apperrors"
    "go-api-template/internal/pkg/reason"
)

// Response 统一响应结构
type Response struct {
    Code     reason.Reason          `json:"code"`
    Message  string                 `json:"message"`
    HTTPCode int                    `json:"http_code"`
    Data     any                    `json:"data,omitempty"`
    Details  []apperrors.FieldError `json:"details,omitempty"`
}

// Success 创建成功响应
func Success(data any) *Response {
    return &Response{
        Code:     reason.Success,
        Message:  "操作成功",
        HTTPCode: 200,
        Data:     data,
    }
}

// Error 从 AppError 创建错误响应
// 这是 Response 使用 AppError 的核心方法
func Error(err *apperrors.AppError) *Response {
    return &Response{
        Code:     err.Code,
        Message:  err.Message,
        HTTPCode: err.HTTPCode,
        Details:  err.Details,
    }
}
```

---

## 5. 最佳实践：为什么不使用 gin.H？

### 5.1 问题：gin.H 是什么？

`gin.H` 是 Gin 框架提供的 `map[string]any` 别名，用于快速构造 JSON：

```go
// 新手写法（错误）
c.JSON(200, gin.H{"message": resp.GetMessage()})
```

### 5.2 为什么不应该使用 gin.H？

| 问题                        | 说明                                                                             |
| --------------------------- | -------------------------------------------------------------------------------- |
| **丢失类型安全**      | Map 是弱类型，编译器无法检查字段名拼写错误                                       |
| **重复劳动**          | Proto 已经定义了结构体并带有 `json` 标签，手动拆解再塞进 Map 是多此一举        |
| **违反 Schema First** | 我们的架构核心是"定义优先"，数据结构由 Proto 统一管理，不应在 Handler 层临时构造 |
| **耦合 Gin 框架**     | 在 Handler 层引入 `gin.H` 增加了对框架的依赖                                   |

### 5.3 正确做法：直接传递结构体

在 **Schema First** 架构中，数据流向是严谨的：

```
Proto 定义 → 生成 Go 结构体（带 json 标签）→ Service 返回 → Handler 透传
```

**最佳实践**：

```go
// 正确：直接传递 Service 返回的结构体
resp, err := svc.SayHello(ctx, req)
if err != nil {
    response.ErrorJSON(c, apperrors.Internal("服务处理失败", err))
    return
}
response.SuccessJSON(c, resp)  // resp 是 *v1.SayHelloResponse
```

### 5.4 深入理解：JSON 序列化机制

`resp` 是 Proto 生成的结构体，直接传给 `SuccessJSON` 后，序列化是怎么发生的？

#### 调用链追踪

```go
// 1. Handler 调用
response.SuccessJSON(c, resp)  // resp 是 *v1.SayHelloResponse

// 2. SuccessJSON 内部
func SuccessJSON(c *gin.Context, data any) {
    JSON(c, Success(data))  // 把 resp 传给 Success()
}

// 3. Success 构造 Response
func Success(data any) *Response {
    return &Response{
        Code:     reason.Success,      // "SUCCESS"
        Message:  "操作成功",
        HTTPCode: 200,
        Data:     data,                // resp 被赋值给 Data 字段
    }
}

// 4. JSON 输出
func JSON(c *gin.Context, r *Response) {
    c.JSON(r.HTTPCode, r)  // 这里触发序列化！
}
```

#### 序列化发生的时机和方式

`c.JSON()` 内部调用 Go 标准库的 `encoding/json.Marshal(r)`，对**整个 Response 结构体**进行一次性递归序列化：

```go
// 要序列化的 Response
Response{
    Code:     "SUCCESS",           // string，直接序列化
    Message:  "操作成功",           // string，直接序列化
    HTTPCode: 200,                 // int，直接序列化
    Data:     resp,                // any 类型，值是 *v1.SayHelloResponse
}
```

当 `json.Marshal` 遇到 `Data` 字段时：

1. 看到 `Data` 的类型是 `any`
2. 检查实际值的底层类型是 `*v1.SayHelloResponse`
3. 递归序列化该结构体，使用它的 `json` 标签

Proto 生成的结构体同时有**两种标签**：

```go
// 由 protoc-gen-go 生成
type SayHelloResponse struct {
    Message string `protobuf:"bytes,1,opt,name=message,proto3" json:"message,omitempty"`
}
```

| 标签                         | 用途                               |
| ---------------------------- | ---------------------------------- |
| `protobuf:"..."`           | Protobuf 二进制序列化（用于 gRPC） |
| `json:"message,omitempty"` | JSON 序列化（用于 HTTP）           |

#### 最终输出

```json
{
    "code": "SUCCESS",           // ← Response.Code
    "message": "操作成功",        // ← Response.Message
    "http_code": 200,            // ← Response.HTTPCode
    "data": {                    // ← Response.Data 递归序列化
        "message": "Hello..."    // ← SayHelloResponse.Message
    }
}
```

#### 常见疑问

| 问题                       | 答案                                         |
| -------------------------- | -------------------------------------------- |
| 序列化在哪里发生？         | `c.JSON()` 内部，由 `encoding/json` 完成 |
| 什么时候发生？             | 调用 `c.JSON(code, obj)` 的那一刻          |
| 只对 data 序列化吗？       | 不是，对整个 Response 递归序列化             |
| 用的是 Protobuf 序列化吗？ | 不是，用的是 JSON 序列化（`json` 标签）    |

### 5.5 边缘情况：response.Body

极少数情况下（如快速返回临时数据），如果确实需要构造 Map，使用 `response.Body` 而非 `gin.H`：

```go
// 可接受：使用 response.Body（减少对 Gin 的直接依赖）
response.SuccessJSON(c, response.Body{"id": 123})
```

`response.Body` 定义为 `type Body map[string]any`，功能与 `gin.H` 相同，但避免了在业务代码中引入框架类型。

---

## 6. 使用方式

### 6.1 完整 Handler 示例

```go
func handleSayHello(svc *service.GreeterService) gin.HandlerFunc {
    return func(c *gin.Context) {
        var req dto.SayHelloRequest

        // 验证失败 → 统一错误响应
        if err := c.ShouldBindJSON(&req); err != nil {
            response.ErrorJSON(c, apperrors.FromValidationError(err))
            return
        }

        // 调用 Service
        resp, err := svc.SayHello(c.Request.Context(), req.ToProto())
        if err != nil {
            // 内部错误 → 统一错误响应
            response.ErrorJSON(c, apperrors.Internal("服务处理失败", err))
            return
        }

        // 成功 → 直接传递结构体（最佳实践）
        response.SuccessJSON(c, resp)
    }
}
```

### 6.2 数据流向图

```
┌─────────────────┐
│  客户端请求      │
└────────┬────────┘
         │
         ▼
┌─────────────────┐
│  Handler        │  c.ShouldBindJSON(&req)
│  (验证请求)      │
└────────┬────────┘
         │ 验证失败？
         │ ├─ 是 → apperrors.FromValidationError(err)
         │ │       → response.ErrorJSON(c, appErr)
         │ │       → 返回 {"code":"INVALID_PARAMS",...}
         ▼ 否
┌─────────────────┐
│  Service        │  svc.SayHello(ctx, req)
│  (业务逻辑)      │
└────────┬────────┘
         │ 出错？
         │ ├─ 是 → apperrors.Internal("...", err)
         │ │       → response.ErrorJSON(c, appErr)
         │ │       → 返回 {"code":"INTERNAL_ERROR",...}
         ▼ 否
┌─────────────────┐
│  Response       │  response.SuccessJSON(c, resp)
│  (统一输出)      │  → 返回 {"code":"SUCCESS","data":{...}}
└─────────────────┘
```

---

## 7. 验证测试

启动服务后测试：

### 7.1 测试成功响应

```bash
# PowerShell
Invoke-RestMethod -Uri http://localhost:8080/api/v1/greeter/say-hello `
    -Method Post -ContentType "application/json" `
    -Body '{"name": "World"}'
```

**返回**：

```json
{
    "code": "SUCCESS",
    "message": "操作成功",
    "http_code": 200,
    "data": {
        "message": "Hello, World! You are visitor #1."
    }
}
```

### 7.2 测试错误响应

```bash
# PowerShell
Invoke-RestMethod -Uri http://localhost:8080/api/v1/greeter/say-hello `
    -Method Post -ContentType "application/json" `
    -Body '{}'
```

**返回**：

```json
{
    "code": "INVALID_PARAMS",
    "message": "请求参数验证失败",
    "http_code": 400,
    "details": [
        {"field": "name", "message": "必填字段"}
    ]
}
```

---

## 8. 架构思考：为什么需要 Reason 包？

你可能会问：*为什么不直接把 SUCCESS 放在 AppError 里？*

**原因：**

1. **语义正确性**：AppError 代表"异常流"，Success 代表"正常流"。将 Success 定义在 Error 包中是语义上的混淆。
2. **职责单一**：AppError 负责错误堆栈和处理，Reason 负责业务状态定义。
3. **可读性**：`reason.Success` 比 `apperrors.ErrCodeSuccess` 更自然，因为它不是一个 Error。

这种设计遵循了**领域驱动设计 (DDD)** 中"通用语言"的概念——状态码是业务领域的通用语言，不应被限制在"错误处理"的技术实现细节中。

---

## 9. 总结

本章节实现了统一的 HTTP 响应格式，核心要点：

| 要点                   | 说明                                                           |
| ---------------------- | -------------------------------------------------------------- |
| **Reason 包**    | 定义业务状态字典（Success, InvalidParams 等），中立于成功/失败 |
| **AppError**     | Service 层的错误传递机制，Handler 层捕获后转换为 Response      |
| **Response**     | 统一的 HTTP 输出格式，**使用** AppError 而非替代它       |
| **不使用 gin.H** | 直接传递结构体，保持类型安全和 Schema First 原则               |
