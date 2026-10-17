package convutil

import (
	"math"
	"testing"
)

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
