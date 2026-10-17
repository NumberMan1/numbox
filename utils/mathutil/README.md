# mathutil

DivideCeil 使用整数运算计算向上取整除法，避免浮点精度损失。
零除数保留返回零的历史行为；有符号最小值除以 -1 溢出时 panic。

CompareWithNegOneAsMax 可交给 slices.SortFunc，将所有负值排在非负值之后，
相同符号内按数值排序。名称中的 NegOne 是历史名称，不仅处理 -1。
