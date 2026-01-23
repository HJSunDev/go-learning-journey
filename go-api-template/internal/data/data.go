// Package data 是数据层，负责实现 biz 层定义的 Repository 接口。
// 当前阶段使用内存存储，未来可替换为真实数据库，而 biz 层无需任何改动。
package data

import (
	"sync"
)

// Data 是数据层的核心结构，持有所有数据连接和存储
// 阶段三：使用内存存储
// 阶段五：将替换为真实数据库连接
type Data struct {
	// 内存存储，使用 sync.Map 保证并发安全
	greeterStore *sync.Map
	// 用于生成自增 ID
	idCounter int64
	// 保护 idCounter 的互斥锁
	mu sync.Mutex
}

// NewData 创建并初始化 Data 实例
func NewData() (*Data, error) {
	return &Data{
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
