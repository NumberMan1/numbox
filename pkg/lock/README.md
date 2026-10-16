# lock

带 token、获取超时和持有期限的进程内锁，提供 Locker 与优先处理写意向的 RWLocker。它不是分布式锁。

## 示例

```go
package main

import (
    "fmt"
    "time"
    "github.com/NumberMan1/numbox/pkg/lock"
)

func main() {
    l := lock.NewLocker()
    token := l.Lock(
        lock.WithAcquireTimeout(time.Second),
        lock.WithLockHoldTimeout(5*time.Second),
    )
    defer l.Unlock(token)
    fmt.Println(l.Acquired(token)) // true
}
```

## 接口

- Lock 返回持有者 token，获取超时时 panic ErrLockTimeout，调用方负责自己的 recover 策略。
- Unlock 需要匹配的 token，返回是否成功；重复解锁或 token 不匹配返回 false。
- Acquired 同步检查 token 与到期时间。
- NewRWLocker 额外提供 RLock/RUnlock。存在写意向时阻止新读者进入，多个有效读锁可共存。
- WithLockCallback/WithUnlockCallback 注册同步回调；解锁回调在内部互斥锁释放后执行，可查询或重新获取锁。回调 panic 由调用方处理。

## 持有期限

到期后其他调用可以获取新 token，但原临界区不会被强制停止。因此有期限的锁不能替代 sync.Mutex/sync.RWMutex 对普通共享内存提供的互斥保证。需要这种保证时直接使用标准库锁。

获取阶段使用退避轮询，默认配置应在并发工作启动前设置。获取超时与持有期限是不同的参数，不要把它们混为一次操作的执行截止时间。
