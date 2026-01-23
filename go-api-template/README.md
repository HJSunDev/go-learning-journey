# Go API Template

基于整洁架构 (Clean Architecture) 和领域驱动设计 (DDD) 的 Go API 服务模板。

## 技术栈

- **HTTP 框架**: Gin
- **RPC 框架**: gRPC
- **API 定义**: Protobuf + Buf
- **依赖注入**: Google Wire
- **ORM**: Ent

## 快速启动

```bash
# 安装依赖
go mod tidy

# 运行服务
make run

# 访问健康检查
curl http://localhost:8080/health
```

## 目录结构

```
go-api-template/
├── api/              # Protobuf API 定义
├── cmd/server/       # 程序入口
├── configs/          # 配置文件
├── internal/         # 核心业务代码
│   ├── biz/          # 领域层（实体、仓储接口）
│   ├── data/         # 数据层（仓储实现）
│   ├── service/      # 应用层（API 实现）
│   └── server/       # 传输层（HTTP/gRPC 配置）
├── ent/              # Ent ORM 生成代码
├── third_party/      # 第三方 Proto 文件
└── docs/             # 项目文档
```

## 可用命令

```bash
make build        # 构建可执行文件
make run          # 运行服务
make clean        # 清理构建产物
make proto        # 生成 Proto 代码
make proto-lint   # Proto 文件 lint 检查
make proto-format # Proto 文件格式化
make tidy         # 整理依赖
make test         # 运行测试
```

## 文档

- [架构设计与目录职责](LEARNING.md)
- [文档索引](docs/README.md)

## License

MIT
