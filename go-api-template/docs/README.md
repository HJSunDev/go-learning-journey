# 📚 Go API Template 文档索引

这里是项目学习文档的总入口。所有具体知识点都拆分成独立的章节笔记。

## 阶段零：架构与设计 (起始)

> 项目的架构思想和设计决策存放在根目录的 [LEARNING.md](../LEARNING.md) 中。

- **[LEARNING.md](../LEARNING.md)**
  - Chapter 1: 架构思想与项目设计
  - Clean Architecture、DDD、依赖倒置原则
  - 技术选型 (Gin, gRPC, Wire, Ent)
  - 目录结构与各层职责

## 阶段一：项目构建

> 从零开始构建 API 服务的完整流程

- **[001. 渐进式构建策略](notes/001-progressive-build-strategy.md)**
  - 五阶段构建路线图
  - 阶段一：骨架与传输层搭建
  - 目录创建、go mod 初始化、Gin HTTP Server

## 阶段二：API 定义 (待完成)

> Schema First：使用 Protobuf 定义 API 契约

- **002. Protobuf 与 Buf 工具链** (待创建)

## 阶段三：领域层实现 (待完成)

> 整洁架构核心：实体、用例、仓储接口

- **003. 领域实体与仓储模式** (待创建)

## 阶段四：依赖注入 (待完成)

> 编译期 DI：Google Wire

- **004. Wire 依赖注入实战** (待创建)

## 阶段五：数据持久化 (待完成)

> ORM 选型与数据层实现

- **005. Ent ORM 入门** (待创建)

---

## 维护指南

- 所有详细笔记存放在 `docs/notes/` 目录下
- 命名格式：`SEQ-topic-name.md` (例如 `001-progressive-build-strategy.md`)
- 每次新增笔记后，必须更新本文件的目录
