# collection

泛型集合 `Set[T comparable]`，用于去重与成员查询，不依赖其他 numbox 工具包。

## 示例

```go
package main

import (
    "fmt"
    "github.com/NumberMan1/numbox/pkg/collection"
)

func main() {
    set := collection.NewSet("a", "b", "a")
    set.Add("c")
    fmt.Println(set.Contains("a")) // true
    fmt.Println(set.ContainsAny("x", "b")) // true
    set.Remove("b")
    fmt.Println(len(set)) // 2
}
```

`All()` 返回独立的成员切片，顺序不固定。通过 `NewSet` 初始化后再 Add；未初始化的 nil Set 不能写入。集合不提供并发同步，跨协程访问由调用方加锁。
