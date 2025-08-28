package utils

import (
	"strings"
)

func MarshalMapToString[T1 comparable, T2 any](splitStr string, src map[T1]T2, fn func(k T1, v T2) string) string {
	var res string
	for key, val := range src {
		if res == "" {
			res += fn(key, val)
			continue
		}
		res += splitStr + fn(key, val)
	}
	return res
}

func MarshalSliceToString[T any](splitStr string, src []T, fn func(T) string) string {
	var res string
	for _, val := range src {
		if res == "" {
			res += fn(val)
			continue
		}
		res += splitStr + fn(val)
	}
	return res
}

func UnmarshalStringToSlice[T any](splitStr string, src string, fn func(int, string) (T, bool)) []T {
	var resArr = make([]T, 0)
	if src == "" {
		return resArr
	}

	var valArr = strings.Split(src, splitStr)
	for idx, str := range valArr {
		val, ok := fn(idx, str)
		if !ok {
			continue
		}
		resArr = append(resArr, val)
	}
	return resArr
}

func UnmarshalStringToMap[T1 comparable, T2 any](splitStr string, src string, fn func(string) (T1, T2, bool)) map[T1]T2 {
	var resMap = make(map[T1]T2)
	if src == "" {
		return resMap
	}

	var valArr = strings.Split(src, splitStr)
	for _, str := range valArr {
		k, v, ok := fn(str)
		if !ok {
			continue
		}
		resMap[k] = v
	}
	return resMap
}

func UnmarshalStringToStruct[T any](splitStr string, src string, resPtr *T, fn func(int, string, *T)) {
	valArr := strings.Split(src, splitStr)
	for idx, str := range valArr {
		fn(idx, str, resPtr)
	}
}
