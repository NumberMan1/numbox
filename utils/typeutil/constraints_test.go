package typeutil_test

import (
	"testing"

	"github.com/NumberMan1/numbox/utils/maputil"
	"github.com/NumberMan1/numbox/utils/sliceutil"
)

func TestGeneralizedConstraints(t *testing.T) {
	type key struct {
		Name string
		ID   int
	}
	values := []key{{"a", 1}, {"b", 2}, {"a", 3}}
	byName := maputil.SliceToMapVals(values, func(v key) string { return v.Name })
	if len(byName) != 2 || byName["a"].ID != 3 {
		t.Fatal(byName)
	}
	byKey := maputil.SliceToMapVals(values, func(v key) key { return v })
	if len(byKey) != 3 {
		t.Fatal(byKey)
	}
	type count uint16
	counts := maputil.SliceToCountMap([]string{"a", "b", "a"}, func(string) count { return 2 })
	if counts["a"] != 4 || counts["b"] != 2 {
		t.Fatal(counts)
	}
	type score float64
	if got := sliceutil.MaxValue([]score{-2, 3, 1}); got != 3 {
		t.Fatal(got)
	}
	if got := sliceutil.MaxValue([]string{"a", "z", "b"}); got != "z" {
		t.Fatal(got)
	}
}
