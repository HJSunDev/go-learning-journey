# 018. Go 互斥锁（Mutex）：并发安全的守护者

[返回索引](../README.md) | [查看代码](../../018-mutex/main.go)

## 本章要解决的问题

**当多个协程需要同时读写共享数据时，如何保证数据的正确性？**

---

## 第一步：场景设定

你在开发订单系统的实时统计模块。每当有订单被处理，需要更新统计数据：

```go
// OrderStats 订单统计数据（共享资源）
type OrderStats struct {
    TotalOrders   int     // 总订单数
    TotalAmount   float64 // 总金额
    SuccessOrders int     // 成功订单数
    FailedOrders  int     // 失败订单数
}
```

现在有 1000 个订单需要并发处理，每个协程处理完后都要更新统计数据。

---

## 第二步：没有锁时的痛点 —— 数据竞争

### 问题代码

```go
func processOrderUnsafe(order Order, stats *OrderStats, wg *sync.WaitGroup) {
    defer wg.Done()
  
    time.Sleep(10 * time.Millisecond)
  
    // 多个协程同时执行这些代码
    stats.TotalOrders++
    stats.TotalAmount += order.Amount
    stats.SuccessOrders++
}
```

### 运行结果

```
期望订单数: 1000
实际统计数: 987   // 每次运行结果不同！
丢失订单数: 13
```

1000 个订单，统计出来却少了 13 个！这就是 **数据竞争（Race Condition）**。

### 问题根源

`stats.TotalOrders++` 看起来是一行代码，但在 CPU 层面是三个操作：

```
1. 读取 stats.TotalOrders 的当前值 → 放入寄存器
2. 将寄存器中的值加 1
3. 将新值写回 stats.TotalOrders
```

当两个协程同时执行时：

```
时间线          协程 A                 协程 B
────────────────────────────────────────────────
T1         读取值 = 5                
T2                                 读取值 = 5
T3         计算 5 + 1 = 6            
T4                                 计算 5 + 1 = 6
T5         写入值 = 6                
T6                                 写入值 = 6
────────────────────────────────────────────────
结果：两个订单只增加了 1，丢失一次计数
```

这种"读取-修改-写入"操作不是原子的，中间可能被其他协程打断。

---

## 第三步：理解锁的概念

### 什么是锁

锁是一种 **同步机制**，用于控制多个执行单元（协程/线程/进程）对共享资源的访问。

核心思想：在任意时刻，只允许一个执行单元访问被保护的资源。

### 锁的分类

"悲观锁"和"乐观锁"是按 **并发控制策略** 分类的，而互斥锁是按 **实现机制** 分类的。它们是不同维度的概念：

| 分类维度           | 类型   | 说明                               |
| ------------------ | ------ | ---------------------------------- |
| **并发策略** | 悲观锁 | 假设会发生冲突，先加锁再操作       |
|                    | 乐观锁 | 假设不会冲突，操作后检查并处理冲突 |
| **实现机制** | 互斥锁 | 通过锁变量实现互斥访问             |
|                    | 读写锁 | 允许多读单写                       |
|                    | 自旋锁 | 忙等待，不释放 CPU                 |

### 互斥锁属于哪一类

**互斥锁是悲观锁的一种实现方式**。

对比你提到的例子：

| 类型                  | 实现方式                  | 适用场景           |
| --------------------- | ------------------------- | ------------------ |
| 基于数据库表的悲观锁  | `SELECT ... FOR UPDATE` | 跨进程、分布式系统 |
| 基于 Redis 的分布式锁 | `SETNX` + 过期时间      | 跨服务、分布式系统 |
| Go 的 `sync.Mutex`  | 内存中的锁变量            | 单进程内多协程     |

**关键区别**：

- 数据库锁和 Redis 锁是 **分布式锁**，用于跨进程/跨服务的同步
- Go 的 `sync.Mutex` 是 **进程内锁**，只能保护同一进程内的多个协程

### 乐观锁的例子

乐观锁不使用锁机制，而是通过版本号或 CAS 操作：

```go
// 乐观锁思路（伪代码）
for {
    oldValue := stats.TotalOrders          // 读取当前值
    newValue := oldValue + 1               // 计算新值
  
    // CAS: 如果当前值仍然是 oldValue，才更新为 newValue
    if atomic.CompareAndSwap(&stats.TotalOrders, oldValue, newValue) {
        break  // 成功
    }
    // 失败说明被其他协程修改了，重试
}
```

Go 的 `sync/atomic` 包提供了原子操作，适合简单的计数场景，但复杂场景仍需互斥锁。

---

## 第四步：互斥锁的工作原理

### sync.Mutex 是什么

`sync.Mutex` 是 Go 标准库提供的互斥锁类型：

- **Mutex** = Mutual Exclusion（互斥）的缩写
- 保证同一时刻只有一个协程能进入临界区
- 零值可用，不需要初始化

### 核心方法

| 方法         | 作用                                     |
| ------------ | ---------------------------------------- |
| `Lock()`   | 获取锁。如果锁已被持有，当前协程阻塞等待 |
| `Unlock()` | 释放锁。让其他等待的协程可以获取锁       |

### 工作流程图解

```
协程A                     协程B                     锁状态
──────                   ──────                   ──────
Lock()                                            🔓→🔒 A持有
  修改数据                                       
                         Lock()                   B等待...
                         (阻塞)                 
Unlock()                                          🔒→🔓 释放
                         (获取锁)                  🔓→🔒 B持有
                           修改数据
                         Unlock()                 🔒→🔓 释放
```

---

## 第五步：互斥锁基础用法

### 基本模式

```go
var mu sync.Mutex

mu.Lock()     // 获取锁
// 临界区：同一时刻只有一个协程能执行这里的代码
// ... 操作共享资源 ...
mu.Unlock()   // 释放锁
```

### 完整示例

```go
// SafeOrderStats 线程安全的订单统计
type SafeOrderStats struct {
    mu sync.Mutex  // 互斥锁，保护下面的字段
  
    TotalOrders   int
    TotalAmount   float64
    SuccessOrders int
}

func processOrderSafe(order Order, stats *SafeOrderStats, wg *sync.WaitGroup) {
    defer wg.Done()
  
    time.Sleep(10 * time.Millisecond)
  
    // 获取锁
    stats.mu.Lock()
  
    // === 临界区 ===
    stats.TotalOrders++
    stats.TotalAmount += order.Amount
    stats.SuccessOrders++
  
    // 释放锁
    stats.mu.Unlock()
}
```

### 语法逐行拆解

```go
type SafeOrderStats struct {
    mu sync.Mutex  // 将锁和数据放在同一结构体中
```

- `sync.Mutex` 是值类型，零值可用
- 惯用做法：将锁作为结构体的第一个字段，或紧挨着它保护的字段

```go
stats.mu.Lock()
```

- 调用 `Lock()` 方法获取锁
- 如果锁已被其他协程持有，当前协程会在这里 **阻塞等待**
- 直到锁被释放，当前协程才能继续执行

```go
stats.mu.Unlock()
```

- 释放锁，让其他等待的协程可以获取
- **必须** 与 `Lock()` 配对使用

### 运行结果

```
期望订单数: 1000
实际统计数: 1000  // 完全正确！
丢失订单数: 0
```

---

## 第六步：defer 解锁模式（推荐）

### 为什么需要 defer

手动 `Unlock()` 可能遗漏：

```go
func processOrder(order Order, stats *SafeOrderStats) error {
    stats.mu.Lock()
  
    if order.Amount <= 0 {
        return errors.New("invalid amount")  // 忘记解锁！
    }
  
    stats.TotalOrders++
    stats.mu.Unlock()
    return nil
}
```

函数提前返回时，锁未释放，后续所有尝试获取锁的协程都会永远阻塞。

### defer 解决方案

```go
func processOrder(order Order, stats *SafeOrderStats) error {
    stats.mu.Lock()
    defer stats.mu.Unlock()  // 函数返回时自动解锁
  
    if order.Amount <= 0 {
        return errors.New("invalid amount")  // 会自动解锁
    }
  
    stats.TotalOrders++
    return nil  // 会自动解锁
}
```

`defer` 的特点：

- 函数返回时执行，无论是正常 return、提前 return 还是 panic
- 确保锁一定会被释放
- 代码更简洁，加锁解锁成对出现

### 最佳实践模式

```go
stats.mu.Lock()
defer stats.mu.Unlock()

// 临界区代码
// ...
```

将 `Lock()` 和 `defer Unlock()` 紧挨着写，形成固定模式。

---

## 第七步：读写锁 RWMutex

### 场景延续：订单统计需要提供查询接口

回到我们的订单统计系统。现在有新需求：

- **写操作**：每处理完一个订单，更新统计（相对较少）
- **读操作**：前端页面要展示统计数据，多个用户可能同时查看（非常频繁）

### 问题：普通 Mutex 让读操作也互斥

```go
// 使用普通 Mutex 的查询
func (s *StatsWithMutex) Query() (int, float64) {
    s.mu.Lock()
    defer s.mu.Unlock()
  
    // 模拟查询耗时 100ms
    time.Sleep(100 * time.Millisecond)
    return s.TotalOrders, s.TotalAmount
}
```

假设 3 个用户同时查询：

```
用户1 查询 ──────────▶ 100ms
                      用户2 等待 ──────────▶ 100ms
                                            用户3 等待 ──────────▶ 100ms
总耗时: 300ms
```

用户2 和用户3 明明只是查数据，不会修改任何东西，为什么要排队等待？

### RWMutex 的解决方案

```go
// 使用 RWMutex 的查询
func (s *StatsWithRWMutex) Query() (int, float64) {
    s.mu.RLock()        // 读锁，多个协程可同时持有
    defer s.mu.RUnlock()
  
    time.Sleep(100 * time.Millisecond)
    return s.TotalOrders, s.TotalAmount
}
```

3 个用户同时查询：

```
用户1 查询 ──────────▶
用户2 查询 ──────────▶  同时进行！
用户3 查询 ──────────▶
总耗时: 100ms
```

### RWMutex 的规则

| 锁类型         | 方法                        | 特点               |
| -------------- | --------------------------- | ------------------ |
| 读锁（共享锁） | `RLock()` / `RUnlock()` | 多个协程可同时持有 |
| 写锁（独占锁） | `Lock()` / `Unlock()`   | 只能一个协程持有   |

互斥关系：

```
        读锁      写锁
读锁    ✓ 共存    ✗ 互斥
写锁    ✗ 互斥    ✗ 互斥
```

- 读 + 读 = 可以同时
- 读 + 写 = 必须等待
- 写 + 写 = 必须等待

### 完整示例

```go
type StatsWithRWMutex struct {
    mu          sync.RWMutex
    TotalOrders int
    TotalAmount float64
}

// Query 查询统计（读操作）
func (s *StatsWithRWMutex) Query() (int, float64) {
    s.mu.RLock()         // R = Read，获取读锁
    defer s.mu.RUnlock()
  
    return s.TotalOrders, s.TotalAmount
}

// Update 更新统计（写操作）
func (s *StatsWithRWMutex) Update(amount float64) {
    s.mu.Lock()          // 获取写锁（独占）
    defer s.mu.Unlock()
  
    s.TotalOrders++
    s.TotalAmount += amount
}
```

### 语法拆解

```go
s.mu.RLock()   // R = Read，获取读锁
```

- 如果当前没有写锁，立即获取成功
- 如果有写锁被持有，阻塞等待

```go
s.mu.RUnlock()  // 释放读锁
```

```go
s.mu.Lock()     // 获取写锁（与普通 Mutex 的 Lock 用法相同）
```

- 必须等待所有读锁和写锁都释放后才能获取

### 何时使用 RWMutex

| 场景               | 推荐                            |
| ------------------ | ------------------------------- |
| 读操作远多于写操作 | `RWMutex`                     |
| 读写频率相近       | `Mutex`（RWMutex 有额外开销） |
| 只有写操作         | `Mutex`                       |

---

## 第八步：常见陷阱

### 陷阱 1：重复加锁（死锁）

```go
var mu sync.Mutex

mu.Lock()
mu.Lock()  // 死锁！同一协程再次加锁会永远阻塞
mu.Unlock()
```

Go 的 Mutex 不是可重入锁，同一协程不能对同一锁加锁两次。

### 陷阱 2：忘记解锁

```go
func doSomething() {
    mu.Lock()
    if someCondition {
        return  // 忘记解锁！
    }
    mu.Unlock()
}
```

解决：使用 `defer mu.Unlock()`

### 陷阱 3：锁拷贝

```go
mu := sync.Mutex{}
mu2 := mu       // 错误！拷贝了锁

mu.Lock()
mu2.Unlock()    // 解锁的是另一个锁！
```

`Mutex` 是值类型，拷贝后是两个独立的锁。传递时要用指针。

### 陷阱 4：交叉锁定（死锁）

场景：系统中有两把锁，某个操作需要同时持有两把锁。

```go
var lockA, lockB sync.Mutex

// 协程 A：先锁 A，再锁 B
go func() {
    lockA.Lock()
    time.Sleep(1 * time.Millisecond)
    lockB.Lock()    // 想获取 lockB，但被协程B持有，等待...
    // ...
    lockB.Unlock()
    lockA.Unlock()
}()

// 协程 B：先锁 B，再锁 A（顺序相反！）
go func() {
    lockB.Lock()
    time.Sleep(1 * time.Millisecond)
    lockA.Lock()    // 想获取 lockA，但被协程A持有，等待...
    // ...
    lockA.Unlock()
    lockB.Unlock()
}()
```

死锁过程：

1. 协程A 获取了 lockA，然后想获取 lockB
2. 协程B 获取了 lockB，然后想获取 lockA
3. 协程A 等 lockB 释放，但 lockB 被协程B 持有
4. 协程B 等 lockA 释放，但 lockA 被协程A 持有
5. 双方互相等待，谁也无法继续，程序卡死

**解决**：所有协程都按相同顺序获取锁（比如都先锁A再锁B）。

### 陷阱 5：在临界区内执行耗时操作

```go
mu.Lock()
result := callExternalAPI()  // 网络请求可能很慢
processResult(result)
mu.Unlock()
```

锁持有时间过长会阻塞其他协程。

**解决**：只在必要时持有锁。

```go
result := callExternalAPI()  // 锁外执行耗时操作

mu.Lock()
processResult(result)        // 只锁住必要的部分
mu.Unlock()
```

---

## 第九步：最佳实践

### 1. 将锁封装在结构体内

核心思想：**调用者不需要知道锁的存在，只管调用方法**。

```go
// 不好的写法：调用者需要自己加锁
mu.Lock()
counter++
mu.Unlock()

// 好的写法：锁封装在方法里
counter.Add(1)
```

完整示例：

```go
// Counter 线程安全的计数器
type Counter struct {
    mu    sync.Mutex  // 锁
    value int         // 被保护的数据
}

// Add 增加计数（调用者不需要关心锁）
func (c *Counter) Add(n int) {
    c.mu.Lock()
    defer c.mu.Unlock()
    c.value += n
}

// Value 获取当前值
func (c *Counter) Value() int {
    c.mu.Lock()
    defer c.mu.Unlock()
    return c.value
}
```

使用时直接调用方法：

```go
counter := &Counter{}
counter.Add(1)       // 不需要关心锁
fmt.Println(counter.Value())
```

### 2. 始终使用 defer 解锁

```go
s.mu.Lock()
defer s.mu.Unlock()
// 临界区代码
```

### 3. 最小化临界区

只锁住必须保护的代码：

```go
// 不好：锁住了不必要的代码
mu.Lock()
result := expensiveComputation()  // 不涉及共享数据
sharedData += result
mu.Unlock()

// 好：只锁住必要部分
result := expensiveComputation()
mu.Lock()
sharedData += result
mu.Unlock()
```

### 4. 避免在持有锁时调用外部函数

```go
// 不好：持有锁时调用可能阻塞的函数
mu.Lock()
callExternalService()  // 如果阻塞，锁一直被持有
mu.Unlock()

// 好：先获取数据，解锁后再调用外部服务
mu.Lock()
data := sharedData
mu.Unlock()
callExternalService(data)
```

### 5. 使用 go vet 检测数据竞争

```bash
go run -race main.go
```

`-race` 标志会检测数据竞争，开发时建议开启。

---

## 总结

| 概念       | 语法                              | 说明                         |
| ---------- | --------------------------------- | ---------------------------- |
| 互斥锁类型 | `sync.Mutex`                    | 保证同一时刻只有一个协程访问 |
| 加锁       | `mu.Lock()`                     | 获取锁，阻塞等待             |
| 解锁       | `mu.Unlock()`                   | 释放锁                       |
| 读写锁类型 | `sync.RWMutex`                  | 允许多读单写                 |
| 读锁       | `mu.RLock()` / `mu.RUnlock()` | 多个协程可同时持有           |
| 写锁       | `mu.Lock()` / `mu.Unlock()`   | 独占访问                     |
| 竞态检测   | `go run -race`                  | 开发时检测数据竞争           |

### 核心原则

1. **有共享数据被多个协程修改时，必须使用锁**
2. **锁的粒度要小**：只保护必要的代码
3. **使用 defer 解锁**：防止遗漏
4. **读多写少用 RWMutex**：提高并发性能
5. **避免死锁**：不重复加锁、按顺序加锁、不在临界区阻塞
