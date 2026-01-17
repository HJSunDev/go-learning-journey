# 007. Go 映射（Map）：键值对的艺术

本章核心目标：**掌握 Map 的创建与操作，理解 comma-ok 模式和键类型要求**。

---

## 1. 概述

**Map（映射）** 是 Go 内置的键值对数据结构，类似其他语言中的 Dictionary、HashMap。

| 特性       | 说明                                   |
| ---------- | -------------------------------------- |
| 类型       | **引用类型**（赋值共享底层数据） |
| 元素顺序   | **无序**（遍历顺序随机）         |
| 键的唯一性 | 每个键只能出现一次                     |
| 零值       | `nil`（nil map 可读不可写）          |
| 键类型要求 | 必须是**可比较类型**             |

**一句话理解**：Map 就是"用任意类型查值"的数据结构。

```go
// 语法：map[KeyType]ValueType
scores := map[string]int{
    "Alice": 95,
    "Bob":   87,
}
```

---

## 2. Map 的创建

### 2.1 字面量创建

最常用的创建方式：

```go
// 带初始值
capitals := map[string]string{
    "China":  "Beijing",
    "Japan":  "Tokyo",
    "France": "Paris",
}

// 空 map
emptyMap := map[string]int{}
```

### 2.2 make 创建与 make 函数详解

`make` 是 Go 的内置函数，专门用于创建 **Slice（切片）**、**Map（映射）** 和 **Channel（通道）**。它会根据第一个参数的类型，决定后续参数的含义。

#### make 的参数差异

`make` 在创建 Slice 和 Map 时参数不同：

| 类型            | 完整语法                | 参数解释                                                                         |
| :-------------- | :---------------------- | :------------------------------------------------------------------------------- |
| **Slice** | `make([]T, len, cap)` | **3个参数**：类型、初始长度、底层容量。`<br>`例如 `make([]int, 0, 10)` |
| **Map**   | `make(map[K]V, cap)`  | **2个参数**：类型、初始容量提示。`<br>`例如 `make(map[string]int, 10)` |

#### Map 的容量（Capacity）代表什么？

对于 Map，`make` 的第二个参数代表**初始建议的键值对数量**（Initial Capacity Hint）。

```go
// 意思是：请预分配足够容纳 1000 个键值对的内存空间
userMap := make(map[string]int, 1000)
```

**核心解惑**：

> "Key Value 大小不确定，怎么设置容量？容量又代表什么？"

1. **容量是"个数"，不是"字节数"**：
   这里的 `1000` 仅仅代表你打算往 Map 里存 **1000 对**数据。Go 运行时不管你的 Key 是小整数还是大字符串，它主要关心的是要准备多少个**哈希桶（Buckets）**来存放这些数据，以减少哈希冲突。
2. **为什么要设置容量？**
   Map 在数据量增长时会自动扩容。扩容是一个昂贵的操作，需要重新申请更大的内存，并把旧数据重新计算哈希值搬迁到新内存中（Rehash）。

   * **不设容量**：存入第 1 个元素时分配一点，存到第 9 个扩容一次，存到第 100 个又扩容... 产生多次搬迁开销。
   * **设置容量**：如果你预知有 1000 个用户，直接 `make(..., 1000)`，Go 会一步到位分配好足够容纳 1000 个元素的桶。后续存入这 1000 个元素时，**不会发生任何扩容**，性能最高。
3. **注意：Map 没有 cap()**

   * Slice 可以用 `cap(s)` 查看容量。
   * Map **不能**用 `cap(m)` 查看容量（编译错误）。Go 隐藏了 Map 的底层扩容细节。

**总结**：如果你知道大概有多少数据，请在 `make` 时传入这个数字；如果不知道，忽略第二个参数即可（Go 会自动管理）。

### 2.3 nil Map（零值）

**什么是 nil Map？**

在 Go 中，所有变量在声明但未赋值时，都会被初始化为**零值**。

- `int` 的零值是 `0`
- `string` 的零值是 `""`
- **引用类型（如 Map、Slice、Pointer）的零值是 `nil`**

所以，当你声明一个 map 但不初始化时，它就是 `nil`：

```go
var m map[string]int  // 声明变量 m，它的值是 nil
// 注意：m 只是个变量名，nilMap 只是为了演示取的变量名，不是关键字！
```

这和 `var s []int`（nil 切片）是一样的道理。

**nil Map 的用途**

虽然 nil map **不能写入**（会 panic），但它在以下场景非常有用：

1. **作为只读的空容器**：
   nil map 可以安全地进行读取、删除（静默忽略）、获取长度（返回 0）、遍历（不执行循环体）。这意味着你不需要在使用 map 前判断 `if m != nil`。

   ```go
   var m map[string]int
   fmt.Println(len(m)) // 0
   fmt.Println(m["a"]) // 0
   ```
2. **延迟初始化**：

   `var m map[...]` 仅仅声明了一个变量（本质是一个 nil 指针），几乎不占用内存。
   `make(...)` 才会真正分配底层哈希表所需的内存空间。

   如果不确定是否会用到 map，可以先声明 `var`，只有在真正需要写入数据时才调用 `make`。

   ```go
   var m map[string]int // 1. 声明：此时 m 是 nil，只占一个指针大小

   // ... 一些业务逻辑 ...

   if needToStoreData {
       m = make(map[string]int) // 2. 初始化：分配底层哈希表内存，并将地址赋值给 m
       m["key"] = 100           // 3. 写入：现在 m 指向真实的 map，可以写入了
   }
   ```

   **优势**：如果 `needToStoreData` 为 false，就完全避免了哈希表的内存分配开销。
3. **作为函数返回值表示"无数据"**：
   当函数出错或没有数据可返回时，返回 `nil` 是惯用做法。

   ```go
   func getUsers() map[string]int {
       if dbError {
           return nil // 返回 nil map，调用方可以正常读取 len() 得到 0
       }
       return make(map[string]int)
   }
   ```
4. **JSON 序列化差异**：

   - `nil` map 序列化为 JSON 的 `null`
   - 空 map (`make(...)`) 序列化为 JSON 的 `{}`

| 操作   | nil Map 行为        |
| ------ | ------------------- |
| 读取   | ✅ 返回值类型的零值 |
| 写入   | ❌**panic**   |
| len()  | ✅ 返回 0           |
| delete | ✅ 静默忽略         |

**关键记忆**：nil map 可读不可写！

```go
var m map[string]int
fmt.Println(m["key"])  // ✅ 输出 0
m["key"] = 1           // ❌ panic: assignment to entry in nil map
```

---

## 3. Map 基本操作

### 3.1 添加 / 修改元素

```go
inventory := map[string]int{}

inventory["apple"] = 50   // 添加
inventory["banana"] = 30  // 添加
inventory["apple"] = 100  // 修改（键已存在）
```

### 3.2 读取元素

```go
count := inventory["apple"]  // 50
```

读取不存在的键会返回**值类型的零值**：

```go
fmt.Println(inventory["grape"])  // 0（int 的零值）
```

### 3.3 删除元素

使用内置函数 `delete(map, key)`：

```go
delete(inventory, "banana")  // 删除存在的键
delete(inventory, "mango")   // 删除不存在的键，不报错
```

### 3.4 获取长度

```go
fmt.Println(len(inventory))  // 元素数量
```

### 3.5 清空 Map（Go 1.21+）

```go
clear(inventory)  // 清空所有元素
```

---

## 4. 键存在性检查（comma-ok 模式）⭐

这是 Map 操作中**最重要**的模式。

### 4.1 问题：无法区分"键不存在"和"值为零值"

```go
scores := map[string]int{
    "Alice": 95,
    "Bob":   0,   // Bob 的分数确实是 0
}

fmt.Println(scores["Bob"])   // 0（真实值）
fmt.Println(scores["Eve"])   // 0（键不存在，返回零值）
```

### 4.2 解决方案：comma-ok 模式

Go 提供了双返回值语法来检查键是否存在：

```go
value, ok := map[key]
```

- `value`：键对应的值（键不存在时为零值）
- `ok`：布尔值，`true` 表示键存在，`false` 表示不存在

```go
// 检查并获取值
if score, ok := scores["Alice"]; ok {
    fmt.Println("Alice 的分数:", score)
} else {
    fmt.Println("Alice 不存在")
}
```

### 4.3 使用 `_` 只检查存在性

当不需要值，只想知道键是否存在时：

```go
// 使用空白标识符 _ 忽略值
if _, exists := scores["Alice"]; exists {
    fmt.Println("Alice 在名单中")
}

if _, exists := scores["Eve"]; !exists {
    fmt.Println("Eve 不在名单中")
}
```

**核心记忆**：

- `value, ok := m[key]` → 需要值
- `_, ok := m[key]` → 只检查存在性

---

## 5. Map 遍历

### 5.1 遍历键值对

```go
for key, value := range capitals {
    fmt.Printf("%s → %s\n", key, value)
}
```

### 5.2 只遍历键

```go
for key := range capitals {
    fmt.Println(key)
}
```

### 5.3 只遍历值

```go
for _, value := range capitals {
    fmt.Println(value)
}
```

### 5.4 ⚠️ 遍历顺序是随机的

Go **故意**使 Map 遍历顺序随机化，每次运行可能不同：

```go
for k := range m {
    fmt.Println(k)  // 顺序每次可能不同
}
```

**如需有序遍历**，先收集键并排序：

```go
import "slices"

keys := make([]string, 0, len(m))
for k := range m {
    keys = append(keys, k)
}
slices.Sort(keys)

for _, k := range keys {
    fmt.Println(k, m[k])
}
```

---

## 6. Map 键类型要求

Map 的键必须是**可比较类型**（可以用 `==` 比较的类型）。

### ✅ 可作为键的类型

| 类型     | 示例                                          |
| -------- | --------------------------------------------- |
| 基本类型 | `int`, `string`, `float64`, `bool`    |
| 数组     | `[3]int`, `[2]string`（长度是类型一部分） |
| 结构体   | `struct{ X, Y int }`（所有字段可比较）      |
| 指针     | `*int`, `*User`                           |
| 接口     | `interface{}`（动态值必须可比较）           |

```go
// 数组作为键
pointNames := map[[2]int]string{
    {0, 0}: "origin",
    {1, 2}: "point A",
}

// 结构体作为键
type Coord struct{ X, Y int }
locations := map[Coord]string{
    {0, 0}: "起点",
    {5, 5}: "终点",
}
```

### ❌ 不能作为键的类型

| 类型 | 原因     |
| ---- | -------- |
| 切片 | 不可比较 |
| Map  | 不可比较 |
| 函数 | 不可比较 |

```go
// ❌ 编译错误
// invalidMap := map[[]int]string{}
```

---

## 7. maps 标准库（Go 1.21+）

Go 1.21 引入了 `maps` 包，提供常用的 Map 操作函数。

```go
import "maps"
```

### 7.1 maps.Equal - 比较两个 Map

Map **不能**用 `==` 直接比较，必须使用 `maps.Equal`：

```go
m1 := map[string]int{"a": 1, "b": 2}
m2 := map[string]int{"a": 1, "b": 2}

maps.Equal(m1, m2)  // true

// m1 == m2  // ❌ 编译错误！Map 只能与 nil 比较
```

### 7.2 maps.Clone - 深拷贝

```go
original := map[string]int{"x": 10, "y": 20}
cloned := maps.Clone(original)

cloned["x"] = 999
// original["x"] 仍是 10
```

### 7.3 maps.Copy - 复制到目标 Map

```go
dest := map[string]int{"a": 1}
src := map[string]int{"b": 2, "c": 3}

maps.Copy(dest, src)
// dest: {"a": 1, "b": 2, "c": 3}
```

### 7.4 maps.DeleteFunc - 按条件删除

```go
numbers := map[string]int{"one": 1, "two": 2, "three": 3, "four": 4}

// 删除所有偶数值
maps.DeleteFunc(numbers, func(k string, v int) bool {
    return v % 2 == 0
})
// numbers: {"one": 1, "three": 3}
```

---

## 8. Map 陷阱与最佳实践

### 陷阱 1：nil Map 写入会 panic

```go
var m map[string]int
m["key"] = 1  // ❌ panic!
```

**解决方案**：始终使用 `make` 或字面量初始化。

### 陷阱 2：Map 不能用 == 比较

```go
m1 := map[string]int{"a": 1}
m2 := map[string]int{"a": 1}
// m1 == m2  // ❌ 编译错误
```

**解决方案**：使用 `maps.Equal(m1, m2)`。

### 陷阱 3：遍历顺序不固定

Go 故意随机化遍历顺序，不要依赖遍历顺序。

**解决方案**：如需有序，先排序键。

### 陷阱 4：并发读写会 panic

多个 goroutine 同时读写同一个 Map 会导致 panic。

**解决方案**：

- 使用 `sync.RWMutex` 保护
- 使用 `sync.Map`（适合读多写少场景）

### 陷阱 5：不能对 Map 元素取地址

```go
m := map[string]int{"a": 1}
// ptr := &m["a"]  // ❌ 编译错误
```

**原因**：Map 内部可能重新分配内存，地址会失效。

---

## 9. 实际使用场景

### 场景 1：词频统计

```go
words := []string{"apple", "banana", "apple", "cherry"}
count := make(map[string]int)

for _, word := range words {
    count[word]++  // 利用零值特性
}
// {"apple": 2, "banana": 1, "cherry": 1}
```

### 场景 2：切片去重

```go
numbers := []int{1, 2, 2, 3, 3, 3}
seen := make(map[int]bool)
unique := []int{}

for _, n := range numbers {
    if !seen[n] {
        seen[n] = true
        unique = append(unique, n)
    }
}
// unique: [1 2 3]
```

### 场景 3：分组

```go
words := []string{"apple", "banana", "apricot", "cherry"}
groups := make(map[byte][]string)

for _, word := range words {
    firstChar := word[0]
    groups[firstChar] = append(groups[firstChar], word)
}
// {'a': ["apple", "apricot"], 'b': ["banana"], 'c': ["cherry"]}
```

### 场景 4：集合（Set）实现

Go 没有内置 Set，用 `map[T]struct{}` 模拟：

```go
type Set map[string]struct{}

set := make(Set)
set["apple"] = struct{}{}
set["banana"] = struct{}{}

// 检查存在性
_, exists := set["apple"]  // true
```

**为什么用 `struct{}` 而不是 `bool`？** `struct{}` 占用 0 字节内存。

### 场景 5：缓存

```go
cache := make(map[string]string)

func getData(key string) string {
    if val, ok := cache[key]; ok {
        return val  // 缓存命中
    }
    // 缓存未命中，从数据源获取
    result := fetchFromDB(key)
    cache[key] = result
    return result
}
```

---

## 10. 本章小结

| 概念          | 要点                                                  |
| ------------- | ----------------------------------------------------- |
| 创建方式      | 字面量 `map[K]V{}` 或 `make(map[K]V)`             |
| 基本操作      | 读取、添加/修改、`delete()`、`len()`、`clear()` |
| comma-ok 模式 | `value, ok := m[key]` 检查键是否存在                |
| 空白标识符    | `_, ok := m[key]` 只检查存在性                      |
| 遍历          | `for k, v := range m`（顺序随机）                   |
| 键类型要求    | 必须是可比较类型（不能是切片、Map、函数）             |
| maps 包       | `Equal`、`Clone`、`Copy`、`DeleteFunc`        |
| nil Map       | 可读不可写，写入会 panic                              |

**核心记忆**：

1. Map 是**引用类型**，赋值共享数据
2. `value, ok := m[key]` 是检查键存在的**标准模式**
3. 遍历顺序**不固定**，不要依赖顺序
4. nil Map 可读不可写，**必须先初始化**
5. 用 `maps.Equal` 比较 Map，不能用 `==`

---

## 11. 动手练习

1. 创建一个 `map[string]int`，存储 5 种水果的库存，练习增删改查操作
2. 使用 comma-ok 模式检查某个水果是否存在
3. 统计一段文本中每个单词出现的次数
4. 使用 Map 实现切片去重
5. 创建一个使用结构体作为键的 Map，存储二维坐标点的名称
