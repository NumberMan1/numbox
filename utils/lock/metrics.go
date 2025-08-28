package lock

var metrics Metrics

func GetMetrics() Metrics {
	return metrics
}

type Metrics struct {
	LTimeOutTimes  int64 // 锁超时次数
	RWTimeOutTimes int64 // 读写锁超时次数
}
