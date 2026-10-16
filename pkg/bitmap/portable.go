// Adapted from github.com/NumberMan1/orbit (81cbe7b).
// Copyright (c) 2025 NumberMan1. See LICENSE in the repository root.

package bitmap

import (
	"bytes"
	"errors"
	"math/bits"
)

// Bitmap contains indices [0,size). It is owner-managed, not concurrent-safe.
// Its byte representation is little-bit-order within each byte.
type Bitmap struct {
	data        []byte
	size, count uint32
}

func NewBitmap(size uint32) *Bitmap {
	return &Bitmap{data: make([]byte, (uint64(size)+7)/8), size: size}
}
func (b *Bitmap) Size() uint32          { return b.size }
func (b *Bitmap) Len() uint32           { return b.count }
func (b *Bitmap) Has(index uint32) bool { return index < b.size && b.data[index/8]&(1<<(index%8)) != 0 }
func (b *Bitmap) Add(index uint32) bool {
	if index >= b.size || b.Has(index) {
		return false
	}
	b.data[index/8] |= 1 << (index % 8)
	b.count++
	return true
}
func (b *Bitmap) Remove(index uint32) bool {
	if !b.Has(index) {
		return false
	}
	b.data[index/8] &^= 1 << (index % 8)
	b.count--
	return true
}
func (b *Bitmap) Clear()        { clear(b.data); b.count = 0 }
func (b *Bitmap) Bytes() []byte { return bytes.Clone(b.data) }

// Load validates length and unused bits before replacing the state; it copies data.
func (b *Bitmap) Load(data []byte) error {
	if len(data) != len(b.data) {
		return errors.New("util: bitmap size mismatch")
	}
	if tail := b.size % 8; tail != 0 && data[len(data)-1]>>tail != 0 {
		return errors.New("util: nonzero bitmap padding")
	}
	var count uint32
	for _, v := range data {
		count += uint32(bits.OnesCount8(v))
	}
	b.data = bytes.Clone(data)
	b.count = count
	return nil
}
