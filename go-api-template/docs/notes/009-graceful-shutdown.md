# 009. 优雅关闭 (Graceful Shutdown)

## 1. 问题背景

### 1.1 什么是优雅关闭？

当服务需要停止时（部署更新、扩缩容、故障处理），有两种方式：

| 方式               | 行为                                     | 后果                                       |
| ------------------ | ---------------------------------------- | ------------------------------------------ |
| **强制关闭** | 立即终止进程，丢弃所有连接               | 正在处理的请求失败，客户端收到连接重置错误 |
| **优雅关闭** | 停止接受新连接，等待现有请求完成后再关闭 | 客户端无感知，请求正常完成                 |

### 1.2 为什么需要优雅关闭？

**生产环境常见场景**：

1. **滚动更新**：Kubernetes 部署新版本时，旧 Pod 需要优雅退出
2. **水平扩缩容**：缩容时需要等待实例处理完请求
3. **手动重启**：运维人员更新配置后重启服务
4. **资源清理**：确保数据库连接、文件句柄等资源正确释放

**没有优雅关闭的后果**：

- 用户请求失败（HTTP 502/504）
- 数据库事务未提交，数据不一致
- 连接池中的连接泄漏
- 日志丢失

## 2. 核心概念

在看具体实现之前，必须先理解几个关键概念。

### 为什么不能直接用 gin.Engine.Run()？

之前我们启动服务器的代码是这样的：

```go
// 之前的写法
engine := gin.Default()
engine.GET("/health", ...)
engine.Run(":8080")  // 阻塞，直到程序退出
```

这种写法有一个问题：**无法优雅关闭**。

`gin.Engine.Run()` 内部创建了 `http.Server` 并调用 `ListenAndServe()`，但它没有把 `http.Server` 的引用返回给我们。而优雅关闭需要调用 `http.Server.Shutdown()` 方法，我们拿不到 `http.Server`，就没法调用 `Shutdown()`。

**解决方案**：自己创建 `http.Server`，把 `gin.Engine` 作为 Handler 传入，这样我们就有了 `http.Server` 的引用，可以在需要时调用 `Shutdown()`。

下面介绍实现优雅关闭需要理解的几个核心概念。

### 2.1 net/http 包与 http.Server

Go 标准库 `net/http` 是用来构建 HTTP 服务器和客户端的包。其中 `http.Server` 是 HTTP 服务器的核心结构体。

```go
// http.Server 是 Go 标准库中的 HTTP 服务器
type Server struct {
    Addr         string        // 监听地址，如 ":8080"
    Handler      http.Handler  // 处理请求的对象（Gin 引擎实现了这个接口）
    ReadTimeout  time.Duration // 读取请求的超时时间
    WriteTimeout time.Duration // 写入响应的超时时间
    // ... 其他字段
}
```

**http.Server 的两个关键方法**：

```go
// ListenAndServe 启动服务器，开始监听端口并处理请求
// 这个方法会阻塞，直到服务器停止
func (srv *Server) ListenAndServe() error

// Shutdown 优雅关闭服务器
// 1. 停止接受新连接
// 2. 等待正在处理的请求完成
// 3. 返回
func (srv *Server) Shutdown(ctx context.Context) error
```

### 2.2 Gin 与 http.Server 的关系

**http.Server 和 gin.Engine 各自负责什么？**

```
请求到达
    │
    ▼
┌─────────────────────────────────────────────────────────
│  http.Server 的职责：    
│  1. 监听 TCP 端口（如 :8080）  
│  2. 接受客户端连接           
│  3. 读取 HTTP 请求数据         
│  4. 把请求交给 Handler 处理      
│  5. 把 Handler 返回的响应发送给客户端  
│  6. 管理连接的生命周期（保活、超时、关闭）   
└─────────────────────────────────────────────────────────
    │
    │ 调用 Handler.ServeHTTP(w, r)
    ▼
┌─────────────────────────────────────────────────────────
│  gin.Engine 的职责（实现了 Handler 接口）：  
│  1. 根据 URL 匹配路由（/api/v1/users → handleUsers）  
│  2. 执行中间件（日志、认证、CORS）             
│  3. 调用你写的处理函数                           
│  4. 生成响应内容（JSON、HTML）                      
└─────────────────────────────────────────────────────────
```

**gin.Engine 实现了 http.Handler 接口**：

```go
// http.Handler 接口定义
type Handler interface {
    ServeHTTP(ResponseWriter, *Request)
}

// gin.Engine 实现了这个接口
// 当请求到来时，http.Server 调用 gin.Engine.ServeHTTP()
// gin.Engine 负责路由匹配、执行中间件、调用你的处理函数
```

**gin.Engine.Run() 的真面目**：

```go
// gin.Engine.Run() 的简化实现（源码在 gin/gin.go）
func (engine *Engine) Run(addr string) error {
    // 其实内部就是创建 http.Server 并调用 ListenAndServe
    server := &http.Server{
        Addr:    addr,
        Handler: engine,  // gin.Engine 作为 Handler
    }
    return server.ListenAndServe()
}
```

**问题在于**：`gin.Engine.Run()` 没有返回 `http.Server` 的引用，我们无法调用 `Shutdown()` 方法。

**解决方案**：自己创建 `http.Server`，把 `gin.Engine` 作为 Handler 传入：

```go
engine := gin.Default()
// 配置路由...

server := &http.Server{
    Addr:    ":8080",
    Handler: engine,  // gin.Engine 作为请求处理器
}

// 现在我们有了 server 的引用，可以调用 Shutdown()
server.ListenAndServe()
// ... 收到关闭信号后 ...
server.Shutdown(ctx)
```

### 2.3 context 包：超时控制与取消机制

**context 在本次场景中的作用**：

在优雅关闭中，我们会调用 `server.Shutdown(ctx)`。这个 `ctx` 参数的作用是**控制 Shutdown 操作的超时时间**。

```go
// Shutdown 的函数签名
func (srv *Server) Shutdown(ctx context.Context) error
```

Shutdown 会等待所有正在处理的请求完成。但如果有请求一直不结束（比如 WebSocket 长连接），Shutdown 就会一直等下去。`ctx` 的作用就是设置一个"最长等待时间"，超过这个时间就强制关闭。

**为什么需要 context？看一个真实问题：**

```go
// 没有超时控制的代码
err := server.Shutdown(context.Background())
// 问题：如果有一个 WebSocket 连接永远不关闭，Shutdown 会永远等下去
// 程序卡死，无法退出
```

```go
// 有超时控制的代码
ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
defer cancel()
err := server.Shutdown(ctx)
// Shutdown 最多等 10 秒
// 超过 10 秒后，ctx 超时，Shutdown 立即返回，强制关闭剩余连接
```

**context.WithTimeout 的执行过程**：

```go
// 第 0 秒：创建 context，启动内部计时器
ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
defer cancel()

// 第 0 秒：调用 Shutdown，开始等待
err := server.Shutdown(ctx)

// 第 0-10 秒：Shutdown 在等待请求完成
// 同时，ctx 内部的计时器在倒计时

// 情况A：第 3 秒所有请求完成
// → Shutdown 正常返回，err == nil

// 情况B：第 10 秒还有请求没完成
// → ctx 超时，Shutdown 被迫返回，err == context.DeadlineExceeded
// → 未完成的请求被强制中断
```

**关于 cancel 函数**：

```go
ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
```

这行代码返回两个值：

1. `ctx`：带超时的 context，内部有一个计时器，10 秒后自动标记为"已取消"
2. `cancel`：一个函数，调用它可以提前取消 context（不用等 10 秒）

`defer cancel()` 的作用：

```go
ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
defer cancel()  // 函数结束时调用 cancel()

// 为什么需要 defer cancel()？
// context.WithTimeout 内部启动了一个 goroutine 来管理计时器
// 调用 cancel() 会停止这个 goroutine，释放资源
// 如果不调用 cancel()，这个 goroutine 会一直存在直到超时
// 这是 Go 的惯用写法，确保资源被正确释放
```

### 2.4 channel 与 select：等待多个事件

**channel 是 goroutine 之间通信的管道**：

```go
ch := make(chan int)  // 创建一个传递 int 的管道

// goroutine A：发送数据
go func() {
    ch <- 42  // 发送 42 到管道
}()

// goroutine B：接收数据
value := <-ch  // 从管道接收，如果管道为空，会阻塞等待
fmt.Println(value)  // 输出 42
```

**select 用于同时等待多个 channel**：

```go
select {
case v := <-ch1:
    // ch1 有数据时执行
case v := <-ch2:
    // ch2 有数据时执行
}
// select 会阻塞，直到某个 case 触发
```

**在我们的代码中**：

```go
select {
case err := <-errChan:
    // 服务器启动失败时触发
case sig := <-quit:
    // 收到关闭信号时触发
}
// select 会一直阻塞在这里，等待其中一个事件发生
// 后面的代码不会执行，直到 select 完成
```

### 2.5 信号（Signal）：操作系统通知进程的方式

当你按 Ctrl+C 时，操作系统不是直接杀死程序，而是发送一个 **信号** 给程序。

**信号监听需要三步，每一步的作用不同**：

```go
// 第 1 步：创建一个通道
quit := make(chan os.Signal, 1)
// 这只是创建了一个空的通道，什么都没发生
// 此时按 Ctrl+C，程序会被直接杀死（默认行为）

// 第 2 步：注册信号监听（关键！）
signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
// signal.Notify 的作用是告诉 Go 运行时：
// "当操作系统发送 SIGINT 或 SIGTERM 信号时，不要执行默认行为（杀死程序），
//  而是把这个信号发送到 quit 通道"
// 这一步之后，按 Ctrl+C 就不会直接杀死程序了

// 第 3 步：从通道读取，等待信号到来
sig := <-quit  // 阻塞，直到 quit 通道收到信号
fmt.Println("收到信号:", sig)
// 只有当用户按 Ctrl+C 时，quit 通道才会收到数据
// 这行代码才会继续执行
```

**执行时序**：

```
程序启动
    │
    ▼
quit := make(chan os.Signal, 1)     ← 创建空通道
    │
    ▼
signal.Notify(quit, SIGINT, SIGTERM) ← 注册监听，此后 Ctrl+C 会发到 quit
    │
    ▼
sig := <-quit                        ← 阻塞在这里，等待通道有数据
    │
    │ （程序在这里等待，可能等几秒、几分钟、几小时...）
    │
    │ 用户按 Ctrl+C
    │     ↓
    │ 操作系统发送 SIGINT 信号
    │     ↓
    │ Go 运行时收到信号，发送到 quit 通道
    │     ↓
    ▼
sig := <-quit 返回，sig 的值是 syscall.SIGINT
    │
    ▼
继续执行后面的代码
```

**常见信号**：

| 信号        | 触发方式         | 含义     |
| ----------- | ---------------- | -------- |
| `SIGINT`  | Ctrl+C           | 中断     |
| `SIGTERM` | `kill <pid>`   | 终止     |
| `SIGKILL` | `kill -9 <pid>` | 强制杀死（无法捕获） |

## 3. 技术方案

### 3.1 整体流程

```
程序启动
    │
    ▼
┌─────────────────────────────────────────────────────
│  httpServer.Start()  
│  在独立 goroutine 中启动服务器  
│  返回 errChan（用于接收启动错误）   
└─────────────────────────────────────────────────────
    │
    ▼
┌─────────────────────────────────────────────────────
│  select { ... }    
│  主 goroutine 阻塞在这里，等待两种事件之一：  
│  1. errChan 收到错误（服务器启动失败）  
│  2. quit 收到信号（用户按 Ctrl+C）   
└─────────────────────────────────────────────────────
    │
    │ （假设用户按了 Ctrl+C）
    ▼
┌─────────────────────────────────────────────────────
│  收到关闭信号，select 结束        
│  注意：此时服务器还在运行！只是收到了"请关闭"的通知   
└─────────────────────────────────────────────────────
    │
    ▼
┌─────────────────────────────────────────────────────
│  httpServer.Stop(ctx)                  
│  调用 http.Server.Shutdown()，这才是真正关闭服务器   
│  - 停止接受新连接                            
│  - 等待正在处理的请求完成（最多等 10 秒）       
│  - 关闭所有连接                                  
└─────────────────────────────────────────────────────
    │
    ▼
程序退出
```

**关键理解**：

- `errChan` 不是"关闭服务的通道"，而是"服务器启动失败的错误通道"
- `quit` 是"用户请求关闭的信号通道"
- 收到信号时服务器还在运行，`Stop()` 才是真正执行关闭的动作

### 3.2 http.Server.Shutdown() 的行为

```go
func (srv *Server) Shutdown(ctx context.Context) error
```

Shutdown 做了什么：

1. **立即停止接受新连接**：关闭监听的端口
2. **等待正在处理的请求完成**：已经在处理中的请求会继续处理
3. **关闭空闲连接**：没有正在处理请求的连接立即关闭
4. **超时控制**：如果 ctx 超时，立即返回，未完成的请求被强制中断

**关于 gin.Engine**：不需要单独关闭。gin.Engine 只是请求处理器，当 http.Server 关闭后，它自然就不会再收到请求了。

## 4. 代码实现详解

### 4.1 配置结构 (`internal/conf/conf.go`)

```go
// ServerConfig HTTP 服务器配置
type ServerConfig struct {
    ShutdownTimeout time.Duration `mapstructure:"shutdown_timeout"`
    ReadTimeout     time.Duration `mapstructure:"read_timeout"`
    WriteTimeout    time.Duration `mapstructure:"write_timeout"`
}

// GetShutdownTimeout 获取优雅关闭超时时间，提供默认值
func (c *ServerConfig) GetShutdownTimeout() time.Duration {
    if c.ShutdownTimeout <= 0 {
        return 10 * time.Second
    }
    return c.ShutdownTimeout
}
```

**为什么要判断 `<= 0`？**

这是**防御性编程**，处理配置缺失或错误的情况：

```yaml
# 情况1：配置文件中没有这个字段
server:
  # shutdown_timeout 没有配置

# 情况2：配置错误
server:
  shutdown_timeout: "abc"  # 无法解析为 Duration

# 情况3：用户手误
server:
  shutdown_timeout: 0s     # 配成了 0
```

在这些情况下，`ShutdownTimeout` 的值会是 Go 的零值 `0`。如果我们直接使用 0 秒作为超时时间，`Shutdown` 会立即返回，根本不等待请求完成——这不是我们想要的。

所以我们检查：如果值不合理（<= 0），就使用默认值 10 秒。

### 4.2 HTTPServer 封装 (`internal/server/server.go`)

```go
type HTTPServer struct {
    server *http.Server  // 底层 HTTP 服务器（负责网络监听）
    engine *gin.Engine   // Gin 引擎（负责路由和请求处理）
}
```

**Start 方法详解**：

```go
func (s *HTTPServer) Start() <-chan error {
    // 创建一个容量为 1 的 error 管道
    errChan := make(chan error, 1)
  
    // 启动一个新的 goroutine 来运行服务器
    go func() {
        // ListenAndServe 会阻塞，直到：
        // 1. 发生错误（端口被占用等）
        // 2. 调用 Shutdown() 关闭服务器
        err := s.server.ListenAndServe()
    
        // 判断错误类型
        // http.ErrServerClosed 是正常关闭时返回的错误，不算真正的错误
        if err != nil && !errors.Is(err, http.ErrServerClosed) {
            errChan <- err  // 只有真正的错误才发送到管道
        }
        close(errChan)  // 关闭管道，表示服务器已停止
    }()
  
    // 立即返回管道引用（不等待服务器启动完成）
    // 主函数可以通过这个管道知道服务器是否启动失败
    return errChan
}
```

**为什么返回 `<-chan error` 而不是 `chan error`？**

`<-chan error` 表示"只读管道"，调用者只能从中读取，不能写入。这是一种良好的封装。

**Stop 方法**：

```go
func (s *HTTPServer) Stop(ctx context.Context) error {
    // 直接调用标准库的 Shutdown 方法
    return s.server.Shutdown(ctx)
}
```

### 4.3 构建 http.Server (`internal/server/http.go`)

```go
func buildHTTPServer(cfg *conf.Config, engine *gin.Engine) *http.Server {
    return &http.Server{
        Addr:         fmt.Sprintf(":%d", cfg.App.Port),  // 监听地址
        Handler:      engine,                             // Gin 作为请求处理器
        ReadTimeout:  cfg.Server.GetReadTimeout(),        // 读取超时
        WriteTimeout: cfg.Server.GetWriteTimeout(),       // 写入超时
    }
}

func NewHTTPServer(cfg *conf.Config, greeterSvc *service.GreeterService) *HTTPServer {
    engine := gin.Default()
    // ... 配置路由 ...
  
    // 构建 http.Server，把 gin.Engine 作为 Handler
    httpServer := buildHTTPServer(cfg, engine)
  
    return &HTTPServer{
        server: httpServer,
        engine: engine,
    }
}
```

### 4.4 主函数实现 (`cmd/server/main.go`)

```go
func main() {
    // ... 加载配置、初始化应用 ...
  
    // ========================================
    // 步骤 1：启动服务器（非阻塞）
    // ========================================
    errChan := httpServer.Start()
    // Start() 会启动一个 goroutine 运行服务器，然后立即返回
    // 此时服务器可能还在启动中，也可能已经启动完成
    // 如果启动失败，错误会发送到 errChan
  
    // ========================================
    // 步骤 2：准备信号监听
    // ========================================
    quit := make(chan os.Signal, 1)
    signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
    // 告诉操作系统：当收到 SIGINT 或 SIGTERM 时，发送到 quit 管道
  
    // ========================================
    // 步骤 3：阻塞等待（这是关键！）
    // ========================================
    select {
    case err := <-errChan:
        // 情况A：服务器启动失败
        // 例如：端口被占用、权限不足等
        log.Fatalf("Server error: %v", err)
    
    case sig := <-quit:
        // 情况B：收到关闭信号（用户按 Ctrl+C）
        // 注意：此时服务器还在正常运行！
        // 我们只是收到了"请关闭"的通知
        log.Printf("Received signal: %v", sig)
    }
    // select 会一直阻塞在这里
    // 只有当 errChan 或 quit 有数据时，才会执行对应的 case
    // 然后 select 结束，继续执行后面的代码
  
    // ========================================
    // 步骤 4：优雅关闭
    // ========================================
    log.Println("Shutting down server...")
  
    // 创建带超时的 context
    // 含义：最多等待 shutdownTimeout 时间
    shutdownTimeout := cfg.Server.GetShutdownTimeout()  // 默认 10 秒
    ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
    defer cancel()  // 确保 context 资源被释放
  
    // 真正执行关闭
    // Stop 会调用 http.Server.Shutdown()
    // Shutdown 会等待所有正在处理的请求完成，但最多等待 10 秒
    if err := httpServer.Stop(ctx); err != nil {
        // 超时了，有些请求没处理完就被强制关闭了
        log.Printf("Server forced to shutdown: %v", err)
    } else {
        // 所有请求都正常完成了
        log.Println("Server gracefully stopped")
    }
}
```

## 5. 常见疑问解答

### Q1：errChan 和 quit 区别？

不一样：

| 管道    | 触发时机       | 含义               |
| ------- | -------------- | ------------------ |
| errChan | 服务器启动失败 | "出错了，无法启动" |
| quit    | 用户按 Ctrl+C  | "请关闭服务器"     |

正常流程是：

1. 服务器启动成功，errChan 没有数据
2. 服务器正常运行中，处理请求
3. 用户按 Ctrl+C，quit 收到信号
4. 执行优雅关闭

异常流程是：

1. 服务器启动失败（端口被占用），errChan 收到错误
2. 程序直接退出

### Q2：select 是阻塞的吗？

是的，`select` 会阻塞当前 goroutine，直到某个 case 触发。

```go
select {
case err := <-errChan:
    // ...
case sig := <-quit:
    // ...
}
// 这行代码不会执行，直到上面的 select 完成
log.Println("这要等 select 完成后才执行")
```

### Q3：收到信号时服务器不是已经关闭了吗？

不是！信号只是"通知"，不是"动作"。

```
用户按 Ctrl+C
    │
    ▼
操作系统发送 SIGINT 信号给程序
    │
    ▼
quit 管道收到信号
    │
    ▼
select 结束，程序继续执行
    │
    ▼
此时服务器还在运行！还在处理请求！
    │
    ▼
调用 httpServer.Stop(ctx)
    │
    ▼
这才是真正关闭服务器
```

### Q4：Shutdown 和 gin.Engine 的关系？

`http.Server.Shutdown()` 关闭的是底层网络连接。`gin.Engine` 只是请求处理器，它不管网络连接。

关闭顺序：

1. `Shutdown()` 停止接受新连接
2. 等待正在处理的请求完成（`gin.Engine` 还在工作）
3. 所有请求完成后，`gin.Engine` 自然没有请求可处理了
4. 服务器关闭完成

不需要单独关闭 `gin.Engine`。

## 6. 配置说明

`configs/config.yaml` 中的相关配置：

```yaml
server:
  # 优雅关闭超时时间
  # 收到关闭信号后，等待正在处理的请求完成的最大时间
  # 超过此时间后，强制关闭所有连接
  shutdown_timeout: 10s
  
  # 读取请求的超时时间
  # 防止慢速客户端攻击（Slowloris）
  read_timeout: 30s
  
  # 写入响应的超时时间
  write_timeout: 30s
```

**配置建议**：

| 环境       | shutdown_timeout | 说明                     |
| ---------- | ---------------- | ------------------------ |
| 开发环境   | 5-10s            | 快速迭代                 |
| 生产环境   | 25s              | K8s 默认 30s，留 5s 余量 |
| 长连接服务 | 根据业务         | WebSocket 可能需要更长   |

## 7. 测试优雅关闭

### 7.1 手动测试

1. 启动服务：`make run`
2. 按 Ctrl+C
3. 观察日志：
   - 应该看到 `Received signal: interrupt`
   - 应该看到 `Shutting down server...`
   - 应该看到 `Server gracefully stopped`

### 7.2 模拟慢请求测试

可以临时添加一个测试端点来验证优雅关闭：

```go
// 在 http.go 的 NewHTTPServer 中添加
engine.GET("/slow", func(c *gin.Context) {
    time.Sleep(5 * time.Second)  // 模拟慢请求
    c.JSON(200, gin.H{"message": "done"})
})
```

测试步骤：

1. 启动服务
2. 在浏览器访问 `http://localhost:8080/slow`
3. 立即按 Ctrl+C
4. 观察：
   - 慢请求应该正常完成（返回 "done"）
   - 然后服务器才关闭

## 8. 生产环境注意事项

### 8.1 Kubernetes 配置

```yaml
apiVersion: apps/v1
kind: Deployment
spec:
  template:
    spec:
      # K8s 等待 Pod 优雅关闭的时间
      # 应该大于你的 shutdown_timeout
      terminationGracePeriodSeconds: 30
      containers:
      - name: app
        lifecycle:
          preStop:
            exec:
              # 等待负载均衡器从后端列表中摘除此 Pod
              command: ["/bin/sh", "-c", "sleep 5"]
```

### 8.2 健康检查配合

关闭时应该让健康检查失败，使负载均衡器提前摘除实例：

```go
var isShuttingDown atomic.Bool

engine.GET("/health", func(c *gin.Context) {
    if isShuttingDown.Load() {
        c.JSON(503, gin.H{"status": "shutting_down"})
        return
    }
    c.JSON(200, gin.H{"status": "ok"})
})
```

### 8.3 资源清理顺序

优雅关闭时，资源清理的顺序很重要：

1. 停止接受新请求（HTTP Server Shutdown）
2. 等待正在处理的请求完成
3. 关闭数据库连接池
4. 关闭 Redis 连接池
5. 刷新日志缓冲区

## 9. 总结

### 核心概念

| 概念        | 说明                                                     |
| ----------- | -------------------------------------------------------- |
| http.Server | Go 标准库的 HTTP 服务器，负责网络监听和连接管理          |
| gin.Engine  | 请求处理器，实现了 http.Handler 接口，负责路由和业务逻辑 |
| context     | 控制超时和取消的机制，传递给需要控制生命周期的函数       |
| channel     | goroutine 之间通信的管道                                 |
| select      | 同时等待多个 channel，阻塞直到某个 case 触发             |
| signal      | 操作系统通知进程的方式（Ctrl+C 发送 SIGINT）             |

### 实现要点

| 要点             | 说明                                            |
| ---------------- | ----------------------------------------------- |
| 使用 http.Server | 不要使用 gin.Run()，它不暴露 Server 引用        |
| 非阻塞启动       | Start() 在 goroutine 中运行，主函数可以继续执行 |
| select 等待      | 同时监听启动错误和关闭信号                      |
| 超时控制         | 使用 context.WithTimeout 限制关闭等待时间       |
| 防御性默认值     | 配置缺失时提供合理的默认值                      |
