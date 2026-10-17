package sliceutil

import (
	"cmp"

	"github.com/NumberMan1/numbox/utils/maputil"
)

func DeduplicateAddInt32[T comparable](arr []T, elements ...T) []T {
	// 使用 map 来记录已经存在的元素，key 为元素值，value 为 true 表示存在
	exists := make(map[T]bool)
	result := make([]T, 0)

	// 先将原数组中的元素添加到结果数组中（去重）
	for _, num := range arr {
		if !exists[num] {
			result = append(result, num)
			exists[num] = true
		}
	}

	// 添加新元素到结果数组中（确保新元素不重复）
	for _, element := range elements {
		if !exists[element] {
			result = append(result, element)
			exists[element] = true
		}
	}

	return result
}

func SliceContains[T comparable](arr []T, item T) bool {
	for _, v := range arr {
		if v == item {
			return true
		}
	}
	return false
}

func SliceAllContains[T comparable](arr []T, items ...T) bool {
	arrColl := maputil.SliceToMap(arr, func(val T) (T, bool) {
		return val, true
	})
	for _, item := range items {
		if !arrColl[item] {
			return false
		}
	}
	return true
}

func MaxValue[T cmp.Ordered](list []T) T {
	if len(list) == 0 {
		var zero T
		return zero
	}
	maxVal := list[0]
	for _, v := range list[1:] {
		if v > maxVal {
			maxVal = v
		}
	}
	return maxVal
}

func BatchSliceGroups[T any](all []T, size int, fn func([]T)) {
	n := len(all)
	if n == 0 {
		return
	}
	if size <= 0 || size >= n {
		fn(all)
		return
	}
	// 逐批次遍历执行（非并发），避免额外边界判断开销
	for start := 0; start < n; {
		end := start + size
		if end > n {
			end = n
		}
		fn(all[start:end])
		if end == n {
			break
		}
		start = end
	}
}
