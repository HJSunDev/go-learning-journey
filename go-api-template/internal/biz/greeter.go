package biz

import (
	"context"
	"fmt"
	"time"

	"github.com/google/wire"
)

// GreeterProviderSet 是 Greeter 模块的依赖提供者集合
var GreeterProviderSet = wire.NewSet(NewGreeterUsecase)

// Greeter 是领域实体，表示一条问候记录
type Greeter struct {
	ID        int64     // 唯一标识
	Name      string    // 被问候者名称
	Message   string    // 问候消息
	CreatedAt time.Time // 创建时间
}

// GreeterRepo 定义了问候数据的存储接口
// 这是依赖倒置的关键：接口定义在领域层，实现在数据层
type GreeterRepo interface {
	// Save 保存一条问候记录
	Save(ctx context.Context, g *Greeter) (*Greeter, error)
	// GetByName 根据名称获取最近的问候记录
	GetByName(ctx context.Context, name string) (*Greeter, error)
	// Count 获取问候总数
	Count(ctx context.Context) (int64, error)
}

// GreeterUsecase 是问候业务用例，包含核心业务逻辑
type GreeterUsecase struct {
	repo GreeterRepo
}

// NewGreeterUsecase 创建 GreeterUsecase 实例
// repo 参数通过依赖注入传入，Usecase 不知道也不关心具体实现
func NewGreeterUsecase(repo GreeterRepo) *GreeterUsecase {
	return &GreeterUsecase{repo: repo}
}

// SayHello 执行问候业务逻辑
// 核心逻辑：创建问候记录并返回个性化消息
func (uc *GreeterUsecase) SayHello(ctx context.Context, name string) (*Greeter, error) {
	// 获取当前问候总数，用于生成个性化消息
	count, err := uc.repo.Count(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get count: %w", err)
	}

	// 构建问候消息,sprintf用于格式化字符串，不输出到控制台，而是返回一个字符串
	message := fmt.Sprintf("Hello, %s! You are visitor #%d.", name, count+1)

	// 创建问候记录,greeter 是问候记录的结构体，且没有 ID
	greeter := &Greeter{
		Name:      name,
		Message:   message,
		CreatedAt: time.Now(),
	}

	// 保存到存储,saved 是保存后的问候记录，就是 greeter 的副本，但是有 ID
	saved, err := uc.repo.Save(ctx, greeter)
	if err != nil {
		return nil, fmt.Errorf("failed to save greeter: %w", err)
	}

	return saved, nil
}
