package mathutil

import (
	"testing"
)

func TestDivideCeilPrecisionAndSigns(t *testing.T) {
	for _, c := range []struct{ a, b, w int64 }{{9007199254740993, 1, 9007199254740993}, {5, 2, 3}, {-5, 2, -2}, {5, -2, -2}, {-5, -2, 3}, {1, 0, 0}} {
		if got := DivideCeil(c.a, c.b); got != c.w {
			t.Fatalf("%+v: %d", c, got)
		}
	}
}
