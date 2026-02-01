# 014. 中间件

## 1. 为什么需要中间件

### 1.1 回顾最初的目标

在开始"统一错误处理、请求验证、响应格式"这个大目标时，我们将其拆分为多个章节：

| 章节          | 内容                    | 状态             |
| ------------- | ----------------------- | ---------------- |
| 010           | Struct Tag 与 Validator | ✅ 完成          |
| 011           | DTO 模式与请求验证实践  | ✅ 完成          |
| 012           | 统一错误处理            | ✅ 完成          |
| 013           | 统一响应格式            | ✅ 完成          |
| **014** | **中间件**        | **本章节** |

### 1.2 修改前的现状

在 012 和 013 章节完成后，我们的错误处理和响应格式**在 Handler 内部**已经统一：

```go
// Handler 内部 - 已统一
func handleSayHello(svc *service.GreeterService) gin.HandlerFunc {
    return func(c *gin.Context) {
        var req dto.SayHelloRequest
        if err := c.ShouldBindJSON(&req); err != nil {
            // ✅ 验证错误 → 统一格式
            response.ErrorJSON(c, apperrors.FromValidationError(err))
            return
        }
        resp, err := svc.SayHello(c.Request.Context(), req.ToProto())
        if err != nil {
            // ✅ 业务错误 → 统一格式
            response.ErrorJSON(c, apperrors.Internal("服务处理失败", err))
            return
        }
        // ✅ 成功响应 → 统一格式
        response.SuccessJSON(c, resp)
    }
}
```

但是，**Handler 外部的错误没有被统一处理**。

### 1.3 已解决与未覆盖场景对比

| 场景                       | 发生位置             | 012/013 是否覆盖 | 当前返回格式                         |
| -------------------------- | -------------------- | ---------------- | ------------------------------------ |
| 请求验证失败               | Handler 内           | ✅ 已覆盖        | 统一 JSON                            |
| 业务逻辑错误               | Handler 内           | ✅ 已覆盖        | 统一 JSON                            |
| 成功响应                   | Handler 内           | ✅ 已覆盖        | 统一 JSON                            |
| **路由不存在 (404)** | **Handler 外** | ❌ 未覆盖        | `404 page not found`（纯文本）     |
| **方法不允许 (405)** | **Handler 外** | ❌ 未覆盖        | `405 method not allowed`（纯文本） |
| **Handler panic**    | **Handler 外** | ❌ 未覆盖        | 空响应或 HTML                        |

### 1.4 问题演示

**问题 1：访问不存在的路由**

```bash
curl http://localhost:8080/api/v1/not-exist
```

修改前返回：

```
404 page not found
```

这是纯文本，不是 JSON。前端需要写额外的逻辑来处理这种情况。

**问题 2：Handler 内发生 panic**

假设某个 Handler 有 bug：

```go
func handleBuggy() gin.HandlerFunc {
    return func(c *gin.Context) {
        var data []int
        fmt.Println(data[10])  // 数组越界，触发 panic
    }
}
```

修改前，Gin 内置的 `Recovery` 会捕获 panic，但返回的格式是 HTML 或纯文本，不是我们定义的统一 JSON 格式。

### 1.5 核心结论

**中间件是"统一错误处理"和"统一响应格式"在 Handler 外的补完。**

```
┌─────────────────────────────────────────────────────────
│                      请求到达   
└───────────────────────────┬─────────────────────────────
                            │
                            ▼
┌─────────────────────────────────────────────────────────
│  第一道防线：中间件层（本章节实现）  
│  - 404 Not Found → 统一 JSON 格式   
│  - 405 Method Not Allowed → 统一 JSON 格式  
│  - Panic Recovery → 统一 JSON 格式  
│  - 请求 ID 注入（链路追踪）      
└───────────────────────────┬─────────────────────────────
                            │
                            ▼
┌─────────────────────────────────────────────────────────
│  第二道防线：Handler 层（012/013 已完成）  
│  - 验证错误 → apperrors.FromValidationError()   
│  - 业务错误 → apperrors.NotFound() / Internal()   
│  - 成功响应 → response.SuccessJSON()        
└─────────────────────────────────────────────────────────
```

---

## 2. 中间件基础概念

### 2.1 什么是中间件

中间件（Middleware）是在请求到达 Handler 之前和响应返回客户端之前执行的代码。它用于处理**横切关注点**（Cross-cutting Concerns）——那些与具体业务无关但每个请求都需要的功能。

常见的中间件功能：

| 功能     | 说明                             |
| -------- | -------------------------------- |
| 日志记录 | 记录每个请求的 URL、耗时、状态码 |
| 异常恢复 | 捕获 panic，防止服务崩溃         |
| 认证授权 | 验证 JWT Token，检查权限         |
| 请求追踪 | 生成请求 ID，用于日志关联        |
| 限流     | 防止接口被刷                     |
| CORS     | 处理跨域请求                     |

### 2.2 Gin 中间件的执行模型：洋葱模型

Gin 的中间件采用"洋葱模型"，请求像穿过洋葱一样，先从外层进入，到达 Handler，再从内层返回：

```
请求进入
    │
    ▼
┌──────────────────────────────────────────────────────
│  中间件 A（前半部分）   
│      ↓   
│  ┌──────────────────────────────────────────────
│  │  中间件 B（前半部分）   
│  │      ↓   
│  │  ┌──────────────────────────────────────┐   
|  |  |  Handler                             |
│  │  └──────────────────────────────────────┘   
│  │      ↓  
│  │  中间件 B（后半部分）   
│  └────────────────────────────────────────────── 
│      ↓   
│  中间件 A（后半部分）   
└──────────────────────────────────────────────────────
    │
    ▼
响应返回
```

### 2.3 中间件的代码结构

```go
func MyMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        // ======== 前半部分：请求到达时执行 ========
        fmt.Println("请求进入")
  
        // 调用下一个中间件或 Handler
        c.Next()
  
        // ======== 后半部分：响应返回时执行 ========
        fmt.Println("响应返回")
    }
}
```

| 方法                  | 作用                                               |
| --------------------- | -------------------------------------------------- |
| `c.Next()`          | 继续执行后续的中间件和 Handler                     |
| `c.Abort()`         | 终止后续执行，直接返回响应                         |
| `c.Set(key, value)` | 在 Context 中存储数据，供后续中间件和 Handler 使用 |
| `c.Get(key)`        | 从 Context 中获取数据                              |

### 2.4 中间件的注册与执行：两个不同的阶段

看到这段代码：

```go
engine := gin.New()
engine.Use(RequestID())   // 注册中间件
engine.Use(Recovery())
```

疑问：*`engine.Use()` 只执行了一次，那每个请求都会有不同的 RequestID 吗？*

答案是：会的。每个 HTTP 请求都会执行所有注册的中间件。

#### 注册阶段 vs 执行阶段

| 阶段               | 时机           | 发生次数                     | 做了什么                                 |
| ------------------ | -------------- | ---------------------------- | ---------------------------------------- |
| **注册阶段** | 服务启动时     | **1 次**               | 把中间件函数添加到 Engine 的中间件列表中 |
| **执行阶段** | 每个请求到达时 | **N 次**（N = 请求数） | 按顺序执行列表中的每个中间件             |

#### 原理详解

```go
// ====== 注册阶段（服务启动时，只执行一次）======
engine := gin.New()

// engine.Use() 做的事情：
// 把 RequestID() 返回的函数存入 engine.handlers 列表
// 此时 RequestID() 被调用一次，返回一个 gin.HandlerFunc
engine.Use(RequestID())

// ====== 执行阶段（每个请求到达时，都会执行）======
// 当 HTTP 请求到达时，Gin 的处理流程是：
// 1. 创建一个新的 gin.Context（每个请求独立）
// 2. 把 engine.handlers 列表中的所有函数依次执行
// 3. 每个函数接收的是这个请求独有的 *gin.Context

// 伪代码：
for _, handler := range engine.handlers {
    handler(ctx)  // ctx 是当前请求的 Context
}
```

#### 关键理解

```go
func RequestID() gin.HandlerFunc {
    // 这一层在【注册阶段】执行，只执行一次
    // 可以在这里做一些初始化工作
  
    return func(c *gin.Context) {
        // 这一层在【执行阶段】执行，每个请求都会执行
        // c 是当前请求独有的 Context
        requestID := uuid.New().String()  // 每次都生成新的 UUID
        c.Set("request_id", requestID)    // 存入当前请求的 Context
        c.Next()
    }
}
```

| 代码位置                                   | 执行时机       | 执行次数 |
| ------------------------------------------ | -------------- | -------- |
| `func RequestID()` 外层函数体            | 服务启动时     | 1 次     |
| `return func(c *gin.Context)` 内层函数体 | 每个请求到达时 | N 次     |

**类比**：`engine.Use()` 就像把一个"印章"放进工具箱（注册），而每个请求到达时都会从工具箱取出印章盖一次（执行）。

### 2.5 gin.Default() vs gin.New()

```go
// gin.Default() = gin.New() + Logger + Recovery
engine := gin.Default()

// gin.New() = 空白引擎，无任何中间件
engine := gin.New()
```

| 方法              | 内置中间件        | 适用场景                       |
| ----------------- | ----------------- | ------------------------------ |
| `gin.Default()` | Logger + Recovery | 快速开发、原型验证             |
| `gin.New()`     | 无                | 生产环境，需要自定义中间件行为 |

**为什么我们选择 `gin.New()`？**

因为 `gin.Default()` 内置的 `Recovery` 返回的是 HTML/纯文本格式，不是我们定义的统一 JSON 格式。我们需要用自己的 `Recovery` 替换它。

---

## 3. 设计方案

### 3.1 文件结构

```
internal/server/
├── http.go                    # 重构：使用 gin.New() + 手动注册中间件
└── middleware/
    ├── middleware.go          # 中间件聚合器，统一注册入口
    ├── request_id.go          # 请求 ID 中间件
    ├── recovery.go            # Panic 恢复中间件
    └── routes.go              # 404/405 处理器
```

### 3.2 中间件注册顺序

```go
// middleware.go
var middlewareChain = []gin.HandlerFunc{
    RequestID(),   // [0] 最先执行
    Recovery(),    // [1] 捕获后续所有 panic
    gin.Logger(),  // [2] 记录请求日志
}

func Register(engine *gin.Engine) {
    engine.Use(middlewareChain...)
}
```

**顺序说明**：

| 索引 | 中间件    | 原因                                                              |
| ---- | --------- | ----------------------------------------------------------------- |
| [0]  | RequestID | 最先执行，确保后续所有代码（包括 Recovery 的日志）都能获取请求 ID |
| [1]  | Recovery  | 第二执行，捕获后续所有代码的 panic（包括 Logger 和 Handler）      |
| [2]  | Logger    | 第三执行，记录请求信息                                            |

### 3.3 与已有代码的集成

中间件复用 012/013 章节实现的 `apperrors` 和 `response` 包：

```go
// recovery.go
func Recovery() gin.HandlerFunc {
    return func(c *gin.Context) {
        defer func() {
            if err := recover(); err != nil {
                // 复用 apperrors 和 response
                response.ErrorJSON(c, apperrors.Internal("服务器内部错误", nil))
                c.Abort()
            }
        }()
        c.Next()
    }
}
```

这样，中间件层的错误输出与 Handler 层**完全一致**。

---

## 4. 代码实现

### 4.1 RequestID 中间件

```go
// internal/server/middleware/request_id.go
package middleware

import (
    "github.com/gin-gonic/gin"
    "github.com/google/uuid"
)

const (
    HeaderXRequestID    = "X-Request-ID"
    ContextKeyRequestID = "request_id"
)

// RequestID 返回请求 ID 中间件
// 职责：
//   - 从请求头提取 X-Request-ID，如果没有则生成 UUID
//   - 将请求 ID 存入 gin.Context，供后续处理使用
//   - 将请求 ID 写入响应头，方便客户端关联请求和响应
func RequestID() gin.HandlerFunc {
    return func(c *gin.Context) {
        requestID := c.GetHeader(HeaderXRequestID)
        if requestID == "" {
            requestID = uuid.New().String()
        }
        c.Set(ContextKeyRequestID, requestID)
        c.Header(HeaderXRequestID, requestID)
        c.Next()
    }
}

// GetRequestID 从 gin.Context 中获取请求 ID
func GetRequestID(c *gin.Context) string {
    if id, exists := c.Get(ContextKeyRequestID); exists {
        if requestID, ok := id.(string); ok {
            return requestID
        }
    }
    return ""
}
```

**为什么需要 RequestID？**

| 场景       | 没有 RequestID               | 有 RequestID             |
| ---------- | ---------------------------- | ------------------------ |
| 日志排查   | 并发请求的日志交错，无法区分 | 按 ID 过滤，完整链路清晰 |
| 用户反馈   | "我刚才报错了"，无法定位     | 用户提供 ID，直接定位    |
| 分布式追踪 | 需要额外机制                 | 可作为 TraceID 的基础    |

#### 两个常量的用途

```go
const (
    HeaderXRequestID    = "X-Request-ID"     // HTTP 头名称
    ContextKeyRequestID = "request_id"       // Context 键名
)
```

| 常量                    | 用在哪里         | 用途                                                           |
| ----------------------- | ---------------- | -------------------------------------------------------------- |
| `ContextKeyRequestID` | 服务端内部       | 存入 `gin.Context`，供 Handler、Service 层获取，用于日志记录 |
| `HeaderXRequestID`    | HTTP 请求/响应头 | 客户端与服务端之间传递，用于问题排查                           |

#### X-Request-ID 响应头：客户端什么时候需要？

**场景 1：用户反馈问题**

用户说"我刚才提交订单失败了"，但没有更多信息。你如何在海量日志中找到那一条？

如果前端在请求失败时展示 RequestID 给用户：

```javascript
// 前端代码示例
try {
    const response = await fetch('/api/v1/orders', { method: 'POST', ... });
    const data = await response.json();
    if (data.code !== 'SUCCESS') {
        // 从响应头获取 RequestID
        const requestId = response.headers.get('X-Request-ID');
        alert(`操作失败，请联系客服，错误编号：${requestId}`);
    }
} catch (error) {
    // ...
}
```

用户反馈时可以说："错误编号是 `550e8400-e29b-41d4-a716-446655440000`"。后端直接用这个 ID 搜索日志：

```bash
grep "550e8400-e29b-41d4-a716-446655440000" /var/log/app.log
```

**场景 2：微服务链路追踪**

在微服务架构中，一个用户请求可能跨越多个服务：

```
用户 → API Gateway → 订单服务 → 库存服务 → 支付服务
```

如果每个服务都把收到的 `X-Request-ID` 传递给下游，整条链路就能串起来：

```bash
# 在任何一个服务的日志中，用同一个 ID 都能找到
grep "550e8400..." api-gateway.log
grep "550e8400..." order-service.log
grep "550e8400..." inventory-service.log
```

**场景 3：客户端主动指定 RequestID**

有时客户端想自己控制 RequestID（比如用于幂等性检查）：

```javascript
// 前端生成 UUID，作为请求的唯一标识
const myRequestId = crypto.randomUUID();

fetch('/api/v1/orders', {
    method: 'POST',
    headers: {
        'X-Request-ID': myRequestId,  // 客户端指定
        'Content-Type': 'application/json'
    },
    body: JSON.stringify({ ... })
});
```

服务端收到后，会复用这个 ID 而不是生成新的：

```go
requestID := c.GetHeader(HeaderXRequestID)
if requestID == "" {
    requestID = uuid.New().String()  // 客户端没传才生成
}
```

**你作为前端没遇到过的原因**

很多项目没有做好这个设计，出问题时只能靠"时间范围 + 用户 ID"去搜日志，效率很低。

RequestID 是**可观测性**（Observability）的基础设施。不是"必须有"，而是"有了之后排查问题效率翻倍"。

### 4.2 Recovery 中间件

#### 前置知识：Go 的 defer、panic、recover

在理解 Recovery 中间件之前，需要先理解 Go 语言的三个内置机制。

##### panic：程序崩溃

`panic` 是 Go 的内置函数，用于表示"不可恢复的错误"。调用 `panic` 会：

1. 立即停止当前函数的执行
2. 逐层向上返回（展开调用栈）
3. 如果没有被 `recover` 捕获，程序会崩溃并打印堆栈

```go
func main() {
    fmt.Println("开始")
    panic("出大问题了！")  // 程序在这里崩溃
    fmt.Println("结束")    // 这行永远不会执行
}
```

常见的隐式 panic：

```go
var data []int
fmt.Println(data[10])  // 数组越界，触发 panic

var m map[string]int
m["key"] = 1           // nil map 写入，触发 panic

var p *int
fmt.Println(*p)        // 空指针解引用，触发 panic
```

##### defer：延迟执行

`defer` 是 Go 的关键字，用于延迟函数的执行。被 `defer` 的函数会在**当前函数返回之前**执行。

```go
func example() {
    defer fmt.Println("3. 最后执行")  // defer 注册
    fmt.Println("1. 先执行")
    fmt.Println("2. 再执行")
    // 函数返回前，执行所有 defer
}

// 输出：
// 1. 先执行
// 2. 再执行
// 3. 最后执行
```

**关键特性**：即使函数发生 panic，defer 也会执行！

```go
func example() {
    defer fmt.Println("这行会执行！")
    panic("崩溃了")
    fmt.Println("这行不会执行")
}

// 输出：
// 这行会执行！
// panic: 崩溃了
```

这就是为什么 Recovery 中间件要用 `defer`——它确保即使 Handler panic，recover 逻辑也会执行。

##### recover：捕获 panic

`recover` 是 Go 的内置函数，用于"拦截"正在发生的 panic，让程序不崩溃。

**重要限制**：`recover` 只能在 `defer` 函数内调用才有效。

```go
func safeCall() {
    defer func() {
        if err := recover(); err != nil {
            // panic 被捕获，程序继续运行
            fmt.Println("捕获到 panic:", err)
        }
    }()
  
    panic("出问题了")
    // 如果没有上面的 recover，程序会崩溃
}

func main() {
    safeCall()
    fmt.Println("程序继续运行")  // 这行会执行
}
```

| 函数/关键字 | 来源        | 作用                       |
| ----------- | ----------- | -------------------------- |
| `panic`   | Go 内置函数 | 触发程序崩溃               |
| `defer`   | Go 关键字   | 延迟执行，确保清理逻辑运行 |
| `recover` | Go 内置函数 | 捕获 panic，阻止程序崩溃   |

它们不需要 import，就像 `len()`、`make()`、`append()` 一样，是语言内置的。

##### defer + recover 的组合模式

```go
func Recovery() gin.HandlerFunc {
    return func(c *gin.Context) {
        // defer 确保这个函数在 c.Next() 发生 panic 时也会执行
        defer func() {
            // recover() 只能在 defer 内调用
            // 如果没有 panic，recover() 返回 nil
            // 如果有 panic，recover() 返回 panic 的值，并阻止崩溃
            if err := recover(); err != nil {
                // 处理 panic：记录日志、返回错误响应
            }
        }()
  
        c.Next()  // 执行后续中间件和 Handler，可能 panic
    }
}
```

#### c.Abort() 的作用

`c.Abort()` 是 Gin 提供的方法，用于**终止后续中间件和 Handler 的执行**。

```go
defer func() {
    if err := recover(); err != nil {
        response.ErrorJSON(c, ...)  // 返回错误响应
        c.Abort()                   // 终止后续执行
    }
}()
c.Next()  // 执行后续中间件和 Handler
```

**问：服务阻塞了？还是请求中断了？**

**答：请求中断了，服务正常运行。**

| 概念           | 说明                                       |
| -------------- | ------------------------------------------ |
| `c.Abort()`  | 当前请求的处理链被中断，后续中间件不再执行 |
| 服务状态       | 服务本身不受影响，继续处理其他请求         |
| 已执行的 defer | 会正常执行（响应已经写入）                 |

**执行流程对比**：

```
正常情况（无 panic）：
RequestID → Recovery → Logger → Handler → Logger后 → Recovery后 → RequestID后

发生 panic 时（有 c.Abort()）：
RequestID → Recovery → Logger → Handler[panic!] 
                              ↓
                         Recovery 的 defer 捕获
                              ↓
                         response.ErrorJSON() 写入响应
                              ↓
                         c.Abort() 终止后续
                              ↓
                         响应返回客户端
```

如果不调用 `c.Abort()`，Gin 可能会继续执行后续的中间件，导致响应被覆盖或重复写入。

#### Recovery 中间件代码

```go
// internal/server/middleware/recovery.go
package middleware

import (
    "fmt"
    "log"
    "runtime/debug"

    "github.com/gin-gonic/gin"

    "go-api-template/internal/pkg/apperrors"
    "go-api-template/internal/server/response"
)

// Recovery 返回 Panic 恢复中间件
// 职责：
//   - 捕获 Handler 中发生的 panic，防止服务崩溃
//   - 记录错误堆栈到日志（用于调试和告警）
//   - 返回统一格式的 500 错误响应（不暴露内部细节）
func Recovery() gin.HandlerFunc {
    return func(c *gin.Context) {
        defer func() {
            if err := recover(); err != nil {
                requestID := GetRequestID(c)
                stack := debug.Stack()
      
                // 记录到日志（生产环境应接入日志系统/告警）
                log.Printf("[PANIC RECOVERED] request_id=%s error=%v\n%s",
                    requestID, err, string(stack))

                // 返回统一格式的错误响应
                // 注意：不暴露 panic 的具体信息给客户端（安全性）
                response.ErrorJSON(c, apperrors.Internal("服务器内部错误", fmt.Errorf("%v", err)))
                c.Abort()
            }
        }()
        c.Next()
    }
}
```

**关键设计点**：

| 要点     | 说明                                                   |
| -------- | ------------------------------------------------------ |
| 安全性   | 返回通用消息"服务器内部错误"，不暴露 panic 细节        |
| 可观测性 | 堆栈信息写入日志，包含 RequestID                       |
| 格式一致 | 复用 `response.ErrorJSON`，与 Handler 层格式完全一致 |

### 4.3 路由错误处理

```go
// internal/server/middleware/routes.go
package middleware

import (
    "github.com/gin-gonic/gin"

    "go-api-template/internal/pkg/apperrors"
    "go-api-template/internal/pkg/reason"
    "go-api-template/internal/server/response"
)

// HandleNoRoute 返回 404 路由不存在的处理函数
func HandleNoRoute() gin.HandlerFunc {
    return func(c *gin.Context) {
        appErr := apperrors.New(reason.NotFound, "请求的资源不存在")
        response.ErrorJSON(c, appErr)
    }
}

// HandleNoMethod 返回 405 方法不允许的处理函数
func HandleNoMethod() gin.HandlerFunc {
    return func(c *gin.Context) {
        appErr := apperrors.New(reason.InvalidParams, "请求方法不允许")
        appErr.HTTPCode = 405
        response.ErrorJSON(c, appErr)
    }
}
```

**注意**：`NoRoute` 和 `NoMethod` 不是中间件，而是 Gin 的特殊 Handler。但它们的职责是处理"框架层错误"，逻辑上属于本章节。

### 4.4 中间件聚合器

顺序很重要，按顺序使用中间件：

```go
// 不推荐：多次调用 engine.Use()
func Register(engine *gin.Engine) {
    engine.Use(RequestID())
    engine.Use(Recovery())
    engine.Use(gin.Logger())
}
```

这样写有几个问题：

| 问题       | 说明                           |
| ---------- | ------------------------------ |
| 顺序不直观 | 需要从上往下读代码才能理解顺序 |
| 修改麻烦   | 调整顺序需要剪切粘贴代码行     |
| 不够声明式 | 命令式的"做这个、做那个"风格   |

#### 最佳实践：使用切片声明中间件链

```go
// internal/server/middleware/middleware.go
package middleware

import (
    "github.com/gin-gonic/gin"
)

// middlewareChain 定义中间件链
// 顺序很重要，遵循"洋葱模型"：
//
//     请求进入 → [0] → [1] → [2] → Handler
//     响应返回 ← [0] ← [1] ← [2] ← Handler
//
// 使用切片声明的优势：
//  1. 顺序一目了然，修改只需调整数组
//  2. 符合声明式编程风格
//  3. 避免多次调用 engine.Use() 的冗余
var middlewareChain = []gin.HandlerFunc{
    RequestID(),   // [0] 最先执行，确保后续中间件都能获取请求 ID
    Recovery(),    // [1] 捕获后续所有代码的 panic
    gin.Logger(),  // [2] 记录请求日志
}

// Register 注册所有中间件到 Gin 引擎
func Register(engine *gin.Engine) {
    engine.Use(middlewareChain...)
}

// RegisterRouteHandlers 注册路由级别的错误处理
func RegisterRouteHandlers(engine *gin.Engine) {
    engine.NoRoute(HandleNoRoute())
    engine.HandleMethodNotAllowed = true
    engine.NoMethod(HandleNoMethod())
}
```

#### 关键语法：`middlewareChain...`

`...` 是 Go 的**展开操作符**（Spread Operator），用于将切片展开为多个参数：

```go
// 以下两种写法等价
engine.Use(middlewareChain...)
engine.Use(middlewareChain[0], middlewareChain[1], middlewareChain[2])
```

#### 未来扩展

如果需要根据环境动态调整中间件，可以这样做：

```go
func buildMiddlewareChain(cfg *conf.Config) []gin.HandlerFunc {
    chain := []gin.HandlerFunc{
        RequestID(),
        Recovery(),
    }
  
    // 只在开发环境启用详细日志
    if !cfg.IsProduction() {
        chain = append(chain, gin.Logger())
    }
  
    return chain
}
```

### 4.5 重构 http.go

**修改前**：

```go
engine := gin.Default()  // 内置 Logger + Recovery（格式不统一）
```

**修改后**：

```go
// 使用 gin.New() 创建空白引擎
engine := gin.New()

// 注册我们的中间件
middleware.Register(engine)

// 注册 404/405 处理
middleware.RegisterRouteHandlers(engine)
```

---

## 5. 验证测试

启动服务后进行测试：

### 5.1 测试 404 响应

```bash
# PowerShell
Invoke-RestMethod -Uri http://localhost:8080/api/v1/not-exist -Method Get
```

**预期返回**（统一 JSON 格式）：

```json
{
    "code": "NOT_FOUND",
    "message": "请求的资源不存在",
    "http_code": 404
}
```

### 5.2 测试 405 响应

```bash
# 对只支持 POST 的接口发送 DELETE 请求
Invoke-RestMethod -Uri http://localhost:8080/api/v1/greeter/say-hello -Method Delete
```

**预期返回**：

```json
{
    "code": "INVALID_PARAMS",
    "message": "请求方法不允许",
    "http_code": 405
}
```

### 5.3 测试请求 ID

```bash
# 检查响应头中的 X-Request-ID
$response = Invoke-WebRequest -Uri http://localhost:8080/health -Method Get
$response.Headers["X-Request-ID"]
```

**预期**：返回一个 UUID 格式的字符串，如 `550e8400-e29b-41d4-a716-446655440000`

### 5.4 测试正常请求（确保不影响原有功能）

```bash
Invoke-RestMethod -Uri http://localhost:8080/api/v1/greeter/say-hello `
    -Method Post -ContentType "application/json" `
    -Body '{"name": "World"}'
```

**预期返回**（与之前一致）：

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

---

## 6. 总结

### 6.1 本章节解决的问题

| 问题           | 修改前                           | 修改后                    |
| -------------- | -------------------------------- | ------------------------- |
| 404 路由不存在 | `404 page not found`（纯文本） | 统一 JSON 格式            |
| 405 方法不允许 | `method not allowed`（纯文本） | 统一 JSON 格式            |
| Handler panic  | 空响应或 HTML                    | 统一 JSON 格式 + 日志记录 |
| 请求追踪       | 无                               | 每个请求有唯一 ID         |

### 6.2 新增文件

| 文件                                         | 职责                         |
| -------------------------------------------- | ---------------------------- |
| `internal/server/middleware/middleware.go` | 中间件聚合器，统一注册入口   |
| `internal/server/middleware/request_id.go` | 请求 ID 生成与传递           |
| `internal/server/middleware/recovery.go`   | Panic 捕获，返回统一错误格式 |
| `internal/server/middleware/routes.go`     | 404/405 处理器               |

### 6.3 架构回顾

完成本章节后，"统一错误处理、请求验证、响应格式"的目标已经**完整实现**：

```
┌─────────────────────────────────────────────────────────────────
│                        请求到达  
└───────────────────────────────┬─────────────────────────────────
                                │
                                ▼
┌─────────────────────────────────────────────────────────────────
│  中间件层  
│  ┌─────────────┬─────────────┬─────────────  
│  │ RequestID   │ Recovery    │ Logger  
│  │ 请求追踪     │ Panic→JSON  │ 请求日志   
│  └─────────────┴─────────────┴───────────── 
│  ┌─────────────┬─────────────
│  │ NoRoute     │ NoMethod   
│  │ 404→JSON    │ 405→JSON  
│  └─────────────┴─────────────
└───────────────────────────────┬─────────────────────────────────
                                │
                                ▼
┌─────────────────────────────────────────────────────────────────
│  Handler 层  
│  ┌─────────────────────────────────────────────────────────  
│  │ 验证错误 → apperrors.FromValidationError() → JSON  
│  │ 业务错误 → apperrors.NotFound/Internal() → JSON  
│  │ 成功响应 → response.SuccessJSON() → JSON   
│  └─────────────────────────────────────────────────────────   
└─────────────────────────────────────────────────────────────────
```

### 6.4 与其他章节的关系

| 章节                 | 关系                                |
| -------------------- | ----------------------------------- |
| 010 Struct Tag       | 定义验证规则，供 011 使用           |
| 011 DTO 模式         | Handler 内的请求验证                |
| 012 统一错误处理     | AppError 类型，本章节复用           |
| 013 统一响应格式     | Response 输出，本章节复用           |
| **014 中间件** | **补完 Handler 外的错误处理** |
