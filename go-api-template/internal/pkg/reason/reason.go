// Package reason 定义业务状态原因
// 这是一个纯字典包，只定义字符串常量，不包含错误处理逻辑
// 这里的定义既用于错误响应，也用于成功响应
package reason

// Reason 业务状态原因
// 使用字符串而非整数，便于阅读和调试
type Reason string

const (
	// ==================== 成功状态 (2xx) ====================

	// Success 操作成功
	// 用于所有成功响应的状态码
	Success Reason = "SUCCESS"

	// ==================== 客户端错误 (4xx) ====================

	// InvalidParams 请求参数验证失败
	// 用于 binding tag 验证失败、参数格式错误等场景
	InvalidParams Reason = "INVALID_PARAMS"

	// Unauthorized 未授权
	// 用于未登录、token 无效等场景
	Unauthorized Reason = "UNAUTHORIZED"

	// Forbidden 禁止访问
	// 用于已登录但无权限访问的场景
	Forbidden Reason = "FORBIDDEN"

	// NotFound 资源不存在
	// 用于请求的资源未找到的场景
	NotFound Reason = "NOT_FOUND"

	// ==================== 服务端错误 (5xx) ====================

	// InternalError 内部错误
	// 用于未预期的服务端错误，不应暴露具体信息给客户端
	InternalError Reason = "INTERNAL_ERROR"

	// ServiceUnavailable 服务不可用
	// 用于依赖服务故障、维护等场景
	ServiceUnavailable Reason = "SERVICE_UNAVAILABLE"
)

// codeHTTPStatus Reason 到 HTTP 状态码的映射
var codeHTTPStatus = map[Reason]int{
	Success:            200,
	InvalidParams:      400,
	Unauthorized:       401,
	Forbidden:          403,
	NotFound:           404,
	InternalError:      500,
	ServiceUnavailable: 503,
}

// HTTPStatus 返回 Reason 对应的 HTTP 状态码
func (r Reason) HTTPStatus() int {
	if status, ok := codeHTTPStatus[r]; ok {
		return status
	}
	// 未知 Reason 默认返回 500
	return 500
}
