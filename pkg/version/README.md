# version

项目数字版本 `major.minor.patch[.build][-suffix]` 的格式化与比较。可带 `v` 前缀，不实现完整 SemVer 的预发布优先级。

## 示例

```go
package main

import (
    "fmt"
    "github.com/NumberMan1/numbox/pkg/version"
)

func main() {
    server := version.Version("v1.2.3.260814-release")
    fmt.Println(server.Tag()) // 1.2.3-260814
    fmt.Println(server.PkgSuffix()) // 1.2.3-260814
    fmt.Println(server.CompatibleVersion()) // v1.2.3.260814
    fmt.Println(server.Compatibility("v1.2.3.260813")) // true
    fmt.Println(server.Compatibility("v1.2.2.260814")) // false
    fmt.Println(version.Version("1.2").CompatibleVersion()) // 1.2.0
}
```

## 规则

- `Compatibility(other)` 将接收者作为服务端版本：前三段相同且服务端 build 不小于对方才兼容。缺失段按零参与数值比较。
- 数值编码前三段各允许 `[0,999]`，build 允许 `[0,999999]`。非法数字或越界时 Number 返回 0，Compatibility 返回 false；Number 返回 0 不能单独用于校验，因为合法零版本也编码为 0。
- `TrimSuffix` 从第一个 `-` 截断。Tag/PkgSuffix 去掉 `v`，最多格式化四段，短版本补到三段；PkgSuffix 缺少 build 时使用 `-0`。
- `Equal` 比较原字符串；`IsLowerThan` 比较 Number。不同前缀、后缀或缺失零段的字符串可能数值相等而 Equal 为 false。
- 格式化方法不代替输入校验，后缀不参与数值排序。
