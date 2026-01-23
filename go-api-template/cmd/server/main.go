package main

import (
	"log"

	"go-api-template/internal/biz"
	"go-api-template/internal/data"
	"go-api-template/internal/server"
	"go-api-template/internal/service"
)

func main() {
	// ========================================
	// 手动依赖注入（阶段四将使用 Wire 自动化）
	// ========================================
	// 依赖组装顺序遵循整洁架构的依赖流向：
	// Data -> Repo -> Usecase -> Service -> Server

	// 1. 初始化数据层
	// 当前使用内存存储，阶段五将替换为真实数据库
	dataLayer, err := data.NewData()
	if err != nil {
		log.Fatalf("Failed to create data layer: %v", err)
	}

	// 2. 创建 Repository（数据层实现 biz 层定义的接口）
	greeterRepo := data.NewGreeterRepo(dataLayer)

	// 3. 创建业务用例（领域层，依赖 Repository 接口）
	greeterUsecase := biz.NewGreeterUsecase(greeterRepo)

	// 4. 创建服务（应用层，实现 proto 定义的接口）
	greeterService := service.NewGreeterService(greeterUsecase)

	// 5. 创建 HTTP 服务器（传输层）
	httpServer := server.NewHTTPServer(greeterService)

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
