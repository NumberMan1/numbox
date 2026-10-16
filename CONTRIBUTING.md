# 贡献指南

欢迎修复缺陷、补充文档、改进测试与提交通用工具。交流请遵守 [行为准则](CODE_OF_CONDUCT.md)；
使用问题见 [SUPPORT.md](SUPPORT.md)，漏洞见 [SECURITY.md](SECURITY.md)。

## 开始之前

先搜索已有 Issue。新增工具、公开接口调整或新依赖请先说明场景、预期 API 和兼容性影响。
Bug 报告请附上版本、Go 版本、运行环境及不依赖私有服务的最小复现。
中文和英文均可。

## 本地开发

要求 Go 1.25.1 或以上版本，模块路径为 `github.com/NumberMan1/numbox`。

`pkg/go-zset` 通过 CGo 编译 C 跳表，完整检查需要启用 CGo 和可用的 C 编译器。
该 C 实现使用 POSIX `random()`，目前 CI 在 Ubuntu 和 macOS 上运行；
Windows 完整构建建议使用 WSL。其他工具包可按需单独导入，不会因此编译 go-zset。

在仓库根目录运行：

```sh
go mod download all
go test -race -count=1 ./...
go vet ./...
git diff --check
```

添加或删除依赖时运行 `go mod tidy`，一并提交 `go.mod` 与 `go.sum`；
普通改动不要顺带升级无关依赖。没有 C 编译器时可以先检查所修改的纯 Go 包，
但这不能代替 Ubuntu 上的完整 CI。

## 修改约定

- 工具保持通用，避免引入私有业务依赖。
- 使用 `gofmt` 格式化修改的 Go 文件，导出符号应有准确文档注释。
- 说明边界值、错误、并发与资源释放约定；修复缺陷时补充可复现问题的行为测试。
- 涉及锁、共享状态或后台任务的修改须运行竞态检查；随机算法测试应使用固定种子或合理统计容差。
- 不提交日志、覆盖率、IDE 配置、凭据或真实个人数据。
- 引用外部代码或数据须记录来源与版本并保留原始许可，参见 [第三方说明](THIRD_PARTY_NOTICES.md)。
- API、依赖要求与使用方式变更同步更新 README、相关子包说明和 CHANGELOG。

## Pull Request

从 fork 创建分支，并以 `develop` 为 PR 目标，一次 PR 聚焦一个问题。`master` 由维护者明确授权后从 `develop` 合入。
描述解决的问题、行为变化、兼容性和实际验证结果；
未执行或失败的检查要明确说明原因。提交前自查 diff，避免无关格式化。

维护者评审通过后合并。提交贡献表示你有权提供这些内容，
并同意原创贡献采用项目 [MIT 许可证](LICENSE)；第三方内容保留其自身许可。
