# 工具子包

此目录只作索引，没有根 `utils` Go 包。按功能导入具体子包，例如
`github.com/NumberMan1/numbox/utils/convutil`。

| 子包 | 职责 |
| --- | --- |
| [convutil](convutil/README.md) | 整数/布尔转换、数值转换、分隔文本编码解析 |
| [sliceutil](sliceutil/README.md) | 切片去重、包含判断、最大值和批处理 |
| [maputil](maputil/README.md) | Map 键值提取、索引、归并、切片与 Map 转换 |
| [jsonutil](jsonutil/README.md) | JSON 编解码及带类型的 JSON 转换 |
| [randutil](randutil/README.md) | 随机整数和洗牌，可复现种子 |
| [timeutil](timeutil/README.md) | 中文时长解析；时间戳与定时器在 pkg/time |
| [mathutil](mathutil/README.md) | 整数向上取整除法、特殊排序比较 |
| [strutil](strutil/README.md) | 字符串显示长度规则 |
| [ioutil](ioutil/README.md) | 完整写入、短写与零进展检查 |
| [assertutil](assertutil/README.md) | Must、Assert、WithRecover |
| [defaultutil](defaultutil/README.md) | 可选参数默认值 |
| [typeutil](typeutil/README.md) | Integer、Number 泛型约束 |
| [env](env/README.md) | 环境变量读取 |
| [percent](percent/README.md) | 千分比、百分数字符串及概率计算 |
| uuid | UUID 生成 |

整数、切片及 Map 工具不会依赖 JSON、正则或随机库。JSON 的动态类型处理限制在
jsonutil；中文时长正则只由 timeutil 初始化。权重配置解析属于 [pkg/weight](../pkg/weight/README.md)。

旧导入路径及函数归属见 [迁移说明](../docs/utils-migration.md)。
