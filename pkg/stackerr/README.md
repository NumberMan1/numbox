# stackerr

为错误添加或提取堆栈，同时保留标准库 `errors.Is/errors.As` 的错误链。适用于返回错误和 recover 边界。

## 示例

```go
package main

import (
    "errors"
    "fmt"
    "github.com/NumberMan1/numbox/pkg/stackerr"
)

func main() {
    cause := errors.New("request failed")
    err := stackerr.WithStack(cause)
    fmt.Println(errors.Is(err, cause)) // true
    trace, ok := stackerr.TraceFrom(err)
    fmt.Println(ok, len(trace) > 0) // true true
    fmt.Printf("%+v\n", err) // 消息与堆栈
}
```

## 接口

| 接口 | 行为 |
| --- | --- |
| `WithStack(err)` | nil 保持 nil；错误链已有堆栈则复用原错误，否则添加当前堆栈 |
| `FromPanic(value)` | 将 recover 值转为带堆栈的错误；原值为 error 时保留错误链 |
| `TraceFrom(err)` | 沿错误链查找第一个实现 `Tracer` 的错误 |
| `Capture(skipFrames)` | 捕获调用堆栈；负 skip 按零处理，超过深度返回空堆栈 |
| `NewStackError/FormatStackError/Callers` | 现有带堆栈错误和捕获入口；包装错误支持 Unwrap |

本包不执行 recover、不记录日志、不处理业务错误码。堆栈在调用处捕获；与 panic 配合时应在 deferred recover 中调用 `FromPanic`，不能在协程退出后再补原始现场。
