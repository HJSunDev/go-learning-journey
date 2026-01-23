package main

import (
	"log"
)

func main() {
	// 使用 Wire 生成的 wireApp 函数初始化所有依赖
	// wireApp 定义在 wire.go，实现代码由 Wire 自动生成在 wire_gen.go
	httpServer, err := wireApp()
	if err != nil {
		log.Fatalf("Failed to initialize application: %v", err)
	}

	// ========================================
	// 启动服务
	// ========================================
	log.Println("Starting server on :8080")
	log.Println("API endpoints:")
	log.Println("  GET  /health                       - Health check")
	log.Println("  GET  /                             - Service info")
	log.Println("  POST /api/v1/greeter/say-hello     - Say hello (JSON body)")
	log.Println("  GET  /api/v1/greeter/say-hello/:name - Say hello (URL param)")

	if err := httpServer.Run(":8080"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
