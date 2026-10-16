# AGENTS.md

此文件为在本仓库中工作的 AI 编程助手提供通用指导，适用于整个仓库。

## 分支与操作授权

- 开发修改位于 `develop` 分支。
- 完成修改与验证后先汇报，只有获得用户对本次操作的明确授权才执行 commit、push 或 PR 合并。
- 对 `master` 的提交、推送或合并需要用户明确授权；历史任务的授权不自动延续到新的修改。

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

- `pkg/sensitive-word` → `pkg/matcher`（词表匹配）、`pkg/utils`（JSON 反序列化）
- `pkg/encryp` → `pkg/env`（获取 JWT 默认密钥）

### 关键设计模式

**函数选项模式（Functional Options）** 在 `pkg/lock` 中用于 `Lock()`/`Unlock()` 的配置（`LockOption`、`UnlockOption`）。为其他包添加可配置行为时，请遵循此模式。

**锁包（`pkg/lock`）** 提供两种实现：
- `Locker` — 带 token 的进程内互斥锁，持有权不会自动过期，必须显式释放。`Lock` 等待超时时 panic `ErrLockTimeout`；`LockContext` 通过 error 返回超时或取消。
- `RWLocker` — 读写锁，通过 `writeIntention` 计数器实现写者优先。多个读者可以同时持锁；当有写意向时，新的读锁请求会被阻塞。

锁获取阶段使用可取消、有上限的退避轮询（`waitForLock`）。`Do`/`DoRead` 在锁内同步执行任务，任务超时只取消传入的 context；必须等待任务真正退出后才解锁。不得通过后台 goroutine 加超时返回来提前释放锁。任务函数还须等待自己启动的临界区子任务结束后才能返回。

`WithTaskTimeout` 只用于 `Do`/`DoRead`，与获取阶段的 `WithAcquireTimeout` 分离。旧 `WithLockHoldTimeout` 是任务超时选项的弃用别名，`DefaultHoldLockExpired` 不再生效。迁移规则见 `pkg/lock/README.md`。

**权重池（`pkg/weight`）** 基于累计权重的二分查找实现加权随机选取（O(log N)）。`Item` 接口要求实现 `Weight() int` 方法。泛型便捷函数（`PickOneFromItems[T Item]`、`PickManyFromFairPool[T]` 等）封装了常见使用场景。`PickRandom()` 为不放回抽样——若需放回抽样，使用 `PickRandomAndPutBack()`。

**时间轮（`pkg/timewheel`）** 采用 Orbit 的泛型主动推进实现。由同一个 owner 调用 `Schedule/Cancel/Reset/Clear/Advance`，不提供隐式锁、后台 goroutine 或旧版 `Start/Stop/Task` 接口。回调串行执行，截止时间向上取整，panic 默认通过 slog 记录。外部协程应投递到同一执行队列。详见 `pkg/timewheel/README.md`。

**go-zset（`pkg/go-zset`）** 通过 cgo 封装了 C 实现的跳表——需要 `CGO_ENABLED=1`。如果需要纯 Go 方案，可使用 `pkg/skiplist`（基于 `big.Float` 键，适合任意精度的场景）。

### 错误处理模式

- `errors/` — 自定义 `CodeError` 类型，携带错误码（`Code int32`）以及可选的内部错误和堆栈信息。使用 `errors.ToError(code, data)` 将任意值转换为 CodeError。
- `pkg/stackerr` — `StackError` 通过 `pkg/errors` 为错误包裹堆栈跟踪信息。`FormatStackError()` 是捕获 panic 并附加上下文的主要入口。
- `pkg/utils` — `Must(err)` 和 `Assert(isOk, err)` 是遇错即 panic 的辅助函数。`WithRecover()` 用 recover 块包裹函数执行。

### Go 版本及依赖

- **Go 1.25.1**，代码中广泛使用了泛型
- 测试框架：`github.com/stretchr/testify`（assert 断言）
- 主要依赖：`github.com/google/uuid`、`github.com/dgrijalva/jwt-go`、`golang.org/x/exp`、`github.com/pkg/errors`
