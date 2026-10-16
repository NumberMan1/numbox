// Adapted from github.com/NumberMan1/orbit (81cbe7b).
// Copyright (c) 2025 NumberMan1. See LICENSE in the repository root.

package timewheel

import (
	"errors"
	"math"
	"reflect"
	"testing"
	"time"
)

func wheel(t testing.TB) (*TimeWheel[int], time.Time) {
	t.Helper()
	start := time.Unix(0, 0)
	w, err := NewTimeWheel[int](start, WheelOptions{Tick: 10 * time.Millisecond, Slots: 4})
	if err != nil {
		t.Fatal(err)
	}
	return w, start
}
func TestWheelBoundariesRoundsAndNoEarlyFire(t *testing.T) {
	w, start := wheel(t)
	var fired []int
	for id, delay := range []time.Duration{0, time.Millisecond, 10 * time.Millisecond, 11 * time.Millisecond, 40 * time.Millisecond, 41 * time.Millisecond, 80 * time.Millisecond} {
		if err := w.Schedule(id, start.Add(delay), 0, func(now time.Time) {
			if now.Before(start.Add(delay)) {
				t.Error("early callback")
			}
			fired = append(fired, id)
		}); err != nil {
			t.Fatal(err)
		}
	}
	if _, err := w.Advance(start.Add(9*time.Millisecond), 10); err != nil {
		t.Fatal(err)
	}
	if len(fired) != 0 {
		t.Fatal(fired)
	}
	for _, step := range []struct {
		at   time.Duration
		want []int
	}{
		{10 * time.Millisecond, []int{0, 1, 2}}, {20 * time.Millisecond, []int{0, 1, 2, 3}},
		{40 * time.Millisecond, []int{0, 1, 2, 3, 4}}, {50 * time.Millisecond, []int{0, 1, 2, 3, 4, 5}},
		{80 * time.Millisecond, []int{0, 1, 2, 3, 4, 5, 6}},
	} {
		if _, err := w.Advance(start.Add(step.at), 10); err != nil {
			t.Fatal(err)
		}
		if !reflect.DeepEqual(fired, step.want) {
			t.Fatalf("%v: %v", step.at, fired)
		}
	}
	if w.Len() != 0 {
		t.Fatal(w.Len())
	}
}
func TestWheelCallbackMutationsAndPanic(t *testing.T) {
	w, start := wheel(t)
	var fired []int
	add := func(id int, fn func(time.Time)) {
		t.Helper()
		if err := w.Schedule(id, start.Add(10*time.Millisecond), 0, fn); err != nil {
			t.Fatal(err)
		}
	}
	add(1, func(now time.Time) {
		w.Cancel(2)
		if ok, err := w.Reset(3, now.Add(20*time.Millisecond)); !ok || err != nil {
			t.Fatal(ok, err)
		}
		panic("isolated")
	})
	add(2, func(time.Time) { t.Error("canceled task ran") })
	add(3, func(time.Time) { fired = append(fired, 3) })
	add(4, func(time.Time) {
		fired = append(fired, 4)
		if _, err := w.Advance(start, 1); err == nil {
			t.Error("recursive advance allowed")
		}
	})
	w.Advance(start.Add(10*time.Millisecond), 10)
	if !reflect.DeepEqual(fired, []int{4}) {
		t.Fatal(fired)
	}
	w.Advance(start.Add(30*time.Millisecond), 10)
	if !reflect.DeepEqual(fired, []int{4, 3}) {
		t.Fatal(fired)
	}
}
func TestWheelCatchupCoalescesAndClearCancelsReadyBatch(t *testing.T) {
	w, start := wheel(t)
	calls := 0
	if err := w.Schedule(1, start.Add(time.Millisecond), 10*time.Millisecond, func(time.Time) { calls++ }); err != nil {
		t.Fatal(err)
	}
	now := start.Add(time.Second)
	for steps := 0; ; steps++ {
		caught, err := w.Advance(now, 3)
		if err != nil {
			t.Fatal(err)
		}
		if caught {
			break
		}
		if steps > 40 {
			t.Fatal("catchup stalled")
		}
	}
	if calls != 1 {
		t.Fatal("replayed missed periodic callbacks", calls)
	}
	w.Advance(now.Add(10*time.Millisecond), 2)
	if calls != 2 {
		t.Fatal(calls)
	}
	w.Clear()
	w.Schedule(2, now.Add(20*time.Millisecond), 0, func(time.Time) { w.Clear() })
	w.Schedule(3, now.Add(20*time.Millisecond), 0, func(time.Time) { t.Error("Clear failed to suppress ready task") })
	w.Advance(now.Add(20*time.Millisecond), 2)
	if w.Len() != 0 {
		t.Fatal(w.Len())
	}
	if _, err := w.Advance(start, 1); !errors.Is(err, ErrClockBackwards) {
		t.Fatal(err)
	}
}
func TestWheelValidationAndLargePopulation(t *testing.T) {
	if _, err := NewTimeWheel[int](time.Time{}, WheelOptions{Tick: time.Millisecond, Slots: 0}); err == nil {
		t.Fatal("invalid slots")
	}
	w, start := wheel(t)
	calls := 0
	for i := range 10000 {
		if err := w.Schedule(i, start.Add(time.Hour), 0, func(time.Time) { calls++ }); err != nil {
			t.Fatal(err)
		}
	}
	if err := w.Schedule(1, start, 0, func(time.Time) {}); !errors.Is(err, ErrTimerExists) {
		t.Fatal(err)
	}
	for i := range 10000 {
		if !w.Cancel(i) || w.Cancel(i) {
			t.Fatal("cancel identity", i)
		}
	}
	if w.Len() != 0 || calls != 0 {
		t.Fatal(w.Len(), calls)
	}
}
func BenchmarkWheelScheduleCancel(b *testing.B) {
	w, start := wheel(b)
	fn := func(time.Time) {}
	for i := range 10000 {
		w.Schedule(i, start.Add(time.Hour), 0, fn)
	}
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		w.Cancel(i % 10000)
		w.Schedule(i%10000, start.Add(time.Hour), 0, fn)
	}
}
func BenchmarkRuntimeScheduleCancel(b *testing.B) {
	timers := make([]*time.Timer, 10000)
	fn := func() {}
	for i := range timers {
		timers[i] = time.AfterFunc(time.Hour, fn)
	}
	b.Cleanup(func() {
		for _, timer := range timers {
			timer.Stop()
		}
	})
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		timers[i%10000].Stop()
		timers[i%10000] = time.AfterFunc(time.Hour, fn)
	}
}

func TestWheelDynamicKeysAndValidation(t *testing.T) {
	w, err := NewTimeWheel[any](time.Unix(0, 0), WheelOptions{Tick: time.Millisecond, Slots: 4})
	if err != nil {
		t.Fatal(err)
	}
	for _, key := range []any{[]int{1}, map[string]int{}, struct{ V any }{[]int{1}}, math.NaN()} {
		if err := w.Schedule(key, time.Unix(0, 0), 0, func(time.Time) {}); !errors.Is(err, ErrInvalidKey) {
			t.Fatal(key, err)
		}
		if w.Cancel(key) {
			t.Fatal("invalid key canceled")
		}
		if _, err := w.Reset(key, time.Unix(0, 0)); !errors.Is(err, ErrInvalidKey) {
			t.Fatal(err)
		}
	}
	if err := w.Schedule(nil, time.Unix(0, 0), 0, func(time.Time) {}); err != nil || !w.Cancel(nil) {
		t.Fatal("nil generic key should work", err)
	}
	if _, err := w.Advance(time.Unix(0, 0), 0); err == nil {
		t.Fatal("zero budget accepted")
	}
}
func TestWheelIdleAdvanceAndRepeatSelfCancel(t *testing.T) {
	w, start := wheel(t)
	later := start.Add(24 * time.Hour)
	if caught, err := w.Advance(later, 1); !caught || err != nil {
		t.Fatal(caught, err)
	}
	calls := 0
	if err := w.Schedule(1, later.Add(time.Millisecond), time.Millisecond, func(time.Time) { calls++; w.Cancel(1) }); err != nil {
		t.Fatal(err)
	}
	if _, err := w.Advance(later.Add(100*time.Millisecond), 100); err != nil {
		t.Fatal(err)
	}
	if calls != 1 || w.Len() != 0 {
		t.Fatal(calls, w.Len())
	}
}
