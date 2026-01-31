# 010. Struct Tag 与 Validator 请求验证

## 1. 从一个前端场景说起

假设你在做一个注册页面，用户填写表单后提交。作为前端，你会怎么验证？

```javascript
// 前端表单验证（以 Element UI 为例）
const rules = {
  username: [
    { required: true, message: '请输入用户名', trigger: 'blur' },
    { min: 2, max: 20, message: '长度在 2 到 20 个字符', trigger: 'blur' }
  ],
  email: [
    { required: true, message: '请输入邮箱', trigger: 'blur' },
    { type: 'email', message: '请输入正确的邮箱格式', trigger: 'blur' }
  ],
  age: [
    { type: 'number', min: 0, max: 150, message: '年龄必须在 0-150 之间' }
  ]
}
```

**后端也需要做同样的验证**（永远不要相信前端传来的数据）。Go 语言怎么做呢？

答案是：**把验证规则直接写在数据结构的定义上**。

```go
type RegisterRequest struct {
    Username string `binding:"required,min=2,max=20"`
    Email    string `binding:"required,email"`
    Age      int    `binding:"gte=0,lte=150"`
}
```

看到 `binding:"required,min=2,max=20"` 了吗？这就是 Go 的**验证规则声明方式**。

但要理解这个写法，我们需要先了解一个 Go 的基础概念：**Struct Tag**。

---

## 2. 什么是 Struct Tag

### 2.1 先看 Struct 是什么

在 Go 里，`struct` （结构体）就是**一组字段的集合**，类似于：

- JavaScript 的对象 `{ name: "张三", age: 18 }`
- Go中的结构体是 值类型，属于原始类型，值传递，修改需要使用指针： &结构体

**定义一个结构体** = 声明这个结构体有哪些字段，每个字段是什么类型：

```go
// 定义一个 User 结构体
// 意思是：User 这种类型的数据，包含 Name（字符串）和 Age（整数）两个字段
type User struct {
    Name string  // 字段名 类型
    Age  int     // 字段名 类型
}
```

**创建一个结构体变量**（也叫"实例化"）：

```go
// 方式一：直接赋值
user := User{
    Name: "张三",
    Age:  18,
}

// 方式二：先声明，再赋值
var user2 User          // 声明一个 User 类型的变量，叫 user2
user2.Name = "李四"      // 给字段赋值
user2.Age = 20
```

**关于 `var req LoginRequest` 这行代码**：

```go
var req LoginRequest
// 拆解：
// var           = 声明一个变量
// req           = 变量的名字
// LoginRequest  = 变量的类型

// 类比 TypeScript：
// let req: LoginRequest;
```

声明之后，`req` 里的每个字段都是"零值"：

- 字符串字段是空字符串 `""`
- 数字字段是 `0`
- 布尔字段是 `false`

### 2.2 什么是 Tag

Tag 是**附加在字段后面的一段文字**，用反引号 `` ` `` 包裹。

```go
type User struct {
    Name string `这里就是 Tag`
    Age  int    `这里也是 Tag`
}
```

**Tag 本身不会影响代码的运行**，它只是一段"备注信息"。但是，**某些库会读取这些备注信息，然后根据备注做特定的事情**。

### 2.3 Tag 的格式

Tag 的标准格式是 `key:"value"`，多个 key 用空格分隔：

```go
type User struct {
    Name string `json:"name" binding:"required"`
    //                        ^^^^^^^^^^^  ^^^^^^^^^^^^^^^^^
    //                        第一个 Tag       第二个 Tag
    //                        key=json          key=binding
    //                        value=name   value=required
}
```

### 2.4 不同的库读取不同的 Tag

| Tag 名称          | 哪个库在用                  | 作用                       |
| ----------------- | --------------------------- | -------------------------- |
| `json:"xxx"`    | Go 标准库 `encoding/json` | 控制 JSON 序列化时的字段名 |
| `binding:"xxx"` | Gin 框架                    | 请求数据绑定时的验证规则   |
| `yaml:"xxx"`    | yaml 库                     | 控制 YAML 序列化时的字段名 |
| `gorm:"xxx"`    | GORM 库                     | 数据库字段映射和约束       |

**你可以把 Tag 理解为"给不同库看的配置说明"**。

---

## 3. json Tag 详解

`json` Tag 是最常用的，它告诉 Go 的 JSON 库：**这个字段在 JSON 里叫什么名字**。

### 3.1 先解释两个词：Marshal 和 Unmarshal

这两个词你会经常看到：

| 术语                | 含义                               | 类比JavaScript          |
| ------------------- | ---------------------------------- | ----------------------- |
| **Marshal**   | 把 Go 数据结构 → 转成 JSON 字符串 | `JSON.stringify(obj)` |
| **Unmarshal** | 把 JSON 字符串 → 转成 Go 数据结构 | `JSON.parse(str)`     |

```go
// Marshal = 序列化 = stringify
user := User{Name: "张三"}
jsonBytes, _ := json.Marshal(user)  // 得到 []byte 类型的 JSON 数据
// jsonBytes 的内容是：{"name":"张三"}

// Unmarshal = 反序列化 = parse
var user2 User
json.Unmarshal(jsonBytes, &user2)   // 把 JSON 填充到 user2 里
// 现在 user2.Name = "张三"
```

### 3.2 为什么用 &user 而不是 user？

这是 Go 的一个重要概念：**指针**。

```go
var user User
json.Unmarshal(jsonBytes, &user)
//                                                            ^ 这个 & 是什么？
```

**简单解释**：

- `user` 是变量本身
- `&user` 是变量的"地址"（指针）

**为什么需要传地址？**

```go
// 错误写法（不会报错，但不起作用）
var user User
json.Unmarshal(jsonBytes, user)
// Unmarshal 收到的是 user 的"副本"
// 它修改的是副本，原来的 user 没变

// 正确写法
var user User
json.Unmarshal(jsonBytes, &user)
// Unmarshal 收到的是 user 的"地址"
// 它通过地址找到原来的 user，直接修改它
```

**类比前端**：

```javascript
// JavaScript 中，对象是引用传递，所以不需要 &
function fillData(obj) {
    obj.name = "张三";  // 直接修改原对象
}
let user = {};
fillData(user);  // user.name 变成了 "张三"

// 但如果是基本类型，就修改不了
function fillNumber(n) {
    n = 100;  // 修改的是副本
}
let num = 0;
fillNumber(num);  // num 还是 0
```

Go 的 struct 默认是"值传递"（像 JS 的基本类型），所以需要用 `&` 来传地址。

**记住这个规则**：当你想让一个函数"填充"或"修改"你的变量时，传 `&变量`。

### 3.3 基本用法

```go
type User struct {
    UserName string `json:"user_name"`
    Age      int    `json:"age"`
}
```

当你把这个 struct 转成 JSON 时：

```go
user := User{UserName: "张三", Age: 18}
jsonBytes, _ := json.Marshal(user)   // Marshal = Go对象 转 JSON
fmt.Println(string(jsonBytes))
// 输出: {"user_name":"张三","age":18}
// 注意：是 "user_name"，不是 "UserName"
```

当你把 JSON 转成 struct 时：

```go
jsonStr := `{"user_name":"李四","age":20}`
var user User
json.Unmarshal([]byte(jsonStr), &user)  // Unmarshal = JSON 转 Go对象
// 注意：传的是 &user，因为要让函数"填充"这个变量
fmt.Println(user.UserName) // 输出: 李四
```

### 3.4 为什么需要 json Tag？

**命名习惯不同**：

- Go 语言习惯用大驼峰：`UserName`、`CreatedAt`
- JSON/JavaScript 习惯用小驼峰或下划线：`userName`、`created_at`

json Tag 就是用来做这个转换的。

### 3.5 常用的 json Tag 选项

```go
type User struct {
    // 基本映射：JSON 字段名为 "name"
    Name string `json:"name"`
  
    // omitempty：如果字段是零值（空字符串、0、nil 等），JSON 中不包含这个字段
    Email string `json:"email,omitempty"`
  
    // -：完全忽略这个字段，不参与 JSON 转换
    Password string `json:"-"`
}
```

**omitempty 的效果**：

```go
user := User{Name: "张三", Email: "", Password: "123456"}
jsonBytes, _ := json.Marshal(user)
// 输出: {"name":"张三"}
// Email 是空字符串，加了 omitempty 所以不输出
// Password 加了 -，完全不输出
```

---

## 4. binding Tag 与 Validator 库

### 4.1 Gin 和 Validator 的关系

**Validator** 是 Go 生态中最流行的验证库，全名是 `go-playground/validator`。

**Gin 框架内部集成了 Validator**，当你调用 `c.ShouldBindJSON()` 时，Gin 会：

1. 把请求的 JSON 数据填充到你的 struct
2. 读取 `binding` Tag
3. 调用 Validator 库进行验证
4. 验证失败就返回错误

所以 `binding` Tag 里写的规则，实际上是 Validator 库的规则。

### 4.2 c.ShouldBindJSON 是什么？

在 Gin 的 Handler 函数里，`c` 是 Gin 提供的"上下文对象"，它包含了这次请求的所有信息。

```go
func handleLogin(c *gin.Context) {
    // c 就是这次请求的上下文
    // c.ShouldBindJSON() 是 c 的一个方法
}
```

**ShouldBindJSON 做了什么？**

```go
var req LoginRequest
err := c.ShouldBindJSON(&req)
```

这一行代码做了三件事：

| 步骤          | 说明                              | 类比前端             |
| ------------- | --------------------------------- | -------------------- |
| 1. 读取请求体 | 从 HTTP 请求中读取 JSON 数据      | `request.body`     |
| 2. 解析 JSON  | 把 JSON 字符串转成 Go struct      | `JSON.parse(body)` |
| 3. 验证数据   | 根据 `binding` Tag 验证每个字段 | 表单验证             |

**为什么传 `&req`？**

和前面解释的一样：ShouldBindJSON 需要"填充"req 这个变量，所以要传它的地址。

```go
var req LoginRequest           // req 现在是空的
c.ShouldBindJSON(&req)         // 把请求数据填充到 req 里
// 现在 req.Username 和 req.Password 有值了
```

**err 是什么？**

如果 JSON 解析失败，或者验证不通过，`err` 就不是 `nil`：

```go
if err := c.ShouldBindJSON(&req); err != nil {
    // err != nil 表示出错了
    // err.Error() 可以获取错误信息
    c.JSON(400, gin.H{"error": err.Error()})
    return
}
// 走到这里说明验证通过
```

### 4.3 第一个验证示例

```go
// 定义请求结构体，带验证规则
type LoginRequest struct {
    Username string `json:"username" binding:"required"`
    Password string `json:"password" binding:"required,min=6"`
}

// 在 Handler 中使用
func handleLogin(c *gin.Context) {
    var req LoginRequest
  
    // ShouldBindJSON 会自动验证
    if err := c.ShouldBindJSON(&req); err != nil {
        // 验证失败，err 包含具体哪个字段、什么规则没通过
        c.JSON(400, gin.H{"error": err.Error()})
        return
    }
  
    // 验证通过，req 里的数据是可信的
    fmt.Println(req.Username, req.Password)
}
```

**测试效果**：

```bash
# 缺少 username
curl -X POST -d '{"password":"123456"}' http://localhost:8080/login
# 返回: {"error":"Key: 'LoginRequest.Username' Error:Field validation for 'Username' failed on the 'required' tag"}

# 密码太短
curl -X POST -d '{"username":"admin","password":"123"}' http://localhost:8080/login
# 返回: {"error":"Key: 'LoginRequest.Password' Error:Field validation for 'Password' failed on the 'min' tag"}

# 验证通过
curl -X POST -d '{"username":"admin","password":"123456"}' http://localhost:8080/login
# 正常处理
```

---

## 5. Validator 验证规则大全

### 5.1 必填与可选

| 规则          | 含义                                           | 示例                          |
| ------------- | ---------------------------------------------- | ----------------------------- |
| `required`  | 必填，不能是零值                               | `binding:"required"`        |
| `omitempty` | 如果是零值，跳过**该字段**的后续验证规则 | `binding:"omitempty,min=2"` |

**什么是零值？**

- 字符串的零值是 `""`（空字符串）
- 数字的零值是 `0`
- 布尔的零值是 `false`
- 指针/切片/map 的零值是 `nil`

**omitempty 的作用范围**：

`omitempty` 只影响**当前这个字段**，不影响其他字段的验证。

```go
type Request struct {
    // 必填：不能是空字符串
    Name string `binding:"required"`
  
    // 可选字段：
    // - 如果 Nickname 是空字符串，跳过 min=2,max=20 的验证，直接通过
    // - 如果 Nickname 有值，才执行 min=2,max=20 的验证
    // - 无论 Nickname 怎样，Name 的 required 验证照常执行
    Nickname string `binding:"omitempty,min=2,max=20"`
}
```

**举例说明**：

```go
// 请求1: {"name": "张三"}
// Name 验证通过（有值）
// Nickname 是空字符串，omitempty 生效，跳过 min/max 验证 → 整体通过

// 请求2: {"name": "张三", "nickname": "A"}
// Name 验证通过
// Nickname 有值，执行 min=2 验证 → 失败（长度只有1）

// 请求3: {"name": "张三", "nickname": "小明"}
// Name 验证通过
// Nickname 有值，执行 min=2,max=20 验证 → 通过
```

### 5.2 字符串长度

| 规则      | 含义     | 示例                  |
| --------- | -------- | --------------------- |
| `min=n` | 最小长度 | `binding:"min=2"`   |
| `max=n` | 最大长度 | `binding:"max=100"` |
| `len=n` | 精确长度 | `binding:"len=11"`  |

```go
type Request struct {
    // 用户名：2-20 个字符
    Username string `binding:"required,min=2,max=20"`
  
    // 手机号：必须 11 位
    Phone string `binding:"required,len=11"`
}
```

### 5.3 数值范围

| 规则      | 含义                             | 示例                  |
| --------- | -------------------------------- | --------------------- |
| `gt=n`  | 大于 (greater than)              | `binding:"gt=0"`    |
| `gte=n` | 大于等于 (greater than or equal) | `binding:"gte=0"`   |
| `lt=n`  | 小于 (less than)                 | `binding:"lt=100"`  |
| `lte=n` | 小于等于 (less than or equal)    | `binding:"lte=100"` |
| `eq=n`  | 等于                             | `binding:"eq=1"`    |
| `ne=n`  | 不等于                           | `binding:"ne=0"`    |

```go
type Request struct {
    // 年龄：0-150
    Age int `binding:"gte=0,lte=150"`
  
    // 数量：必须大于 0
    Quantity int `binding:"gt=0"`
  
    // 折扣：0.1 到 1.0
    Discount float64 `binding:"gte=0.1,lte=1.0"`
}
```

### 5.4 格式验证

| 规则              | 含义         | 示例                              |
| ----------------- | ------------ | --------------------------------- |
| `email`         | 邮箱格式     | `binding:"email"`               |
| `url`           | URL 格式     | `binding:"url"`                 |
| `uri`           | URI 格式     | `binding:"uri"`                 |
| `uuid`          | UUID 格式    | `binding:"uuid"`                |
| `ip`            | IP 地址      | `binding:"ip"`                  |
| `ipv4`          | IPv4 地址    | `binding:"ipv4"`                |
| `ipv6`          | IPv6 地址    | `binding:"ipv6"`                |
| `datetime=格式` | 日期时间格式 | `binding:"datetime=2006-01-02"` |

```go
type Request struct {
    Email     string `binding:"required,email"`
    Website   string `binding:"omitempty,url"`
    Birthday  string `binding:"required,datetime=2006-01-02"`
    // 日期格式说明：Go 用 2006-01-02 15:04:05 作为格式模板
    // 这是 Go 的特殊约定，记住就行
}
```

### 5.5 枚举值

| 规则            | 含义             | 示例                            |
| --------------- | ---------------- | ------------------------------- |
| `oneof=a b c` | 必须是指定值之一 | `binding:"oneof=male female"` |

```go
type Request struct {
    // 性别：只能是 male 或 female
    Gender string `binding:"required,oneof=male female"`
  
    // 状态：只能是 0 1 2
    Status int `binding:"oneof=0 1 2"`
  
    // 排序：只能是 asc 或 desc
    Order string `binding:"omitempty,oneof=asc desc"`
}
```

### 5.6 字符串内容

| 规则               | 含义             | 示例                          |
| ------------------ | ---------------- | ----------------------------- |
| `alpha`          | 只包含字母       | `binding:"alpha"`           |
| `alphanum`       | 只包含字母和数字 | `binding:"alphanum"`        |
| `numeric`        | 只包含数字       | `binding:"numeric"`         |
| `lowercase`      | 只包含小写字母   | `binding:"lowercase"`       |
| `uppercase`      | 只包含大写字母   | `binding:"uppercase"`       |
| `contains=xxx`   | 包含指定字符串   | `binding:"contains=@"`      |
| `startswith=xxx` | 以指定字符串开头 | `binding:"startswith=http"` |
| `endswith=xxx`   | 以指定字符串结尾 | `binding:"endswith=.com"`   |

```go
type Request struct {
    // 用户名：只能是字母和数字
    Username string `binding:"required,alphanum,min=4,max=20"`
  
    // 邀请码：只能是大写字母
    InviteCode string `binding:"required,uppercase,len=6"`
}
```

### 5.7 跨字段验证

| 规则               | 含义                 | 示例                              |
| ------------------ | -------------------- | --------------------------------- |
| `eqfield=字段名` | 必须等于另一个字段   | `binding:"eqfield=Password"`    |
| `nefield=字段名` | 必须不等于另一个字段 | `binding:"nefield=OldPassword"` |
| `gtfield=字段名` | 必须大于另一个字段   | `binding:"gtfield=StartDate"`   |
| `ltfield=字段名` | 必须小于另一个字段   | `binding:"ltfield=EndDate"`     |

```go
type RegisterRequest struct {
    Password        string `binding:"required,min=6"`
    ConfirmPassword string `binding:"required,eqfield=Password"`
    // ConfirmPassword 必须等于 Password
}

type ChangePasswordRequest struct {
    OldPassword string `binding:"required"`
    NewPassword string `binding:"required,min=6,nefield=OldPassword"`
    // NewPassword 必须不等于 OldPassword
}

type DateRangeRequest struct {
    StartDate string `binding:"required,datetime=2006-01-02"`
    EndDate   string `binding:"required,datetime=2006-01-02,gtfield=StartDate"`
    // EndDate 必须大于 StartDate
}
```

### 5.8 切片/数组验证

| 规则      | 含义               | 示例                        |
| --------- | ------------------ | --------------------------- |
| `min=n` | 最少 n 个元素      | `binding:"min=1"`         |
| `max=n` | 最多 n 个元素      | `binding:"max=10"`        |
| `dive`  | 验证切片内每个元素 | `binding:"dive,required"` |

```go
type Request struct {
    // 标签列表：至少 1 个，最多 5 个
    Tags []string `binding:"required,min=1,max=5"`
  
    // ID 列表：每个 ID 都必须大于 0
    IDs []int `binding:"required,dive,gt=0"`
  
    // 邮箱列表：每个都必须是有效邮箱
    Emails []string `binding:"omitempty,dive,email"`
}
```

---

## 6. 组合使用示例

### 6.1 用户注册

```go
type RegisterRequest struct {
    // 用户名：必填，2-20字符，只能字母数字
    Username string `json:"username" binding:"required,min=2,max=20,alphanum"`
  
    // 邮箱：必填，邮箱格式
    Email string `json:"email" binding:"required,email"`
  
    // 密码：必填，8-72字符
    Password string `json:"password" binding:"required,min=8,max=72"`
  
    // 确认密码：必填，必须等于密码
    ConfirmPassword string `json:"confirm_password" binding:"required,eqfield=Password"`
  
    // 年龄：可选，如果填了必须 0-150
    Age *int `json:"age" binding:"omitempty,gte=0,lte=150"`
  
    // 性别：可选，如果填了必须是 male/female/other
    Gender string `json:"gender" binding:"omitempty,oneof=male female other"`
}
```

### 6.2 创建文章

```go
type CreateArticleRequest struct {
    // 标题：必填，2-100字符
    Title string `json:"title" binding:"required,min=2,max=100"`
  
    // 内容：必填，至少 10 字符
    Content string `json:"content" binding:"required,min=10"`
  
    // 分类：必填，必须是指定值之一
    Category string `json:"category" binding:"required,oneof=tech life travel food"`
  
    // 标签：可选，最多 5 个，每个标签 1-20 字符
    Tags []string `json:"tags" binding:"omitempty,max=5,dive,min=1,max=20"`
  
    // 是否发布：可选
    Published bool `json:"published"`
}
```

### 6.3 分页查询

```go
type ListRequest struct {
    // 页码：可选，默认会是 0，但如果传了必须 >= 1
    Page int `json:"page" binding:"omitempty,gte=1"`
  
    // 每页数量：可选，1-100
    PageSize int `json:"page_size" binding:"omitempty,gte=1,lte=100"`
  
    // 排序字段：可选，必须是指定值之一
    SortBy string `json:"sort_by" binding:"omitempty,oneof=created_at updated_at name"`
  
    // 排序方向：可选，必须是 asc 或 desc
    SortOrder string `json:"sort_order" binding:"omitempty,oneof=asc desc"`
  
    // 搜索关键词：可选，最多 100 字符
    Keyword string `json:"keyword" binding:"omitempty,max=100"`
}
```

---

## 7. 处理验证错误

当验证失败时，`ShouldBindJSON` 返回的错误信息是英文的、格式化的，不适合直接返回给用户。

### 7.1 原始错误信息

```go
if err := c.ShouldBindJSON(&req); err != nil {
    fmt.Println(err.Error())
    // 输出类似：
    // Key: 'RegisterRequest.Username' Error:Field validation for 'Username' failed on the 'required' tag
}
```

### 7.2 解析验证错误

Validator 的错误可以被解析成结构化的信息：

```go
import "github.com/go-playground/validator/v10"

if err := c.ShouldBindJSON(&req); err != nil {
    // 尝试转换为 validator 的错误类型
    var validationErrors validator.ValidationErrors
    if errors.As(err, &validationErrors) {
        // 遍历每个字段的错误
        for _, fieldError := range validationErrors {
            fmt.Printf("字段: %s, 规则: %s, 值: %v\n",
                fieldError.Field(),  // 字段名，如 "Username"
                fieldError.Tag(),    // 验证规则，如 "required"
                fieldError.Value(),  // 实际值
            )
        }
    }
}
```

### 7.3 生成友好的错误信息

```go
func translateValidationError(err error) []string {
    var messages []string
  
    var validationErrors validator.ValidationErrors
    if errors.As(err, &validationErrors) {
        for _, e := range validationErrors {
            field := e.Field()
            tag := e.Tag()
            param := e.Param()
      
            var msg string
            switch tag {
            case "required":
                msg = fmt.Sprintf("%s 不能为空", field)
            case "min":
                msg = fmt.Sprintf("%s 长度不能小于 %s", field, param)
            case "max":
                msg = fmt.Sprintf("%s 长度不能大于 %s", field, param)
            case "email":
                msg = fmt.Sprintf("%s 必须是有效的邮箱地址", field)
            case "oneof":
                msg = fmt.Sprintf("%s 必须是以下值之一: %s", field, param)
            default:
                msg = fmt.Sprintf("%s 验证失败", field)
            }
            messages = append(messages, msg)
        }
    }
  
    return messages
}
```

---

## 8. 总结

### 核心概念

| 概念              | 说明                                  |
| ----------------- | ------------------------------------- |
| Struct            | Go 的数据结构，类似 JS 的对象         |
| Struct Tag        | 写在字段后面的元信息，用反引号包裹    |
| `json:"xxx"`    | 告诉 JSON 库这个字段的 JSON 名称      |
| `binding:"xxx"` | 告诉 Gin/Validator 这个字段的验证规则 |

### 常用验证规则速查

| 类别       | 规则                                            |
| ---------- | ----------------------------------------------- |
| 必填       | `required`                                    |
| 字符串长度 | `min=n`, `max=n`, `len=n`                 |
| 数值范围   | `gt=n`, `gte=n`, `lt=n`, `lte=n`        |
| 格式       | `email`, `url`, `uuid`, `datetime=格式` |
| 枚举       | `oneof=a b c`                                 |
| 内容       | `alpha`, `alphanum`, `numeric`            |
| 跨字段     | `eqfield=字段`, `nefield=字段`              |
| 数组       | `dive`（验证每个元素）                        |
| 可选       | `omitempty`（为空时跳过验证）                 |

### 与前端对比

| 前端 (Element UI)           | Go (Validator)                |
| --------------------------- | ----------------------------- |
| `{ required: true }`      | `binding:"required"`        |
| `{ min: 2, max: 20 }`     | `binding:"min=2,max=20"`    |
| `{ type: 'email' }`       | `binding:"email"`           |
| `{ pattern: /^[a-z]+$/ }` | `binding:"alpha,lowercase"` |
| 写在单独的 rules 对象里     | 写在 struct 字段的 Tag 里     |
