# 002b. Makefile 入门：项目命令统一入口

本章介绍 Makefile 是什么、为什么需要它、以及如何在本项目中使用它。

---

## 1. Makefile 是什么

### 1.1 定义

Makefile 是一个**命令快捷方式文件**，把一长串命令简化成一个短单词：

```
make run    =  go run ./cmd/server
make proto  =  buf lint api/ && buf generate api/
make build  =  go build -o bin/server ./cmd/server
```

### 1.2 为什么不直接敲命令

| 直接敲命令 | 使用 Makefile |
|------------|---------------|
| 需要记住完整命令 | 只需记住短单词 |
| 每次都要敲一长串 | `make xxx` 搞定 |
| 新人不知道该执行什么 | `make help` 列出所有可用命令 |
| 命令可能写错 | 命令固化在文件中，不会出错 |
| 多个命令要分开执行 | 可以组合成一个目标 |

---

## 2. 基本概念

### 2.1 Makefile 的结构

```makefile
# 这是注释

# 目标名:
目标名:
	命令1
	命令2
```

**注意**：命令前面必须是 **Tab 键**，不能是空格！

### 2.2 示例

```makefile
# 运行服务
run:
	go run ./cmd/server

# 构建二进制文件
build:
	go build -o bin/server ./cmd/server
```

执行方式：

```powershell
make run     # 执行 run 目标下的命令
make build   # 执行 build 目标下的命令
```

### 2.3 依赖关系

目标可以依赖其他目标：

```makefile
# proto 依赖 proto-lint，会先执行 proto-lint
proto: proto-lint
	buf generate api/

proto-lint:
	buf lint api/
```

执行 `make proto` 时，会先执行 `proto-lint`，再执行 `proto` 本身的命令。

### 2.4 .PHONY 声明

```makefile
.PHONY: all build run clean
```

`.PHONY` 声明这些目标是"伪目标"，不是真实的文件名。这样即使目录下有同名文件，`make` 也会正常执行命令。

---

## 3. 本项目的 Makefile

### 3.1 可用命令一览

| 命令 | 作用 | 阶段 |
|------|------|------|
| `make help` | 显示所有可用命令 | - |
| `make run` | 启动服务 | 阶段一+ |
| `make build` | 编译生成可执行文件 | 阶段一+ |
| `make clean` | 清理编译产物 | 阶段一+ |
| `make proto` | 生成 Proto 代码（含 lint） | 阶段二+ |
| `make proto-lint` | 检查 Proto 文件规范 | 阶段二+ |
| `make proto-format` | 格式化 Proto 文件 | 阶段二+ |
| `make tidy` | 整理 Go 依赖 | 阶段一+ |
| `make test` | 运行测试 | 阶段三+ |
| `make wire` | 生成依赖注入代码 | 阶段四 |

### 3.2 常用场景

#### 场景 1：日常开发

```powershell
# 启动服务进行调试
make run
```

#### 场景 2：修改了 Proto 文件

```powershell
# 重新生成代码
make proto
```

#### 场景 3：准备提交代码

```powershell
# 整理依赖 + 运行测试
make tidy
make test
```

#### 场景 4：构建发布版本

```powershell
# 编译二进制文件
make build

# 生成的文件在 bin/server
```

#### 场景 5：忘记有哪些命令

```powershell
make help
# 或直接
make
```

---

## 4. Windows 用户注意事项

### 4.1 安装 Make

Windows 默认没有 `make` 命令，需要安装：

**方法 1：通过 Chocolatey 安装**

```powershell
choco install make
```

**方法 2：通过 Scoop 安装**

```powershell
scoop install make
```

**方法 3：使用 Git Bash**

如果安装了 Git for Windows，Git Bash 自带 `make`。

### 4.2 验证安装

```powershell
make --version
# GNU Make 4.x.x
```

### 4.3 替代方案

如果不想安装 `make`，可以直接执行原始命令：

```powershell
# 代替 make run
go run ./cmd/server

# 代替 make proto
buf lint api/
buf generate api/

# 代替 make build
go build -o bin/server ./cmd/server
```

---

## 5. 完整的 Makefile 解读

```makefile
# 声明所有伪目标（避免与同名文件冲突）
.PHONY: all build run clean proto proto-lint proto-format wire help

# 默认目标：执行 make 不带参数时运行
all: help

# 构建可执行文件到 bin/ 目录
build:
	go build -o bin/server ./cmd/server

# 直接运行服务（开发时使用）
run:
	go run ./cmd/server

# 清理编译产物
clean:
	rm -rf bin/

# 生成 Proto 代码（先检查规范，再生成）
proto: proto-lint
	buf generate api/

# Proto 文件规范检查
proto-lint:
	buf lint api/

# Proto 文件格式化（-w 表示直接修改文件）
proto-format:
	buf format -w api/

# 生成 Wire 依赖注入代码（阶段四实现）
wire:
	@echo "Wire generation not configured yet. See Phase 4."

# 整理 Go 模块依赖
tidy:
	go mod tidy

# 运行所有测试
test:
	go test -v ./...

# 显示帮助信息（@ 表示不打印命令本身）
help:
	@echo "Available targets:"
	@echo "  build        - Build the server binary"
	@echo "  run          - Run the server"
	@echo "  clean        - Remove build artifacts"
	@echo "  proto        - Generate code from proto files (with lint check)"
	@echo "  proto-lint   - Lint check proto files"
	@echo "  proto-format - Format proto files"
	@echo "  wire         - Generate dependency injection code (Phase 4)"
	@echo "  tidy         - Run go mod tidy"
	@echo "  test         - Run tests"
```

---

## 6. 小结

| 概念 | 说明 |
|------|------|
| Makefile | 项目命令的统一入口，简化操作 |
| 目标 (target) | 快捷命令的名称，如 `run`、`build` |
| 依赖 | 目标可以依赖其他目标，按顺序执行 |
| .PHONY | 声明伪目标，避免与文件名冲突 |
| `make help` | 查看所有可用命令 |

**核心价值**：统一团队操作方式，降低沟通成本，新人看一眼 `make help` 就知道该怎么做。
