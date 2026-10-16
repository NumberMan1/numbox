# numbox

[![Go CI](https://github.com/NumberMan1/numbox/actions/workflows/ci.yml/badge.svg)](https://github.com/NumberMan1/numbox/actions/workflows/ci.yml)
[![Go Reference](https://pkg.go.dev/badge/github.com/NumberMan1/numbox.svg)](https://pkg.go.dev/github.com/NumberMan1/numbox)

NumberMan 的 Go 工具箱，包含错误码、集合、随机抽样、时间、环境变量和常用数据结构。
模块路径为 `github.com/NumberMan1/numbox`，按需导入各子包。

## 安装与环境

要求 **Go 1.25.1 或以上版本**。在使用方的 Go module 中执行：

```sh
go get github.com/NumberMan1/numbox@latest
```

可将 `@latest` 替换为 [已有 tag](https://github.com/NumberMan1/numbox/tags) 以固定版本。
项目仍处于 v0 阶段，升级前请阅读 [变更记录](CHANGELOG.md) 与版本之间的差异。

`pkg/go-zset` 依赖 **CGo 和 C 编译器**，C 实现使用 POSIX `random()`。
完整 CI 使用 Ubuntu 和 macOS；Windows 完整开发建议使用 WSL。
不使用该包时，可以正常按需导入其他纯 Go 工具包。

## 快速开始

```go
package main

import (
    "fmt"

    "github.com/NumberMan1/numbox/pkg/collection"
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
| [`pkg/utils`](pkg/utils/README.md) | 字符串、切片、JSON、随机数等通用函数 |
| `pkg/bitmap` | 位图 |
| `pkg/collection` | 泛型集合 |
| [`pkg/dag`](pkg/dag/README.md) | 有向无环图 |
| [`pkg/encryp`](pkg/encryp/README.md) | JWT token 工具 |
| [`pkg/env`](pkg/env/README.md) | 环境变量读取 |
| `pkg/go-zset` | 基于 CGo 的有序集合 |
| `pkg/idcard` | 身份证校验与解析 |
| [`pkg/lock`](pkg/lock/README.md) | 不自动过期的 token 锁与 context 协作取消任务 |
| [`pkg/percent`](pkg/percent/README.md) | 千分比工具 |
| `pkg/sensitive-word` | 内嵌词表校验，包名为 `sensitiveword` |
| [`pkg/skiplist`](pkg/skiplist/README.md) | 以 `big.Float` 为 key 的跳表 |
| [`pkg/time`](pkg/time/README.md) | 毫秒时间戳 |
| `pkg/timewheel` | 时间轮调度 |
| `pkg/uuid` | UUID 工具 |
| `pkg/weight` | 权重池、随机抽样与保底机制 |
| [`pkg/matcher`](pkg/matcher/README.md) | 字符串匹配 |
| [`pkg/safego`](pkg/safego/README.md) | goroutine 启动与 panic 处理 |
| [`pkg/stackerr`](pkg/stackerr/README.md) | 带堆栈的错误 |
| [`pkg/version`](pkg/version/README.md) | 版本解析与比较 |
| [`pkg/event-bus`](pkg/event-bus/README.md) | 事件总线类型定义 |

新版工具统一位于 `pkg/*`，旧版 `utils/*` 导入路径需要迁移；公开函数的更名见各包 README。
锁不再自动过期，迁移规则见 [lock 说明](pkg/lock/README.md)。

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

CI 配置见 [.github/workflows/ci.yml](.github/workflows/ci.yml)。
贡献流程、依赖维护与验证要求见 [CONTRIBUTING.md](CONTRIBUTING.md)。

## 社区与许可

- [获取帮助](SUPPORT.md) · [报告问题](https://github.com/NumberMan1/numbox/issues)
- [贡献指南](CONTRIBUTING.md) · [行为准则](CODE_OF_CONDUCT.md)
- [安全政策](SECURITY.md) · [维护者与发布流程](MAINTAINERS.md)
- [变更记录](CHANGELOG.md) · [第三方来源与许可说明](THIRD_PARTY_NOTICES.md)

项目原创内容采用 [MIT 许可证](LICENSE)。第三方代码、依赖和数据的权利以其原始许可为准；
尚待确认的历史来源见第三方说明。
