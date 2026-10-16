# numbox

[![Go CI](https://github.com/NumberMan1/numbox/actions/workflows/go.yml/badge.svg)](https://github.com/NumberMan1/numbox/actions/workflows/go.yml)
[![Go Reference](https://pkg.go.dev/badge/github.com/NumberMan1/numbox.svg)](https://pkg.go.dev/github.com/NumberMan1/numbox)

NumberMan 的 Go 工具箱，包含错误码、集合、随机抽样、时间、环境变量和常用数据结构。
模块路径为 `github.com/NumberMan1/numbox`，按需导入各子包。

## 安装与环境

要求 **Go 1.24.0 或以上版本**。在使用方的 Go module 中执行：

```sh
go get github.com/NumberMan1/numbox@latest
```

可将 `@latest` 替换为 [已有 tag](https://github.com/NumberMan1/numbox/tags) 以固定版本。
项目仍处于 v0 阶段，升级前请阅读 [变更记录](CHANGELOG.md) 与版本之间的差异。

`utils/go-zset` 依赖 **CGo 和 C 编译器**，C 实现使用 POSIX `random()`。
完整构建与 CI 使用 Ubuntu；Windows 完整开发建议使用 WSL。
不使用该包时，可以正常按需导入其他纯 Go 工具包。

## 快速开始

```go
package main

import (
    "fmt"

    "github.com/NumberMan1/numbox/utils/collection"
)

func main() {
    names := collection.NewSet("alice", "bob")
    names.Add("alice", "carol")
    names.Remove("bob")

    fmt.Println(names.Contains("alice")) // true
    fmt.Println(names.Contains("bob"))   // false
}
```

## 工具索引

| 导入路径（模块路径后的部分） | 用途 |
| --- | --- |
| `errors` | 错误码、自定义错误与标准错误判断 |
| `utils` | 字符串、切片、JSON、随机数等通用函数 |
| `utils/bitmap` | 位图 |
| `utils/collection` | 泛型集合 |
| [`utils/dag`](utils/dag/README.md) | 有向无环图 |
| [`utils/encryp`](utils/encryp/README.md) | JWT token 工具 |
| [`utils/env`](utils/env/README.md) | 环境变量读取 |
| `utils/go-zset` | 基于 CGo 的有序集合 |
| `utils/idcard` | 身份证校验与解析 |
| `utils/lock` | 带 token、超时与过期选项的锁 |
| [`utils/percent`](utils/percent/README.md) | 千分比工具 |
| `utils/sensitive-word` | 内嵌词表校验，包名为 `sensitiveword` |
| [`utils/skiplist`](utils/skiplist/README.md) | 以 `big.Float` 为 key 的跳表 |
| [`utils/time`](utils/time/README.md) | 毫秒时间戳 |
| `utils/timewheel` | 时间轮调度 |
| `utils/uuid` | UUID 工具 |
| `utils/weight` | 权重池、随机抽样与保底机制 |

具体 API 见 [Go 文档](https://pkg.go.dev/github.com/NumberMan1/numbox) 和对应包源码。
不同工具的并发与生命周期约定各异，例如 `bitmap.BitMap` 不支持并发修改，
使用前应阅读对应接口说明。

## 开发与验证

在已启用 CGo 并安装 C 编译器的环境中执行：

```sh
go mod download all
go test -race -count=1 ./...
go vet ./...
git diff --check
```

CI 配置见 [.github/workflows/go.yml](.github/workflows/go.yml)。
贡献流程、依赖维护与验证要求见 [CONTRIBUTING.md](CONTRIBUTING.md)。

## 社区与许可

- [获取帮助](SUPPORT.md) · [报告问题](https://github.com/NumberMan1/numbox/issues)
- [贡献指南](CONTRIBUTING.md) · [行为准则](CODE_OF_CONDUCT.md)
- [安全政策](SECURITY.md) · [维护者与发布流程](MAINTAINERS.md)
- [变更记录](CHANGELOG.md) · [第三方来源与许可说明](THIRD_PARTY_NOTICES.md)

项目原创内容采用 [MIT 许可证](LICENSE)。第三方代码、依赖和数据的权利以其原始许可为准；
尚待确认的历史来源见第三方说明。
