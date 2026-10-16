# weight

基于累计权重的抽样工具，提供普通权重池、公平抽样和按未命中次数触发的保底池。显式传入种子可复现相同操作序列。

## 示例

```go
package main

import (
    "fmt"
    "github.com/NumberMan1/numbox/pkg/weight"
)

type prize struct {
    name string
    weight int
}
func (p prize) Weight() int { return p.weight }

func main() {
    pool := weight.NewWeightPool(42)
    pool.Add(prize{"common", 90})
    pool.Add(prize{"rare", 10})
    picked, ok := pool.PickRandomAndPutBack()
    fmt.Println(picked, ok)
    fmt.Println(pool.Length()) // 2，放回抽样不改变池
    fmt.Println(len(pool.PickManyRandom(2))) // 2，批量操作也不改变原池
    _, _ = pool.PickRandom() // 单次不放回抽样会移除命中项
    fmt.Println(pool.Length()) // 1
}
```

## 普通池

| 接口 | 行为 |
| --- | --- |
| `PickRandom()` | 一次不放回抽样，命中后移除 |
| `PickRandomAndPutBack()` | 一次放回抽样，不修改物品集合 |
| `PickManyRandom(count)` | 在副本上执行最多 count 次不放回尝试，返回独立切片，原池不变 |
| `PickRandomManyAndPutBack(count)` | count 次放回尝试 |
| `SetTotalWeight(total)` | total 可大于有效权重和，额外区间表示不命中；小于有效权重和时 panic |

零权重项不会普通命中，零总权重返回未命中。Add 拒绝 nil 接口、负数或溢出权重。批量次数不代表必定返回的数量：抽空或未命中都会减少结果。二分选择为 O(log N)，移除元素并维护索引为 O(N)。

通过构造函数初始化池。不要直接修改 `Items` 切片或已加入物品的权重。Pool 不提供并发同步；需要多协程访问时自行加锁。Copy 复制集合和索引，但按初始种子重建随机源，不复制当前随机流位置；PickManyRandom 使用该副本语义，因此相同池状态可重复得到相同批量结果。

## 保底池

`GuaranteedItem` 在 Weight 之外提供 `GuaranteedTimes() int` 和 `MissTimes() int`。保底阈值大于零且内部未命中数达到阈值时，下一次抽取优先从满足条件的项中随机选择，并移除命中项。普通未命中会增加所有剩余项的计数。

零权重项可通过保底选中；阈值小于等于零表示没有保底。`PickManyFromGuaranteedItems` 使用临时池，连续调用需要保留状态时应管理 `NewGuaranteedPool` 返回的池。随机工具不适用于安全令牌。
