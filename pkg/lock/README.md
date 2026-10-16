# lock

带 token 的进程内互斥锁和读写锁。**锁不会自动过期，任务完成后才释放**，不提供租约或分布式锁能力。
只需要基础互斥时优先使用标准库 `sync.Mutex` / `sync.RWMutex`；本包额外提供 token 校验、可取消的获取过程及托管任务接口。

## 手动管理

```go
l := lock.NewLocker()
token, err := l.LockContext(ctx, lock.WithAcquireTimeout(time.Second))
if err != nil {
    return err
}
defer l.Unlock(token)
// 在此同步访问共享数据；ctx 取消不会自动释放 token。
```

`Lock` / `RLock` 保留等待失败时 panic 的行为。新代码可以使用返回 `(Token, error)` 的
`LockContext` / `RLockContext`：

- `WithAcquireTimeout` 仅控制获取锁的等待时间，默认 1 秒，配置到期返回 `ErrLockTimeout`。
- 传入 context 取消或到期返回 `context.Canceled` / `context.DeadlineExceeded`。
- 获取期限非正值表示不设独立等待上限，仍响应 context。
- token 在同一锁实例内生成；错误、过时 token 或重复释放返回 false。零 token 无效。
- `Acquired` 仅表示查询当时是否持有锁，不能用作替代加锁的并发保护。
- 同一个锁不可重入，也不支持读锁升级为写锁。读写锁优先处理等待中的写者，不承诺 FIFO。

## 托管任务与协作取消

```go
package main

import (
    "context"
    "fmt"
    "time"

    "github.com/NumberMan1/numbox/pkg/lock"
)

func main() {
    l := lock.NewLocker()
    err := l.Do(context.Background(), func(taskCtx context.Context) error {
        // 下游 I/O 使用 taskCtx；计算循环应主动检查 taskCtx.Err()。
        <-taskCtx.Done()
        return taskCtx.Err()
    }, lock.WithAcquireTimeout(time.Second),
       lock.WithTaskTimeout(20*time.Millisecond))
    fmt.Println(err) // context deadline exceeded
}
```

`Do` 在写锁下运行任务；`RWLocker.DoRead` 在读锁下运行任务。
获取锁后才为任务创建执行期限，获取阶段的独立等待期限不会泄漏到任务 context。
调用方 context 的截止时间仍作用于等待和任务两个阶段。

任务在调用 goroutine 中**同步执行**，不通过后台 goroutine 模拟超时：

1. 获取锁失败，任务不执行。
2. 任务执行超时，context 发出取消信号，任务须主动响应。
3. 任务返回后才释放锁；任务忽略取消或尚在清理时，`Do` 会继续等待。
4. 任务返回错误时优先保留该错误；返回 nil 时返回任务 context 的错误。
5. 任务 panic 时先释放锁，再继续传播 panic；nil 任务返回 `ErrNilTask`。

`WithTaskTimeout` 非正值表示不设独立任务期限，对手动获取接口不生效。
Go 的取消信号不会强制终止 goroutine，因而不能承诺忽略取消的任务按时返回。
任务若自行启动访问同一临界区的子 goroutine，必须在返回前等待它们全部结束。

## 回调与指标

- 加锁回调在成功取得 token 后同步执行，不持有内部状态互斥锁，但此时业务锁仍被占用。
  回调不得等待重新获取同一把业务锁；回调 panic 时会先释放尚未交付的 token。
- 解锁回调仅在成功释放后执行，可以查询状态或重新获取锁。此时其他任务也可能已取得锁。
- nil 回调被忽略；解锁回调 panic 仍向上传播，此时 token 已释放。
- `GetMetrics` 并发安全地读取普通锁、读写锁的获取超时次数；
  调用方取消和任务执行超时不计入获取超时。两个计数是分别读取的快照。
- `DefaultAcquireTimeout` 等包级配置应在并发工作启动前设置。

## 从自动过期版本迁移

这次调整改变了旧的过期行为：

| 旧接口或行为 | 新行为与迁移 |
| --- | --- |
| 持有期限到期，其他任务可取得锁 | 所有读写 token 都不会自动过期，必须由持有者释放 |
| `WithLockHoldTimeout(d)` | 弃用，作为 `WithTaskTimeout(d)` 的兼容别名；仅在 `Do` / `DoRead` 中生效 |
| `DefaultHoldLockExpired` | 仅保留符号兼容，值不再被读取；任务超时需显式传入选项 |
| 手动获取后依赖超时释放 | 使用 `defer Unlock(token)` / `defer RUnlock(token)`，或改用托管任务接口 |
| 自行实现 `Locker` / `RWLocker` | 接口新增 context 获取和任务执行方法，自定义实现须补齐对应方法 |

旧业务若依赖自动过期来恢复死锁，须先修复任务退出和释放路径。
不要在超时后提前返回并解锁，让旧任务继续在后台访问共享数据。

## 验证

```sh
go test -race -count=1 ./pkg/lock
go vet ./pkg/lock
```
