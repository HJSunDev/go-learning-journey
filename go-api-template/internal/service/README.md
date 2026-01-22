# service - 应用服务层 (Application Layer)

## 职责

这是整洁架构的**应用层**，负责：
- **实现 api (proto) 定义的服务接口**
- **DTO 转换**：将 API 请求转换为 biz 层能理解的格式
- **编排业务逻辑**：调用 biz 层的 Use Case
- **不包含业务规则**：业务规则属于 biz 层

## 依赖规则

此层**依赖 biz 层和 api 层**。

- 允许 import `internal/biz`
- 允许 import `api/` (proto 生成的代码)
- 不允许 import `internal/data`
- 不允许 import `internal/server`

## 依赖方向

```
service --> biz (调用业务逻辑)
service --> api (实现 proto 接口)
```

## 示例结构

```
service/
├── service.go     # Wire Provider 注册
├── user.go        # 实现 api.UserServiceServer
└── order.go       # 实现 api.OrderServiceServer
```

## 关键模式

service 层是 API 和业务逻辑之间的桥梁：

```go
// service/user.go
type UserService struct {
    userUseCase *biz.UserUseCase  // 注入 biz 层
}

func (s *UserService) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.GetUserResponse, error) {
    // 1. 调用 biz 层获取业务数据
    user, err := s.userUseCase.GetByID(ctx, req.Id)
    
    // 2. 转换为 API 响应格式
    return &pb.GetUserResponse{
        User: convertToProto(user),
    }, nil
}
```
