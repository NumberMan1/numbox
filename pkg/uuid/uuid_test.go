package uuid

import "testing"

func TestUUID(t *testing.T) {
	repeated := make(map[string]struct{})
	for i := 0; i < 1000000; i++ {
		id := NewUUID()
		if _, ok := repeated[id]; ok {
			panic(i)
		}
		repeated[id] = struct{}{}
	}
}

func BenchmarkNewUUID(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = NewUUID()
	}
}
