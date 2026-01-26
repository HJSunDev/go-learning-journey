package server

import (
	"net/http"

	"github.com/gin-gonic/gin"

	v1 "go-api-template/api/helloworld/v1"
	"go-api-template/internal/conf"
	"go-api-template/internal/service"
)

// NewHTTPServer 创建并配置 HTTP 服务器
// cfg 提供服务器配置（端口、环境等）
// greeterSvc 是通过依赖注入传入的服务实例
func NewHTTPServer(cfg *conf.Config, greeterSvc *service.GreeterService) *HTTPServer {
	// 根据环境设置 Gin 模式
	setGinMode(cfg)

	engine := gin.Default()

	// 健康检查端点
	engine.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
		})
	})

	// 服务信息端点
	engine.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"name":    cfg.App.Name,
			"version": "0.1.0",
			"env":     cfg.App.Env,
			"message": "Welcome to Go API Template",
		})
	})

	// 注册 Greeter 服务的 HTTP 路由
	registerGreeterRoutes(engine, greeterSvc)

	// 构建 http.Server（支持优雅关闭和超时配置）
	httpServer := buildHTTPServer(cfg, engine)

	return &HTTPServer{
		server: httpServer,
		engine: engine,
	}
}

// registerGreeterRoutes 注册 Greeter 服务的 HTTP 路由
// 将 gRPC 风格的服务暴露为 RESTful HTTP 端点
func registerGreeterRoutes(engine *gin.Engine, svc *service.GreeterService) {
	// API v1 路由组
	v1Group := engine.Group("/api/v1")
	{
		// POST /api/v1/greeter/say-hello
		// 请求体: {"name": "World"}
		// 响应体: {"message": "Hello, World! You are visitor #1."}
		v1Group.POST("/greeter/say-hello", handleSayHello(svc))

		// GET /api/v1/greeter/say-hello/:name
		// 便捷的 GET 端点，name 作为 URL 参数
		v1Group.GET("/greeter/say-hello/:name", handleSayHelloByPath(svc))
	}
}

// handleSayHello 处理 POST 请求
func handleSayHello(svc *service.GreeterService) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 使用 proto 生成的请求类型
		var req v1.SayHelloRequest

		// 绑定并验证请求
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "name is required",
			})
			return
		}

		if req.GetName() == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "name is required",
			})
			return
		}

		// 调用服务
		resp, err := svc.SayHello(c.Request.Context(), &req)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": resp.GetMessage(),
		})
	}
}

// handleSayHelloByPath 处理 GET 请求，name 从 URL 路径获取
func handleSayHelloByPath(svc *service.GreeterService) gin.HandlerFunc {
	return func(c *gin.Context) {
		name := c.Param("name")
		if name == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "name is required",
			})
			return
		}

		// 构造 proto 请求
		req := &v1.SayHelloRequest{Name: name}

		// 调用服务
		resp, err := svc.SayHello(c.Request.Context(), req)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": resp.GetMessage(),
		})
	}
}
