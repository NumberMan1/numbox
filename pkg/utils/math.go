package utils

import "math"

func DivideCeil[T ~int | ~int8 | ~int32 | ~int64](dividend, divisor T) T {
	if dividend == 0 || divisor == 0 {
		return T(0)
	}
	return T(math.Ceil(float64(dividend) / float64(divisor)))
}
