# safego

启动异步协程，在 panic 时 recover 并报告带堆栈的错误。默认使用标准库 slog，也可注入自己的上报函数。

## 示例

```go
package main

import (
    "fmt"
    "github.com/NumberMan1/numbox/pkg/safego"
)

func main() {
    reported := make(chan struct{})
    safego.GoWithReporter(func() {
        panic("background task failed")
    }, func(err error) {
        fmt.Printf("%+v\n", err)
        close(reported)
    })
    <-reported // 示例自行协调完成；这不是 safego 的等待接口
}
```

- `Go(work)` 异步启动，发生 panic 时通过 slog 记录错误及堆栈。
- `GoWithReporter(work, reporter)` 使用指定 reporter；nil reporter 使用默认上报。reporter 在工作协程内调用，自身不应 panic。
- work 正常返回时不调用 reporter；reporter 收到的错误保留原 error 的错误链。
- 本包只隔离 panic，不提供数据竞争保护、重试、取消、等待或关闭排空。需要生命周期管理的工作由调用方传入 context，并自行协调结束。
