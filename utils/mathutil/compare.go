package mathutil

import (
	"cmp"
)

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
