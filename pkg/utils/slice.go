package utils

import (
	"cmp"
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
		}
	}

	return result
}

type Numeric interface {
	int | int32 | int64 | float32 | float64
}

func MapVals[T1 comparable, T2 any](src map[T1]T2) []T2 {
	var res = make([]T2, 0, len(src))
	for _, v := range src {
		res = append(res, v)
	}
	return res
}

func MapKeys[T1 comparable, T2 any](src map[T1]T2) []T1 {
	var res = make([]T1, 0, len(src))
	for k := range src {
		res = append(res, k)
	}
	return res
}

func SliceToMapKeys[T1 comparable, T2 any](src []T1, getValFn func(T1) T2) map[T1]T2 {
	var res = make(map[T1]T2, len(src))
	for _, v := range src {
		res[v] = getValFn(v)
	}
	return res
}

func SliceToMapVals[T1 any, T2 Numeric](src []T1, getKeyFn func(T1) T2) map[T2]T1 {
	var res = make(map[T2]T1, len(src))
	for _, v := range src {
		res[getKeyFn(v)] = v
	}
	return res
}

func SliceToMapWithIgnore[T1 any, T2 comparable, T3 any](src []T1, getKVFn func(T1) (T2, T3, bool)) map[T2]T3 {
	var res = make(map[T2]T3, len(src))
	for _, item := range src {
		key, val, ok := getKVFn(item)
		if !ok {
			continue
		}
		res[key] = val
	}
	return res
}

func SliceToMap[T1 any, T2 comparable, T3 any](src []T1, getKVFn func(T1) (T2, T3)) map[T2]T3 {
	var res = make(map[T2]T3, len(src))
	for _, item := range src {
		key, val := getKVFn(item)
		res[key] = val
	}
	return res
}

func SliceToCountMap[T1 comparable, T2 Numeric](arr []T1, getValFn func(T1) T2) map[T1]T2 {
	var res = make(map[T1]T2, len(arr))
	for _, v := range arr {
		res[v] += getValFn(v)
	}
	return res
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
	arrColl := SliceToMap(arr, func(val T) (T, bool) {
		return val, true
	})
	for _, item := range items {
		if !arrColl[item] {
			return false
		}
	}
	return true
}

// MergeSliceMaps 合并两个映射
func MergeSliceMaps[T1, T2 comparable](srcMap, anotherSrcMap map[T1][]T2) map[T1][]T2 {
	// 创建一个新的映射，容量为两个源映射的总和
	mergedMap := make(map[T1][]T2, len(srcMap)+len(anotherSrcMap))

	// 将第一个映射的键值对添加到合并映射中
	for key, values := range srcMap {
		mergedMap[key] = append([]T2{}, values...) // 复制切片
	}

	// 将第二个映射的键值对添加到合并映射中
	for key, values := range anotherSrcMap {
		// 如果键已存在，合并切片
		if existingValues, ok := mergedMap[key]; ok {
			mergedMap[key] = append(existingValues, values...) // 合并切片
		} else {
			mergedMap[key] = append([]T2{}, values...) // 复制切片
		}
	}

	return mergedMap
}

func MaxValue[T Numeric](list []T) T {
	var maxVal T
	for _, v := range list {
		if v > maxVal {
			maxVal = v
		}
	}
	return maxVal
}

// CompareWithNegOneAsMax 对比两个整数值，
// 按升序排序但将负数（如 -1）视作最大值。
// 适配 slices.SortFunc，返回值为：
// -1 表示 vi < vj
//
//	0 表示 vi == vj
//	1 表示 vi > vj
func CompareWithNegOneAsMax[T int | int8 | int32 | int64 | uint | uint8 | uint32 | uint64](vi T, vj T) int {
	if vi == vj {
		return 0
	}

	ineg, jneg := vi < 0, vj < 0

	if ineg && !jneg {
		return 1 // vi 为负数，视为无穷大，所以 vi > vj
	}
	if !ineg && jneg {
		return -1 // vj 为负数，视为无穷大，所以 vi < vj
	}

	// 符号相同时（同为正数或同为负数），按正常的数值大小进行升序比较
	return cmp.Compare(vi, vj)
}

func FormatIntTypeMap[T1 ~int32 | ~int64, T2 int | int8 | int32 | int64 | uint | uint8 | uint32 | uint64](src map[T1]T2) map[int32]int64 {
	var res = make(map[int32]int64, len(src))
	for k, v := range src {
		res[int32(k)] = int64(v)
	}
	return res
}

func ParseIntTypeMap[T1 ~int32 | ~int64, T2 int | int8 | int32 | int64 | uint | uint8 | uint32 | uint64](src map[int32]int64) map[T1]T2 {
	var res = make(map[T1]T2, len(src))
	for k, v := range src {
		res[T1(k)] = T2(v)
	}
	return res
}

func MapToSlice[T1 comparable, T2, T any](src map[T1]T2, fn func(k T1, v T2) (T, bool)) []T {
	var res = make([]T, 0, len(src))
	for k, v := range src {
		val, ok := fn(k, v)
		if !ok {
			continue
		}
		res = append(res, val)
	}
	return res
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
