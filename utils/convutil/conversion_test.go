package convutil

import (
	"errors"
	"math"
	"slices"
	"strconv"
	"testing"

	"github.com/NumberMan1/numbox/utils/typeutil"
)

type namedSigned int16

type namedUnsigned uint64

func checkIntegerRange[T typeutil.Integer](t *testing.T, minText, maxText string, minimum, maximum T, invalid ...string) {
	t.Helper()
	for _, c := range []struct {
		text string
		want T
	}{{minText, minimum}, {maxText, maximum}, {"0", 0}} {
		got, err := ParseInteger[T](c.text)
		if err != nil || got != c.want {
			t.Fatalf("ParseInteger(%q) = %v, %v; want %v", c.text, got, err, c.want)
		}
		if FormatIntString(got) != c.text {
			t.Fatalf("format roundtrip: %q", c.text)
		}
	}
	for _, text := range invalid {
		got, err := ParseInteger[T](text)
		if err == nil || got != 0 {
			t.Fatalf("accepted %q: %v, %v", text, got, err)
		}
		var numErr *strconv.NumError
		if !errors.As(err, &numErr) {
			t.Fatalf("not a strconv error: %v", err)
		}
	}
}

func TestParseIntegerRanges(t *testing.T) {
	t.Run("int8", func(t *testing.T) { checkIntegerRange(t, "-128", "127", int8(-128), int8(127), "128", "-129") })
	t.Run("int16", func(t *testing.T) {
		checkIntegerRange(t, "-32768", "32767", int16(-32768), int16(32767), "32768", "-32769")
	})
	t.Run("int32", func(t *testing.T) {
		checkIntegerRange(t, "-2147483648", "2147483647", int32(math.MinInt32), int32(math.MaxInt32), "2147483648", "-2147483649")
	})
	t.Run("int64", func(t *testing.T) {
		checkIntegerRange(t, "-9223372036854775808", "9223372036854775807", int64(math.MinInt64), int64(math.MaxInt64), "9223372036854775808", "-9223372036854775809")
	})
	t.Run("int", func(t *testing.T) {
		checkIntegerRange(t, strconv.Itoa(math.MinInt), strconv.Itoa(math.MaxInt), math.MinInt, math.MaxInt, "18446744073709551616")
	})
	t.Run("uint8", func(t *testing.T) { checkIntegerRange(t, "0", "255", uint8(0), uint8(math.MaxUint8), "256", "-1") })
	t.Run("uint16", func(t *testing.T) {
		checkIntegerRange(t, "0", "65535", uint16(0), uint16(math.MaxUint16), "65536", "-1")
	})
	t.Run("uint32", func(t *testing.T) {
		checkIntegerRange(t, "0", "4294967295", uint32(0), uint32(math.MaxUint32), "4294967296", "-1")
	})
	t.Run("uint64", func(t *testing.T) {
		checkIntegerRange(t, "0", "18446744073709551615", uint64(0), uint64(math.MaxUint64), "18446744073709551616", "-1")
	})
	t.Run("uint", func(t *testing.T) {
		checkIntegerRange(t, "0", strconv.FormatUint(uint64(^uint(0)), 10), uint(0), ^uint(0), "18446744073709551616", "-1")
	})
	t.Run("named", func(t *testing.T) {
		checkIntegerRange(t, "-32768", "32767", namedSigned(-32768), namedSigned(32767), "32768")
		checkIntegerRange(t, "0", "18446744073709551615", namedUnsigned(0), namedUnsigned(math.MaxUint64), "-1")
	})
	for _, text := range []string{"", " ", " 1", "1 ", "1.5", "0x10", "1_000", "hello"} {
		if _, err := ParseInteger[int](text); !errors.Is(err, strconv.ErrSyntax) {
			t.Fatalf("%q: %v", text, err)
		}
	}
	if got, err := ParseInteger[int]("+12"); got != 12 || err != nil {
		t.Fatal(got, err)
	}
	if _, err := ParseInteger[int8]("128"); !errors.Is(err, strconv.ErrRange) {
		t.Fatal(err)
	}
}

func compareLooseParsing[T typeutil.Integer](t *testing.T) {
	t.Helper()
	for _, text := range []string{"", "bad", "-1", "+1", "128", "256", "65536", "9223372036854775808", "18446744073709551615", "18446744073709551616", "-9223372036854775809"} {
		if got, want := ParseIntString[T](text), parseIntReflect[T](text); got != want {
			t.Fatalf("%q: got %v, want legacy %v", text, got, want)
		}
	}
}

func TestLooseIntegerCompatibility(t *testing.T) {
	compareLooseParsing[int](t)
	compareLooseParsing[int8](t)
	compareLooseParsing[int16](t)
	compareLooseParsing[int32](t)
	compareLooseParsing[int64](t)
	compareLooseParsing[uint](t)
	compareLooseParsing[uint8](t)
	compareLooseParsing[uint16](t)
	compareLooseParsing[uint32](t)
	compareLooseParsing[uint64](t)
	compareLooseParsing[namedSigned](t)
	compareLooseParsing[namedUnsigned](t)
}

func TestIntegerSliceConversions(t *testing.T) {
	signed := []int64{math.MinInt64, -1, 0, 42, math.MaxInt64}
	unsigned := []namedUnsigned{0, math.MaxUint64}
	if got, err := ParseIntegerSlice[int64](FormatIntSliceString(signed), ","); err != nil || !slices.Equal(got, signed) {
		t.Fatal(got, err)
	}
	if got, err := ParseIntegerSlice[namedUnsigned](FormatIntSliceString(unsigned), ","); err != nil || !slices.Equal(got, unsigned) {
		t.Fatal(got, err)
	}
	if got, err := ParseIntegerSlice[int16]("1::2", "::"); err != nil || !slices.Equal(got, []int16{1, 2}) {
		t.Fatal(got, err)
	}
	for _, text := range []string{"1,,2", ",1", "1,", "1,128", "1,bad"} {
		if got, err := ParseIntegerSlice[int8](text, ","); err == nil || got != nil {
			t.Fatalf("accepted %q: %v, %v", text, got, err)
		}
	}
	if _, err := ParseIntegerSlice[int]("12", ""); err == nil {
		t.Fatal("empty separator accepted")
	}
	if got, err := ParseIntegerSlice[int]("", ","); err != nil || got == nil || len(got) != 0 {
		t.Fatal(got, err)
	}
	if got := ParseIntSliceString[int]("1,,bad,", ","); !slices.Equal(got, []int{1, 0, 0, 0}) {
		t.Fatal(got)
	}
	if got := ParseIntSliceString[int]("12", ""); !slices.Equal(got, []int{1, 2}) {
		t.Fatal(got)
	}
	if got := FormatIntSliceToStringSlice([]int{}); got != nil {
		t.Fatal("empty slice changed")
	}
	if got := ParseStringSliceToIntSlice[int](nil); got == nil || len(got) != 0 {
		t.Fatal("empty slice changed")
	}
}
