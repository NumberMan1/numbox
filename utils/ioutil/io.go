// Adapted from github.com/NumberMan1/orbit (81cbe7b).
// Copyright (c) 2025 NumberMan1. See LICENSE in the repository root.

package ioutil

import (
	"errors"
	"io"
)

// WriteFull retries short writes, but rejects a writer making no progress or
// reporting impossible byte counts. It never treats a partial write as success.
func WriteFull(writer io.Writer, data []byte) error {
	if writer == nil {
		return errors.New("util: writer required")
	}
	for len(data) > 0 {
		n, err := writer.Write(data)
		if n < 0 || n > len(data) {
			return errors.New("util: invalid write count")
		}
		if err != nil {
			return err
		}
		if n == 0 {
			return io.ErrNoProgress
		}
		data = data[n:]
	}
	return nil
}
