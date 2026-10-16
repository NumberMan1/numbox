# 变更记录

记录面向使用者的变更。只有实际发布后才填写版本与发布日期；历史详情以
[Tags](https://github.com/NumberMan1/numbox/tags) 和
[提交记录](https://github.com/NumberMan1/numbox/commits) 为准，不在此推测历史发布日期。

## 未发布

### 不兼容变更

- 工具包从 `utils/*` 统一迁移到 `pkg/*`，并调整部分公开函数名称；升级时须更新导入路径并参照各包 README 迁移。
- 进程内锁不再自动过期。新增 `LockContext`、`RLockContext`、`Do`、`DoRead` 和 `WithTaskTimeout`，任务取消后等待任务退出才释放锁；旧持锁期限选项弃用，详见 `pkg/lock/README.md`。
- 最低 Go 版本提升至 1.25.1；时间轮改为由调用者主动推进，接口和生命周期约定见 `pkg/timewheel/README.md`。

### 修复

- 修复普通锁 `Acquired` 与加锁、解锁之间的数据竞争；解锁回调在内部互斥锁释放后执行，允许回调查询持锁状态。
- 补齐 `go.sum` 并启用 CI 依赖缓存。

### 新增

- MIT 许可证、贡献指南、行为准则、安全政策、支持说明与维护发布约定。
- Issue / PR 模板、CODEOWNERS、Dependabot、编辑器与换行配置。
- Ubuntu 上的 Go 测试、竞态检测和静态检查工作流。

### 依赖

- 升级 `actions/checkout`、`actions/setup-go` 至 v7，并升级 `testify` 至 1.12.1。

### 文档

- 补充安装、快速开始、工具索引与 CGo / POSIX 构建要求。
- 记录第三方依赖以及历史内嵌代码和词表的来源核对事项。
