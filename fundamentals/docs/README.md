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
- **[010. Go 闭包：自带状态的函数](notes/010-closures.md)**
  - 闭包定义、变量捕获机制、内存管理与销毁（GC）、生成器模式、状态隔离
- **[011. Go 指针：直接操作内存的钥匙](notes/011-pointers.md)**
  - 指针基础（& 和 *）、值传递 vs 指针传递、修改结构体、引用类型特殊性、nil 指针
- **[012. Go 结构体：构建自定义类型的基石](notes/012-structs.md)**
  - 结构体定义与创建、方法与接收器、值接收器 vs 指针接收器、构造函数模式、匿名结构体、结构体组合（嵌入）
- **[013. Go 接口：定义行为的契约](notes/013-interfaces.md)**
  - 接口定义与隐式实现、io.Writer 经典案例、接口值内部结构、类型断言与类型选择、空接口 any、接口组合、设计原则
- **[014. Go 枚举：用 const 和 iota 构建类型安全的常量集合](notes/014-enums.md)**
  - const 常量定义、iota 自动递增、自定义类型增加类型安全、枚举方法（String/业务逻辑）、iota 进阶用法
- **[015. Go 泛型（Generics）：一次编写，多种类型](notes/015-generics.md)**
  - 泛型函数、类型参数与约束、comparable/any/自定义约束、泛型类型、Filter/Map/Reduce 集合操作

## 🚀 阶段三：并发编程

> 协程、通道与并发模式

- **[016. Go 协程（Goroutine）：轻量级并发的基石](notes/016-goroutines.md)**
  - 协程本质与优势、go 关键字、sync.WaitGroup 等待组、闭包变量捕获陷阱、sync.Mutex 互斥锁
- **[017. Go 通道（Channel）：协程间通信的桥梁](notes/017-channels.md)**
  - 无缓冲/有缓冲通道、发送/接收/关闭、for range 遍历、select 多路复用、死锁、工作池模式
- **[018. Go 互斥锁（Mutex）：并发安全的守护者](notes/018-mutex.md)**
  - 数据竞争问题、锁的概念（悲观锁/乐观锁）、Mutex 与 RWMutex、死锁陷阱、临界区最小化

## 📝 维护指南

- 所有的详细笔记存放在 `docs/notes/` 目录下。
- 命名格式：`SEQ-topic-name.md` (例如 `002-controller-basics.md`)。
- 每次新增笔记后，必须更新本文件的目录。
