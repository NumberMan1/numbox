# 变更记录

记录面向使用者的变更。只有实际发布后才填写版本与发布日期；历史详情以
[Tags](https://github.com/NumberMan1/numbox/tags) 和
[提交记录](https://github.com/NumberMan1/numbox/commits) 为准，不在此推测历史发布日期。

## 未发布

### 不兼容变更

- 移除 Numeric；纯键值索引使用 comparable，数值运算使用 utils/typeutil.Number，最大值使用 cmp.Ordered。外部自行引用 Numeric 的代码需要调整约束。
- 工具按职责拆至 `utils/convutil`、`sliceutil`、`maputil`、`jsonutil`、`randutil`、`timeutil`、`mathutil`、`strutil`、`ioutil`、`assertutil`、`defaultutil`、`typeutil`；env、percent、uuid 也位于 `utils` 子目录。根 `utils` 不提供 Go 包，权重配置解析归入 `pkg/weight`，详见 `docs/utils-migration.md`。
- 移除 `pkg/sensitive-word` 和内嵌词表，避免导入时自动构建全局匹配树。`pkg/matcher` 保留，由应用显式传入词表和校验策略。
- 进程内锁不再自动过期。新增 `LockContext`、`RLockContext`、`Do`、`DoRead` 和 `WithTaskTimeout`，任务取消后等待任务退出才释放锁；旧持锁期限选项弃用，详见 `pkg/lock/README.md`。
- 最低 Go 版本提升至 1.25.1；时间轮改为由调用者主动推进，接口和生命周期约定见 `pkg/timewheel/README.md`。

### 修复

- 整数转换改用无反射的泛型符号判断，保留宽松解析兼容行为；切片转换预分配，整数切片格式化直接拼接字节，解析避免中间字符串切片。
- 时间轮 key 校验改用无反射比较，继续拒绝不可比较的接口值及 NaN。
- 修复普通锁 `Acquired` 与加锁、解锁之间的数据竞争；解锁回调在内部互斥锁释放后执行，允许回调查询持锁状态。
- 补齐 `go.sum` 并启用 CI 依赖缓存。

### 新增

- 新增 Integer、Number 约束，以及检查目标类型范围的 ParseInteger / ParseIntegerSlice；补充转换边界回归测试和旧实现性能对照基准。
- MIT 许可证、贡献指南、行为准则、安全政策、支持说明与维护发布约定。
- Issue / PR 模板、CODEOWNERS、Dependabot、编辑器与换行配置。
- Ubuntu 上的 Go 测试、竞态检测和静态检查工作流。

### 依赖

- 升级 `actions/checkout`、`actions/setup-go` 至 v7，并升级 `testify` 至 1.12.1。

### 文档

- 补充安装、快速开始、工具索引与 CGo / POSIX 构建要求。
- 记录第三方依赖及历史内嵌代码的来源核对事项，并记录敏感词表的移除。
