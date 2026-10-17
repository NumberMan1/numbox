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

## 历史内嵌内容：来源待确认

以下项目需要维护者结合原始引入记录确认来源、版本与授权，再补入相应版权和许可证全文：

| 路径 | 当前证据与需要确认的内容 |
| --- | --- |
| `pkg/skiplist/skiplist.go`、`type.go` | 常量、结构和注释与 [sean-public/fast-skiplist](https://github.com/sean-public/fast-skiplist) 存在相似之处，这是来源线索，尚未确认实际引入版本或中间 fork。 |

代码相似性不能单独证明来源或授权，不据此给这些内容补写猜测的作者、年份或许可。
完成来源核对前，不应宣称上述历史内嵌内容已统一获得 MIT 授权。
后续引入的第三方代码和数据必须在引入时记录上游地址、版本、修改范围与许可。

## 已移除的内容

`pkg/go-zset/skiplist.c`、`skiplist.h` 已由纯 Go 实现替换。历史 C 源码说明其实现类似 Redis，
与 [xjdrew/lua-zset](https://github.com/xjdrew/lua-zset) 的代码存在来源线索；历史引入链、版本
及授权仍未确认。替换不代表历史版本的来源与授权已经完成核对。

`pkg/sensitive-word/sensitive_word.json` 及其加载组件已移除。该历史词表的原始来源、
生成方式和许可未得到确认；移除不代表已确认历史版本的再分发授权。当前版本不再包含该词表。
