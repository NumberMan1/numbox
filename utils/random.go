package utils

import (
	"math/rand/v2"
)

// RandomInt 生成一个指定范围的随机整数
func RandomInt[T int | int8 | int32 | int64](seed int64, min, max T) T {
	if min >= max {
		return min // 如果 min >= max，直接返回 min
	}

	// 创建随机数生成器
	r := rand.New(rand.NewPCG(uint64(seed), uint64(seed)))

	// 计算随机值
	randomVal := int64(max - min)
	return T(r.Int64N(randomVal)) + min
}

// RandInt 生成一个在 [min, max] 范围内的随机整数 (包含min和max)。
func RandInt(min, max int) int {
	return min + rand.IntN(max-min+1)
}
