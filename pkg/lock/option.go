package lock

import "time"

type callback []func()

func (c *callback) append(cb func()) {
	if cb != nil {
		*c = append(*c, cb)
	}
}

func (c *callback) invoke() {
	for _, cb := range *c {
		cb()
	}
}

func newConfig(opt ...LockOption) *lockConfig {
	conf := &lockConfig{acquireTimeout: DefaultAcquireTimeout, cb: new(callback)}
	for _, option := range opt {
		option(conf)
	}
	return conf
}

type lockConfig struct {
	acquireTimeout time.Duration
	taskTimeout    time.Duration
	cb             *callback
}

// LockOption 配置一次获取锁或执行任务的行为；同类选项以最后一次设置为准。
type LockOption func(*lockConfig)

// WithAcquireTimeout 设置等待获取锁的期限；非正值不设独立期限。
// 此期限不包含任务执行时间，任务阶段始终使用调用方的原始 context。
func WithAcquireTimeout(duration time.Duration) LockOption {
	return func(c *lockConfig) { c.acquireTimeout = duration }
}

// WithTaskTimeout 设置 Do/DoRead 在获取锁后创建的任务 context 期限。
// 非正值不设独立期限。任务必须主动响应取消，函数实际退出后才解锁。
// 对手动管理的 Lock、LockContext、RLock、RLockContext 不生效。
func WithTaskTimeout(duration time.Duration) LockOption {
	return func(c *lockConfig) { c.taskTimeout = duration }
}

// WithLockHoldTimeout 是 WithTaskTimeout 的兼容别名，不再使锁自动过期。
//
// Deprecated: 改用 WithTaskTimeout，并使用 Do/DoRead 管理任务生命周期。
func WithLockHoldTimeout(duration time.Duration) LockOption {
	return WithTaskTimeout(duration)
}

// WithLockCallback 注册成功获取锁后同步执行的回调，nil 被忽略。
// 回调在内部互斥锁外执行；若回调 panic，归还尚未交给调用方的 token 后继续 panic。
func WithLockCallback(cb func()) LockOption {
	return func(c *lockConfig) { c.cb.append(cb) }
}

func newUnlockConfig(opt ...UnlockOption) *unlockConfig {
	conf := &unlockConfig{cb: new(callback)}
	for _, option := range opt {
		option(conf)
	}
	return conf
}

type unlockConfig struct{ cb *callback }

// UnlockOption 配置一次释放锁的行为。
type UnlockOption func(config *unlockConfig)

// WithUnlockCallback 注册成功释放 token 后同步执行的回调，nil 被忽略。
// 回调执行时锁已释放，可查询状态或重新获取；回调 panic 向调用方传播。
func WithUnlockCallback(cb func()) UnlockOption {
	return func(c *unlockConfig) { c.cb.append(cb) }
}
