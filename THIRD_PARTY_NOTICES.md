# 第三方来源与许可说明

根目录 [LICENSE](LICENSE) 用于本项目有权授权的原创内容，不替代第三方内容的原始许可。
本文件记录已知依赖和历史来源核对事项，不是完整的许可审计结论。

## Go module 依赖

依赖与版本以 [go.mod](go.mod) 和 go.sum 为准，依赖模块保留各自的版权与许可证。
发布包含依赖源码的压缩包、vendor 或其他再分发物时，应一并保留对应版本中的许可与通知。
直接依赖包括：

- [dgrijalva/jwt-go](https://github.com/dgrijalva/jwt-go)
- [google/uuid](https://github.com/google/uuid)
- [guanguans/id-validator](https://github.com/guanguans/id-validator)
- [json-iterator/go](https://github.com/json-iterator/go)
- [stretchr/testify](https://github.com/stretchr/testify)
- [golang.org/x/exp](https://go.googlesource.com/exp)

## 内嵌开源组件

### pkg/skiplist

- 上游项目：[sean-public/fast-skiplist](https://github.com/sean-public/fast-skiplist)。
- 本地派生源码：`pkg/skiplist/skiplist.go`、`pkg/skiplist/type.go`。
- 上游许可：MIT，版权声明 `Copyright (c) 2017 sean`；全文保留于
  [pkg/skiplist/LICENSE](pkg/skiplist/LICENSE)，派生源码保留相应来源说明。
- 2026-10-18 核对参考：上游 commit
  [`d7f7945b944e5e211d906cbc4b08a07153d66ecb`](https://github.com/sean-public/fast-skiplist/tree/d7f7945b944e5e211d906cbc4b08a07153d66ecb)，
  [对应 LICENSE](https://github.com/sean-public/fast-skiplist/blob/d7f7945b944e5e211d906cbc4b08a07153d66ecb/LICENSE)。
  历史最初引入的 commit 与可能的中间 fork 未记录，此 commit 仅用于本次来源和许可核对。
- 本地定制：将 `float64` 键改为 `*big.Float` 并使用 `Cmp` 比较；迁移至
  `math/rand/v2` / PCG；增加 `RangeVals` 与 `RankIndex` 辅助方法。
  具体说明见 [包文档](pkg/skiplist/README.md)。

本包是上游开源实现的本地定制版本，使用与再分发时须保留上游版权及许可通知。
后续引入的第三方代码和数据必须在引入时记录上游地址、版本、修改范围与许可。

## 已移除的内容

`pkg/go-zset/skiplist.c`、`skiplist.h` 已由纯 Go 实现替换。历史 C 源码说明其实现类似 Redis，
与 [xjdrew/lua-zset](https://github.com/xjdrew/lua-zset) 的代码存在来源线索；历史引入链、版本
及授权仍未确认。替换不代表历史版本的来源与授权已经完成核对。

`pkg/sensitive-word/sensitive_word.json` 及其加载组件已移除。该历史词表的原始来源、
生成方式和许可未得到确认；移除不代表已确认历史版本的再分发授权。当前版本不再包含该词表。
