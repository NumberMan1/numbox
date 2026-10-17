# 变更记录

记录面向使用者的变更。只有实际发布后才填写版本与发布日期；历史详情以
[Tags](https://github.com/NumberMan1/numbox/tags) 和
[提交记录](https://github.com/NumberMan1/numbox/commits) 为准，不在此推测历史发布日期。

## 未发布

## v0.4.1 - 2026-10-18

### 文档与许可

- 明确 `pkg/skiplist` 基于 `sean-public/fast-skiplist` 定制，补充上游地址、核对参考 commit、原始 MIT 许可全文及源码版权说明。
- 记录已有的 `big.Float` 键、`math/rand/v2` / PCG 及范围和排名辅助方法等本地调整；同步第三方说明。本次不改变运行时行为或公开接口。

## v0.4.0 - 2026-10-18

### 不兼容变更

- ZSet 明确拒绝 NaN 分数（Add 在修改前 panic），分数查询中的 NaN 边界返回空结果；负数保留数量按零处理。
- 移除 Numeric；纯键值索引使用 comparable，数值运算使用 utils/typeutil.Number，最大值使用 cmp.Ordered。外部自行引用 Numeric 的代码需要调整约束。
- 工具按职责拆至 `utils/convutil`、`sliceutil`、`maputil`、`jsonutil`、`randutil`、`timeutil`、`mathutil`、`strutil`、`ioutil`、`assertutil`、`defaultutil`、`typeutil`；env、percent、uuid 也位于 `utils` 子目录。根 `utils` 不提供 Go 包，权重配置解析归入 `pkg/weight`，详见 `docs/utils-migration.md`。
- 移除 `pkg/sensitive-word` 和内嵌词表，避免导入时自动构建全局匹配树。`pkg/matcher` 保留，由应用显式传入词表和校验策略。

### 修复

- `pkg/go-zset` 改为纯 Go 的成员索引与带排名跨度跳表，移除 C 源码、CGo 和 finalizer，支持 Windows 及关闭 CGo 的构建。保留公开接口与 1-based 排名语义，同分按成员字节序排序。
- ZSet 计数与逆序查询纳入锁保护，整个 Limit/RevLimit 裁剪在单次写锁内完成；限制范围查询的分配容量，分数更新复用节点。
- 整数转换改用无反射的泛型符号判断，保留宽松解析兼容行为；切片转换预分配，整数切片格式化直接拼接字节，解析避免中间字符串切片。
- 时间轮 key 校验改用无反射比较，继续拒绝不可比较的接口值及 NaN。

### 新增

- ZSet 排序模型、跳表跨度/反向链、边界、并发及模糊测试；CI 增加 Windows 和无 CGo 构建验证。
- 新增 Integer、Number 约束，以及检查目标类型范围的 ParseInteger / ParseIntegerSlice；补充转换边界回归测试和旧实现性能对照基准。

### 文档

- 补充纯 Go ZSet 的用法、并发与边界约定，以及与旧 CGo 实现的性能和进程内存对比，详见 `docs/zset-performance.md`。
- 更新 Windows 和无 CGo 构建说明，保留已移除 C 源码及词表的历史来源核对事项。

## v0.3.0 - 2026-10-17

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
