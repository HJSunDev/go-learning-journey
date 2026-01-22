# 📚 Go API Template 文档索引

这里是项目学习文档的总入口。所有具体知识点都拆分成独立的章节笔记。

---

## 001. 架构与设计

> 项目的架构思想和设计决策存放在根目录的 [LEARNING.md](../LEARNING.md) 中。

- **[LEARNING.md](../LEARNING.md)**
  - 整洁架构 (Clean Architecture)、DDD、依赖倒置原则
  - 技术选型 (Gin, gRPC, Wire, Ent)
  - 目录结构与各层职责

---

## 阶段一：项目构建

> 从零开始，渐进式构建一个完整的 Go API 服务

### 一、骨架与传输层

- **[002. 渐进式构建策略](notes/002-progressive-build-strategy.md)**
  - 五阶段构建路线图
  - 目录结构创建、go mod 初始化
  - Gin HTTP Server 启动

### 二、API 定义

- **[003. Protobuf 与 Buf 工具链](notes/003-protobuf-and-buf-toolchain.md)**
  - Buf CLI、protoc-gen-go、protoc-gen-go-grpc 介绍
  - 开发环境 vs 部署环境
  - 配置文件详解与最佳实践

- **004. 编写 Proto 与生成代码** (待创建)

### 三、领域层与模拟数据

- **005. 领域实体与仓储模式** (待创建)

### 四、依赖注入

- **006. Wire 依赖注入实战** (待创建)

### 五、持久化

- **007. Ent ORM 入门** (待创建)

---

## 维护指南

- 所有详细笔记存放在 `docs/notes/` 目录下
- 命名格式：`SEQ-topic-name.md` (例如 `002-progressive-build-strategy.md`)
- 每次新增笔记后，必须更新本文件的目录
