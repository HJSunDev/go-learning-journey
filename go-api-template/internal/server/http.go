package server

import (
	"net/http"

	"github.com/gin-gonic/gin"

	v1 "go-api-template/api/helloworld/v1"
	"go-api-template/internal/conf"
	"go-api-template/internal/pkg/apperrors"
	"go-api-template/internal/server/dto"
	"go-api-template/internal/server/middleware"
	"go-api-template/internal/server/response"
	"go-api-template/internal/service"

	// 导入生成的 Swagger 文档包（空导入，执行 init 函数注册规范）
	_ "go-api-template/internal/swagger"
)

// NewHTTPServer 创建并配置 HTTP 服务器
// cfg 提供服务器配置（端口、环境等）
// greeterSvc 是通过依赖注入传入的服务实例
func NewHTTPServer(cfg *conf.Config, greeterSvc *service.GreeterService) *HTTPServer {
	// 根据环境设置 Gin 模式
	setGinMode(cfg)

	// 使用 gin.New() 创建空白引擎，手动控制中间件
	// 不使用 gin.Default()，因为它内置的 Recovery 返回非 JSON 格式
	engine := gin.New()

	// 注册中间件（顺序重要）
	// 1. RequestID - 请求追踪
	// 2. Recovery - Panic 恢复，返回统一 JSON 格式
	// 3. Logger - 请求日志
	middleware.Register(engine)

	// 注册路由级别的错误处理（404、405）
	middleware.RegisterRouteHandlers(engine)

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

	// 注册 Swagger UI（非生产环境）
	registerSwagger(engine, cfg.App.Env)

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
// 使用 DTO 接收请求，Validator 自动验证，然后转换为 Proto 类型调用 Service
//
// @Summary      发送问候
// @Description  向指定用户发送问候消息，返回问候语和访问计数
// @Tags         greeter
// @Accept       json
// @Produce      json
// @Param        request body     dto.SayHelloRequest true "问候请求参数"
// @Success      200     {object} response.Response{data=v1.SayHelloResponse} "成功"
// @Failure      400     {object} response.Response "请求参数错误"
// @Failure      500     {object} response.Response "服务内部错误"
// @Router       /greeter/say-hello [post]
func handleSayHello(svc *service.GreeterService) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 使用 DTO 接收请求（DTO 有 binding tag，会自动验证）
		var req dto.SayHelloRequest

		// ShouldBindJSON 会：
		// 1. 把 JSON 填充到 req
		// 2. 根据 binding tag 验证（required, min=1, max=100）
		// 3. 验证失败返回错误,err!=nil 表示验证失败
		if err := c.ShouldBindJSON(&req); err != nil {
			// 使用统一响应：将 validator 错误转换为 AppError，再输出
			response.ErrorJSON(c, apperrors.FromValidationError(err))
			return
		}

		// DTO 转 Proto，调用 Service
		resp, err := svc.SayHello(c.Request.Context(), req.ToProto())
		if err != nil {
			// 使用统一响应：包装内部错误
			response.ErrorJSON(c, apperrors.Internal("服务处理失败", err))
			return
		}

		// 使用统一响应：成功响应
		// 最佳实践：直接传递结构体（DTO 或 Proto），避免手动构造 map
		response.SuccessJSON(c, resp)
	}
}

// handleSayHelloByPath 处理 GET 请求，name 从 URL 路径获取
//
// @Summary      发送问候（URL参数）
// @Description  通过 URL 路径参数向指定用户发送问候消息
// @Tags         greeter
// @Accept       json
// @Produce      json
// @Param        name path     string true "用户名称" minlength(1) maxlength(100)
// @Success      200  {object} response.Response{data=v1.SayHelloResponse} "成功"
// @Failure      400  {object} response.Response "请求参数错误"
// @Failure      500  {object} response.Response "服务内部错误"
// @Router       /greeter/say-hello/{name} [get]
func handleSayHelloByPath(svc *service.GreeterService) gin.HandlerFunc {
	return func(c *gin.Context) {
		name := c.Param("name")
		if name == "" {
			// 使用统一响应
			response.ErrorJSON(c, apperrors.InvalidParams("name 参数是必填的"))
			return
		}

		// 构造 proto 请求
		req := &v1.SayHelloRequest{Name: name}

		// 调用服务
		resp, err := svc.SayHello(c.Request.Context(), req)
		if err != nil {
			// 使用统一响应：包装内部错误
			response.ErrorJSON(c, apperrors.Internal("服务处理失败", err))
			return
		}

		// 使用统一响应：成功响应
		// 灵活用法：使用 response.Body (map[string]any) 构造临时数据
		response.SuccessJSON(c, response.Body{
			"message": resp.GetMessage(),
		})
	}
}
