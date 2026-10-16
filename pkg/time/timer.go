// Adapted from github.com/NumberMan1/orbit (81cbe7b).
// Copyright (c) 2025 NumberMan1. See LICENSE in the repository root.

package time

import (
	"context"
	"errors"
	"time"
)

// Every runs callbacks serially until cancellation. Slow callbacks coalesce ticks;
// they never overlap. The caller owns the goroutine and any panic policy.
func Every(ctx context.Context, interval time.Duration, fn func(context.Context, time.Time) error) error {
	if interval <= 0 || fn == nil {
		return errors.New("util: invalid periodic callback")
	}
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case now := <-ticker.C:
			if err := ctx.Err(); err != nil {
				return err
			}
			if err := fn(ctx, now); err != nil {
				return err
			}
		}
	}
}

// After waits once. Cancellation never waits for the delay to expire.
func After(ctx context.Context, delay time.Duration, fn func(context.Context) error) error {
	if delay < 0 || fn == nil {
		return errors.New("util: invalid delayed callback")
	}
	timer := time.NewTimer(delay)
	defer timer.Stop()
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-timer.C:
		if err := ctx.Err(); err != nil {
			return err
		}
		return fn(ctx)
	}
}
