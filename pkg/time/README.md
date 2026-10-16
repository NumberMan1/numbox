# time

毫秒时间戳、日历边界、可调业务时钟，以及从 Orbit 提取的 After/Every/PollTimer。导入时不会修改进程时区。

## 时间戳与日历

```go
package main

import (
    "fmt"
    "time"
    ntime "github.com/NumberMan1/numbox/pkg/time"
)

func main() {
    now := ntime.NewTimestamp(time.Now())
    start, end := ntime.GetTodayStartEnd(now)
    fmt.Println(now.IsValid()) // true
    fmt.Println(start.Before(end)) // true
    duration, err := ntime.ParseDuration("1h30m")
    fmt.Println(duration, err) // 1h30m0s <nil>
}
```

Timestamp 以 int64 毫秒表示，IsValid 判断是否大于零。NewTimestampWithUnix 的输入是秒，NewTimestampWithUnixMilli 的输入是毫秒；Time/Unix/UnixNano 负责转换。

日/周边界使用 `[start,end)`，周一为一周起点；按应用的 time.Local 与日历计算，支持夏令时。需要固定上海时区的应用应在启动时明确设置 time.Local 或 TZ，工具包不会代为修改全局配置。

## 时钟选择

| 场景 | 使用方式 |
| --- | --- |
| 可调业务日历，如每日刷新或活动时间 | Now、NewTimestampNow、Since、SinceDays |
| 超时、耗时统计、日志、时间轮 | 标准库 time.Now/time.Since |
| 模拟和更新循环 | 显式向 PollTimer 或时间轮传入一致的模拟时钟 |

SetOffsetMillis 原子设置全局业务时钟偏移，OffsetMillis 读取，零表示真实时间；超出 time.Duration 可表示范围时 panic。调整只影响本包的业务当前时间，不调整操作系统时钟。

SinceDays 按 24 小时折算并加一表示含当天，不是自然日差；时钟回拨时 Since/SinceDays 可能为负。DurationToNextInterval 返回下一个固定时长边界的剩余时长；无效时间戳或非正 interval 返回零。ParseDuration 保留 `(Duration,error)`。

## After/Every/PollTimer

```go
package main

import (
    "context"
    "fmt"
    "time"
    ntime "github.com/NumberMan1/numbox/pkg/time"
)

func main() {
    if err := ntime.After(context.Background(), 0, func(context.Context) error {
        fmt.Println("once")
        return nil
    }); err != nil { panic(err) }

    start := time.Unix(0, 0)
    var timer ntime.PollTimer
    if err := timer.Reset(start, time.Second, 0); err != nil { panic(err) }
    fmt.Println(timer.Poll(start)) // false
    fmt.Println(timer.Poll(start.Add(time.Second))) // true
    fmt.Println(timer.Active()) // false，单次触发后失活
}
```

After/Every 与 Orbit 实现一致，是阻塞函数，使用标准库 timer/ticker，不受业务偏移时钟影响。context 取消时返回 context 错误，回调返回错误时结束；Every 回调串行，慢回调合并 tick，不重叠执行。调用方决定是否启动 goroutine，并负责 panic 策略。

PollTimer 同样采用 Orbit 实现：零值失活，Reset 的 interval=0 为单次，正值为周期；Pause/Resume 保存剩余时长，Cancel 清空。迟到的周期 Poll 合并错过次数，从传入 now 计算下一次截止时间。该计时器由单个 owner 调用，不提供并发同步，调用方需保持一致的时钟域。
