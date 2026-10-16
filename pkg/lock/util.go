package lock

import (
	"context"
	"errors"
	"time"
)

// Token 是同一个锁实例内的持有权编号；零值无效，不是分布式 fencing token。
type Token int

var (
	// ErrLockTimeout 表示配置的获取锁等待期限已到；context 取消返回对应的 context 错误。
	ErrLockTimeout = errors.New("lock_timeout")
	// ErrNilTask 表示 Do 或 DoRead 收到了 nil 任务函数。
	ErrNilTask = errors.New("lock: nil task")
	// DefaultAcquireTimeout 是默认获取锁等待期限，应在并发工作启动前配置。
	// 非正值表示不设独立等待期限；LockContext 仍受 ctx 取消控制。
	DefaultAcquireTimeout = time.Second
	// DefaultHoldLockExpired 仅保留以兼容旧代码，其值不再被读取。
	//
	// Deprecated: 锁不再自动过期。任务执行期限使用 WithTaskTimeout 和 Do/DoRead。
	DefaultHoldLockExpired = time.Second * 10
)

func waitForLock(ctx context.Context, timeout time.Duration, try func() Token) (Token, error) {
	waitCtx := ctx
	if timeout > 0 {
		var cancel context.CancelFunc
		waitCtx, cancel = context.WithTimeout(ctx, timeout)
		defer cancel()
	}
	delay := time.Millisecond
	for {
		if err := waitCtx.Err(); err != nil {
			if err := ctx.Err(); err != nil {
				return 0, err
			}
			return 0, ErrLockTimeout
		}
		if token := try(); token != 0 {
			return token, nil
		}
		timer := time.NewTimer(delay)
		select {
		case <-waitCtx.Done():
			timer.Stop()
		case <-timer.C:
		}
		if delay < 10*time.Millisecond {
			delay *= 2
			if delay > 10*time.Millisecond {
				delay = 10 * time.Millisecond
			}
		}
	}
}

// If a lock callback panics or calls Goexit, the caller never receives its token.
func finishAcquire(token Token, cb *callback, unlock func(Token, ...UnlockOption) bool) Token {
	delivered := false
	defer func() {
		if !delivered {
			unlock(token)
		}
	}()
	cb.invoke()
	delivered = true
	return token
}

func runTask(
	ctx context.Context, fn func(context.Context) error, conf *lockConfig,
	acquire func(context.Context, *lockConfig) (Token, error),
	unlock func(Token, ...UnlockOption) bool,
) error {
	if fn == nil {
		return ErrNilTask
	}
	token, err := acquire(ctx, conf)
	if err != nil {
		return err
	}
	defer unlock(token)
	taskCtx := ctx
	if conf.taskTimeout > 0 {
		var cancel context.CancelFunc
		taskCtx, cancel = context.WithTimeout(ctx, conf.taskTimeout)
		defer cancel()
	}
	if err := taskCtx.Err(); err != nil {
		return err
	}
	// Run synchronously: cancellation must never let a second task enter early.
	if err := fn(taskCtx); err != nil {
		return err
	}
	return taskCtx.Err()
}
