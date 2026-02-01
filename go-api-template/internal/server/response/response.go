// Package response 提供统一的 HTTP 响应格式
// 确保所有 API 响应具有一致的结构，方便前端处理
package response

import (
	"go-api-template/internal/pkg/apperrors"
	"go-api-template/internal/pkg/reason"

	"github.com/gin-gonic/gin"
)

// Body 响应数据别名，避免在 Handler 中频繁引用 gin.H
// 如果不想定义结构体，可以直接使用 response.Body{"key": "value"}
type Body map[string]any

// Response 统一响应结构
// 所有 HTTP 响应都使用此结构，确保格式一致
type Response struct {
	Code     reason.Reason          `json:"code"`              // 业务状态码
	Message  string                 `json:"message"`           // 状态消息
	HTTPCode int                    `json:"http_code"`         // HTTP 状态码
	Data     any                    `json:"data,omitempty"`    // 成功时的数据
	Details  []apperrors.FieldError `json:"details,omitempty"` // 失败时的字段错误详情
}

// ==================== 成功响应构造函数 ====================

// Success 创建成功响应
// data 参数会被放入 Response.Data 字段
func Success(data any) *Response {
	return &Response{
		Code:     reason.Success,
		Message:  "操作成功",
		HTTPCode: 200,
		Data:     data,
	}
}

// SuccessWithMessage 创建带自定义消息的成功响应
func SuccessWithMessage(message string, data any) *Response {
	return &Response{
		Code:     reason.Success,
		Message:  message,
		HTTPCode: 200,
		Data:     data,
	}
}

// ==================== 错误响应构造函数 ====================

// Error 从 AppError 创建错误响应
// 将内部错误类型转换为统一的响应格式
func Error(err *apperrors.AppError) *Response {
	return &Response{
		Code:     err.Code,
		Message:  err.Message,
		HTTPCode: err.HTTPCode,
		Details:  err.Details,
	}
}

// ==================== Gin 响应输出 ====================

// JSON 输出统一响应到 gin.Context
// 使用 Response 的 HTTPCode 作为 HTTP 状态码
func JSON(c *gin.Context, r *Response) {
	c.JSON(r.HTTPCode, r)
}

// SuccessJSON 快捷方法：输出成功响应
func SuccessJSON(c *gin.Context, data any) {
	JSON(c, Success(data))
}

// ErrorJSON 快捷方法：输出错误响应
func ErrorJSON(c *gin.Context, err *apperrors.AppError) {
	JSON(c, Error(err))
}
