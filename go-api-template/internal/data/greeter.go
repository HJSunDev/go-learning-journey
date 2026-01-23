package data

import (
	"context"
	"sync/atomic"

	"go-api-template/internal/biz"
)

// greeterRepo 实现 biz.GreeterRepo 接口
// 当前使用内存 Map 存储，未来可替换为数据库实现
type greeterRepo struct {
	data *Data
}

// NewGreeterRepo 创建 GreeterRepo 实例
// 返回接口类型，隐藏实现细节
func NewGreeterRepo(data *Data) biz.GreeterRepo {
	return &greeterRepo{data: data}
}

// Save 保存问候记录到内存存储
func (r *greeterRepo) Save(ctx context.Context, g *biz.Greeter) (*biz.Greeter, error) {
	// 分配新 ID
	g.ID = r.data.NextID()

	// 存储到 sync.Map，使用 ID 作为 key
	r.data.greeterStore.Store(g.ID, g)

	// 同时以 name 为 key 存储，便于按名称查询
	r.data.greeterStore.Store("name:"+g.Name, g)

	return g, nil
}

// GetByName 根据名称获取最近的问候记录
func (r *greeterRepo) GetByName(ctx context.Context, name string) (*biz.Greeter, error) {
	value, ok := r.data.greeterStore.Load("name:" + name)
	if !ok {
		return nil, nil
	}
	return value.(*biz.Greeter), nil
}

// Count 获取问候记录总数
func (r *greeterRepo) Count(ctx context.Context) (int64, error) {
	// idCounter 就是当前的记录总数
	return atomic.LoadInt64(&r.data.idCounter), nil
}
