# sensitive-word

使用内嵌词表判断名称是否含敏感子串，同时拒绝 Unicode 标点和非法 UTF-8。底层匹配由不可变 Matcher 实现，可并发调用。

## 示例

```go
package main

import (
    "fmt"
    sensitiveword "github.com/NumberMan1/numbox/pkg/sensitive-word"
)

func main() {
    fmt.Println(sensitiveword.IsStringValid("普通玩家123")) // true
    fmt.Println(sensitiveword.IsStringValid("普通玩家!")) // false
}
```

词表来自本包 `sensitive_word.json`，通过 go:embed 编入程序；更换内嵌词表需要重新构建。IsStringValid 同时应用词表匹配和标点规则，不提供长度检查、空白限制、大小写折叠或 Unicode 归一化，空字符串也不会单独被拒绝。

需要动态或自定义词表、查询命中区间或执行替换时，使用 [matcher](../matcher/README.md)；业务身份、名称长度和允许字符集由调用方校验。
