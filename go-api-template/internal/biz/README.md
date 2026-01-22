# biz - 领域层 (Domain Layer)

## 职责

这是整洁架构的**核心层**，包含：
- **领域实体 (Entities)**：核心业务对象的定义
- **业务接口 (Repository Interfaces)**：数据访问的抽象接口
- **业务逻辑 (Use Cases)**：核心业务规则的实现

## 依赖规则

**此层不依赖任何外部层。** 它是整个架构的中心。

- 不允许 import `internal/data`
- 不允许 import `internal/service`
- 不允许 import `internal/server`
- 不允许 import 任何数据库驱动 (如 gorm, ent)
- 不允许 import 任何 HTTP 框架 (如 gin)

## 允许的依赖

- 标准库
- 纯工具库 (如 uuid 生成器)
- 本项目的 `api` 层定义 (可选，用于复用 DTO)

## 示例结构

```
biz/
├── biz.go         # Wire Provider 注册
├── user.go        # User 实体 & UserRepo 接口
└── order.go       # Order 实体 & OrderRepo 接口
```
