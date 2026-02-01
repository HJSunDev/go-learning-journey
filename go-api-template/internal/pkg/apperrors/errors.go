package apperrors

import (
	"fmt"
	"strings"

	"go-api-template/internal/pkg/reason"

	"github.com/go-playground/validator/v10"
)

// FieldError 字段级别的错误信息
// 用于验证错误时，提供每个字段的具体错误
type FieldError struct {
	Field   string `json:"field"`   // 字段名（JSON 格式）
	Message string `json:"message"` // 错误描述
}

// AppError 应用错误类型
// 统一的错误结构，包含错误码、消息、HTTP 状态码等信息
type AppError struct {
	Code     reason.Reason `json:"code"`              // 业务错误码
	Message  string        `json:"message"`           // 错误消息（面向用户）
	HTTPCode int           `json:"http_code"`         // HTTP 状态码
	Details  []FieldError  `json:"details,omitempty"` // 字段级错误详情
	Cause    error         `json:"-"`                 // 原始错误（用于日志，不暴露给客户端）
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

// ==================== 构造函数 ====================

// New 创建一个新的 AppError
func New(code reason.Reason, message string) *AppError {
	return &AppError{
		Code:     code,
		Message:  message,
		HTTPCode: code.HTTPStatus(),
	}
}

// Wrap 包装一个已有错误，添加业务错误码
// 用于将底层错误（如数据库错误）转换为业务错误
func Wrap(code reason.Reason, message string, cause error) *AppError {
	return &AppError{
		Code:     code,
		Message:  message,
		HTTPCode: code.HTTPStatus(),
		Cause:    cause,
	}
}

// WithDetails 添加字段级错误详情
func (e *AppError) WithDetails(details []FieldError) *AppError {
	e.Details = details
	return e
}

// ==================== 常用错误快捷构造 ====================

// InvalidParams 创建参数验证失败错误
func InvalidParams(message string) *AppError {
	return New(reason.InvalidParams, message)
}

// InvalidParamsWithDetails 创建带字段详情的参数验证失败错误
func InvalidParamsWithDetails(message string, details []FieldError) *AppError {
	return New(reason.InvalidParams, message).WithDetails(details)
}

// NotFound 创建资源不存在错误
func NotFound(message string) *AppError {
	return New(reason.NotFound, message)
}

// Internal 创建内部错误
// cause 参数用于记录原始错误（日志用），不会暴露给客户端
func Internal(message string, cause error) *AppError {
	return Wrap(reason.InternalError, message, cause)
}

// Unauthorized 创建未授权错误
func Unauthorized(message string) *AppError {
	return New(reason.Unauthorized, message)
}

// Forbidden 创建禁止访问错误
func Forbidden(message string) *AppError {
	return New(reason.Forbidden, message)
}

// ==================== 验证错误转换 ====================

// FromValidationError 将 validator 库的错误转换为 AppError
// 提取每个字段的验证错误，生成友好的错误消息
func FromValidationError(err error) *AppError {
	// 尝试断言为 validator.ValidationErrors
	validationErrors, ok := err.(validator.ValidationErrors)
	if !ok {
		// 不是验证错误，可能是 JSON 解析错误等
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

	return InvalidParamsWithDetails("请求参数验证失败", details)
}

// toJSONFieldName 将结构体字段名转换为 JSON 字段名（小写驼峰）
// 简单实现：首字母小写
func toJSONFieldName(field string) string {
	if len(field) == 0 {
		return field
	}
	return strings.ToLower(field[:1]) + field[1:]
}

// translateValidationError 将 validator 的错误转换为友好的中文消息
func translateValidationError(fe validator.FieldError) string {
	// 根据验证 tag 返回对应的友好消息
	switch fe.Tag() {
	case "required":
		return "必填字段"
	case "min":
		return fmt.Sprintf("最小长度为 %s", fe.Param())
	case "max":
		return fmt.Sprintf("最大长度为 %s", fe.Param())
	case "len":
		return fmt.Sprintf("长度必须为 %s", fe.Param())
	case "email":
		return "邮箱格式不正确"
	case "url":
		return "URL 格式不正确"
	case "numeric":
		return "必须是数字"
	case "alpha":
		return "只能包含字母"
	case "alphanum":
		return "只能包含字母和数字"
	case "gt":
		return fmt.Sprintf("必须大于 %s", fe.Param())
	case "gte":
		return fmt.Sprintf("必须大于等于 %s", fe.Param())
	case "lt":
		return fmt.Sprintf("必须小于 %s", fe.Param())
	case "lte":
		return fmt.Sprintf("必须小于等于 %s", fe.Param())
	case "oneof":
		return fmt.Sprintf("必须是以下值之一: %s", fe.Param())
	default:
		// 未知的验证规则，返回原始错误
		return fmt.Sprintf("验证失败: %s", fe.Tag())
	}
}
