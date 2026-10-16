# matcher

构造后不可变的 Unicode trie，词表由调用方注入，查询可并发执行。提供精确子串匹配、重叠区间查询和替换，不包含默认敏感词或名称校验规则。

## 示例

```go
package main

import (
    "fmt"
    "github.com/NumberMan1/numbox/pkg/matcher"
)

func main() {
    m := matcher.NewMatcher([]string{"aba", "ba", "坏词"})
    fmt.Println(m.Contains("末尾坏词")) // true
    fmt.Println(m.Find("ababa")) // [{0 3} {1 3} {2 5} {3 5}]
    fmt.Println(m.Replace("  坏词 & ababa!  ", '*')) //   ** & *****!
}
```

- Match.Start/End 是 rune 下标，表示 `[Start,End)`，不是字节偏移。
- Find 返回所有重叠匹配，按起始位置递增，同一起点按最长匹配优先。
- Replace 遮盖所有命中范围的并集，保留未命中的空格、标点及其他字符。
- 大小写敏感，不进行 Unicode 归一化。构造时忽略空词和非法 UTF-8 词条。
- 对非法 UTF-8 文本，Contains 返回 false，Find 返回 nil，Replace 返回原文本；文本合法性需由调用方单独校验。
