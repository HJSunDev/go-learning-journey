// Package conf 定义应用配置结构体和配置加载逻辑。
// 采用 Viper 库实现配置文件读取和环境变量覆盖。
package conf

import (
	"fmt"
	"strings"
	"time"

	"github.com/google/wire"
	"github.com/spf13/viper"
)

// ProviderSet 配置模块的依赖注入 Provider
var ProviderSet = wire.NewSet(LoadConfig)

// Config 应用根配置，聚合所有配置模块
type Config struct {
	App      AppConfig      `mapstructure:"app"`
	Server   ServerConfig   `mapstructure:"server"`
	Log      LogConfig      `mapstructure:"log"`
	Database DatabaseConfig `mapstructure:"database"`
	Redis    RedisConfig    `mapstructure:"redis"`
	JWT      JWTConfig      `mapstructure:"jwt"`
}

// AppConfig 应用基础配置
type AppConfig struct {
	// 应用名称
	Name string `mapstructure:"name"`
	// 运行环境：development | production
	Env string `mapstructure:"env"`
	// HTTP 服务监听端口
	Port int `mapstructure:"port"`
}

// ServerConfig HTTP 服务器配置
type ServerConfig struct {
	// 优雅关闭超时时间
	// 收到关闭信号后，等待正在处理的请求完成的最大时间
	ShutdownTimeout time.Duration `mapstructure:"shutdown_timeout"`
	// 读取请求的超时时间（包括请求头和请求体）
	ReadTimeout time.Duration `mapstructure:"read_timeout"`
	// 写入响应的超时时间
	WriteTimeout time.Duration `mapstructure:"write_timeout"`
}

// GetShutdownTimeout 获取优雅关闭超时时间，提供默认值
func (c *ServerConfig) GetShutdownTimeout() time.Duration {
	if c.ShutdownTimeout <= 0 {
		return 10 * time.Second
	}
	return c.ShutdownTimeout
}

// GetReadTimeout 获取读取超时时间，提供默认值
func (c *ServerConfig) GetReadTimeout() time.Duration {
	if c.ReadTimeout <= 0 {
		return 30 * time.Second
	}
	return c.ReadTimeout
}

// GetWriteTimeout 获取写入超时时间，提供默认值
func (c *ServerConfig) GetWriteTimeout() time.Duration {
	if c.WriteTimeout <= 0 {
		return 30 * time.Second
	}
	return c.WriteTimeout
}

// LogConfig 日志配置
type LogConfig struct {
	// 日志级别：debug | info | warn | error
	Level string `mapstructure:"level"`
	// 日志格式：text | json
	Format string `mapstructure:"format"`
}

// DatabaseConfig 数据库配置
type DatabaseConfig struct {
	// 数据库驱动：postgres | mysql | sqlite
	Driver string `mapstructure:"driver"`
	// 数据库主机地址
	Host string `mapstructure:"host"`
	// 数据库端口
	Port int `mapstructure:"port"`
	// 数据库名称
	Database string `mapstructure:"database"`
	// 数据库用户名
	Username string `mapstructure:"username"`
	// 数据库密码（敏感信息，建议通过环境变量覆盖）
	Password string `mapstructure:"password"`

	// 连接池配置
	// 最大空闲连接数
	MaxIdleConns int `mapstructure:"max_idle_conns"`
	// 最大打开连接数
	MaxOpenConns int `mapstructure:"max_open_conns"`
	// 连接最大存活时间
	ConnMaxLifetime time.Duration `mapstructure:"conn_max_lifetime"`
}

// DSN 生成数据库连接字符串
func (c *DatabaseConfig) DSN() string {
	switch c.Driver {
	case "postgres":
		return fmt.Sprintf(
			"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
			c.Host, c.Port, c.Username, c.Password, c.Database,
		)
	case "mysql":
		return fmt.Sprintf(
			"%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
			c.Username, c.Password, c.Host, c.Port, c.Database,
		)
	case "sqlite":
		return c.Database
	default:
		return ""
	}
}

// RedisConfig Redis 配置
type RedisConfig struct {
	// Redis 主机地址
	Host string `mapstructure:"host"`
	// Redis 端口
	Port int `mapstructure:"port"`
	// Redis 密码（敏感信息，建议通过环境变量覆盖）
	Password string `mapstructure:"password"`
	// 数据库索引
	DB int `mapstructure:"db"`
}

// Addr 返回 Redis 连接地址
func (c *RedisConfig) Addr() string {
	return fmt.Sprintf("%s:%d", c.Host, c.Port)
}

// JWTConfig JWT 认证配置
type JWTConfig struct {
	// JWT 签名密钥（敏感信息，必须通过环境变量覆盖）
	Secret string `mapstructure:"secret"`
	// Token 过期时间
	ExpiresIn time.Duration `mapstructure:"expires_in"`
}

// LoadConfig 加载应用配置
// 配置加载优先级（从低到高）：
// 1. 配置文件默认值
// 2. 配置文件中的值
// 3. 环境变量覆盖
//
// 环境变量命名规则：将配置路径中的 "." 替换为 "_"，并全部大写
// 例如：database.password -> DATABASE_PASSWORD
func LoadConfig(configPath string) (*Config, error) {
	v := viper.New()

	// 设置配置文件路径
	v.SetConfigFile(configPath)

	// 设置环境变量前缀（可选，避免与系统环境变量冲突）
	// 设置后，环境变量需要加上前缀：GO_API_DATABASE_PASSWORD
	// 如果不需要前缀，可以注释掉这行
	// v.SetEnvPrefix("GO_API")

	// 自动绑定环境变量
	v.AutomaticEnv()

	// 将配置路径中的 "." 替换为 "_"，以便环境变量可以覆盖嵌套配置
	// 例如：database.password 可以被 DATABASE_PASSWORD 覆盖
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	// 读取配置文件
	if err := v.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	// 将配置映射到结构体
	var config Config
	if err := v.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return &config, nil
}

// IsDevelopment 判断是否为开发环境
func (c *Config) IsDevelopment() bool {
	return c.App.Env == "development"
}

// IsProduction 判断是否为生产环境
func (c *Config) IsProduction() bool {
	return c.App.Env == "production"
}
