// Go API Template - Swagger 全局信息配置
// swag 工具会扫描这些注释生成 OpenAPI 规范

// @title           Go API Template
// @version         1.0
// @description     一个基于整洁架构的 Go API 服务模板，用于学习和实践 Go Web 开发

// @contact.name   开发者
// @contact.email  dev@example.com

// @license.name  MIT
// @license.url   https://opensource.org/licenses/MIT

// @host      localhost:8080
// @BasePath  /api/v1

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description 输入格式: Bearer {token}

package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"

	"go-api-template/internal/conf"
)

// 命令行参数
var configPath string

func init() {
	// 支持通过命令行参数指定配置文件路径
	// 默认值为 configs/config.yaml（相对于项目根目录）
	flag.StringVar(&configPath, "config", "configs/config.yaml", "config file path")
}

func main() {
	flag.Parse()

	// ========================================
	// 加载配置
	// ========================================
	cfg, err := conf.LoadConfig(configPath)
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	log.Printf("Loaded config: env=%s, port=%d", cfg.App.Env, cfg.App.Port)

	// ========================================
	// 初始化应用
	// ========================================
	// 使用 Wire 生成的 wireApp 函数初始化所有依赖
	// wireApp 定义在 wire.go，实现代码由 Wire 自动生成在 wire_gen.go
	// wireApp 可以直接调用，因为它在同一个 package main 中
	httpServer, err := wireApp(cfg)
	if err != nil {
		log.Fatalf("Failed to initialize application: %v", err)
	}

	// ========================================
	// 启动服务
	// ========================================
	log.Printf("Starting server on %s", httpServer.Addr())
	log.Println("API endpoints:")
	log.Println("  GET  /health                         - Health check")
	log.Println("  GET  /                               - Service info")
	log.Println("  POST /api/v1/greeter/say-hello       - Say hello (JSON body)")
	log.Println("  GET  /api/v1/greeter/say-hello/:name - Say hello (URL param)")
	log.Println("  GET  /swagger/*                      - Swagger UI")

	// 启动 HTTP 服务器（非阻塞）
	errChan := httpServer.Start()

	// ========================================
	// 等待关闭信号
	// ========================================
	// 创建关闭服务信号通道
	quit := make(chan os.Signal, 1)
	// 告诉操作系统：当收到 SIGINT 或 SIGTERM 时，发送到 quit 管道
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// 阻塞等待：服务器错误 或 关闭信号
	select {
	case err := <-errChan:
		// 服务器启动失败或意外停止
		log.Fatalf("Server error: %v", err)
	case sig := <-quit:
		// 收到关闭信号
		log.Printf("Received signal: %v", sig)
	}

	// ========================================
	// 优雅关闭
	// ========================================
	log.Println("Shutting down server...")

	// 创建带超时的 context 用于优雅关闭
	// 超时时间从配置读取，确保不会无限等待
	shutdownTimeout := cfg.Server.GetShutdownTimeout()
	ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	// 执行优雅关闭
	// Shutdown 会：
	// 1. 停止接受新连接
	// 2. 等待正在处理的请求完成（或直到 context 超时）
	// 3. 关闭所有空闲连接
	if err := httpServer.Stop(ctx); err != nil {
		log.Printf("Server forced to shutdown: %v", err)
	} else {
		log.Println("Server gracefully stopped")
	}
}
