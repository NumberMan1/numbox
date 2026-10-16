// Adapted from github.com/NumberMan1/orbit (81cbe7b).
// Copyright (c) 2025 NumberMan1. See third_party/ORBIT_LICENSE.

package time

import (
	"testing"
	"time"
)

func TestPollTimerPauseResumeAndRepeat(t *testing.T) {
	var timer PollTimer
	start := time.Unix(0, 0)
	if err := timer.Reset(start, 10*time.Second, 0); err != nil {
		t.Fatal(err)
	}
	if !timer.Pause(start.Add(3*time.Second)) || timer.Poll(start.Add(time.Hour)) {
		t.Fatal("pause failed")
	}
	if !timer.Resume(start.Add(time.Hour)) || timer.Poll(start.Add(time.Hour+6*time.Second)) {
		t.Fatal("remaining duration lost")
	}
	if !timer.Poll(start.Add(time.Hour+7*time.Second)) || timer.Poll(start.Add(2*time.Hour)) {
		t.Fatal("once timer repeated")
	}
	timer.Reset(start, 0, time.Second)
	if !timer.Poll(start) || !timer.Poll(start.Add(time.Hour)) || timer.Poll(start.Add(time.Hour)) {
		t.Fatal("periodic coalescing failed")
	}
	timer.Cancel()
	if timer.Active() || timer.Poll(start.Add(2*time.Hour)) {
		t.Fatal("cancel failed")
	}
}
