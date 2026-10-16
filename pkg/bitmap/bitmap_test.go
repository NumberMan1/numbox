package bitmap

import "testing"

func TestLegacyBitmapSerialization(t *testing.T) {
	b := NewBitMap(9)
	b.Add(9)
	data := b.Bits()
	data[1] = 0
	if !b.IsExist(9) {
		t.Fatal("Bits aliases state")
	}
	input := []byte{1, 0}
	b.SetBits(input)
	input[0] = 0
	if !b.IsExist(0) || b.Count() != 1 {
		t.Fatal("SetBits aliases input")
	}
	for _, bad := range [][]byte{nil, {0}, {0, 4}, {0, 0, 0}} {
		if b.Load(bad) == nil || !b.IsExist(0) || b.Count() != 1 {
			t.Fatal("invalid input changed state", bad)
		}
	}
	if err := NewBitMap(7).Load([]byte{255}); err != nil {
		t.Fatal(err)
	}
}
