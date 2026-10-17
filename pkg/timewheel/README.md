# timewheel

采用 Orbit 的泛型、单 owner 主动推进实现。任务以 key 管理，支持取消、重设和周期执行；不创建后台 goroutine 或 timer。

## 示例

```go
package main

import (
    "fmt"
    "time"
    "github.com/NumberMan1/numbox/pkg/timewheel"
)

func main() {
    start := time.Unix(0, 0)
    wheel, err := timewheel.NewTimeWheel[string](start, timewheel.WheelOptions{
        Tick: 10 * time.Millisecond,
        Slots: 256,
    })
    if err != nil { panic(err) }
    if err := wheel.Schedule("job", start.Add(time.Second), 0, func(now time.Time) {
        fmt.Println(now.Sub(start)) // 1s
    }); err != nil { panic(err) }
    caught, err := wheel.Advance(start.Add(time.Second), 128)
    if err != nil { panic(err) }
    fmt.Println(caught, wheel.Len()) // true 0
}
```

实际服务在同一执行队列或更新循环内调用 Advance，传入标准库 time.Now 或单调模拟时钟；caught=false 时后续继续追赶。不要使用可回拨的业务日历时钟。

## 配置与接口

| 接口 | 语义 |
| --- | --- |
| NewTimeWheel(start, options) | Tick 必须为正，Slots 至少为 2，Logger 默认为 slog.Default |
| Schedule(key, deadline, interval, callback) | 绝对截止时间；interval=0 单次，正数为固定延迟周期；重复 key 返回 ErrTimerExists |
| Cancel(key) | 移除尚未执行的任务，返回是否仍待执行 |
| Reset(key, deadline) | 保留回调与周期，更新截止时间；不存在则返回 false |
| Clear() | 清除待执行任务，包括当前批次中还未调用的任务 |
| Advance(now, maxTicks) | 按预算处理槽位，返回是否追到 now；maxTicks 必须为正 |
| Len() | 当前待执行任务数量 |

截止时间向上取整到 tick，不提前触发；零/过去的截止时间在下一 tick 执行。同一批次按插入顺序执行符合条件的任务，前面的回调可取消或重设后面的任务。周期按本次 Advance 的 now 计算下一次时间，合并错过次数。

## 执行边界

- 所有方法必须来自同一 owner；跨协程先投递到同一执行队列。回调串行，慢回调会阻塞推进。
- maxTicks 限制处理槽数，不限制同槽任务数或回调耗时。空轮可直接跳到目标 tick；有任务时按预算追赶。
- 回调 panic 通过配置的 slog Logger 记录并隔离；递归 Advance 和时钟回拨返回错误。
- key 必须可比较且等于自身；any 持有切片/映射或 NaN key 会被拒绝，nil 接口 key 可以使用。
- key 校验使用泛型相等比较，局部捕获接口内部不可比较值引起的 panic；不使用 reflect，仍检查嵌套接口及 NaN。
- 单层槽位加圈数，取消通常为 O(1) 定位移除；远期任务每圈仍需检查，需按任务分布设置 tick/slots。
- Clear 不会停止外部队列或已开始的回调，owner 的启动和退出由应用负责。旧版 Start/Stop/Task 接口已移除。

## 基准

```sh
go test ./pkg/timewheel -run '^$' -bench 'Benchmark(Wheel|Runtime)ScheduleCancel' -benchmem
go test ./pkg/timewheel -run '^$' -bench BenchmarkKeyValidation -benchmem -count=3
```

该基准只比较 10,000 个远期任务的取消和重建，不包含到期执行、并发或整体调度延迟，不能代表全部业务负载。

## 许可证

本库采用 [MIT License](../../LICENSE)。
