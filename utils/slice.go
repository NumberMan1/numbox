package utils

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
