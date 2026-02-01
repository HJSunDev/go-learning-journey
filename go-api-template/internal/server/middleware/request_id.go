// Package middleware 提供 HTTP 中间件
// 处理横切关注点：请求追踪、异常恢复、日志等
package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// 请求 ID 相关常量
const (
	// HeaderXRequestID 请求 ID 的 HTTP 头名称
	// 客户端可以通过此头传入自定义的请求 ID
	HeaderXRequestID = "X-Request-ID"

	// ContextKeyRequestID 请求 ID 在 gin.Context 中的键名
	// 用于在 Handler 和 Service 层获取请求 ID
	ContextKeyRequestID = "request_id"
)

// RequestID 返回请求 ID 中间件
// 职责：
//   - 从请求头提取 X-Request-ID，如果没有则生成 UUID
//   - 将请求 ID 存入 gin.Context，供后续处理使用
//   - 将请求 ID 写入响应头，方便客户端关联请求和响应
func RequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 尝试从请求头获取，支持客户端传入自定义 ID
		requestID := c.GetHeader(HeaderXRequestID)

		// 如果客户端没有传，生成新的 UUID
		if requestID == "" {
			requestID = uuid.New().String()
		}

		// 存入 Context，供 Handler/Service 层使用
		c.Set(ContextKeyRequestID, requestID)

		// 写入响应头，方便客户端追踪
		c.Header(HeaderXRequestID, requestID)

		c.Next()
	}
}

// GetRequestID 从 gin.Context 中获取请求 ID
// 如果不存在则返回空字符串
func GetRequestID(c *gin.Context) string {
	if id, exists := c.Get(ContextKeyRequestID); exists {
		if requestID, ok := id.(string); ok {
			return requestID
		}
	}
	return ""
}
