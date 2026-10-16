package lock

import "sync/atomic"

var metrics Metrics

// GetMetrics 返回超时计数的并发安全快照；两个计数分别读取，不是事务快照。
func GetMetrics() Metrics {
	return Metrics{
		LTimeOutTimes:  atomic.LoadInt64(&metrics.LTimeOutTimes),
		RWTimeOutTimes: atomic.LoadInt64(&metrics.RWTimeOutTimes),
	}
}

// Metrics 统计配置的获取锁期限导致的失败，不包括调用方取消和任务执行超时。
type Metrics struct {
	// LTimeOutTimes 是普通锁获取超时次数。
	LTimeOutTimes int64
	// RWTimeOutTimes 是读写锁获取读或写 token 超时次数。
	RWTimeOutTimes int64
}
