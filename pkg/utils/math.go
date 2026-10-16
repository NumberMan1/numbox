package utils

// DivideCeil computes ceil(dividend/divisor) without losing integer precision.
// A zero dividend or divisor preserves the legacy result of zero.
func DivideCeil[T ~int | ~int8 | ~int32 | ~int64](dividend, divisor T) T {
	if dividend == 0 || divisor == 0 {
		return 0
	}
	q, r := dividend/divisor, dividend%divisor
	if dividend < 0 && divisor == -1 && q == dividend {
		panic("utils: integer division overflow")
	}
	if r != 0 && (r > 0) == (divisor > 0) {
		q++
	}
	return q
}
