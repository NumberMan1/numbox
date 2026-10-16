// Adapted from github.com/NumberMan1/orbit (81cbe7b).
// Copyright (c) 2025 NumberMan1. See third_party/ORBIT_LICENSE.

package bitmap

import "testing"

func TestBitmapCountBoundsAndOwnership(t *testing.T) {
	b := NewBitmap(10)
	if !b.Add(9) || b.Add(9) || b.Add(10) || b.Len() != 1 {
		t.Fatal("bad addition")
	}
	if b.Remove(8) || !b.Remove(9) || b.Remove(9) || b.Len() != 0 {
		t.Fatal("bad removal")
	}
	if err := b.Load([]byte{255, 3}); err != nil || b.Len() != 10 {
		t.Fatal(err, b.Len())
	}
	copy := b.Bytes()
	copy[0] = 0
	if !b.Has(0) {
		t.Fatal("aliased bytes")
	}
	if err := b.Load([]byte{0, 4}); err == nil || b.Len() != 10 {
		t.Fatal("padding validation not atomic")
	}
	input := []byte{1, 0}
	b.Load(input)
	input[0] = 0
	if !b.Has(0) || b.Len() != 1 {
		t.Fatal("aliased input")
	}
	b.Clear()
	if b.Len() != 0 {
		t.Fatal(b.Len())
	}
	if err := NewBitmap(0).Load(nil); err != nil {
		t.Fatal(err)
	}
}
