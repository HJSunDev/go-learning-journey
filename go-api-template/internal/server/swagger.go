// Package server Swagger UI 路由配置
package server

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// registerSwagger 注册 Swagger UI 路由
// 仅在非生产环境启用，避免暴露 API 文档给外部
func registerSwagger(engine *gin.Engine, env string) {
	// 生产环境不启用 Swagger UI
	if env == "production" {
		return
	}

	// 注册 Swagger UI 路由
	// 访问 /swagger/index.html 可查看 API 文档
	engine.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// 访问 /swagger 时自动重定向到 /swagger/index.html
	engine.GET("/swagger", func(c *gin.Context) {
		c.Redirect(301, "/swagger/index.html")
	})
}
