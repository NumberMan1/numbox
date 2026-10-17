# SkipList

以 `*big.Float` 为键的有序索引，导入路径为 `github.com/NumberMan1/numbox/pkg/skiplist`。

## 开源来源与许可

本包基于 [sean-public/fast-skiplist](https://github.com/sean-public/fast-skiplist)
进行定制，保留其跳表节点结构、前驱节点缓存及概率表设计。
上游采用 MIT 许可，版权声明为 `Copyright (c) 2017 sean`。
本目录 [LICENSE](LICENSE) 保留了上游许可证全文，派生源码也注明了来源与本地修改。

本次对照的上游版本为
[`d7f7945b944e5e211d906cbc4b08a07153d66ecb`](https://github.com/sean-public/fast-skiplist/tree/d7f7945b944e5e211d906cbc4b08a07153d66ecb)。
这是 2026-10-18 的核对参考；历史最初引入的 commit 和可能的中间 fork 未记录，
不将此参考版本作为历史引入版本。

## 本地调整

- 将上游 `float64` 键改为 `*big.Float`，通过 `Cmp` 比较大小，并同步调整查找、插入、删除及 `Element.Key()` 的类型。
- 将随机数源迁移至 `math/rand/v2` 和 PCG，调整初始化与随机数读取。
- 增加 `RangeVals(begin, end)`，按位置返回左闭右开的值区间。
- 增加 `RankIndex(key)` 排名辅助方法。

接口与行为以本包源码为准，上游的性能数据不直接代表本地定制版本。

成员分数排序与带跨度的排名查询由独立的 [`pkg/go-zset`](../go-zset/README.md) 提供。
更多第三方内容说明见 [THIRD_PARTY_NOTICES.md](../../THIRD_PARTY_NOTICES.md)。
