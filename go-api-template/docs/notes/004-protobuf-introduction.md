# 004. Protobuf 入门：跨语言的接口契约

本章从零开始介绍 Protocol Buffers（简称 Protobuf 或 Proto），帮助你理解它是什么、为什么需要它、以及如何使用它。

---

## 1. Proto 文件是什么

### 1.1 一句话定义

`.proto` 文件是一种**接口定义文件**，用来描述数据的结构和服务的接口。

### 1.2 类比理解

| 类比对象                | 说明                     |
| ----------------------- | ------------------------ |
| JSON Schema             | 定义 JSON 数据的结构规范 |
| TypeScript 的 `.d.ts` | 定义类型，不包含实现     |
| 数据库的 DDL            | 定义表结构，不包含数据   |
| OpenAPI/Swagger         | 定义 REST API 接口规范   |

**Proto 文件就像是一份"契约"**：它只描述"数据长什么样"和"有哪些方法"，不包含任何业务逻辑代码。

### 1.3 编辑器配置

默认情况下，VS Code / Cursor 打开 `.proto` 文件时，所有内容都是白色的（无语法高亮），因为编辑器不认识这种文件格式。

**解决方法**：安装 `vscode-proto3` 插件。

1. 打开 VS Code / Cursor
2. 按 `Ctrl+Shift+X` 打开扩展面板
3. 搜索 `vscode-proto3`
4. 点击安装

安装后，`.proto` 文件会有：
- 语法高亮（关键字、类型、注释不同颜色）
- 代码格式化
- 语法错误提示
- 自动补全

### 1.4 一个简单的例子

```protobuf
// 定义一个"用户"数据结构
message User {
  int64 id = 1;           // 用户ID
  string name = 2;        // 用户名
  string email = 3;       // 邮箱
}

// 定义一个"用户服务"接口
service UserService {
  rpc GetUser(GetUserRequest) returns (User);      // 根据ID获取用户
  rpc CreateUser(CreateUserRequest) returns (User); // 创建新用户
}
```

这就是一个 proto 文件的基本样子：定义数据结构（`message`）和服务接口（`service`）。

---

## 2. 为什么需要 Proto

### 2.1 核心价值：一份定义，多端共享

假设你在开发一个系统，有以下技术栈：

- 后端：Go
- 前端：TypeScript
- 移动端：Kotlin/Swift
- 数据分析：Python

**传统方式的问题**：

```
前端定义：{ userId: number, userName: string }
后端定义：type User struct { UserID int64; UserName string }
移动端定义：data class User(val user_id: Long, val user_name: String)
```

三个端各自定义，字段名不一致（`userId` vs `UserID` vs `user_id`），类型也可能不匹配。一旦接口变更，需要手动同步修改所有端。

**使用 Proto 的方式**：

```
                    ┌─────────────────┐
                    │  user.proto     │  唯一的接口定义
                    │  (定义一次)      │
                    └────────┬────────┘
                             │
              ┌──────────────┼──────────────┬──────────────┐
              ▼              ▼              ▼              ▼
         ┌────────┐    ┌────────┐    ┌────────┐    ┌────────┐
         │ Go 代码 │   │ TS代码  │    │ Kotlin │    │Python  │
         │ 自动生成│    │ 自动生成│    │自动生成│    │ 自动生成│
         └────────┘    └────────┘    └────────┘    └────────┘
```

**所有端的字段名、类型绝对一致**，因为它们都是从同一份 proto 文件生成的。

### 2.2 Proto 的三大优势

| 优势             | 说明                                                                            |
| ---------------- | ------------------------------------------------------------------------------- |
| **跨语言** | 一份 proto 可以生成 Go、TypeScript、Python、Java、C++、Rust 等 20+ 种语言的代码 |
| **强类型** | 编译期就能发现类型错误，不用等到运行时                                          |
| **高性能** | 二进制序列化格式，比 JSON 体积更小、解析更快                                    |

### 2.3 谁在用 Proto

- Google 内部几乎所有服务都使用 Protobuf
- gRPC（Google 的 RPC 框架）的默认序列化格式
- Kubernetes、Envoy、Istio 等云原生项目
- 微服务架构中的标准接口定义方式

---

## 3. Proto 语法详解

Proto 有自己的语法，但非常简单，10 分钟就能掌握核心内容。

### 3.1 文件头部

```protobuf
syntax = "proto3";                                    // 声明使用 proto3 语法
package helloworld.v1;                                // 包名（命名空间）
option go_package = "go-api-template/api/helloworld/v1;v1";  // Go 包路径
```

| 元素                  | 说明                                             |
| --------------------- | ------------------------------------------------ |
| `syntax = "proto3"` | 必须放在第一行，声明使用 proto3 版本（目前主流） |
| `package`           | 命名空间，防止不同 proto 文件的类型名冲突        |
| `option go_package` | Go 专用配置，指定生成代码的包路径                |

### 3.2 定义数据结构：message

`message` 是 proto 的核心，用来定义数据结构（类似 Go 的 struct）。

```protobuf
message User {
  int64 id = 1;           // 字段类型 字段名 = 字段编号;
  string name = 2;
  string email = 3;
  bool is_active = 4;
  repeated string roles = 5;  // repeated 表示数组/切片
}
```

**字段编号的作用**：

- 每个字段必须有唯一的编号（1, 2, 3...）
- 编号用于二进制序列化，一旦使用就不能修改
- 1-15 占用 1 字节，16-2047 占用 2 字节，所以常用字段放前面

### 3.3 基本数据类型

| Proto 类型 | Go 类型     | TypeScript 类型 | 说明            |
| ---------- | ----------- | --------------- | --------------- |
| `double` | `float64` | `number`      | 双精度浮点      |
| `float`  | `float32` | `number`      | 单精度浮点      |
| `int32`  | `int32`   | `number`      | 32位有符号整数  |
| `int64`  | `int64`   | `bigint`      | 64位有符号整数  |
| `uint32` | `uint32`  | `number`      | 32位无符号整数  |
| `uint64` | `uint64`  | `bigint`      | 64位无符号整数  |
| `bool`   | `bool`    | `boolean`     | 布尔值          |
| `string` | `string`  | `string`      | 字符串（UTF-8） |
| `bytes`  | `[]byte`  | `Uint8Array`  | 字节数组        |

### 3.4 复合类型

#### 数组/切片：repeated

```protobuf
message User {
  repeated string roles = 1;       // 字符串数组 → Go: []string
  repeated Address addresses = 2;  // 对象数组 → Go: []*Address
}
```

#### 字典/Map：map

```protobuf
message User {
  map<string, string> metadata = 1;  // → Go: map[string]string
  map<int64, Order> orders = 2;      // → Go: map[int64]*Order
}
```

#### 嵌套 Message

```protobuf
message Order {
  int64 id = 1;
  repeated OrderItem items = 2;  // 引用另一个 message
}

message OrderItem {
  string product_name = 1;
  int32 quantity = 2;
  double price = 3;
}
```

#### 枚举：enum

```protobuf
enum UserStatus {
  USER_STATUS_UNSPECIFIED = 0;  // 第一个值必须是 0
  USER_STATUS_ACTIVE = 1;
  USER_STATUS_INACTIVE = 2;
  USER_STATUS_BANNED = 3;
}

message User {
  int64 id = 1;
  UserStatus status = 2;  // 使用枚举类型
}
```

### 3.5 定义服务接口：service

`service` 定义 RPC 服务接口（类似 Go 的 interface）。

```protobuf
service UserService {
  // 一元 RPC：一个请求，一个响应
  rpc GetUser(GetUserRequest) returns (GetUserResponse);
  
  // 服务端流：一个请求，多个响应（流式返回）
  rpc ListUsers(ListUsersRequest) returns (stream User);
  
  // 客户端流：多个请求，一个响应
  rpc UploadUsers(stream User) returns (UploadUsersResponse);
  
  // 双向流：多个请求，多个响应
  rpc Chat(stream ChatMessage) returns (stream ChatMessage);
}
```

**命名规范**：

- Service 名称：`PascalCase`，以 `Service` 结尾（如 `UserService`）
- RPC 方法名：`PascalCase`（如 `GetUser`、`CreateUser`）
- 每个方法有独立的 Request 和 Response 类型

### 3.6 完整示例

```protobuf
syntax = "proto3";

package user.v1;

option go_package = "myproject/api/user/v1;v1";

// 用户状态枚举
enum UserStatus {
  USER_STATUS_UNSPECIFIED = 0;
  USER_STATUS_ACTIVE = 1;
  USER_STATUS_INACTIVE = 2;
}

// 用户实体
message User {
  int64 id = 1;
  string name = 2;
  string email = 3;
  UserStatus status = 4;
  repeated string roles = 5;
  map<string, string> metadata = 6;
}

// 获取用户请求
message GetUserRequest {
  int64 id = 1;
}

// 获取用户响应
message GetUserResponse {
  User user = 1;
}

// 创建用户请求
message CreateUserRequest {
  string name = 1;
  string email = 2;
}

// 创建用户响应
message CreateUserResponse {
  User user = 1;
}

// 用户服务
service UserService {
  rpc GetUser(GetUserRequest) returns (GetUserResponse);
  rpc CreateUser(CreateUserRequest) returns (CreateUserResponse);
}
```

---

## 4. Proto、Buf 和生成的文件

### 4.1 工具链关系图

```
┌─────────────────────────────────────────────────────────────────────────┐
│                             开发流程                                     │
│                                                                         │
│   你编写 .proto 文件                                                     │
│         │                                                               │
│         ▼                                                               │
│   ┌─────────────┐                                                       │
│   │   Buf CLI   │  ← 协调者：读取配置，调用插件                           │
│   └──────┬──────┘                                                       │
│          │                                                              │
│    ┌─────┴─────┐                                                        │
│    │           │                                                        │
│    ▼           ▼                                                        │
│ ┌────────────────┐  ┌────────────────────┐                              │
│ │ protoc-gen-go  │  │ protoc-gen-go-grpc │                              │
│ │ (消息结构体插件) │  │ (服务接口插件)      │                              │
│ └───────┬────────┘  └─────────┬──────────┘                              │
│         │                     │                                         │
│         ▼                     ▼                                         │
│   ┌───────────┐         ┌─────────────────┐                             │
│   │ *.pb.go   │         │ *_grpc.pb.go    │                             │
│   │ 消息结构体 │         │ gRPC 服务接口   │                             │
│   └───────────┘         └─────────────────┘                             │
└─────────────────────────────────────────────────────────────────────────┘
```

### 4.2 Buf 是什么

Buf 是 Protobuf 生态的**现代化管理工具**，它不是生成代码的工具，而是一个协调者：

| Buf 的职责      | 说明                              |
| --------------- | --------------------------------- |
| 管理 proto 文件 | 类似 npm 管理 JS 包               |
| 调用生成插件    | 根据配置调用 protoc-gen-go 等插件 |
| Lint 检查       | 检查 proto 文件是否符合规范       |
| Breaking 检测   | 检测接口是否有破坏性变更          |

### 4.3 两个生成的文件

当你运行 `buf generate` 后，每个 `.proto` 文件会生成两个 Go 文件：

#### 文件 1：`*.pb.go`（消息结构体）

由 `protoc-gen-go` 插件生成，包含：

```go
// proto 中的 message → Go 中的 struct
type User struct {
    Id    int64  `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
    Name  string `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
    Email string `protobuf:"bytes,3,opt,name=email,proto3" json:"email,omitempty"`
}

// 自动生成的 Getter 方法
func (x *User) GetId() int64 { ... }
func (x *User) GetName() string { ... }
func (x *User) GetEmail() string { ... }

// 序列化/反序列化方法（内部使用）
func (x *User) ProtoReflect() protoreflect.Message { ... }
```

**这个文件的作用**：定义数据结构，用于在代码中创建和操作数据对象。

#### 文件 2：`*_grpc.pb.go`（服务接口）

由 `protoc-gen-go-grpc` 插件生成，包含：

```go
// 客户端接口（用于调用远程服务）
type UserServiceClient interface {
    GetUser(ctx context.Context, in *GetUserRequest, opts ...grpc.CallOption) (*GetUserResponse, error)
    CreateUser(ctx context.Context, in *CreateUserRequest, opts ...grpc.CallOption) (*CreateUserResponse, error)
}

// 服务端接口（你需要实现这个）
type UserServiceServer interface {
    GetUser(context.Context, *GetUserRequest) (*GetUserResponse, error)
    CreateUser(context.Context, *CreateUserRequest) (*CreateUserResponse, error)
    mustEmbedUnimplementedUserServiceServer()
}

// 未实现的默认结构体（用于向前兼容）
type UnimplementedUserServiceServer struct{}
```

**这个文件的作用**：定义服务接口，你需要实现 `Server` 接口来提供具体的业务逻辑。

---

## 5. 使用场景与工作流程

### 5.1 核心原则：关注 Proto，忽略生成的代码

```
重要程度排序：

1. greeter.proto          ← 你唯一需要编写和维护的文件
2. greeter.pb.go          ← 自动生成，不要手动修改
3. greeter_grpc.pb.go     ← 自动生成，不要手动修改
```

**生成的 `*.pb.go` 和 `*_grpc.pb.go` 文件你不需要阅读、理解或修改**。它们是工具的产物，就像 TypeScript 编译出的 JavaScript 文件一样。

### 5.2 标准开发流程

```
步骤 1：编写/修改 .proto 文件
         │
         ▼
步骤 2：运行 buf generate 重新生成代码
         │
         ▼
步骤 3：在你的业务代码中使用生成的类型和接口
         │
         ▼
步骤 4：接口有变更？回到步骤 1
```

### 5.3 实际使用示例

假设你有这个 proto 定义：

```protobuf
service GreeterService {
  rpc SayHello(SayHelloRequest) returns (SayHelloResponse);
}

message SayHelloRequest {
  string name = 1;
}

message SayHelloResponse {
  string message = 1;
}
```

在你的业务代码中这样使用：

```go
package service

import (
    "context"
  
    v1 "go-api-template/api/helloworld/v1"  // 导入生成的包
)

// 实现生成的 GreeterServiceServer 接口
type GreeterService struct {
    v1.UnimplementedGreeterServiceServer  // 嵌入以实现向前兼容
}

// 实现 SayHello 方法
func (s *GreeterService) SayHello(ctx context.Context, req *v1.SayHelloRequest) (*v1.SayHelloResponse, error) {
    // req.Name 和 SayHelloResponse 都是从 proto 生成的
    return &v1.SayHelloResponse{
        Message: "Hello, " + req.Name + "!",
    }, nil
}
```

**你只需要**：

1. 导入生成的包
2. 使用生成的类型（`SayHelloRequest`、`SayHelloResponse`）
3. 实现生成的接口（`GreeterServiceServer`）

### 5.4 什么时候需要重新生成

| 场景                      | 是否需要重新生成                    |
| ------------------------- | ----------------------------------- |
| 修改了 proto 文件         | ✅ 需要                             |
| 修改了业务代码            | ❌ 不需要                           |
| 升级了 protoc-gen-go 版本 | ⚠️ 建议重新生成                   |
| 换了一台电脑              | ❌ 不需要（生成的代码已提交到 Git） |

---

## 6. 常见问题

### 6.1 生成的代码要不要提交到 Git？

**建议提交**。原因：

- 其他开发者无需安装工具链就能直接编译
- Code Review 时能看到接口变更的实际影响

### 6.2 Proto 和 JSON 的关系？

Proto 定义的结构体同时支持 JSON 序列化：

```go
user := &v1.User{Id: 1, Name: "Alice"}

// Protobuf 二进制格式（更小更快）
data, _ := proto.Marshal(user)

// JSON 格式（可读性好）
jsonData, _ := protojson.Marshal(user)
// 输出: {"id":"1","name":"Alice"}
```

### 6.3 只用 HTTP 不用 gRPC，还需要 Proto 吗？

可以只用 Proto 定义数据结构，不定义 service。生成的 `*.pb.go` 中的结构体可以直接用于 HTTP API：

```go
func HandleGetUser(c *gin.Context) {
    user := &v1.User{Id: 1, Name: "Alice"}
    c.JSON(200, user)  // 自动序列化为 JSON
}
```

---

## 7. 小结

| 概念             | 说明                                       |
| ---------------- | ------------------------------------------ |
| Proto 文件       | 接口定义文件，描述数据结构和服务接口       |
| 核心价值         | 一份定义生成多语言代码，保证各端一致       |
| 语法要点         | `message` 定义数据，`service` 定义接口 |
| Buf              | 协调工具，调用插件生成代码                 |
| `*.pb.go`      | 消息结构体，由 protoc-gen-go 生成          |
| `*_grpc.pb.go` | 服务接口，由 protoc-gen-go-grpc 生成       |
| 工作流程         | 编写 proto → 生成代码 → 实现接口         |

**记住**：你只需要关注 `.proto` 文件，生成的代码是工具的产物，不需要阅读和修改。
