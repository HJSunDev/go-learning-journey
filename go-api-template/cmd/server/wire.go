//go:build wireinject
// +build wireinject

// wireinject 构建标签告诉 Go 编译器：这个文件只在 Wire 代码生成时使用
// 正常编译时会被忽略（由 wire_gen.go 替代）

package main

import (
	"github.com/gin-gonic/gin"
	"github.com/google/wire"

	"go-api-template/internal/biz"
	"go-api-template/internal/data"
	"go-api-template/internal/server"
	"go-api-template/internal/service"
)

// wireApp 是 Wire 的注入器函数（Injector）
// Wire 会分析这个函数，根据 ProviderSet 中的构造函数自动生成依赖组装代码
//
// 函数签名说明：
// - 返回值：最终需要的依赖（*gin.Engine）和可能的错误
// - 函数体：调用 wire.Build 并传入所有 ProviderSet
// - panic：占位符，Wire 生成代码后会被替换
func wireApp() (*gin.Engine, error) {
	// wire.Build 声明所有需要的 Provider
	// Wire 会分析依赖关系，按正确顺序调用构造函数
	wire.Build(
		data.ProviderSet,    // Data -> GreeterRepo
		biz.ProviderSet,     // GreeterUsecase
		service.ProviderSet, // GreeterService
		server.ProviderSet,  // HTTPServer
	)

	// 占位返回，Wire 会替换整个函数体
	return nil, nil
}
