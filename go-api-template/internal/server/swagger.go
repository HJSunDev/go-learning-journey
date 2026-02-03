// Package server Swagger UI 路由配置
package server

import (
	"net/http"

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

	// 缓存 Handler，避免每次请求都创建新实例
	swaggerHandler := ginSwagger.WrapHandler(swaggerFiles.Handler)

	// 注册 Swagger UI 路由
	// 包装 Handler 以处理根路径重定向
	// 当访问 /swagger/ 时，*any 为 "/"，gin-swagger 无法处理，手动重定向到 index.html
	engine.GET("/swagger/*any", func(c *gin.Context) {
		any := c.Param("any")
		if any == "/" {
			c.Redirect(http.StatusMovedPermanently, "/swagger/index.html")
			return
		}
		swaggerHandler(c)
	})

	// '/swagger'（无尾部斜杠）不匹配 /swagger/*any，需单独注册
	engine.GET("/swagger", func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, "/swagger/index.html")
	})
}
