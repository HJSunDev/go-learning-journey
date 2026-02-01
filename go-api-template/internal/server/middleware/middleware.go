package middleware

import (
	"github.com/gin-gonic/gin"
)

// middlewareChain 定义中间件链
// 顺序很重要，遵循"洋葱模型"：
//
//	请求进入 → RequestID → Recovery → Logger → Handler
//	响应返回 ← RequestID ← Recovery ← Logger ← Handler
//
// 使用切片声明的优势：
//  1. 顺序一目了然，修改只需调整数组
//  2. 符合声明式编程风格
//  3. 避免多次调用 engine.Use() 的冗余
var middlewareChain = []gin.HandlerFunc{
	RequestID(),  // [0] 最先执行，确保后续中间件都能获取请求 ID
	Recovery(),   // [1] 捕获后续所有代码的 panic
	gin.Logger(), // [2] 记录请求日志
}

// Register 注册所有中间件到 Gin 引擎
func Register(engine *gin.Engine) {
	engine.Use(middlewareChain...)
}

// RegisterRouteHandlers 注册路由级别的错误处理
// 包括 404 和 405 处理，确保这些错误也返回统一格式
func RegisterRouteHandlers(engine *gin.Engine) {
	// 处理 404：路由不存在
	engine.NoRoute(HandleNoRoute())

	// 处理 405：方法不允许
	// 需要先启用 HandleMethodNotAllowed
	engine.HandleMethodNotAllowed = true
	engine.NoMethod(HandleNoMethod())
}
