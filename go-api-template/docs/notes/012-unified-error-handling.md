# 012. 统一错误处理

## 0. 前置知识

在阅读本章节前，需要理解以下 Go 语言和 Gin 框架的基础知识。

### 0.1 类型断言（Type Assertion）

在代码中你会看到这样的写法：

```go
validationErrors, ok := err.(validator.ValidationErrors)
```

这**不是方法调用**，而是 Go 语言的**类型断言**语法。

#### 语法解释

```go
value, ok := 变量.(目标类型)
```

| 部分                              | 含义                                                      |
| --------------------------------- | --------------------------------------------------------- |
| `err`                           | 一个 `error` 接口类型的变量                             |
| `.(validator.ValidationErrors)` | 尝试把 `err` 断言为 `validator.ValidationErrors` 类型 |
| `validationErrors`              | 如果断言成功，这就是转换后的值                            |
| `ok`                            | 布尔值，表示断言是否成功                                  |

#### 为什么需要类型断言？

Go 的 `error` 是一个接口类型：

```go
type error interface {
    Error() string
}
```

任何实现了 `Error() string` 方法的类型都可以作为 `error`。但当我们拿到一个 `error` 时，只能调用 `Error()` 方法，无法访问具体类型的其他字段或方法。

```go
// ShouldBindJSON 返回的 err 是 error 接口
if err := c.ShouldBindJSON(&req); err != nil {
    // 此时 err 的类型是 error 接口
    // 但实际上它可能是 validator.ValidationErrors 类型
    // 我们需要"拆开"这个接口，才能访问具体类型的字段
  
    validationErrors, ok := err.(validator.ValidationErrors)
    if ok {
        // 断言成功！现在可以访问 ValidationErrors 的方法了
        for _, fieldErr := range validationErrors {
            fmt.Println(fieldErr.Field())  // 获取字段名
            fmt.Println(fieldErr.Tag())    // 获取验证规则
        }
    }
}
```

#### 两种断言写法

```go
// 写法 1：安全断言（推荐）
value, ok := err.(TargetType)
if ok {
    // 使用 value
}

// 写法 2：直接断言（不推荐，断言失败会 panic）
value := err.(TargetType)  // 如果 err 不是 TargetType，程序会崩溃！
```

**结论**：`err.(validator.ValidationErrors)` 的意思是"尝试把 `err` 接口变量还原为 `validator.ValidationErrors` 具体类型"。

---

### 0.2 Go 包命名规范

为什么包名叫 `apperrors` 而不是 `errors` 或 `AppErrors`？

#### Go 包命名的核心规则

| 规则               | 说明                                 | 示例                                               |
| ------------------ | ------------------------------------ | -------------------------------------------------- |
| **全小写**   | 包名必须全部小写，不使用下划线或驼峰 | ✅`apperrors` ❌ `appErrors` ❌ `app_errors` |
| **简短**     | 包名应该简短，通常是一个单词         | ✅`http` ✅ `json` ✅ `errors`               |
| **有意义**   | 包名应该描述包的功能                 | ✅`apperrors`（应用错误）                        |
| **避免重名** | 不要与标准库包名冲突                 | ❌`errors`（与标准库冲突）                       |

#### 为什么不能叫 `errors`？

Go 标准库有一个 `errors` 包：

```go
import "errors"  // 标准库

err := errors.New("something went wrong")
```

如果我们的包也叫 `errors`：

```go
import "errors"                              // 标准库
import "go-api-template/internal/pkg/errors" // 我们的包

// 冲突了！编译器不知道 errors.New() 是哪个包的
```

虽然可以用别名解决：

```go
import "errors"
import apperrors "go-api-template/internal/pkg/errors"  // 被迫加别名
```

但这很麻烦。**更好的做法是直接把包名改成 `apperrors`**，这样：

```go
import "errors"                                 // 标准库，正常用
import "go-api-template/internal/pkg/apperrors" // 我们的包，无需别名

err := errors.New("standard error")
appErr := apperrors.New(ErrCodeNotFound, "user not found")
```

#### 包名与目录名

在 Go 中，**包名通常与目录名相同**（但不是强制）：

```
internal/pkg/apperrors/     <- 目录名
            ├── codes.go    <- package apperrors
            └── errors.go   <- package apperrors
```

---

### 0.3 Gin 的 `c.JSON` 和 `gin.H`

在 HTTP Handler 中，你会看到两种返回响应的写法：

```go
// 写法 1：使用 gin.H
c.JSON(http.StatusOK, gin.H{
    "message": resp.GetMessage(),
})

// 写法 2：使用 struct
c.JSON(appErr.HTTPCode, appErr)
```

#### `c.JSON` 是什么？

`c` 是 `*gin.Context` 类型，代表当前 HTTP 请求的上下文。

`c.JSON()` 是 Gin 提供的方法，用于返回 JSON 响应：

```go
func (c *Context) JSON(code int, obj any)
```

| 参数     | 含义                                   |
| -------- | -------------------------------------- |
| `code` | HTTP 状态码（如 200, 400, 500）        |
| `obj`  | 任意值，会被序列化为 JSON 返回给客户端 |

`c.JSON` 内部会：

1. 设置响应头 `Content-Type: application/json`
2. 设置 HTTP 状态码
3. 将 `obj` 序列化为 JSON 字符串
4. 写入响应体

#### `gin.H` 是什么？

`gin.H` 的定义非常简单：

```go
// Gin 框架源码
type H map[string]any
```

`H` 是 **"Hash"** 的缩写，代表一个键值对的哈希表（Map）。

它就是 `map[string]any` 的别名，用于快速构造 JSON 对象：

```go
// 这两种写法完全等价
gin.H{"message": "hello", "count": 42}
map[string]any{"message": "hello", "count": 42}
```

使用 `gin.H` 的好处是**写起来更短**。

#### 为什么响应格式看起来不一样？

这是一个很好的观察！让我们对比两种响应：

**成功类型响应**（当前代码）：

```go
c.JSON(http.StatusOK, gin.H{
    "message": resp.GetMessage(),
})
// 输出: {"message": "Hello, World!"}
```

**错误类型响应**（使用 AppError）：

```go
appErr := apperrors.FromValidationError(err)
c.JSON(appErr.HTTPCode, appErr)
// 输出: {"code":"INVALID_PARAMS","message":"请求参数验证失败","http_code":400,"details":[...]}
```

**为什么格式不同？**

因为 `c.JSON` 的第二个参数传的是不同的东西：

| 传入的值                      | 序列化结果                     |
| ----------------------------- | ------------------------------ |
| `gin.H{"message": "hello"}` | `{"message": "hello"}`       |
| `appErr`（AppError struct） | 根据 struct 的 json tag 序列化 |

当你传入一个 struct 时，Go 会根据 struct 的字段和 json tag 自动序列化：

```go
type AppError struct {
    Code     ErrorCode    `json:"code"`
    Message  string       `json:"message"`
    HTTPCode int          `json:"http_code"`
    Details  []FieldError `json:"details,omitempty"`
    Cause    error        `json:"-"`  // 这个字段不会出现在 JSON 中
}

// 当 c.JSON(code, appErr) 时，会生成：
// {"code":"xxx", "message":"xxx", "http_code":400, "details":[...]}
```

#### 这就是为什么需要"统一响应格式"

你发现了一个问题：**成功和失败的响应格式不一致**！

- 成功：`{"message": "..."}`
- 失败：`{"code": "...", "message": "...", "http_code": ..., "details": [...]}`

这确实不好。**013 章节"统一响应格式"就是要解决这个问题**，让所有响应都有统一的结构：

```json
// 成功响应（013 章节目标）
{
    "code": "SUCCESS",
    "message": "操作成功",
    "http_code": 200,
    "data": {
        "message": "Hello, World!"
    }
}

// 失败响应（本章节已实现）
{
    "code": "INVALID_PARAMS",
    "message": "请求参数验证失败",
    "http_code": 400,
    "details": [...]
}
```

---

## 1. 核心问题与概念

### 1.1 解决什么问题

在 011 章节中，我们用 DTO + Validator 实现了自动验证，但错误返回还有问题：

```go
// http.go 中的代码
if err := c.ShouldBindJSON(&req); err != nil {
    c.JSON(http.StatusBadRequest, gin.H{
        "error": err.Error(),  // 问题在这里！
    })
    return
}
```

**问题分析**：

| 问题       | 影响                                                                                                              |
| ---------- | ----------------------------------------------------------------------------------------------------------------- |
| 没有错误码 | 前端无法根据错误类型做不同处理                                                                                    |
| 消息不友好 | Validator 返回的是 `Key: 'SayHelloRequest.Name' Error:Field validation for 'Name' failed on the 'required' tag` |
| 格式不统一 | 不同地方的错误返回格式可能不一致                                                                                  |
| 无法扩展   | 以后加新的错误类型（如权限错误）时无法统一处理                                                                    |

### 1.2 核心概念

| 概念                 | 说明                                                          |
| -------------------- | ------------------------------------------------------------- |
| **AppError**   | 自定义错误类型，包含错误码、消息、HTTP 状态码等               |
| **ErrorCode**  | 业务错误码，字符串类型（如 `INVALID_PARAMS`），便于前端识别 |
| **FieldError** | 字段级错误，用于验证失败时告诉用户具体哪个字段出错            |

### 1.3 目标格式

统一后的错误响应结构：

```json
{
    "code": "INVALID_PARAMS",
    "message": "请求参数验证失败",
    "http_code": 400,
    "details": [
        {"field": "name", "message": "必填字段"},
        {"field": "email", "message": "邮箱格式不正确"}
    ]
}
```

---

## 2. 设计方案

### 2.1 文件结构

```
internal/pkg/apperrors/
├── codes.go    # 错误码定义
└── errors.go   # AppError 类型和工具函数
```

**为什么叫 `apperrors` 而不是 `errors`？**

- 避免与 Go 标准库 `errors` 包重名
- 调用时 `apperrors.New()` 语义清晰，无需别名
- 符合 Go 社区的最佳实践（明确的包名优于通用的包名）

### 2.2 AppError 结构

```go
type AppError struct {
    Code     ErrorCode    `json:"code"`              // 业务错误码
    Message  string        `json:"message"`           // 错误消息（面向用户）
    HTTPCode int           `json:"http_code"`         // HTTP 状态码
    Details  []FieldError  `json:"details,omitempty"` // 字段级错误详情
    Cause    error        `json:"-"`                 // 原始错误（用于日志，不暴露给客户端）
}

type FieldError struct {
    Field   string `json:"field"`   // 字段名
    Message string `json:"message"` // 错误描述
}
```

**设计要点**：

| 字段     | `json` tag                 | 说明                                     |
| -------- | ---------------------------- | ---------------------------------------- |
| Code     | `json:"code"`              | 暴露给前端，用于程序判断                 |
| Message  | `json:"message"`           | 暴露给前端，用于用户展示                 |
| HTTPCode | `json:"http_code"`         | 暴露给前端，方便前端直接获取 HTTP 状态码 |
| Details  | `json:"details,omitempty"` | 有值时才序列化，用于字段级错误           |
| Cause    | `json:"-"`                 | 不暴露给前端，仅用于内部日志             |

### 2.3 错误码设计

```go
type ErrorCode string

const (
    // 客户端错误 (4xx)
    ErrCodeInvalidParams ErrorCode = "INVALID_PARAMS"  // 400
    ErrCodeUnauthorized  ErrorCode = "UNAUTHORIZED"    // 401
    ErrCodeForbidden     ErrorCode = "FORBIDDEN"       // 403
    ErrCodeNotFound      ErrorCode = "NOT_FOUND"       // 404

    // 服务端错误 (5xx)
    ErrCodeInternal           ErrorCode = "INTERNAL_ERROR"      // 500
    ErrCodeServiceUnavailable ErrorCode = "SERVICE_UNAVAILABLE" // 503
)
```

**为什么用字符串而不是数字？**

- 可读性好：`"INVALID_PARAMS"` 比 `10001` 更直观
- 调试友好：日志中直接能看出错误类型
- 兼容性好：不会与 HTTP 状态码混淆

---

## 3. 代码实现

### 3.1 codes.go - 错误码定义

```go
// internal/pkg/apperrors/codes.go
package apperrors

type ErrorCode string

const (
    ErrCodeInvalidParams      ErrorCode = "INVALID_PARAMS"
    ErrCodeUnauthorized       ErrorCode = "UNAUTHORIZED"
    ErrCodeForbidden          ErrorCode = "FORBIDDEN"
    ErrCodeNotFound           ErrorCode = "NOT_FOUND"
    ErrCodeInternal           ErrorCode = "INTERNAL_ERROR"
    ErrCodeServiceUnavailable ErrorCode = "SERVICE_UNAVAILABLE"
)

// codeHTTPStatus 错误码到 HTTP 状态码的映射
var codeHTTPStatus = map[ErrorCode]int{
    ErrCodeInvalidParams:      400,
    ErrCodeUnauthorized:       401,
    ErrCodeForbidden:          403,
    ErrCodeNotFound:           404,
    ErrCodeInternal:           500,
    ErrCodeServiceUnavailable: 503,
}

// HTTPStatus 返回错误码对应的 HTTP 状态码
func (c ErrorCode) HTTPStatus() int {
    if status, ok := codeHTTPStatus[c]; ok {
        return status
    }
    return 500
}
```

### 3.2 errors.go - AppError 实现

```go
// internal/pkg/apperrors/errors.go
package apperrors

import (
    "fmt"
    "strings"

    "github.com/go-playground/validator/v10"
)

type FieldError struct {
    Field   string `json:"field"`
    Message string `json:"message"`
}

type AppError struct {
    Code     ErrorCode    `json:"code"`              // 业务错误码
    Message  string       `json:"message"`           // 面向用户的错误消息
    HTTPCode int          `json:"http_code"`         // HTTP 状态码
    Details  []FieldError `json:"details,omitempty"` // 字段级错误详情
    Cause    error        `json:"-"`                 // 原始错误（用于日志，不暴露给客户端）
}

// Error 实现 error 接口
func (e *AppError) Error() string {
    if e.Cause != nil {
        return fmt.Sprintf("%s: %s (cause: %v)", e.Code, e.Message, e.Cause)
    }
    return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

// Unwrap 支持 errors.Unwrap，用于错误链
func (e *AppError) Unwrap() error {
    return e.Cause
}
```

### 3.3 构造函数

```go
// New 创建一个新的 AppError
func New(code ErrorCode, message string) *AppError {
    return &AppError{
        Code:     code,
        Message:  message,
        HTTPCode: code.HTTPStatus(),
    }
}

// Wrap 包装一个已有错误
func Wrap(code ErrorCode, message string, cause error) *AppError {
    return &AppError{
        Code:     code,
        Message:  message,
        HTTPCode: code.HTTPStatus(),
        Cause:    cause,
    }
}

// 常用错误快捷构造
func InvalidParams(message string) *AppError {
    return New(ErrCodeInvalidParams, message)
}

func NotFound(message string) *AppError {
    return New(ErrCodeNotFound, message)
}

func Internal(message string, cause error) *AppError {
    return Wrap(ErrCodeInternal, message, cause)
}
```

### 3.4 验证错误转换

这是最核心的部分：将 Validator 的错误转换为友好格式。

```go
// FromValidationError 将 validator 库的错误转换为 AppError
func FromValidationError(err error) *AppError {
    validationErrors, ok := err.(validator.ValidationErrors)
    if !ok {
        // 不是验证错误（可能是 JSON 解析错误）
        return InvalidParams(err.Error())
    }

    // 转换每个字段错误
    details := make([]FieldError, 0, len(validationErrors))
    for _, fieldErr := range validationErrors {
        details = append(details, FieldError{
            Field:   toJSONFieldName(fieldErr.Field()),
            Message: translateValidationError(fieldErr),
        })
    }

    return &AppError{
        Code:     ErrCodeInvalidParams,
        Message:  "请求参数验证失败",
        HTTPCode: 400,
        Details:  details,
    }
}

// toJSONFieldName 将 Go 字段名转为 JSON 字段名（首字母小写）
func toJSONFieldName(field string) string {
    if len(field) == 0 {
        return field
    }
    return strings.ToLower(field[:1]) + field[1:]
}

// translateValidationError 将验证错误转为友好消息
func translateValidationError(fe validator.FieldError) string {
    switch fe.Tag() {
    case "required":
        return "必填字段"
    case "min":
        return fmt.Sprintf("最小长度为 %s", fe.Param())
    case "max":
        return fmt.Sprintf("最大长度为 %s", fe.Param())
    case "email":
        return "邮箱格式不正确"
    // ... 更多规则
    default:
        return fmt.Sprintf("验证失败: %s", fe.Tag())
    }
}
```

---

## 4. 使用方式

### 4.1 修改 HTTP Handler

**修改前**：

```go
if err := c.ShouldBindJSON(&req); err != nil {
    c.JSON(http.StatusBadRequest, gin.H{
        "error": err.Error(),
    })
    return
}
```

**修改后**：

```go
import "go-api-template/internal/pkg/apperrors"

if err := c.ShouldBindJSON(&req); err != nil {
    appErr := apperrors.FromValidationError(err)
    c.JSON(appErr.HTTPCode, appErr)
    return
}
```

### 4.2 处理业务错误

```go
// 资源不存在
if user == nil {
    appErr := apperrors.NotFound("用户不存在")
    c.JSON(appErr.HTTPCode, appErr)
    return
}

// 内部错误（不暴露细节给用户）
data, err := repo.GetData(ctx)
if err != nil {
    appErr := apperrors.Internal("获取数据失败", err)
    // Cause 会记录在 appErr 中用于日志，但不会返回给客户端
    c.JSON(appErr.HTTPCode, appErr)
    return
}
```

---

## 5. 验证测试

启动服务后测试：

```bash
# 测试验证错误
curl -X POST http://localhost:8080/api/v1/greeter/say-hello \
  -H "Content-Type: application/json" \
  -d '{}'

# 预期返回：
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

## 6. 最佳实践与注意事项

### 6.1 推荐做法

| 场景       | 做法                                                        |
| ---------- | ----------------------------------------------------------- |
| 验证错误   | 使用 `FromValidationError()` 自动转换                     |
| 资源不存在 | 使用 `NotFound("xx不存在")`                               |
| 内部错误   | 使用 `Internal("操作失败", err)`，原始错误记录在 Cause 中 |
| 权限错误   | 使用 `Unauthorized()` 或 `Forbidden()`                  |

### 6.2 避免做法

| 错误做法                 | 问题               | 正确做法                             |
| ------------------------ | ------------------ | ------------------------------------ |
| 直接返回 `err.Error()` | 可能暴露内部细节   | 使用 `AppError` 包装               |
| 返回数据库错误消息       | 暴露技术细节       | 返回通用消息，原始错误存入 `Cause` |
| 忘记设置 HTTP 状态码     | 客户端无法正确判断 | 使用 `appErr.HTTPCode`             |

### 6.3 关于错误消息的语言

本章节使用中文错误消息作为示例。实际项目中可以：

1. **固定使用中文**：面向国内用户的项目
2. **固定使用英文**：国际化项目
3. **使用 i18n**：根据请求头 `Accept-Language` 返回不同语言

---

## 7. 总结

### 核心文件

| 文件                                 | 职责                                  |
| ------------------------------------ | ------------------------------------- |
| `internal/pkg/apperrors/codes.go`  | 定义错误码和 HTTP 状态码映射          |
| `internal/pkg/apperrors/errors.go` | AppError 类型、构造函数、验证错误转换 |

### 数据流向

```
Validator 错误
      ↓
FromValidationError()
      ↓
AppError (统一格式)
      ↓
c.JSON(appErr.HTTPCode, appErr)
      ↓
{
    "code": "INVALID_PARAMS",
    "message": "请求参数验证失败",
    "http_code": 400,
    "details": [...]
}
```

### 当前的局限性

本章节只统一了**错误响应**的格式，成功响应仍然是简单的 `{"message": "..."}`。

这意味着前端需要用不同的逻辑处理成功和失败的响应，不够优雅。

**013 章节"统一响应格式"将解决这个问题**，让所有响应都有一致的结构。

### 与其他章节的关系

| 章节             | 关系                                   |
| ---------------- | -------------------------------------- |
| 011 DTO 模式     | 提供验证错误，本章节将其转换为友好格式 |
| 013 统一响应格式 | 基于本章节的 AppError 构建统一响应     |
| 014 中间件入门   | 使用本章节的错误处理处理 panic 和 404  |
