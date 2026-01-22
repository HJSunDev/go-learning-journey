# 005. 编写 Proto 与生成代码

本章记录阶段二的完整操作步骤：编写 `.proto` 文件并使用 Buf 生成 Go 代码。

---

## 1. 阶段目标

- 编写 Protobuf 文件定义 API 接口
- 配置 Buf 工具链
- 自动生成 Go 结构体和 gRPC 服务接口
- 理解 Schema First 开发模式

---

## 2. 操作步骤

### 2.1 创建 Proto 文件

在 `api/helloworld/v1/` 目录下创建 `greeter.proto` 文件：

```protobuf
// api/helloworld/v1/greeter.proto

// API 接口定义：Greeter 服务
// 这是 Schema First 开发模式的核心，先定义接口契约，再编写实现代码

syntax = "proto3";

package helloworld.v1;

option go_package = "go-api-template/api/helloworld/v1;v1";

// GreeterService 提供问候相关的服务
service GreeterService {
  // SayHello 向指定用户发送问候
  rpc SayHello(SayHelloRequest) returns (SayHelloResponse);
}

// SayHelloRequest SayHello 方法的请求参数
message SayHelloRequest {
  // 要问候的用户名称
  string name = 1;
}

// SayHelloResponse SayHello 方法的响应结果
message SayHelloResponse {
  // 问候消息
  string message = 1;
}
```

**代码解析**：

| 元素 | 说明 |
|------|------|
| `syntax = "proto3"` | 使用 Protobuf 3 语法 |
| `package helloworld.v1` | Proto 包名，用于命名空间隔离 |
| `option go_package` | 指定生成的 Go 代码的包路径和别名 |
| `service GreeterService` | 定义 gRPC 服务 |
| `rpc SayHello(...)` | 定义 RPC 方法 |
| `message SayHelloRequest` | 定义请求消息结构 |
| `string name = 1` | 字段定义，`1` 是字段编号（用于序列化） |

### 2.2 配置 Buf

#### buf.yaml（api/ 目录下）

```yaml
# api/buf.yaml
version: v1

# 模块名称
name: buf.build/go-api-template/api

# Lint 规则
lint:
  use:
    - DEFAULT

# Breaking change 检测
breaking:
  use:
    - FILE
```

#### buf.gen.yaml（项目根目录）

```yaml
# buf.gen.yaml
version: v1

# 插件配置
plugins:
  # 生成 Go 结构体（使用本地安装的 protoc-gen-go）
  - plugin: go
    out: api
    opt:
      - paths=source_relative

  # 生成 gRPC Go 代码（使用本地安装的 protoc-gen-go-grpc）
  - plugin: go-grpc
    out: api
    opt:
      - paths=source_relative
```

**配置说明**：

| 配置项 | 说明 |
|--------|------|
| `version: v1` | Buf 配置文件版本 |
| `plugin: go` | 使用 protoc-gen-go 插件 |
| `plugin: go-grpc` | 使用 protoc-gen-go-grpc 插件 |
| `out: api` | 输出目录为 api/ |
| `paths=source_relative` | 生成的文件放在与 proto 文件相同的相对路径 |

### 2.3 Lint 检查

运行 lint 检查确保 proto 文件符合规范：

```powershell
buf lint api/
```

如果没有输出，说明检查通过。

### 2.4 生成代码

```powershell
buf generate api/
```

此命令会调用配置的插件，生成以下文件：

```
api/helloworld/v1/
├── greeter.proto          ← 源文件（手动编写）
├── greeter.pb.go          ← 消息结构体（自动生成）
└── greeter_grpc.pb.go     ← gRPC 服务接口（自动生成）
```

### 2.5 添加依赖

生成的代码依赖 protobuf 和 gRPC 包，运行：

```powershell
go mod tidy
```

这会自动下载 `google.golang.org/protobuf` 和 `google.golang.org/grpc` 依赖。

---

## 3. 生成的代码解析

### 3.1 greeter.pb.go（消息结构体）

由 `protoc-gen-go` 生成，包含：

```go
// SayHelloRequest SayHello 方法的请求参数
type SayHelloRequest struct {
    state protoimpl.MessageState `protogen:"open.v1"`
    // 要问候的用户名称
    Name          string `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
    unknownFields protoimpl.UnknownFields
    sizeCache     protoimpl.SizeCache
}

// SayHelloResponse SayHello 方法的响应结果
type SayHelloResponse struct {
    state protoimpl.MessageState `protogen:"open.v1"`
    // 问候消息
    Message       string `protobuf:"bytes,1,opt,name=message,proto3" json:"message,omitempty"`
    unknownFields protoimpl.UnknownFields
    sizeCache     protoimpl.SizeCache
}
```

**特点**：

- proto 中的 `message` 转换为 Go 的 `struct`
- 字段带有 `protobuf` 和 `json` tag
- 自动生成 `GetXxx()` 方法
- 自动生成序列化/反序列化方法

### 3.2 greeter_grpc.pb.go（服务接口）

由 `protoc-gen-go-grpc` 生成，包含：

```go
// 客户端接口
type GreeterServiceClient interface {
    SayHello(ctx context.Context, in *SayHelloRequest, opts ...grpc.CallOption) (*SayHelloResponse, error)
}

// 服务端接口 - 你需要实现这个接口
type GreeterServiceServer interface {
    SayHello(context.Context, *SayHelloRequest) (*SayHelloResponse, error)
    mustEmbedUnimplementedGreeterServiceServer()
}

// 未实现的基础结构 - 用于向前兼容
type UnimplementedGreeterServiceServer struct{}
```

**特点**：

- `GreeterServiceClient`：客户端调用接口
- `GreeterServiceServer`：服务端需要实现的接口
- `UnimplementedGreeterServiceServer`：提供默认的"未实现"响应，用于向前兼容

---

## 4. Makefile 命令

更新后的 Makefile 提供以下 proto 相关命令：

```makefile
# 生成 Proto 代码（包含 lint 检查）
proto: proto-lint
	buf generate api/

# Proto 文件 lint 检查
proto-lint:
	buf lint api/

# Proto 文件格式化
proto-format:
	buf format -w api/
```

使用方式：

```powershell
# 生成代码（先 lint 再生成）
make proto

# 仅检查规范
make proto-lint

# 格式化 proto 文件
make proto-format
```

---

## 5. 验证

### 5.1 编译检查

```powershell
go build ./...
```

如果编译通过，说明生成的代码正确无误。

### 5.2 阶段二完成后的目录结构

```
go-api-template/
├── api/
│   ├── buf.yaml                      ← Buf 模块配置
│   └── helloworld/v1/
│       ├── greeter.proto             ← Proto 定义（手动编写）
│       ├── greeter.pb.go             ← 消息结构体（自动生成）
│       └── greeter_grpc.pb.go        ← gRPC 接口（自动生成）
├── buf.gen.yaml                      ← Buf 生成配置
├── cmd/
│   └── server/
│       └── main.go
├── internal/
│   ├── biz/                          ← 待实现
│   ├── data/                         ← 待实现
│   ├── server/                       ← 待实现
│   └── service/                      ← 待实现
├── go.mod
├── go.sum
└── Makefile
```

---

## 6. Schema First 开发模式总结

阶段二的核心是理解 **Schema First** 开发模式：

```
传统开发：先写代码 → 再补文档 → 接口不一致
Schema First：先定义契约 → 自动生成代码 → 保持一致
```

**优势**：

| 方面 | 说明 |
|------|------|
| 接口一致性 | 客户端和服务端从同一个 proto 文件生成代码 |
| 文档即代码 | proto 文件就是最权威的 API 文档 |
| 类型安全 | 编译期就能发现类型错误 |
| 多语言支持 | 同一个 proto 可以生成 Go、Java、Python 等语言代码 |

---

## 7. 下一步：阶段三

阶段三将实现业务逻辑层：

1. 在 `internal/biz/` 定义领域实体和 Repository 接口
2. 在 `internal/service/` 实现 `GreeterServiceServer` 接口
3. 在 `internal/data/` 使用内存 Map 实现 Repository

目标：体验**依赖倒置**的威力——业务逻辑不关心数据存储在哪里。
