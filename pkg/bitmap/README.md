# bitmap

用连续位存储整数集合，适用于状态标记和成员判断。数据在每字节内按低位到高位排列；两种位图均不提供并发同步。

## 示例

```go
package main

import (
    "fmt"
    "github.com/NumberMan1/numbox/pkg/bitmap"
)

func main() {
    b := bitmap.NewBitmap(10) // 支持 0～9
    fmt.Println(b.Add(9)) // true
    fmt.Println(b.Add(9)) // false，重复添加
    fmt.Println(b.Add(10)) // false，越界
    fmt.Println(b.Len()) // 1
    restored := bitmap.NewBitmap(10)
    if err := restored.Load(b.Bytes()); err != nil { panic(err) }
    fmt.Println(restored.Has(9)) // true
}
```

## 接口与边界

| 类型 | 范围 | 增删及查询 | 序列化 |
| --- | --- | --- | --- |
| Bitmap，NewBitmap(size) | `[0,size)` | Add/Remove 返回是否改变，Has 查询，Len 为 uint32 计数 | Bytes 返回副本，Load 验证并复制 |
| BitMap，NewBitMap(max) | `[0,max]` | Add/Remove 返回越界错误，IsExist 查询，Count 为 int32 计数 | Bits 返回副本，Load 返回验证错误，SetBits 验证失败会 panic |

新代码优先选择 Bitmap。重复增删不会重复更新计数。Load 要求字节长度与容量一致，未使用的尾部位必须为零；验证失败不改变现有状态。Bytes/Bits 和 Load 输入均不与内部存储共享，修改导出数据不会改变位图。

Bitmap 的 Size/Len 分别表示容量和当前成员数，Clear 清空内容。较大的 max 不应使用 BitMap 的 int32 Count 表示全部成员数量。
