# typeutil

Integer 包含有符号/无符号整数（不含 uintptr）及对应命名类型；Number 还包含浮点数及其命名类型。
此包只有编译期约束，没有初始化或运行时类型识别。

Map 键直接使用 comparable；排序/最大值使用 cmp.Ordered。不要把所有约束都换成 Number。
