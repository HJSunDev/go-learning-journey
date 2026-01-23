// Package server 是传输层，配置 HTTP/gRPC 服务器。
// 负责路由注册、中间件配置、请求/响应处理。
package server

import "github.com/google/wire"

// ProviderSet 是 server 层的依赖提供者集合
// server 层按协议划分（HTTP/gRPC），不按业务模块划分
var ProviderSet = wire.NewSet(
	NewHTTPServer,
	// NewGRPCServer, // 未来：gRPC 服务器
)
