# server - 传输层 (Transport Layer)

## 职责

这是整洁架构的**传输层/适配器层**，负责：
- **HTTP Server 配置与启动** (Gin)
- **gRPC Server 配置与启动**
- **路由注册**
- **中间件配置** (日志、认证、限流等)

## 依赖规则

此层**依赖 service 层**，将 service 注册到具体的传输协议。

- 允许 import `internal/service`
- 允许 import HTTP 框架 (gin)
- 允许 import gRPC 库
- 不允许 import `internal/biz`
- 不允许 import `internal/data`

## 依赖方向

```
server --> service (注册服务)
```

## 示例结构

```
server/
├── server.go      # Wire Provider 注册
├── http.go        # Gin HTTP Server 配置
└── grpc.go        # gRPC Server 配置
```

## 关键模式

server 层负责启动服务器并注册路由：

```go
// server/http.go
func NewHTTPServer(userSvc *service.UserService) *gin.Engine {
    r := gin.Default()
    
    // 注册路由
    r.GET("/users/:id", userSvc.GetUser)
    
    return r
}
```
