package randutil

import (
	"math"
	"testing"
)

func TestRandomFullWidthRanges(t *testing.T) {
	for seed := int64(0); seed < 100; seed++ {
		value := RandomInt(seed, int64(math.MinInt64), int64(math.MaxInt64))
		if value == math.MaxInt64 {
			t.Fatal("exclusive maximum")
		}
		if RandomInt(seed, int8(-128), int8(127)) == 127 {
			t.Fatal("exclusive int8 maximum")
		}
		if RandomInt(seed, 5, 6) != 5 {
			t.Fatal("single value")
		}
		if got := RandInt[int](math.MinInt, math.MaxInt); got < math.MinInt || got > math.MaxInt {
			t.Fatal(got)
		}
	}
	defer func() {
		if recover() == nil {
			t.Fatal("reversed bounds accepted")
		}
	}()
	RandomInt(0, 2, 1)
}
