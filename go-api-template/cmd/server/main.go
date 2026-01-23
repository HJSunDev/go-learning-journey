package main

import (
	"flag"
	"log"

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
	httpServer, err := wireApp(cfg)
	if err != nil {
		log.Fatalf("Failed to initialize application: %v", err)
	}

	// ========================================
	// 启动服务
	// ========================================
	log.Printf("Starting server on :%d", cfg.App.Port)
	log.Println("API endpoints:")
	log.Println("  GET  /health                       - Health check")
	log.Println("  GET  /                             - Service info")
	log.Println("  POST /api/v1/greeter/say-hello     - Say hello (JSON body)")
	log.Println("  GET  /api/v1/greeter/say-hello/:name - Say hello (URL param)")

	if err := httpServer.Run(); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
