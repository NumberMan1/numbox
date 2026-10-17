# percent

以千分比表示概率和倍率：`FullPercent=1000`、`HalfPercent=500`、`ZeroPercent=0`。

## 示例

```go
package main

import (
    "fmt"
    "github.com/NumberMan1/numbox/utils/percent"
)

func main() {
    p := percent.Percent(250) // 25%
    fmt.Println(p.Float64()) // 0.25
    fmt.Println(p.MultiplyByInt(20)) // 5
    fmt.Println(percent.GetPercentByDivide(1, 3)) // 333
    fmt.Println(percent.GetPercentByDivideCeil(1, 3)) // 334
    fmt.Println(percent.FullPercent.HitBySeed(42)) // true
}
```

## 语义

- ParsePercentageToPermille 解析带 % 的字符串并四舍五入为千分比，拒绝非有限值或超出 int 范围的结果；不限制在 0%～100%。
- FormatPercentString / ParsePercentString 处理整数形式的千分比；后者保留宽松解析和窄类型转换行为，不用于输入校验。
- `Hit()` 使用共享随机源；`HitBySeed(seed)` 使用独立随机源，同一概率和种子可复现结果。两者对 `p<=0` 总为 false，对 `p>=1000` 总为 true。
- 普通乘除结果转换为整数时向零截断，`Ceil` 系列向上取整。两种除法在除数为零时 panic。
- 算术运算不会自动限制结果到 `[0,1000]`。大整数计算会涉及浮点精度，不能用作精确财务运算。
- 随机判断不用于密码、令牌等安全用途。
