package utils

import "golang.org/x/exp/rand"

// RandomInt returns a reproducible value in [min,max); equal bounds return min.
func RandomInt[T ~int | ~int8 | ~int32 | ~int64](seed int64, min, max T) T {
	if min > max {
		panic("utils: min exceeds max")
	}
	if min == max {
		return min
	}
	r := rand.New(rand.NewSource(uint64(seed)))
	span := uint64(max) - uint64(min)
	// Retain the seeded sequence for ranges already supported by the old API.
	if span <= uint64(^uint(0)>>1) {
		return T(uint64(min) + uint64(r.Intn(int(span))))
	}
	return T(uint64(min) + boundedUint64(r.Uint64, span))
}

// RandInt returns a value in [min,max], including the full machine-int range.
func RandInt[T ~int | ~int8 | ~int32 | ~int64](min, max int) T {
	if min > max {
		panic("utils: min exceeds max")
	}
	span := uint64(max) - uint64(min) + 1
	if span > 0 && span <= uint64(^uint(0)>>1) {
		return T(uint64(min) + uint64(rand.Intn(int(span))))
	}
	return T(uint64(min) + boundedUint64(rand.Uint64, span))
}

// A zero bound represents the full uint64 range. Rejection avoids modulo bias.
func boundedUint64(next func() uint64, bound uint64) uint64 {
	if bound == 0 {
		return next()
	}
	threshold := -bound % bound
	for {
		value := next()
		if value >= threshold {
			return value % bound
		}
	}
}
