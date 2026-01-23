// Package data 是数据层，负责实现 biz 层定义的 Repository 接口。
// 当前阶段使用内存存储，未来可替换为真实数据库，而 biz 层无需任何改动。
package data

import (
	"log"
	"sync"

	"github.com/google/wire"

	"go-api-template/internal/conf"
)

// ProviderSet 聚合 data 层所有模块的 ProviderSet
// NewData 是基础设施（数据库连接等），单独列出
// 各模块的 Repository 在各自文件中定义 ProviderSet
var ProviderSet = wire.NewSet(
	NewData,            // 基础设施：数据库连接
	GreeterProviderSet, // Greeter 模块
	// UserProviderSet,  // 未来：用户模块
	// OrderProviderSet, // 未来：订单模块
)

// Data 是数据层的核心结构，持有所有数据连接和存储
// 阶段三：使用内存存储
// 阶段五：将替换为真实数据库连接（如 Ent Client）
type Data struct {
	// 配置信息
	cfg *conf.Config

	// 内存存储，使用 sync.Map 保证并发安全
	// 阶段五将替换为数据库客户端
	greeterStore *sync.Map
	// 用于生成自增 ID
	idCounter int64
	// 保护 idCounter 的互斥锁
	mu sync.Mutex
}

// NewData 创建并初始化 Data 实例
// cfg 提供数据库连接配置
// 阶段五：此处将初始化数据库连接
func NewData(cfg *conf.Config) (*Data, error) {
	// 记录数据库配置信息（不包含密码）
	log.Printf("Data layer initialized: driver=%s, host=%s, database=%s",
		cfg.Database.Driver, cfg.Database.Host, cfg.Database.Database)

	// 阶段五将在此处：
	// 1. 创建数据库连接
	// 2. 运行数据库迁移
	// 3. 返回包含数据库客户端的 Data 实例

	return &Data{
		cfg:          cfg,
		greeterStore: &sync.Map{},
		idCounter:    0,
	}, nil
}

// NextID 生成下一个自增 ID（并发安全）
func (d *Data) NextID() int64 {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.idCounter++
	return d.idCounter
}

// Close 关闭数据层资源（如数据库连接）
// 阶段五：此处将关闭数据库连接
func (d *Data) Close() error {
	log.Println("Data layer closed")
	return nil
}
