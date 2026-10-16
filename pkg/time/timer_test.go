// Adapted from github.com/NumberMan1/orbit (81cbe7b).
// Copyright (c) 2025 NumberMan1. See third_party/ORBIT_LICENSE.

package time

import (
	"context"
	"errors"
	"testing"
	"testing/synctest"
	"time"
)

func TestTimersCancelAndTick(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		ctx, cancel := context.WithCancel(t.Context())
		n := 0
		done := make(chan error, 1)
		go func() { done <- Every(ctx, time.Second, func(context.Context, time.Time) error { n++; return nil }) }()
		time.Sleep(3500 * time.Millisecond)
		synctest.Wait()
		if n != 3 {
			t.Fatal(n)
		}
		cancel()
		if err := <-done; !errors.Is(err, context.Canceled) {
			t.Fatal(err)
		}
		if err := After(ctx, time.Hour, func(context.Context) error { t.Error("canceled timer fired"); return nil }); !errors.Is(err, context.Canceled) {
			t.Fatal(err)
		}
	})
}
