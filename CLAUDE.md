# CLAUDE.md

此文件为 Claude Code (claude.ai/code) 在此仓库中工作时提供指导。

## 构建与测试

```bash
# 运行所有测试
go test ./...

# 运行指定包的测试
go test ./pkg/lock/...

# 运行单个测试用例
go test -run TestLocker ./pkg/lock/...

# 运行性能基准测试
go test -bench=. ./pkg/lock/...

# 构建 (pkg/go-zset 依赖 CGO)
CGO_ENABLED=1 go build ./...
```

## 架构概览

这是一个 Go 工具库（模块 `github.com/NumberMan1/numbox`），提供可复用的数据结构和辅助函数。各包之间相互独立，没有统一的框架将它们耦合在一起，但存在以下内部依赖关系：

- `pkg/collection` → `pkg/utils`（使用 `MapKeys`）
- `pkg/sensitive-word` → `pkg/collection`（Set）、`pkg/utils`（JSON 反序列化）
- `pkg/encryp` → `pkg/env`（获取 JWT 默认密钥）

### 关键设计模式

**函数选项模式（Functional Options）** 在 `pkg/lock` 中用于 `Lock()`/`Unlock()` 的配置（`LockOption`、`UnlockOption`）。为其他包添加可配置行为时，请遵循此模式。

**锁包（`pkg/lock`）** 提供两种实现：
- `Locker` — 基于 sync.Mutex 的简单锁，具备基于 token 的所有权、获取超时和持锁过期机制。获取锁失败时会 panic `ErrLockTimeout`，调用方需要 `recover()` 来捕获。
- `RWLocker` — 读写锁，通过 `writeIntention` 计数器实现写者优先。多个读者可以同时持锁；当有写意向时，新的读锁请求会被阻塞。

锁等待均使用指数退避轮询（`withTimeout`），而非 channel 机制。

**权重池（`pkg/weight`）** 基于累计权重的二分查找实现加权随机选取（O(log N)）。`Item` 接口要求实现 `Weight() int` 方法。泛型便捷函数（`PickOneFromItems[T Item]`、`PickManyFromFairPool[T]` 等）封装了常见使用场景。`PickRandom()` 为不放回抽样——若需放回抽样，使用 `PickRandomAndPutBack()`。

**时间轮（`pkg/timewheel`）** 是基于 channel 的定时任务调度器。所有操作（添加/删除/滴答推进）都通过单个 goroutine 的 `select` 循环串行化，无需互斥锁即可保证 goroutine 安全。到期任务在独立的 goroutine 中执行，以防止阻塞时间轮的推进。

**go-zset（`pkg/go-zset`）** 通过 cgo 封装了 C 实现的跳表——需要 `CGO_ENABLED=1`。如果需要纯 Go 方案，可使用 `pkg/skiplist`（基于 `big.Float` 键，适合任意精度的场景）。

### 错误处理模式

- `errors/` — 自定义 `CodeError` 类型，携带错误码（`Code int32`）以及可选的内部错误和堆栈信息。使用 `errors.ToError(code, data)` 将任意值转换为 CodeError。
- `pkg/stackerr` — `StackError` 通过 `pkg/errors` 为错误包裹堆栈跟踪信息。`FormatStackError()` 是捕获 panic 并附加上下文的主要入口。
- `pkg/utils` — `Must(err)` 和 `Asset(isOk, err)` 是遇错即 panic 的辅助函数。`WithRecover()` 用 recover 块包裹函数执行。

### Go 版本及依赖

- **Go 1.25.1**，代码中广泛使用了泛型
- 测试框架：`github.com/stretchr/testify`（assert 断言）
- 主要依赖：`github.com/google/uuid`、`github.com/dgrijalva/jwt-go`、`golang.org/x/exp`、`github.com/pkg/errors`
