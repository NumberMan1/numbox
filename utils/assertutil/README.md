# assertutil

Must、Assert、AssertStringError 在条件失败时 panic。
WithRecover 捕获 panic，并把 panic 或函数返回错误打印到标准输出，不把错误返回调用方。
需要结构化日志或异步 panic 处理时使用 pkg/safego。
