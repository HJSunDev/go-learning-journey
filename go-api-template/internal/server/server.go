// Package server 是传输层，配置 HTTP/gRPC 服务器。
// 负责路由注册、中间件配置、请求/响应处理。
package server

import (
	"context"
	"errors"
	"fmt"
	"net/http"

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

// HTTPServer 封装 HTTP 服务器的配置和底层 http.Server
// 使用 http.Server 而非 gin.Engine.Run()，以支持优雅关闭
type HTTPServer struct {
	server *http.Server
	engine *gin.Engine
}

// Start 启动 HTTP 服务器（非阻塞）
// 服务器在独立的 goroutine 中运行，启动错误通过返回的 channel 传递
// 返回的 channel 在服务器停止时会收到错误（正常关闭时为 http.ErrServerClosed）
func (s *HTTPServer) Start() <-chan error {
	errChan := make(chan error, 1)
	go func() {
		// ListenAndServe 会阻塞直到服务器停止
		// 正常关闭时返回 http.ErrServerClosed
		if err := s.server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			errChan <- err
		}
		close(errChan)
	}()
	return errChan
}

// Stop 优雅关闭 HTTP 服务器
// 1. 停止接受新连接
// 2. 等待正在处理的请求完成（或直到 context 超时）
// 3. 关闭所有空闲连接
func (s *HTTPServer) Stop(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}

// Addr 返回服务器监听地址
func (s *HTTPServer) Addr() string {
	return s.server.Addr
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

// buildHTTPServer 构建 http.Server 实例
func buildHTTPServer(cfg *conf.Config, engine *gin.Engine) *http.Server {
	return &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.App.Port),
		Handler:      engine,
		ReadTimeout:  cfg.Server.GetReadTimeout(),
		WriteTimeout: cfg.Server.GetWriteTimeout(),
	}
}
