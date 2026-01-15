# go-learning-journey

一个**分阶段、多模块单体仓库**，用于系统化学习 Go 语言开发。通过 `go.work` 管理独立模块。

## 架构设计

本仓库采用 **Workspace** (`go.work`) 模式，管理多个独立的 Go 模块：

```
go-journey/
├── fundamentals/      # 阶段一：语言基础 & 工程实践
├── web-service/       # 阶段二：框架级后端开发 
└── go.work            # 工作区定义，协调本地模块无缝交互
```

## 快速开始

### 环境要求

- Go 1.18+
- Windows/Linux/macOS

详细的环境配置指南见 [fundamentals/docs/](fundamentals/docs/README.md)。

## 文档

- **[基础阶段文档](fundamentals/docs/README.md)**：阶段一学习笔记与技术指南
- **[环境配置指南](fundamentals/LEARNING.md)**：Windows Go 开发环境完整配置

## 许可证

MIT
