// Adapted from github.com/NumberMan1/orbit (81cbe7b).
// Copyright (c) 2025 NumberMan1. See LICENSE in the repository root.

package ioutil

import (
	"bytes"
	"errors"
	"io"
	"testing"
)

type writerFunc func([]byte) (int, error)

func (f writerFunc) Write(p []byte) (int, error) { return f(p) }

func TestWriteFullShortWritesAndNoProgress(t *testing.T) {
	var output bytes.Buffer
	if err := WriteFull(writerFunc(func(p []byte) (int, error) { return output.Write(p[:min(2, len(p))]) }), []byte("complete")); err != nil || output.String() != "complete" {
		t.Fatal(output.String(), err)
	}
	if err := WriteFull(writerFunc(func([]byte) (int, error) { return 0, nil }), []byte("x")); !errors.Is(err, io.ErrNoProgress) {
		t.Fatal(err)
	}
	failure := errors.New("failed")
	if err := WriteFull(writerFunc(func([]byte) (int, error) { return 1, failure }), []byte("xx")); !errors.Is(err, failure) {
		t.Fatal(err)
	}
}
