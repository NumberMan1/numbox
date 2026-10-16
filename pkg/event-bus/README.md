# event-bus

导入路径为 `github.com/NumberMan1/numbox/pkg/event-bus`，Go 包名为 `eventbus`。本包只定义事件主题和合并契约，不提供订阅、发布、异步投递或请求队列。

## 示例

```go
package main

import (
    "fmt"
    eventbus "github.com/NumberMan1/numbox/pkg/event-bus"
)

func main() {
    topic := eventbus.Topic("score").AddSuffix("changed")
    fmt.Println(topic.String()) // score_changed
}
```

`MergeableEventData` 要求实现 `Merge(other MergeableEventData) MergeableEventData` 和 `MergeKey() string`。队列实现只有在主题与 MergeKey 均相同时才能合并。

MergeKey 应包含允许合并的业务身份，避免将不同用户的事件混合。Merge 的类型匹配、合并规则、投递顺序和并发同步由事件实现与消费队列负责。本包没有默认合并器。
