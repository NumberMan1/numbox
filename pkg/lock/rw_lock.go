package lock

import (
	"context"
	"sync"
	"sync/atomic"
)

type rwLocker struct {
	mutex          sync.Mutex
	write          bool
	writeIntention int32
	token          Token
	readTokens     map[Token]struct{}
}

func (l *rwLocker) lock() Token {
	l.mutex.Lock()
	defer l.mutex.Unlock()
	if l.write || len(l.readTokens) > 0 {
		return 0
	}
	l.write = true
	l.token++
	return l.token
}

func (l *rwLocker) acquire(ctx context.Context, conf *lockConfig) (Token, error) {
	token, err := func() (Token, error) {
		atomic.AddInt32(&l.writeIntention, 1)
		defer atomic.AddInt32(&l.writeIntention, -1)
		return waitForLock(ctx, conf.acquireTimeout, l.lock)
	}()
	return l.finish(token, err, conf, l.Unlock)
}

func (l *rwLocker) finish(token Token, err error, conf *lockConfig, unlock func(Token, ...UnlockOption) bool) (Token, error) {
	if err != nil {
		if err == ErrLockTimeout {
			atomic.AddInt64(&metrics.RWTimeOutTimes, 1)
		}
		return 0, err
	}
	return finishAcquire(token, conf.cb, unlock), nil
}

// Lock 获取写锁；等待失败时 panic，成功后必须显式 Unlock。
func (l *rwLocker) Lock(opt ...LockOption) Token {
	token, err := l.LockContext(context.Background(), opt...)
	if err != nil {
		panic(err)
	}
	return token
}

// LockContext 可取消地等待写锁；取消的写者会撤回写意向。
func (l *rwLocker) LockContext(ctx context.Context, opt ...LockOption) (Token, error) {
	return l.acquire(ctx, newConfig(opt...))
}

// Unlock 释放有效的写 token；用户回调在内部互斥锁外执行。
func (l *rwLocker) Unlock(token Token, opt ...UnlockOption) bool {
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

func (l *rwLocker) rLock() Token {
	l.mutex.Lock()
	defer l.mutex.Unlock()
	if l.write || atomic.LoadInt32(&l.writeIntention) > 0 {
		return 0
	}
	l.token++
	l.readTokens[l.token] = struct{}{}
	return l.token
}

func (l *rwLocker) acquireRead(ctx context.Context, conf *lockConfig) (Token, error) {
	token, err := waitForLock(ctx, conf.acquireTimeout, l.rLock)
	return l.finish(token, err, conf, l.RUnlock)
}

// RLock 获取读锁，等待失败时 panic；不会因持有时间过长而失效。
func (l *rwLocker) RLock(opt ...LockOption) Token {
	token, err := l.RLockContext(context.Background(), opt...)
	if err != nil {
		panic(err)
	}
	return token
}

// RLockContext 可取消地等待读锁，存在写意向时等待写者完成或退出。
func (l *rwLocker) RLockContext(ctx context.Context, opt ...LockOption) (Token, error) {
	return l.acquireRead(ctx, newConfig(opt...))
}

// RUnlock 释放有效的读 token；用户回调在内部互斥锁外执行。
func (l *rwLocker) RUnlock(token Token, opt ...UnlockOption) bool {
	conf := newUnlockConfig(opt...)
	l.mutex.Lock()
	_, ok := l.readTokens[token]
	if ok {
		delete(l.readTokens, token)
	}
	l.mutex.Unlock()
	if ok {
		conf.cb.invoke()
	}
	return ok
}

// Acquired 并发安全地查询有效的读或写 token。
func (l *rwLocker) Acquired(tk Token) bool {
	l.mutex.Lock()
	defer l.mutex.Unlock()
	_, read := l.readTokens[tk]
	return l.write && l.token == tk || read
}

// Do 在写锁下同步执行任务，等待任务退出后解锁。
func (l *rwLocker) Do(ctx context.Context, fn func(context.Context) error, opt ...LockOption) error {
	return runTask(ctx, fn, newConfig(opt...), l.acquire, l.Unlock)
}

// DoRead 在读锁下同步执行任务，等待任务退出后释放读 token。
func (l *rwLocker) DoRead(ctx context.Context, fn func(context.Context) error, opt ...LockOption) error {
	return runTask(ctx, fn, newConfig(opt...), l.acquireRead, l.RUnlock)
}
