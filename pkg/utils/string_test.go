package utils

import (
	"fmt"
	"testing"
)

func TestFormatIntSliceString(t *testing.T) {
	src := FormatIntSliceString([]int32{1})
	fmt.Println(src)

	target := ParseIntSliceString[int32](src)
	fmt.Println(target)

	emptyTarget := ParseIntSliceString[int32]("")
	fmt.Println(emptyTarget)
}
