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
//
// 为什么需要自定义 Recovery？
// Gin 内置的 gin.Recovery() 返回的是纯文本或 HTML 格式，
// 将返回内容统一为 JSON 响应格式和一致的结构体。
func Recovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				// 获取请求 ID（如果有）
				requestID := GetRequestID(c)

				// 获取堆栈信息
				stack := debug.Stack()

				// 记录到日志（生产环境应接入日志系统/告警）
				// 这里使用标准库 log，未来可替换为 zap/zerolog
				log.Printf("[PANIC RECOVERED] request_id=%s error=%v\n%s",
					requestID, err, string(stack))

				// 返回统一格式的错误响应
				// 注意：不暴露 panic 的具体信息给客户端（安全性）
				response.ErrorJSON(c, apperrors.Internal("服务器内部错误", fmt.Errorf("%v", err)))

				// 终止后续中间件和 Handler 的执行
				c.Abort()
			}
		}()

		c.Next()
	}
}
