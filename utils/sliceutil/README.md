# sliceutil

提供 DeduplicateAddInt32、SliceContains、SliceAllContains、MaxValue、BatchSliceGroups。
DeduplicateAddInt32 虽保留历史名称，实际支持任意 comparable 元素，并保留首次出现顺序。

SliceAllContains 检查集合包含关系，不检查重复次数。MaxValue 支持 cmp.Ordered，
空输入返回零值。BatchSliceGroups 同步调用回调，size<=0 时一次处理全部非空输入；
回调收到原切片的视图，需要保留独立数据时自行复制。
