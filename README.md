# numbox 🧰

`numbox` 是 **NumberMan** 的go常用工具箱，包含日常开发中使用的小工具、脚本和实用函数。

---

## 📦 项目内容

`numbox` 目前包含以下工具：
- `errors` 带错误码的自定义错误结构
- `pkg` 业务中通常使用的工具或者函数
- [`pkg/lock`](pkg/lock/README.md) 不自动过期的进程内互斥锁、读写锁，以及支持 context 协作取消的任务执行接口
