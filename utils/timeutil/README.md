# timeutil

ParseCNDuration 完整解析非负整数中文时长，例如 1天2小时3分钟4秒5毫秒。
未知文本、负数、小数和溢出返回 error。天固定表示 24 小时，不表示时区日历日。

此包只处理时长文本；时间戳和业务日历时钟仍在 pkg/time，主动时间轮在 pkg/timewheel。
