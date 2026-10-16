// Adapted from github.com/NumberMan1/orbit (81cbe7b).
// Copyright (c) 2025 NumberMan1. See LICENSE in the repository root.

package time

import (
	"errors"
	"time"
)

// PollTimer is an owner-driven timer for a simulation or update loop. It has no
// goroutine. Pause and Resume preserve remaining time; late periodic polls
// coalesce missed intervals. The zero value is inactive. Not concurrent-safe.
type PollTimer struct {
	deadline            time.Time
	interval, remaining time.Duration
	active, paused      bool
}

func (t *PollTimer) Reset(now time.Time, delay, interval time.Duration) error {
	if delay < 0 || interval < 0 {
		return errors.New("util: negative timer duration")
	}
	*t = PollTimer{deadline: now.Add(delay), interval: interval, active: true}
	return nil
}
func (t *PollTimer) Cancel()      { *t = PollTimer{} }
func (t *PollTimer) Active() bool { return t.active }
func (t *PollTimer) Paused() bool { return t.paused }
func (t *PollTimer) Pause(now time.Time) bool {
	if !t.active || t.paused {
		return false
	}
	t.remaining = max(0, t.deadline.Sub(now))
	t.paused = true
	return true
}
func (t *PollTimer) Resume(now time.Time) bool {
	if !t.active || !t.paused {
		return false
	}
	t.deadline = now.Add(t.remaining)
	t.paused = false
	return true
}
func (t *PollTimer) Poll(now time.Time) bool {
	if !t.active || t.paused || now.Before(t.deadline) {
		return false
	}
	if t.interval == 0 {
		t.active = false
	} else {
		t.deadline = now.Add(t.interval)
	}
	return true
}
