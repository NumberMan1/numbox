package sliceutil

import (
	"reflect"
	"testing"

	"github.com/NumberMan1/numbox/utils/convutil"
)

func TestSliceAndMarshalRegressions(t *testing.T) {
	if got := DeduplicateAddInt32([]int{1, 1}, 2, 2, 3); !reflect.DeepEqual(got, []int{1, 2, 3}) {
		t.Fatal(got)
	}
	if got := MaxValue([]int{-3, -1, -2}); got != -1 {
		t.Fatal(got)
	}
	if got := convutil.MarshalSliceToString(",", []string{"", "x", ""}, func(s string) string { return s }); got != ",x," {
		t.Fatal(got)
	}
}
