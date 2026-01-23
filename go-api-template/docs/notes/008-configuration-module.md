# 008. 配置模块

本章介绍 Go API 服务的配置管理最佳实践

---

## 1. 设计目标

- **类型安全**：使用 Go 结构体映射配置，编译期检查配置字段
- **环境隔离**：区分开发环境和生产环境配置
- **敏感信息保护**：密码、密钥等敏感信息通过环境变量覆盖，不提交到 Git
- **集成 Wire**：配置作为依赖注入的一部分，各层通过构造函数接收配置

---

## 2. 技术选型：Viper

[Viper](https://github.com/spf13/viper) 是 Go 生态中最成熟的配置库。

| 特性         | 说明                               |
| ------------ | ---------------------------------- |
| 多格式支持   | YAML, JSON, TOML, 环境变量         |
| 环境变量覆盖 | 自动将配置路径映射为环境变量名     |
| 配置热重载   | 支持配置文件变更后自动重载（可选） |
| 默认值       | 支持设置配置默认值                 |

安装依赖：

```bash
go get github.com/spf13/viper
```

---

## 3. 目录结构

```
go-api-template/
├── configs/
│   ├── config.yaml          # 开发环境配置（含实际值，不提交 Git）
│   └── config.example.yaml  # 配置模板（提交 Git，供参考）
├── internal/
│   └── conf/
│       └── conf.go          # 配置结构体定义与加载逻辑
```

---

## 4. 配置结构体

配置使用嵌套的 Go 结构体表示，通过 `mapstructure` 标签映射 YAML 字段。

```go
// internal/conf/conf.go

type Config struct {
    App      AppConfig      `mapstructure:"app"`
    Log      LogConfig      `mapstructure:"log"`
    Database DatabaseConfig `mapstructure:"database"`
    Redis    RedisConfig    `mapstructure:"redis"`
    JWT      JWTConfig      `mapstructure:"jwt"`
}

type AppConfig struct {
    Name string `mapstructure:"name"`
    Env  string `mapstructure:"env"`   // development | production
    Port int    `mapstructure:"port"`
}

type DatabaseConfig struct {
    Driver   string `mapstructure:"driver"`   // postgres | mysql | sqlite
    Host     string `mapstructure:"host"`
    Port     int    `mapstructure:"port"`
    Database string `mapstructure:"database"`
    Username string `mapstructure:"username"`
    Password string `mapstructure:"password"`
    // 连接池配置
    MaxIdleConns    int           `mapstructure:"max_idle_conns"`
    MaxOpenConns    int           `mapstructure:"max_open_conns"`
    ConnMaxLifetime time.Duration `mapstructure:"conn_max_lifetime"`
}
```

---

## 5. 配置文件

### 5.1 开发配置 (configs/config.yaml)

```yaml
app:
  name: go-api-template
  env: development
  port: 8080

log:
  level: debug
  format: text

database:
  driver: postgres
  host: localhost
  port: 5432
  database: go_api_template
  username: postgres
  password: ""  # 通过环境变量 DATABASE_PASSWORD 设置
  max_idle_conns: 10
  max_open_conns: 100
  conn_max_lifetime: 1h

redis:
  host: localhost
  port: 6379
  password: ""
  db: 0

jwt:
  secret: dev-secret-change-in-production
  expires_in: 24h
```

### 5.2 配置模板 (configs/config.example.yaml)

与 `config.yaml` 结构相同，但敏感字段为空或使用占位符。此文件提交到 Git，供新开发者参考。

---

## 6. 环境变量覆盖

Viper 支持使用环境变量覆盖配置文件中的值。

### 6.1 命名规则

将配置路径中的 `.` 替换为 `_`，并全部大写：

| 配置路径              | 环境变量              |
| --------------------- | --------------------- |
| `app.port`          | `APP_PORT`          |
| `database.password` | `DATABASE_PASSWORD` |
| `redis.password`    | `REDIS_PASSWORD`    |
| `jwt.secret`        | `JWT_SECRET`        |

### 6.2 使用示例

Windows PowerShell：

```powershell
$env:DATABASE_PASSWORD = "your-secure-password"
$env:APP_ENV = "production"
go run ./cmd/server
```

Linux/macOS：

```bash
export DATABASE_PASSWORD="your-secure-password"
export APP_ENV="production"
go run ./cmd/server
```

---

## 7. 集成到 Wire

配置作为依赖注入的根节点，在 `main.go` 中加载后传入 `wireApp`。

### 7.1 main.go

```go
func main() {
    flag.Parse()

    // 加载配置
    cfg, err := conf.LoadConfig(configPath)
    if err != nil {
        log.Fatalf("Failed to load config: %v", err)
    }

    // 初始化应用（Wire 生成的函数）
    httpServer, err := wireApp(cfg)
    if err != nil {
        log.Fatalf("Failed to initialize application: %v", err)
    }

    // 启动服务
    if err := httpServer.Run(); err != nil {
        log.Fatalf("Failed to start server: %v", err)
    }
}
```

### 7.2 wire.go

```go
func wireApp(c *conf.Config) (*server.HTTPServer, error) {
    wire.Build(
        data.ProviderSet,
        biz.ProviderSet,
        service.ProviderSet,
        server.ProviderSet,
    )
    return nil, nil
}
```

### 7.3 Wire 的"参数即 Provider"机制

**疑问**：`wire.go` 中参数 `c *conf.Config` 并没有显式加入任何 ProviderSet，为什么生成的 `wire_gen.go` 会自动把 `c` 传给 `NewHTTPServer` 和 `NewData`？

**答案**：这是 Wire 的核心设计 —— **Injector 函数的参数自动成为可用的 Provider**。

Wire 的依赖解析过程：

1. **收集所有 Provider**：Wire 扫描 `wire.Build()` 中的所有 ProviderSet，提取每个构造函数能"提供"的类型
2. **参数也是 Provider**：Injector 函数（如 `wireApp`）的参数会被 Wire 视为"已存在的依赖"，无需再构造
3. **类型匹配**：Wire 发现 `NewHTTPServer` 需要 `*conf.Config`，而 `wireApp(c *conf.Config)` 的参数正好提供了这个类型，于是自动匹配

**生成代码对比**：

```go
// wire.go（声明）
func wireApp(c *conf.Config) (*server.HTTPServer, error) {
    wire.Build(...)  // c 没有显式关联任何 Provider
}

// wire_gen.go（生成）
func wireApp(c *conf.Config) (*server.HTTPServer, error) {
    dataData, err := data.NewData(c)          // ← c 被自动传入
    // ...
    httpServer := server.NewHTTPServer(c, greeterService)  // ← c 被自动传入
    return httpServer, nil
}
```

**类比理解**：可以把 Injector 参数想象成"外部注入的依赖"。Wire 不需要知道如何创建它，只需要知道"调用时会有人提供这个值"。

### 7.4 依赖流向

```
main.go
   │
   │ LoadConfig(path)
   ▼
*conf.Config
   │
   │ wireApp(cfg)
   ▼
┌──────────────────────────────────────────────┐
│                   Wire                       │
│                                              │
│  conf.Config ──┬── data.NewData(cfg)         │
│                │                             │
│                ├── server.NewHTTPServer(cfg) │
│                │                             │
│                └── (其他需要配置的层)         │
└──────────────────────────────────────────────┘
```

---

## 8. 运行时行为

启动服务时，会看到配置加载日志：

```
2026/01/23 17:56:59 Loaded config: env=development, port=8080
2026/01/23 17:56:59 Data layer initialized: driver=postgres, host=localhost, database=go_api_template
2026/01/23 17:56:59 Starting server on :8080
```

访问根端点 `GET /` 会返回配置信息：

```json
{
  "name": "go-api-template",
  "version": "0.1.0",
  "env": "development",
  "message": "Welcome to Go API Template"
}
```

---

## 9. 生产环境部署

### 9.1 配置策略

| 配置类型 | 开发环境    | 生产环境               |
| -------- | ----------- | ---------------------- |
| 应用配置 | config.yaml | config.yaml + 环境变量 |
| 敏感信息 | 本地值      | 环境变量覆盖           |
| 日志级别 | debug       | info 或 warn           |
| 日志格式 | text        | json                   |

### 9.2 Docker 部署示例

```dockerfile
FROM golang:1.21 AS builder
WORKDIR /app
COPY . .
RUN go build -o server ./cmd/server

FROM debian:bookworm-slim
WORKDIR /app
COPY --from=builder /app/server .
COPY configs/config.yaml configs/

EXPOSE 8080
CMD ["./server", "-config", "configs/config.yaml"]
```

```bash
docker run -p 8080:8080 \
  -e APP_ENV=production \
  -e DATABASE_PASSWORD=secure-password \
  -e JWT_SECRET=random-secret-key \
  go-api-template
```

### 9.3 Kubernetes ConfigMap + Secret

```yaml
# ConfigMap for non-sensitive config
apiVersion: v1
kind: ConfigMap
metadata:
  name: app-config
data:
  config.yaml: |
    app:
      name: go-api-template
      env: production
      port: 8080
    database:
      driver: postgres
      host: postgres-service
      port: 5432
      database: go_api_template
      username: app_user
---
# Secret for sensitive values
apiVersion: v1
kind: Secret
metadata:
  name: app-secrets
stringData:
  DATABASE_PASSWORD: "your-secure-password"
  JWT_SECRET: "your-jwt-secret"
```

---

## 10. 本章小结

| 完成项     | 说明                                     |
| ---------- | ---------------------------------------- |
| Viper 依赖 | 配置库安装完成                           |
| 配置结构体 | `internal/conf/conf.go` 定义所有配置   |
| 配置文件   | `configs/config.yaml` 开发配置         |
| 配置模板   | `configs/config.example.yaml` 提交 Git |
| Wire 集成  | 配置通过依赖注入传递到各层               |
| .gitignore | 忽略 `configs/config.yaml`             |
| 环境变量   | 支持通过环境变量覆盖敏感配置             |
