# ZSet

纯 Go 的成员索引与带排名跨度的跳表，支持 Windows、Linux、macOS，普通构建无需
CGo 或 C 编译器。导入路径和包名保留为 `github.com/NumberMan1/numbox/pkg/go-zset` / `zset`。

```go
z := zset.New()
z.Add(100, "alice")
z.Add(100, "bob")
z.Add(80, "carol")
rank := z.Rank("alice")       // 2
top := z.RevRange(1, 2)        // ["bob", "alice"]
removed := z.RevLimit(2)       // 1，保留最大的两项
```

## 排序与范围

- 按 float64 分数升序，再按成员字符串的字节序升序。同分排名不共享，每个成员有独立排名。
- Add 更新已有成员的分数，不增加数量；同分更新不改变位置。支持空字符串、Unicode 和含 NUL 的成员。
- Rank / RevRank 从 1 开始，不存在返回 0。Score 返回分数与存在标志。
- Range / RevRange 的两端均包含，排名从 1 开始，端点逆序则反向遍历，不支持负数排名。
  起点越界返回空结果，尾端越界截断；大端点不会导致按该端点大小分配内存。
- RangeByScore 是包含边界的分数区间，起始分数大于结束分数时降序，等分区间返回成员升序。
- Limit 保留最小的 count 项，RevLimit 保留最大的 count 项；返回实际删除数量。
  count<=0 清空集合。
- 分数允许正负无穷；Add 遇到 NaN 会在修改前 panic。RangeByScore 遇到 NaN 边界返回空切片。
- float64 的精度限制保持不变。Dump 只用于诊断，不承诺固定输出格式。

## 并发与复杂度

零值可直接使用，使用后不可复制。所有公开方法均受实例 RWMutex 保护；每次查询及整个
裁剪操作具有一致状态。多个方法调用之间不是事务，跨调用的结果可能被其他协程更新。
返回的成员切片由调用方持有，不借出内部节点。

成员查分、计数为 O(1)；添加、更新、删除和排名预期 O(log N)；范围查询及批量裁剪
预期 O(log N + K)，K 为返回/删除项数。更新分数复用原节点与层数组，不额外分配节点。
新成员的跳表层数按 1/4 概率提升，最多 32 层，不使用反射或外部数据结构依赖。

## 验证

```sh
CGO_ENABLED=0 go test ./pkg/go-zset
go test -race ./pkg/go-zset
go test ./pkg/go-zset -run '^$' -fuzz FuzzPureGoOperations -fuzztime=10s
go test ./pkg/go-zset -run '^$' -bench BenchmarkPureGoZSet -benchmem
```

测试包含排序 Map 对照、每层跨度与反向链校验、随机更新/删除/裁剪和并发读写。
本实现没有 C 内存与 finalizer，生命周期由 Go GC 管理。竞态检测本身仍受 Go 工具链的 CGo 要求约束。

2026-10-18 本地基准：Windows/amd64、Ryzen 7 5800H、10,000 个成员，
`CGO_ENABLED=0 -cpu=1 -benchtime=300ms -count=3`，下表为三次测量中位数。

| 操作 | ns/op | B/op | allocs/op |
| --- | ---: | ---: | ---: |
| 排名查询 | 135.7 | 0 | 0 |
| 已有成员更新分数 | 330.0 | 0 | 0 |
| 返回 20 个连续排名成员 | 479.7 | 320 | 1 |

与旧 CGo 实现的同场景对照、分配及进程内存实验见 [性能对比](../../docs/zset-performance.md)。
