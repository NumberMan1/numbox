package lock

import (
	"context"
	"sync"
	"sync/atomic"
)

type locker struct {
	mutex sync.Mutex
	write bool
	token Token
}

func (l *locker) lock() Token {
	l.mutex.Lock()
	defer l.mutex.Unlock()
	if l.write {
		return 0
	}
	l.write = true
	l.token++
	return l.token
}

func (l *locker) acquire(ctx context.Context, conf *lockConfig) (Token, error) {
	token, err := waitForLock(ctx, conf.acquireTimeout, l.lock)
	if err != nil {
		if err == ErrLockTimeout {
			atomic.AddInt64(&metrics.LTimeOutTimes, 1)
		}
		return 0, err
	}
	return finishAcquire(token, conf.cb, l.Unlock), nil
}

// Lock 获取锁；等待失败时 panic，成功后必须显式 Unlock。
func (l *locker) Lock(opt ...LockOption) Token {
	token, err := l.LockContext(context.Background(), opt...)
	if err != nil {
		panic(err)
	}
	return token
}

// LockContext 获取锁，等待取消或超时返回错误，不改变已经取得的持有权。
func (l *locker) LockContext(ctx context.Context, opt ...LockOption) (Token, error) {
	return l.acquire(ctx, newConfig(opt...))
}

// Unlock 释放有效 token；回调在内部互斥锁释放后同步执行。
func (l *locker) Unlock(token Token, opt ...UnlockOption) bool {
	conf := newUnlockConfig(opt...)
	l.mutex.Lock()
	if l.token != token || !l.write {
		l.mutex.Unlock()
		return false
	}
	l.write = false
	l.mutex.Unlock()
	conf.cb.invoke()
	return true
}

// Acquired 并发安全地查询 token 是否仍持有锁，持有权不会自动过期。
func (l *locker) Acquired(tk Token) bool {
	l.mutex.Lock()
	defer l.mutex.Unlock()
	return l.write && l.token == tk
}

// Do 同步执行任务，等待任务实际退出后解锁；panic 也会先释放锁再向上传播。
func (l *locker) Do(ctx context.Context, fn func(context.Context) error, opt ...LockOption) error {
	return runTask(ctx, fn, newConfig(opt...), l.acquire, l.Unlock)
}
