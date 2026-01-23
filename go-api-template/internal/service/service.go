// Package service 是应用层，负责实现 API 接口定义。
// 这一层编排业务逻辑，进行 DTO 转换，但不包含核心业务规则。
// Package service 是应用层，实现 proto 定义的服务接口。
// 负责接收请求、调用业务用例、转换响应。
package service

import "github.com/google/wire"

// ProviderSet 聚合 service 层所有模块的 ProviderSet
var ProviderSet = wire.NewSet(
	GreeterProviderSet,
	// UserProviderSet,  // 未来：用户模块
	// OrderProviderSet, // 未来：订单模块
)
