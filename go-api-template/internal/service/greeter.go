package service

import (
	"context"

	v1 "go-api-template/api/helloworld/v1"
	"go-api-template/internal/biz"
)

// GreeterService 实现 proto 定义的 GreeterServiceServer 接口
type GreeterService struct {
	// 嵌入 UnimplementedGreeterServiceServer 以保持向前兼容
	v1.UnimplementedGreeterServiceServer

	// 依赖领域层的业务用例，而非数据层
	uc *biz.GreeterUsecase
}

// NewGreeterService 创建 GreeterService 实例
func NewGreeterService(uc *biz.GreeterUsecase) *GreeterService {
	return &GreeterService{uc: uc}
}

// SayHello 实现 GreeterServiceServer.SayHello 方法
// 职责：接收请求 -> 调用业务用例 -> 转换响应
func (s *GreeterService) SayHello(ctx context.Context, req *v1.SayHelloRequest) (*v1.SayHelloResponse, error) {
	// 调用业务用例执行核心逻辑
	greeter, err := s.uc.SayHello(ctx, req.GetName())
	if err != nil {
		return nil, err
	}

	// 将领域对象转换为 API 响应
	return &v1.SayHelloResponse{
		Message: greeter.Message,
	}, nil
}
