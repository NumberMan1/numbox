package convutil

import (
	"reflect"
	"strconv"
	"strings"
	"testing"

	"github.com/NumberMan1/numbox/utils/typeutil"
)

var (
	benchmarkParsed int64
	benchmarkText   string
)

// Keep the previous implementations here to compare the same inputs in one run.
func parseIntReflect[T typeutil.Integer](s string) T {
	var zero T
	if reflect.TypeOf(zero).Kind() >= reflect.Uint {
		value, _ := strconv.ParseUint(s, 10, 64)
		return T(value)
	}
	value, _ := strconv.ParseInt(s, 10, 64)
	return T(value)
}

func BenchmarkIntegerConversion(b *testing.B) {
	b.Run("parse_reflect", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			benchmarkParsed = parseIntReflect[int64]("123456789")
		}
	})
	b.Run("parse_generic", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			benchmarkParsed = ParseIntString[int64]("123456789")
		}
	})
	values := []int64{-9223372036854775808, 0, 42, 123456789, 9223372036854775807}
	b.Run("format_slice_legacy", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			var parts []string
			for _, value := range values {
				parts = append(parts, FormatIntString(value))
			}
			benchmarkText = strings.Join(parts, ",")
		}
	})
	b.Run("format_slice", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			benchmarkText = FormatIntSliceString(values)
		}
	})
}
