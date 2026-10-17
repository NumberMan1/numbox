# maputil

提供 MapKeys、MapVals、SliceToMap 系列、SliceToCountMap、MergeSliceMaps、MapToSlice。
索引键只要求 comparable；需要累加的值使用 utils/typeutil.Number。

SliceToMapVals 支持字符串、结构体和命名类型作为键，重复键由后项覆盖。
Map 遍历顺序不固定。MergeSliceMaps 返回独立的值切片；元素引用本身不会深拷贝。
