package utils

import (
	"golang.org/x/exp/rand"
)

func RandomInt[T int | int8 | ~int32 | ~int64](seed int64, min, max T) T {
	if min == max {
		return min
	}
	randomVal := int(max - min)
	return T(rand.New(rand.NewSource(uint64(seed))).Intn(randomVal)) + min
}

// RandInt 生成一个在 [min, max] 范围内的随机整数 (包含min和max)。
func RandInt[T int | int8 | ~int32 | ~int64](min, max int) T {
	return T(min + rand.Intn(max-min+1))
}
