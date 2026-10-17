# 工具目录与整数转换迁移说明

## 操作范围

| 当前目录 | 目标目录 |
| --- | --- |
| pkg/utils（或根 utils） | 根据函数职责选择下表中的子包 |
| pkg/env | utils/env |
| pkg/percent | utils/percent |
| pkg/uuid | utils/uuid |

锁、时间轮、时间组件、权重池、集合、matcher 等保留在 pkg；errors 保持根目录。
根 `utils` 只有目录索引，没有 Go 包，也不提供聚合兼容层。函数名称保持不变，调用方
按下表更新导入路径和包限定符；不能把旧导入机械替换成根 `utils`。

| 原函数或类型 | 新包（模块路径后的部分） |
| --- | --- |
| ParseInteger、ParseIntString、FormatIntString、整数切片/Map 转换、布尔转换、MarshalSliceToString、UnmarshalStringToSlice、ParseItemsString | utils/convutil |
| DeduplicateAddInt32、SliceContains、SliceAllContains、MaxValue、BatchSliceGroups | utils/sliceutil |
| MapKeys、MapVals、SliceToMap 系列、SliceToCountMap、MergeSliceMaps、MapToSlice | utils/maputil |
| ToJSONString、ToJSONBytes、UnmarshalFromJSON、GetIntSliceFromJSON、GetInt32SliceFromJSON、GetMapFromJSON | utils/jsonutil |
| RandomInt、RandInt、Shuffle、ShuffleRandom | utils/randutil |
| ParseCNDuration | utils/timeutil |
| DivideCeil、CompareWithNegOneAsMax | utils/mathutil |
| GetNameStringLength | utils/strutil |
| WriteFull | utils/ioutil |
| Must、Assert、AssertStringError、WithRecover | utils/assertutil |
| GetParamWithDefaultValue、GetBoolParamWithDefault、GetStringParamWithDefault | utils/defaultutil |
| Integer、Number | utils/typeutil |
| FormatPercentString、ParsePercentString、ParsePercentageToPermille | utils/percent |
| ParseWeightItems、ParseWeightMultiItems、ParseWeightMultiOptions、ParseFixedTotalWeightPool | pkg/weight |

例如 `utils.ParseInteger[int8](text)` 改为导入
`github.com/NumberMan1/numbox/utils/convutil` 并调用 `convutil.ParseInteger[int8](text)`。
本次没有修改 Orbit 或其他仓库的依赖版本。

## 敏感词包

已移除 pkg/sensitive-word 中以下四个文件：

- handle.go
- handle_test.go
- sensitive_word.json（374,480 字节）
- README.md

根 README、依赖说明和第三方说明已同步更新。历史上词表来源未确认的事实仍有记录。

pkg/matcher 保留，应用可显式提供自己的词表。它只提供精确匹配，不自动替代
IsStringValid 的禁止标点策略，使用方必须显式恢复自己的业务校验。

只有导入敏感词包的程序会运行其 init、保留内嵌词表并构建匹配树；
未导入该包的程序不会因本次移除获得同样的内存节省。不承诺未经测量的内存数值。

## 约束与整数转换

- Numeric 已删除。仅作为 Map 键的约束改为 comparable；求和与数值转换使用 Number；
  MaxValue 使用 cmp.Ordered。Integer / Number 支持命名类型，Integer 还包含 int16。
- ParseIntString 等旧入口保持宽松解析，不会因为本次优化自动变成严格校验。
- ParseInteger[T] 新增目标类型范围检查；ParseIntegerSlice[T](src, separator) 严格校验
  每个元素。出错返回零值或 nil 切片和 error；错误可用 errors.Is 检查 strconv.ErrRange。
- 切片转换保留原来的 nil/空切片语义；格式化避免中间字符串切片。

## 性能与验证

生产代码直接使用 reflect 的整数解析和时间轮 key 校验已替换，测试中的 reflect 用于
兼容性断言和旧实现性能对照。JSON 仍依赖 json-iterator，不承诺动态序列化完全没有反射。
可在 utils/convutil 与 pkg/timewheel 运行基准，比较旧实现的耗时、字节与分配次数。
go-zset 现已替换为纯 Go 实现，原 Windows/POSIX 限制已解除；普通构建支持关闭 CGo。

2026-10-18 拆包前的转换优化对照：Windows/amd64、Go 1.27.2、Ryzen 7 5800H，
`-cpu=1 -benchtime=400ms -count=3`，下表为三次测量中位数。

| 场景 | 原实现 | 当前实现 | 原/当前分配次数 |
| --- | --- | --- | --- |
| 单个 int64 解析 | 27.38 ns | 25.48 ns | 0 / 0 |
| 5 个 int64 拼接字符串 | 505.4 ns | 307.3 ns | 7 / 4 |
| 时间轮 int key 校验 | 49.19 ns | 9.028 ns | 1 / 0 |
| 时间轮 any(int) key 校验 | 34.42 ns | 9.148 ns | 1 / 0 |

单值解析两组测量区间重叠，不据此宣称稳定加速；切片样例分配量从 367 B 降至 176 B。
这些微基准不代表完整业务吞吐。测试中保留旧实现便于在实际部署环境重测。
