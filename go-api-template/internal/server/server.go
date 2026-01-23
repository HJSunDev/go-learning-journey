// Package server 是传输层，配置 HTTP/gRPC 服务器。
// 负责路由注册、中间件配置、请求/响应处理。
package server

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/google/wire"

	"go-api-template/internal/conf"
)

// ProviderSet 是 server 层的依赖提供者集合
// server 层按协议划分（HTTP/gRPC），不按业务模块划分
var ProviderSet = wire.NewSet(
	NewHTTPServer,
	// NewGRPCServer, // 未来：gRPC 服务器
)

// HTTPServer 封装 HTTP 服务器的配置和 Gin 引擎
type HTTPServer struct {
	engine *gin.Engine
	port   int
}

// Run 启动 HTTP 服务器
func (s *HTTPServer) Run() error {
	addr := fmt.Sprintf(":%d", s.port)
	return s.engine.Run(addr)
}

// Engine 返回底层的 Gin 引擎（用于测试等场景）
func (s *HTTPServer) Engine() *gin.Engine {
	return s.engine
}

// setGinMode 根据环境设置 Gin 模式
func setGinMode(cfg *conf.Config) {
	if cfg.IsProduction() {
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}
}
