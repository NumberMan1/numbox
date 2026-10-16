// Adapted from github.com/NumberMan1/orbit (81cbe7b).
// Copyright (c) 2025 NumberMan1. See third_party/ORBIT_LICENSE.

package timewheel

import (
	"container/list"
	"errors"
	"log/slog"
	"reflect"
	"runtime/debug"
	"time"
)

var (
	ErrInvalidKey     = errors.New("timewheel: key must be comparable and reflexive")
	ErrTimerExists    = errors.New("timewheel: timer key already exists")
	ErrClockBackwards = errors.New("timewheel: clock moved backwards")
)

type WheelOptions struct {
	Tick   time.Duration
	Slots  int
	Logger *slog.Logger // nil captures slog.Default()
}

type wheelTask[K comparable] struct {
	key      K
	at       time.Time
	interval time.Duration
	fn       func(time.Time)
	rounds   uint64
	slot     int
	element  *list.Element
}

// TimeWheel is a caller-driven, keyed timing wheel with no goroutines or runtime
// timers. All methods must run on one owner (for example an Orbit service.Queue).
// Callbacks may Schedule, Cancel or Reset, but must not recursively Advance.
// Tick precision rounds deadlines UP; callbacks never fire before their deadline.
type TimeWheel[K comparable] struct {
	start, observed time.Time
	tick            time.Duration
	current         uint64
	slots           []list.List
	tasks           map[K]*wheelTask[K]
	logger          *slog.Logger
	advancing       bool
}

func NewTimeWheel[K comparable](start time.Time, options WheelOptions) (*TimeWheel[K], error) {
	if options.Tick <= 0 || options.Slots < 2 {
		return nil, errors.New("timewheel: positive tick and at least two slots required")
	}
	if options.Logger == nil {
		options.Logger = slog.Default()
	}
	return &TimeWheel[K]{start: start, observed: start, tick: options.Tick, slots: make([]list.List, options.Slots), tasks: make(map[K]*wheelTask[K]), logger: options.Logger}, nil
}

// Schedule adds a unique key. interval=0 fires once; interval>0 repeats with a
// fixed delay from the time passed to Advance, coalescing missed repetitions.
// Use time.Now().Add(delay) in real time or your simulation clock in simulations.
func (w *TimeWheel[K]) Schedule(key K, at time.Time, interval time.Duration, fn func(time.Time)) error {
	if !validKey(key) {
		return ErrInvalidKey
	}
	if fn == nil || interval < 0 {
		return errors.New("timewheel: invalid timer callback or interval")
	}
	if _, ok := w.tasks[key]; ok {
		return ErrTimerExists
	}
	if err := w.checkDeadline(at); err != nil {
		return err
	}
	task := &wheelTask[K]{key: key, at: at, interval: interval, fn: fn}
	w.tasks[key] = task
	w.place(task)
	return nil
}

func (w *TimeWheel[K]) checkDeadline(at time.Time) error {
	// Keep tick rounding and duration multiplication within time.Duration.
	if at.After(w.start.Add(time.Duration(1<<63-1) - w.tick)) {
		return errors.New("timewheel: timer exceeds duration horizon")
	}
	return nil
}

func (w *TimeWheel[K]) place(task *wheelTask[K]) {
	var due uint64
	if elapsed := task.at.Sub(w.start); elapsed > 0 {
		due = uint64(elapsed / w.tick)
		if elapsed%w.tick != 0 {
			due++
		}
	}
	if due <= w.current {
		due = w.current + 1
	}
	task.rounds = (due - w.current - 1) / uint64(len(w.slots))
	task.slot = int(due % uint64(len(w.slots)))
	task.element = w.slots[task.slot].PushBack(task)
}

func (w *TimeWheel[K]) Len() int { return len(w.tasks) }

// Cancel returns true if this key was still pending. It cannot undo a callback
// that has begun; a repeating timer can cancel itself from its callback.
func (w *TimeWheel[K]) Cancel(key K) bool {
	if !validKey(key) {
		return false
	}
	task, ok := w.tasks[key]
	if !ok {
		return false
	}
	if task.element != nil {
		w.slots[task.slot].Remove(task.element)
	}
	delete(w.tasks, key)
	return true
}

// Reset preserves the callback and repetition interval. A reset also invalidates
// an expired callback waiting in the current batch. Missing keys return false.
func (w *TimeWheel[K]) Reset(key K, at time.Time) (bool, error) {
	if !validKey(key) {
		return false, ErrInvalidKey
	}
	task, ok := w.tasks[key]
	if !ok {
		return false, nil
	}
	if err := w.checkDeadline(at); err != nil {
		return false, err
	}
	w.Cancel(key)
	return true, w.Schedule(key, at, task.interval, task.fn)
}

// Clear cancels pending callbacks, including the remainder of an expiry batch.
func (w *TimeWheel[K]) Clear() {
	clear(w.tasks)
	for i := range w.slots {
		w.slots[i].Init()
	}
}

// Advance processes at most maxTicks tick buckets, returning whether it caught
// up to now. Call again with the same/later now after a long stall. Work within a
// bucket is proportional to its population; maxTicks is NOT a callback budget.
// Equal-deadline tasks follow insertion order. Panics are logged and isolated.
func (w *TimeWheel[K]) Advance(now time.Time, maxTicks int) (bool, error) {
	if w.advancing {
		return false, errors.New("timewheel: recursive wheel advance")
	}
	if maxTicks <= 0 {
		return false, errors.New("timewheel: positive tick budget required")
	}
	if now.Before(w.observed) {
		return false, ErrClockBackwards
	}
	if err := w.checkDeadline(now); err != nil {
		return false, err
	}
	w.observed = now
	w.advancing = true
	defer func() { w.advancing = false }()
	target := uint64(now.Sub(w.start) / w.tick)
	if len(w.tasks) == 0 {
		w.current = target
		return true, nil
	}
	steps := min(target-w.current, uint64(maxTicks))
	for range steps {
		w.current++
		bucket := &w.slots[w.current%uint64(len(w.slots))]
		var ready []*wheelTask[K]
		for e := bucket.Front(); e != nil; {
			next := e.Next()
			task := e.Value.(*wheelTask[K])
			if task.rounds > 0 {
				task.rounds--
			} else {
				bucket.Remove(e)
				task.element = nil
				ready = append(ready, task)
			}
			e = next
		}
		for _, task := range ready {
			if w.tasks[task.key] != task {
				continue
			}
			delete(w.tasks, task.key)
			if task.interval > 0 {
				task.at = now.Add(task.interval)
				if err := w.checkDeadline(task.at); err == nil {
					w.tasks[task.key] = task
					w.place(task)
				} else {
					w.logger.Error("timer repeat stopped", "error", err)
				}
			}
			w.call(task, now)
		}
	}
	return w.current == target, nil
}

func (w *TimeWheel[K]) call(task *wheelTask[K], now time.Time) {
	defer func() {
		if p := recover(); p != nil {
			w.logger.Error("timer callback panic", "panic", p, "stack", string(debug.Stack()))
		}
	}()
	task.fn(now)
}

// Interface keys can hold non-comparable values despite K's comparable constraint.
// NaN-containing keys are rejected because they cannot be looked up or canceled.
func validKey[K comparable](key K) bool {
	value := reflect.ValueOf(key)
	return (!value.IsValid() || value.Comparable()) && key == key
}
