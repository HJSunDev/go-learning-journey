// Package apperrors 提供统一的错误处理机制
// 本文件定义业务错误码，用于前后端统一的错误识别
package apperrors

// ErrorCode 业务错误码类型
// 使用字符串而非整数，便于阅读和调试
type ErrorCode string

// 通用错误码定义
// 按功能分组，便于管理和扩展
const (
	// ==================== 客户端错误 (4xx) ====================

	// ErrCodeInvalidParams 请求参数验证失败
	// 用于 binding tag 验证失败、参数格式错误等场景
	ErrCodeInvalidParams ErrorCode = "INVALID_PARAMS"

	// ErrCodeUnauthorized 未授权
	// 用于未登录、token 无效等场景
	ErrCodeUnauthorized ErrorCode = "UNAUTHORIZED"

	// ErrCodeForbidden 禁止访问
	// 用于已登录但无权限访问的场景
	ErrCodeForbidden ErrorCode = "FORBIDDEN"

	// ErrCodeNotFound 资源不存在
	// 用于请求的资源未找到的场景
	ErrCodeNotFound ErrorCode = "NOT_FOUND"

	// ==================== 服务端错误 (5xx) ====================

	// ErrCodeInternal 内部错误
	// 用于未预期的服务端错误，不应暴露具体信息给客户端
	ErrCodeInternal ErrorCode = "INTERNAL_ERROR"

	// ErrCodeServiceUnavailable 服务不可用
	// 用于依赖服务故障、维护等场景
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
	// 未知错误码默认返回 500
	return 500
}
