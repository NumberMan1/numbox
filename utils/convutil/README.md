# convutil

整数/布尔转换、整数切片与 Map 的数值转换、分隔文本编码解析。只依赖标准库与 typeutil，
不加载 JSON、随机数、正则或业务词表。

```go
value, err := convutil.ParseInteger[int8]("127")
values, err := convutil.ParseIntegerSlice[uint64]("1,18446744073709551615", ",")
text := convutil.FormatIntSliceString([]int64{-1, 0, 42})
```

- ParseInteger 严格检查目标类型范围，不裁剪空白；失败返回零值和 error。
- ParseIntegerSlice 拒绝空分隔符和空元素；空输入返回非 nil 空切片，失败返回 nil 切片。
- errors.Is 可检查 strconv.ErrSyntax / strconv.ErrRange。
- ParseIntString 系列保持宽松解析：忽略错误，窄类型转换可能截断。不得作为输入校验器。
- Integer 包含 int16 及其他整数命名类型，定义位于 utils/typeutil；整数转换不使用 reflect。
- FormatIntSliceString 直接拼接十进制字节；批量转换预分配，保留 nil/空切片行为。
- MarshalMapToString 顺序不固定；MarshalSliceToString 保留空字段。

```sh
go test ./utils/convutil -run '^$' -bench BenchmarkIntegerConversion -benchmem -count=3
```

基准中的 reflect 是旧实现对照，不进入生产代码。
