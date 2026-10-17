package timewheel

import (
	"reflect"
	"testing"
)

var benchmarkKeyValid bool

func validKeyReflect[K comparable](key K) bool {
	value := reflect.ValueOf(key)
	return (!value.IsValid() || value.Comparable()) && key == key
}

func BenchmarkKeyValidation(b *testing.B) {
	b.Run("int_reflect", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			benchmarkKeyValid = validKeyReflect(i)
		}
	})
	b.Run("int_generic", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			benchmarkKeyValid = validKey(i)
		}
	})
	b.Run("interface_reflect", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			benchmarkKeyValid = validKeyReflect[any](i)
		}
	})
	b.Run("interface_generic", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			benchmarkKeyValid = validKey[any](i)
		}
	})
}
