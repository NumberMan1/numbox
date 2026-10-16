// Package lock 提供不会自动过期的、带 token 的进程内互斥锁和读写锁。
// 手动获取的 token 必须显式释放；Do/DoRead 等待任务真正退出后自动释放。
// context 取消只通知任务停止，不会强制终止 goroutine 或提前放行临界区。
package lock

import "context"

// Locker 是带 token 的进程内互斥锁。持有权不会自动过期，必须显式解锁。
// 不支持重入；持有权属于 token，不绑定到某个 goroutine。
type Locker interface {
	// Acquired 并发安全地查询 token 当前是否仍持有锁。
	Acquired(tk Token) bool
	// Lock 获取锁；等待失败时 panic。获取超时默认为 DefaultAcquireTimeout。
	Lock(opt ...LockOption) Token
	// LockContext 获取锁，取消或超时以 error 返回；成功后由调用方负责 Unlock。
	// ctx 取消不会释放已经取得的锁，ctx 不得为 nil。
	LockContext(ctx context.Context, opt ...LockOption) (Token, error)
	// Unlock 释放匹配 token 的锁；错误 token 或重复释放返回 false。
	Unlock(token Token, opt ...UnlockOption) bool
	// Do 获取锁后同步执行 fn，并在 fn 返回或 panic 展开时解锁。
	// 任务必须主动响应 context；取消后仍会等待 fn 返回，不提前释放锁。
	// fn 必须等待访问同一临界区的子任务结束后才能返回。
	// fn 的错误优先返回；fn 返回 nil 时返回任务 context 的错误。
	Do(ctx context.Context, fn func(context.Context) error, opt ...LockOption) error
}

// RWLocker 允许多个读者或一个写者，并优先处理等待中的写者。
// 所有读写 token 都必须显式释放；不支持重入或读锁升级为写锁。
type RWLocker interface {
	Locker
	// RLock 获取读锁，等待失败时 panic。
	RLock(opt ...LockOption) Token
	// RLockContext 可取消地获取读锁；成功后由调用方负责 RUnlock。
	RLockContext(ctx context.Context, opt ...LockOption) (Token, error)
	// RUnlock 释放匹配的读 token；错误 token 或重复释放返回 false。
	RUnlock(token Token, opt ...UnlockOption) bool
	// DoRead 在读锁下同步执行 fn，取消、错误和 panic 语义与 Do 相同。
	DoRead(ctx context.Context, fn func(context.Context) error, opt ...LockOption) error
}

// NewLocker 创建不会自动过期的进程内互斥锁。
func NewLocker() Locker { return &locker{} }

// NewRWLocker 创建不会自动过期、优先处理写者的进程内读写锁。
func NewRWLocker() RWLocker {
	return &rwLocker{readTokens: make(map[Token]struct{})}
}
