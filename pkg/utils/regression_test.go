package utils

import (
	"math"
	"reflect"
	"testing"
	"time"
)

func TestSliceAndMarshalRegressions(t *testing.T) {
	if got := DeduplicateAddInt32([]int{1, 1}, 2, 2, 3); !reflect.DeepEqual(got, []int{1, 2, 3}) {
		t.Fatal(got)
	}
	if got := MaxValue([]int{-3, -1, -2}); got != -1 {
		t.Fatal(got)
	}
	if got := MarshalSliceToString(",", []string{"", "x", ""}, func(s string) string { return s }); got != ",x," {
		t.Fatal(got)
	}
}
func TestUnsignedStringRoundtrip(t *testing.T) {
	type ID uint64
	const value = ID(math.MaxUint64)
	s := FormatIntString(value)
	if s != "18446744073709551615" {
		t.Fatal(s)
	}
	if got := ParseIntString[ID](s); got != value {
		t.Fatal(got)
	}
	var receiver ID
	ParseIntStringWithReceiver(s, &receiver)
	if receiver != value {
		t.Fatal(receiver)
	}
	if FormatIntString(int64(-42)) != "-42" || ParseIntString[int64]("-42") != -42 {
		t.Fatal("signed conversion changed")
	}
}
func TestDivideCeilPrecisionAndSigns(t *testing.T) {
	for _, c := range []struct{ a, b, w int64 }{{9007199254740993, 1, 9007199254740993}, {5, 2, 3}, {-5, 2, -2}, {5, -2, -2}, {-5, -2, 3}, {1, 0, 0}} {
		if got := DivideCeil(c.a, c.b); got != c.w {
			t.Fatalf("%+v: %d", c, got)
		}
	}
}
func TestChineseDurationValidation(t *testing.T) {
	for _, s := range []string{"", "bad", "1天junk", "-1秒", "1.5秒", "1000000000000000000天"} {
		if _, err := ParseCNDuration(s); err == nil {
			t.Fatalf("accepted %q", s)
		}
	}
	if got, err := ParseCNDuration("1天 2小时3分钟4秒5毫秒"); err != nil || got != 26*time.Hour+3*time.Minute+4*time.Second+5*time.Millisecond {
		t.Fatal(got, err)
	}
}
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

func TestPercentageNonFiniteAndOverflow(t *testing.T) {
	for _, s := range []string{"NaN%", "+Inf%", "-Inf%", "1e100%"} {
		if _, err := ParsePercentageToPermille(s); err == nil {
			t.Fatal(s)
		}
	}
	if value, err := ParsePercentageToPermille("12.34%"); err != nil || value != 123 {
		t.Fatal(value, err)
	}
}
