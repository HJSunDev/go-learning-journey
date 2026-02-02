# 📚 Go API Template 文档索引

这里是项目学习文档的总入口。所有具体知识点都拆分成独立的章节笔记。

## 🚀001. 架构与设计

> 项目的架构思想和设计决策存放在根目录的 [LEARNING.md](../LEARNING.md) 中。

- **[LEARNING.md](../LEARNING.md)**
  - 整洁架构 (Clean Architecture)、DDD、依赖倒置原则；技术选型 (Gin, gRPC, Wire, Ent)；目录结构与各层职责

## 🧩阶段一：项目构建

> 从零开始，渐进式构建一个完整的 Go API 服务

- **[002. 渐进式构建策略与结构搭建](notes/002-progressive-build-strategy.md)**
  - 五阶段构建路线图；目录结构创建、go mod 初始化；Gin HTTP Server 启动
- **[002b. Makefile 入门](notes/002b-makefile-introduction.md)**
  - Makefile 是什么、为什么需要它；本项目可用命令一览；Windows 用户安装指南
- **[002c. Gin 框架入门](notes/002c-gin-framework-introduction.md)**
  - Gin 本质与核心概念；gin.Engine、gin.Context、gin.HandlerFunc 详解；路由、路由组、响应输出；gin.Engine 与 http.Server 的关系（为何需要两者）；与整洁架构的关系
- **[003. Protobuf 与 Buf 工具链](notes/003-protobuf-and-buf-toolchain.md)**
  - Buf CLI、protoc-gen-go、protoc-gen-go-grpc 介绍；开发环境 vs 部署环境；配置文件详解与最佳实践
- **[004. Protobuf 入门](notes/004-protobuf-introduction.md)**
  - Proto 文件是什么、为什么需要它；Proto 语法详解（message、service、类型系统）；生成的两个文件（*.pb.go 和 *_grpc.pb.go）的作用；使用场景与开发工作流程
- **[005. 编写 Proto 与生成代码](notes/005-proto-and-code-generation.md)**
  - Proto 文件编写与命名规范；Buf 配置与代码生成；Schema First 开发模式
- **[006. 领域层与模拟数据](notes/006-domain-layer.md)**
  - 依赖倒置原则 (DIP) 实践；领域实体与 Repository 接口设计；内存存储实现；四层架构组装
- **[007. Wire 依赖注入实战](notes/007-wire-dependency-injection.md)**
  - Google Wire 编译期依赖注入；ProviderSet 与 Injector 模式；自动生成依赖组装代码
- **[008. 配置模块](notes/008-configuration-module.md)**
  - Viper 配置库使用；配置结构体设计；环境变量覆盖机制；生产环境部署策略
- **[009. 优雅关闭](notes/009-graceful-shutdown.md)**
  - http.Server.Shutdown 机制；信号监听（SIGINT/SIGTERM）；超时控制；生产环境注意事项
- **[010. Struct Tag 与 Validator 请求验证](notes/010-struct-tag-and-validator.md)**
  - Struct Tag 是什么；json Tag 详解；binding Tag 与 Validator 库；完整验证规则大全；处理验证错误
- **[011. DTO 模式与请求验证实践](notes/011-dto-pattern.md)**
  - DTO 是什么；DTO 与 Proto 的关系和协作；开发流程；在项目中实践自动验证
- **[012. 统一错误处理](notes/012-unified-error-handling.md)**
  - AppError 类型设计；错误码体系；验证错误转换；友好错误消息
- **[013. 统一响应格式](notes/013-unified-response.md)**
  - Response 结构设计；成功/错误响应统一；SuccessJSON/ErrorJSON 工具函数
- **[014. 中间件](notes/014-middleware.md)**
  - 中间件概念与洋葱模型；RequestID 请求追踪；Recovery Panic 恢复；404/405 统一处理
- **[015. Swagger 接口文档集成](notes/015-swagger-api-documentation.md)**
  - swag + gin-swagger 工具链；全局注释与 API 注释语法；Swagger UI 配置；环境隔离最佳实践
- **016. Ent ORM 入门** (待创建)

---

## 📝维护指南

- 所有详细笔记存放在 `docs/notes/` 目录下
- 命名格式：`SEQ-topic-name.md` (例如 `002-progressive-build-strategy.md`)
- 每次新增笔记后，必须更新本文件的目录
