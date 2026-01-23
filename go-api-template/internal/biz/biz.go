// Package biz 是领域层，包含核心业务逻辑和领域模型。
// 这是整洁架构的核心，不依赖任何外部层。
package biz

import "github.com/google/wire"

// ProviderSet 聚合 biz 层所有模块的 ProviderSet
// 新增模块时，只需在对应文件定义 XxxProviderSet，然后添加到这里
var ProviderSet = wire.NewSet(
	GreeterProviderSet,
	// UserProviderSet,    // 未来：用户模块
	// OrderProviderSet,   // 未来：订单模块
	// ProductProviderSet, // 未来：商品模块
)
