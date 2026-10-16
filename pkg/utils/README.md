# utils

泛型集合转换、字符串与 JSON、概率配置解析、数值计算和完整写入工具。各函数的错误处理不同，使用前应明确它返回 error、忽略错误还是 panic。

## 示例

```go
package main

import (
    "fmt"
    "github.com/NumberMan1/numbox/pkg/utils"
)

func main() {
    fmt.Println(utils.ToJSONString(map[string]int{"count": 2})) // {"count":2}
    fmt.Println(string(utils.ToJSONBytes([]int{1, 2}))) // [1,2]
    fmt.Println(utils.ParseIntString[uint64]("18446744073709551615"))
    fmt.Println(utils.FormatIntString(uint64(18446744073709551615)))
    fmt.Println(utils.DivideCeil(int64(5), int64(2))) // 3
    duration, err := utils.ParseCNDuration("1小时30分钟")
    fmt.Println(duration, err) // 1h30m0s <nil>
}
```

## 常用入口

| 类别 | 接口 |
| --- | --- |
| JSON | ToJSONString、ToJSONBytes、UnmarshalFromJSON、GetIntSliceFromJSON、GetMapFromJSON |
| 整数与字符串 | ParseIntString、FormatIntString、ParseIntTypeString、FormatIntTypeString、ParseIntSliceString |
| 切片与映射 | DeduplicateAddInt32、MapKeys、MapVals、SliceToMap、MaxValue、BatchSliceGroups |
| 数值与随机 | DivideCeil、RandomInt、RandInt、Shuffle、ParsePercentageToPermille |
| 配置文本 | ParseCNDuration、ParseWeightItems、ParseFixedTotalWeightPool |
| 连接与拆分 | MarshalSliceToString、MarshalMapToString、UnmarshalStringToSlice |
| 完整写入 | WriteFull |
| 断言与捕获 | Must、Assert、AssertStringError、WithRecover |

JSON 公开入口统一使用 `JSON` 大写命名，没有旧大小写入口。它们及 JSON 反序列化助手在编解码失败时 panic；ToJSONString(nil) 返回空字符串，ToJSONBytes(nil) 返回 JSON 的 `null`。

## 使用边界

- ParseIntString 系列保持宽松解析：忽略 strconv 错误，窄类型转换可能截断。需要严格校验时使用 strconv 的带错误接口。无符号转换支持完整 uint64 范围。
- ParseCNDuration 完整解析非负整数中文时长，如 `1天2小时3分钟4秒5毫秒`；未知文本、负数、小数和溢出返回 error。
- ParsePercentageToPermille 接受带 `%` 的数字并四舍五入为千分比，拒绝非有限或超出 int 范围的结果，不限制在 0%～100%。
- DivideCeil 用整数运算计算向上取整；除数为零保留返回零的语义，有符号最小值除以 -1 溢出时 panic。
- RandomInt(seed,min,max) 使用 `[min,max)`，相等时返回 min；RandInt 使用 `[min,max]`。调用方应保证结果能装入返回类型。上下界逆序时 panic，随机数不用于密码学。
- MapKeys/MapVals 和映射连接的顺序不固定；切片连接保留空字段。DeduplicateAddInt32 虽保留名称中的 Int32，实际支持任意 comparable 类型。
- WriteFull 补齐短写，拒绝零进展和非法字节计数；发生错误时此前的写入不会撤销。
- WithRecover 将 panic 和返回错误打印到标准输出，不把它们返回给调用方；需要结构化日志和异步 panic 上报时使用 safego。
