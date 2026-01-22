# data - 数据层 (Infrastructure Layer)

## 职责

这是整洁架构的**基础设施层**，负责：
- **实现 biz 定义的 Repository 接口**
- **管理数据库连接**
- **管理缓存连接 (Redis 等)**
- **封装外部 API 调用**

## 依赖规则

此层**依赖 biz 层**，实现 biz 定义的接口。

- 允许 import `internal/biz`
- 允许 import 数据库驱动 (gorm, ent, sqlx)
- 允许 import 缓存驱动 (redis)
- 不允许 import `internal/service`
- 不允许 import `internal/server`

## 依赖方向

```
data --> biz (实现接口)
```

## 示例结构

```
data/
├── data.go        # Wire Provider, 数据库/缓存连接初始化
├── user.go        # userRepo 结构体，实现 biz.UserRepo
└── order.go       # orderRepo 结构体，实现 biz.OrderRepo
```

## 关键模式

data 层通过**依赖注入**将具体实现注入到 biz 层：

```go
// biz/user.go - 定义接口
type UserRepo interface {
    GetByID(ctx context.Context, id int64) (*User, error)
}

// data/user.go - 实现接口
type userRepo struct {
    db *gorm.DB
}

func (r *userRepo) GetByID(ctx context.Context, id int64) (*biz.User, error) {
    // 具体的数据库查询实现
}
```
