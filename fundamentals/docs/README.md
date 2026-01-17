# 📚 Go Learning Journey (索引)

这里是项目学习文档的总入口。为了保持轻量和条理，所有的具体知识点都拆分成了独立的原子笔记。

## 🚀 阶段一：环境与启动 (起始)

> 早期记录均存放在根目录的 [LEARNING.md](../LEARNING.md) 中，作为永久存档。

- **[LEARNING.md](../LEARNING.md)**
  - Chapter 1: Windows Go 语言开发环境配置 (GOROOT, GOPATH, GOPROXY)
- **[002. Go 工程化：项目结构与工作区模式](notes/002-project-layout-and-workspace.md)**
  - go.mod / go.work 的区别与协作、目录规范、核心命令详解
- **[003. Hello World：第一个 Go 程序](notes/003-hello-world.md)**
  - package main、fmt 包、main 函数、go build、go run、路径前缀

## 🧩 阶段二：核心概念

> 新的笔记将自动添加到下方

- **[004. Go 类型系统：值类型与引用类型](notes/004-types-and-values.md)**
  - 基本类型、复合类型、零值机制、值类型 vs 引用类型、类型转换、自定义类型
- **[005. Go 控制流：条件判断与循环](notes/005-control-flow.md)**
  - if-else、switch、for 循环、for range、break/continue、defer 延迟执行
- **[006. Go 切片与数组：动态与固定的选择](notes/006-slices-and-arrays.md)**
  - 切片创建与操作、len/cap、append、slices 标准库、二维切片、数组使用场景
- **[007. Go 映射（Map）：键值对的艺术](notes/007-maps.md)**
  - Map 创建与操作、comma-ok 模式、键存在性检查、遍历、键类型要求、maps 标准库
- **[008. Go Range 迭代器：优雅的遍历之道](notes/008-range-and-strings.md)**
  - range 语法、切片/数组/映射/字符串遍历、strings 标准库、range 陷阱
- **[009. Go 函数：代码复用的基石](notes/009-functions.md)**
  - 函数定义与调用、参数传递、多返回值、可变参数、函数类型、高阶函数、递归、defer

## 📝 维护指南

- 所有的详细笔记存放在 `docs/notes/` 目录下。
- 命名格式：`SEQ-topic-name.md` (例如 `002-controller-basics.md`)。
- 每次新增笔记后，必须更新本文件的目录。
