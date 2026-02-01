package middleware

import (
	"github.com/gin-gonic/gin"

	"go-api-template/internal/pkg/apperrors"
	"go-api-template/internal/pkg/reason"
	"go-api-template/internal/server/response"
)

// HandleNoRoute 返回 404 路由不存在的处理函数
// 当请求的路径不存在时，Gin 会调用此函数
// 替代 Gin 默认的 "404 page not found" 纯文本响应
func HandleNoRoute() gin.HandlerFunc {
	return func(c *gin.Context) {
		appErr := apperrors.New(reason.NotFound, "请求的资源不存在")
		response.ErrorJSON(c, appErr)
	}
}

// HandleNoMethod 返回 405 方法不允许的处理函数
// 当请求的路径存在但 HTTP 方法不匹配时，Gin 会调用此函数
// 例如：路由只注册了 POST，但客户端发送了 GET 请求
//
// 注意：需要设置 engine.HandleMethodNotAllowed = true 才会生效
func HandleNoMethod() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 使用 InvalidParams 而非专门的 MethodNotAllowed
		// 因为 405 本质上是"请求方式错误"，属于参数级别的错误
		appErr := apperrors.New(reason.InvalidParams, "请求方法不允许")
		appErr.HTTPCode = 405 // 覆盖默认的 400
		response.ErrorJSON(c, appErr)
	}
}
