# 003. Protobuf 与 Buf 工具链

本章介绍第二阶段所需的开发工具：Buf CLI、protoc-gen-go、protoc-gen-go-grpc。

---

## 1. 工具链概述

第二阶段「API 定义」需要三个核心工具，它们协同工作完成从 `.proto` 文件到 Go 代码的生成流程。

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                              开发工作流                                      │
│                                                                             │
│   .proto 文件  ──────►  Buf CLI  ──────►  生成的 Go 代码                     │
│   (你编写的)           (协调器)           (自动生成)                          │
│                           │                                                 │
│                           ├──► protoc-gen-go      ──► xxx.pb.go             │
│                           │    (生成消息结构体)                              │
│                           │                                                 │
│                           └──► protoc-gen-go-grpc ──► xxx_grpc.pb.go        │
│                                (生成服务接口)                                │
└─────────────────────────────────────────────────────────────────────────────┘
```

### 1.1 三个工具的职责

| 工具               | 职责                                          | 生成产物             |
| ------------------ | --------------------------------------------- | -------------------- |
| Buf CLI            | 管理 proto 文件、执行 lint 检查、协调代码生成 | 无直接产物，是指挥者 |
| protoc-gen-go      | 将 proto 中的 `message` 转换为 Go 结构体    | `*.pb.go`          |
| protoc-gen-go-grpc | 将 proto 中的 `service` 转换为 Go 接口      | `*_grpc.pb.go`     |

### 1.2 开发环境 vs 部署环境

这三个工具都是**开发时工具**，只在开发机器上需要安装。

| 环境         | 是否需要这些工具 | 原因                                     |
| ------------ | ---------------- | ---------------------------------------- |
| 开发环境     | ✅ 需要          | 编写 proto 文件后需要生成 Go 代码        |
| CI/CD 环境   | ✅ 需要          | 构建流水线中需要重新生成代码以确保一致性 |
| 生产部署环境 | ❌ 不需要        | 只运行编译好的二进制文件，不涉及代码生成 |

**类比理解**：这些工具类似于 TypeScript 的 `tsc` 编译器——开发时需要用它把 `.ts` 编译成 `.js`，但运行时只执行 `.js`，不需要 `tsc`。

---

## 2. Buf CLI

### 2.1 是什么

Buf 是 Protobuf 生态的现代化工具，由 Buf Technologies 公司开发。它替代了传统的 `protoc` 命令行工具，解决了以下痛点：

| 传统 protoc 的问题        | Buf 的解决方案            |
| ------------------------- | ------------------------- |
| 命令行参数冗长复杂        | 使用 YAML 配置文件        |
| 手动管理第三方 proto 依赖 | 内置依赖管理（类似 npm）  |
| 无内置 lint 检查          | 提供三级 lint 规则        |
| 无法检测破坏性变更        | 内置 breaking change 检测 |
| 生成速度慢                | 并行生成，速度更快        |

### 2.2 核心功能

```
buf lint      # 检查 proto 文件是否符合规范
buf format    # 格式化 proto 文件
buf generate  # 生成代码（Go, TypeScript, Python 等）
buf build     # 编译 proto 文件，验证语法正确性
buf breaking  # 检测是否有破坏性的 API 变更
```

### 2.3 安装

```powershell
go install github.com/bufbuild/buf/cmd/buf@latest
```

验证安装：

```powershell
buf --version
# 预期输出类似: 1.28.1
```

### 2.4 配置文件

Buf 需要两个配置文件：

**`buf.yaml`** - 放在 proto 文件目录，定义模块信息和 lint 规则

```yaml
# api/buf.yaml
version: v2
modules:
  - path: helloworld/v1
lint:
  use:
    - STANDARD
breaking:
  use:
    - FILE
```

**`buf.gen.yaml`** - 放在项目根目录，定义代码生成规则

```yaml
# buf.gen.yaml
version: v2
plugins:
  # 生成 Go 消息结构体
  - local: protoc-gen-go
    out: api
    opt:
      - paths=source_relative
  # 生成 gRPC 服务接口
  - local: protoc-gen-go-grpc
    out: api
    opt:
      - paths=source_relative
inputs:
  - directory: api
```

### 2.5 常用命令详解

#### 生成代码

```powershell
# 在项目根目录执行
buf generate
```

此命令读取 `buf.gen.yaml`，找到所有 `.proto` 文件，调用配置的插件生成代码。

#### Lint 检查

```powershell
buf lint api/
```

检查 proto 文件是否符合最佳实践，例如：

- 包名是否正确
- 字段命名是否使用 snake_case
- 服务名是否以 Service 结尾

#### 格式化

```powershell
buf format -w api/
```

`-w` 参数表示直接修改文件（原地格式化）。

---

## 3. protoc-gen-go

### 3.1 是什么

`protoc-gen-go` 是 Google 官方维护的 Protobuf Go 代码生成插件。它的唯一职责是将 `.proto` 文件中的 `message` 定义转换为 Go 结构体。

### 3.2 生成产物示例

**输入（greeter.proto）：**

```protobuf
message HelloRequest {
  string name = 1;
}

message HelloReply {
  string message = 1;
}
```

**输出（greeter.pb.go）：**

```go
type HelloRequest struct {
    state         protoimpl.MessageState
    sizeCache     protoimpl.SizeCache
    unknownFields protoimpl.UnknownFields

    Name string `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
}

type HelloReply struct {
    state         protoimpl.MessageState
    sizeCache     protoimpl.SizeCache
    unknownFields protoimpl.UnknownFields

    Message string `protobuf:"bytes,1,opt,name=message,proto3" json:"message,omitempty"`
}

// 以及序列化/反序列化方法...
```

### 3.3 安装

```powershell
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
```

验证安装：

```powershell
protoc-gen-go --version
# 预期输出类似: protoc-gen-go v1.36.11
```

### 3.4 关键配置选项

| 选项                      | 作用                                      | 推荐值       |
| ------------------------- | ----------------------------------------- | ------------ |
| `paths=source_relative` | 生成的文件放在与 proto 文件相同的相对路径 | ✅ 推荐      |
| `paths=import`          | 根据 Go import 路径决定输出位置           | 复杂项目使用 |

在 `buf.gen.yaml` 中配置：

```yaml
plugins:
  - local: protoc-gen-go
    out: api
    opt:
      - paths=source_relative
```

---

## 4. protoc-gen-go-grpc

### 4.1 是什么

`protoc-gen-go-grpc` 是 gRPC 团队维护的 Go 代码生成插件。它的职责是将 `.proto` 文件中的 `service` 定义转换为 Go 接口。

### 4.2 生成产物示例

**输入（greeter.proto）：**

```protobuf
service Greeter {
  rpc SayHello (HelloRequest) returns (HelloReply);
}
```

**输出（greeter_grpc.pb.go）：**

```go
// 客户端接口
type GreeterClient interface {
    SayHello(ctx context.Context, in *HelloRequest, opts ...grpc.CallOption) (*HelloReply, error)
}

// 服务端接口 - 你需要实现这个接口
type GreeterServer interface {
    SayHello(context.Context, *HelloRequest) (*HelloReply, error)
    mustEmbedUnimplementedGreeterServer()
}

// 未实现的基础结构 - 用于向前兼容
type UnimplementedGreeterServer struct {}

func (UnimplementedGreeterServer) SayHello(context.Context, *HelloRequest) (*HelloReply, error) {
    return nil, status.Errorf(codes.Unimplemented, "method SayHello not implemented")
}
```

### 4.3 安装

```powershell
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
```

验证安装：

```powershell
protoc-gen-go-grpc --version
# 预期输出类似: protoc-gen-go-grpc 1.6.0
```

### 4.4 关于 UnimplementedServer

生成的代码中有一个 `UnimplementedGreeterServer` 结构体，这是 gRPC 的**向前兼容机制**：

```go
type greeterServer struct {
    // 必须嵌入这个结构体
    v1.UnimplementedGreeterServer
}
```

**为什么需要嵌入？**

当 proto 文件新增一个 RPC 方法时，如果你的实现没有更新，嵌入的 `UnimplementedServer` 会返回 "未实现" 错误，而不是编译失败。这让你可以渐进式地实现新方法。

---

## 5. 工具链协作流程

### 5.1 完整的代码生成流程

```
步骤 1：编写 proto 文件
         │
         ▼
步骤 2：执行 buf generate
         │
         ├──► Buf 读取 buf.gen.yaml 配置
         │
         ├──► Buf 解析 api/ 目录下的所有 .proto 文件
         │
         ├──► Buf 调用 protoc-gen-go 生成 *.pb.go
         │
         └──► Buf 调用 protoc-gen-go-grpc 生成 *_grpc.pb.go
         │
         ▼
步骤 3：生成的 Go 文件出现在 api/ 目录
```

### 5.2 文件对应关系

```
api/helloworld/v1/
├── greeter.proto          ← 你编写的 API 定义
├── greeter.pb.go          ← protoc-gen-go 生成的消息结构体
└── greeter_grpc.pb.go     ← protoc-gen-go-grpc 生成的服务接口
```

---

## 6. 最佳实践

### 6.1 Proto 文件组织

```
api/
├── buf.yaml                    # Buf 模块配置
├── helloworld/v1/              # 模块名/版本号
│   └── greeter.proto
└── user/v1/                    # 另一个模块
    └── user.proto
```

**规则**：

- 每个业务模块一个目录
- 目录名使用 `lower_snake_case`
- 始终包含版本号（v1, v2）

### 6.2 Proto 文件命名规范

| 元素       | 命名规则                | 示例                   |
| ---------- | ----------------------- | ---------------------- |
| 文件名     | lower_snake_case.proto  | `user_service.proto` |
| package    | lower_snake_case + 版本 | `helloworld.v1`      |
| message    | PascalCase              | `HelloRequest`       |
| field      | lower_snake_case        | `user_name`          |
| service    | PascalCase + Service    | `GreeterService`     |
| rpc        | PascalCase              | `SayHello`           |
| enum       | PascalCase              | `UserStatus`         |
| enum value | UPPER_SNAKE_CASE        | `USER_STATUS_ACTIVE` |

### 6.3 请求/响应命名

每个 RPC 方法使用独立的请求和响应类型：

```protobuf
// ✅ 推荐：每个方法独立的请求/响应
service GreeterService {
  rpc SayHello (SayHelloRequest) returns (SayHelloResponse);
  rpc SayGoodbye (SayGoodbyeRequest) returns (SayGoodbyeResponse);
}

message SayHelloRequest {
  string name = 1;
}

message SayHelloResponse {
  string message = 1;
}

// ❌ 避免：复用通用类型
service GreeterService {
  rpc SayHello (google.protobuf.Empty) returns (StringValue);  // 不推荐
}
```

### 6.4 生成的代码不要手动修改

```
// Code generated by protoc-gen-go. DO NOT EDIT.
```

生成的 `*.pb.go` 和 `*_grpc.pb.go` 文件顶部都有这行注释。任何手动修改都会在下次 `buf generate` 时被覆盖。

如果需要扩展功能，在**单独的文件**中添加方法：

```go
// greeter_ext.go - 你的扩展代码
package v1

// 为生成的类型添加辅助方法
func (r *HelloRequest) Validate() error {
    if r.Name == "" {
        return errors.New("name is required")
    }
    return nil
}
```

### 6.5 版本控制策略

| 文件类型         | 是否提交到 Git | 原因                                 |
| ---------------- | -------------- | ------------------------------------ |
| `*.proto`      | ✅ 必须提交    | 这是源文件                           |
| `buf.yaml`     | ✅ 必须提交    | 配置文件                             |
| `buf.gen.yaml` | ✅ 必须提交    | 配置文件                             |
| `*.pb.go`      | ⚠️ 建议提交  | 方便其他开发者无需安装工具链即可构建 |
| `*_grpc.pb.go` | ⚠️ 建议提交  | 同上                                 |

**说明**：虽然生成的代码可以通过 `buf generate` 重新生成，但提交到 Git 可以：

- 让没有安装 Buf 工具链的开发者也能直接 `go build`
- 在 Code Review 中看到 API 变更的实际影响

---

## 7. 环境配置检查清单

在开始第二阶段实践前，确认以下环境已就绪：

```powershell
# 检查 Go 版本
go version
# 预期: go version go1.21+ ...

# 检查 Buf
buf --version
# 预期: 1.28.1 或更高

# 检查 protoc-gen-go
protoc-gen-go --version
# 预期: protoc-gen-go v1.36.x

# 检查 protoc-gen-go-grpc
protoc-gen-go-grpc --version
# 预期: protoc-gen-go-grpc 1.6.x

# 检查 GOPATH/bin 是否在 PATH 中
# Windows PowerShell
$env:PATH -split ';' | Select-String 'go'
# 应该能看到类似 C:\Users\xxx\go\bin 的路径
```

如果 `protoc-gen-go` 或 `protoc-gen-go-grpc` 命令找不到，说明 `$GOPATH/bin` 没有加入系统 PATH。

**Windows 添加 PATH 的方法**：

```powershell
# 临时添加（仅当前会话有效）
$env:PATH += ";$(go env GOPATH)\bin"

# 永久添加（需要管理员权限，或通过系统设置手动添加）
# 路径通常是: C:\Users\你的用户名\go\bin
```

---

## 8. 小结

| 工具               | 来源             | 职责                       | 安装命令                                                            |
| ------------------ | ---------------- | -------------------------- | ------------------------------------------------------------------- |
| Buf CLI            | Buf Technologies | 管理 proto、lint、生成协调 | `go install github.com/bufbuild/buf/cmd/buf@latest`               |
| protoc-gen-go      | Google           | 生成消息结构体             | `go install google.golang.org/protobuf/cmd/protoc-gen-go@latest`  |
| protoc-gen-go-grpc | gRPC Team        | 生成服务接口               | `go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest` |

**核心理解**：

- 这三个工具只在开发/构建阶段使用，生产环境不需要
- Buf 是指挥官，protoc-gen-go 和 protoc-gen-go-grpc 是执行者
- 生成的代码不要手动修改，如需扩展请创建单独的文件
