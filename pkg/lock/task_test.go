package lock

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"
)

func taskCases() []struct {
	name string
	lock Locker
	run  func(context.Context, func(context.Context) error, ...LockOption) error
} {
	m := NewLocker()
	w := NewRWLocker()
	r := NewRWLocker()
	return []struct {
		name string
		lock Locker
		run  func(context.Context, func(context.Context) error, ...LockOption) error
	}{{"mutex", m, m.Do}, {"rw-writer", w, w.Do}, {"rw-reader", r, r.DoRead}}
}

func TestTaskCancellationWaitsForExit(t *testing.T) {
	for _, tc := range taskCases() {
		t.Run(tc.name, func(t *testing.T) {
			canceled := make(chan struct{})
			exit := make(chan struct{})
			defer func() {
				select {
				case <-exit:
				default:
					close(exit)
				}
			}()
			result := make(chan error, 1)
			go func() {
				result <- tc.run(context.Background(), func(ctx context.Context) error {
					<-ctx.Done()
					close(canceled)
					<-exit // Simulate task cleanup that has not finished yet.
					return nil
				}, WithTaskTimeout(5*time.Millisecond))
			}()
			select {
			case <-canceled:
			case <-time.After(time.Second):
				t.Fatal("task never observed its timeout")
			}
			select {
			case err := <-result:
				t.Fatalf("Do returned before the task exited: %v", err)
			default:
			}
			if token, err := tc.lock.LockContext(context.Background(), WithAcquireTimeout(5*time.Millisecond)); err != ErrLockTimeout {
				tc.lock.Unlock(token)
				t.Fatalf("another task entered before cleanup: %v", err)
			}
			close(exit)
			select {
			case err := <-result:
				if !errors.Is(err, context.DeadlineExceeded) {
					t.Fatalf("task timeout = %v", err)
				}
			case <-time.After(time.Second):
				t.Fatal("Do did not return after the task exited")
			}
			token := tc.lock.Lock()
			tc.lock.Unlock(token)
		})
	}
}

func TestTaskFailureAndPanicReleaseLock(t *testing.T) {
	for _, tc := range taskCases() {
		t.Run(tc.name, func(t *testing.T) {
			want := errors.New("task failed")
			ctx, cancel := context.WithCancel(context.Background())
			err := tc.run(ctx, func(context.Context) error {
				cancel()
				return want
			})
			if !errors.Is(err, want) {
				t.Fatalf("task error was replaced: %v", err)
			}
			func() {
				defer func() {
					if recover() != "task panic" {
						t.Error("task panic was not propagated")
					}
				}()
				tc.run(context.Background(), func(context.Context) error { panic("task panic") })
			}()
			token, err := tc.lock.LockContext(context.Background(), WithAcquireTimeout(time.Second))
			if err != nil {
				t.Fatalf("task panic leaked the lock: %v", err)
			}
			tc.lock.Unlock(token)
			if err := tc.run(context.Background(), nil); err != ErrNilTask {
				t.Fatalf("nil task = %v", err)
			}
		})
	}
}

func TestTaskIsNotStartedAfterCanceledAcquisition(t *testing.T) {
	for _, tc := range taskCases() {
		t.Run(tc.name, func(t *testing.T) {
			ctx, cancel := context.WithCancel(context.Background())
			cancel()
			called := false
			err := tc.run(ctx, func(context.Context) error { called = true; return nil })
			if !errors.Is(err, context.Canceled) || called {
				t.Fatalf("canceled task started: called=%v err=%v", called, err)
			}
		})
	}
}

func TestTaskContextIsIndependentOfAcquisition(t *testing.T) {
	for _, tc := range taskCases() {
		t.Run(tc.name, func(t *testing.T) {
			// The acquisition context is canceled after acquiring the token.
			// The task must instead get a fresh execution context and deadline.
			calls := 0
			once := LockOption(func(c *lockConfig) {
				calls++
				c.acquireTimeout = time.Second
			})
			err := tc.run(context.Background(), func(ctx context.Context) error {
				if ctx.Err() != nil {
					t.Fatalf("acquisition context leaked into the task: %v", ctx.Err())
				}
				deadline, ok := ctx.Deadline()
				if !ok || time.Until(deadline) < 30*time.Minute {
					t.Fatal("task did not receive a fresh execution deadline")
				}
				return nil
			}, once, WithTaskTimeout(time.Hour))
			if err != nil || calls != 1 {
				t.Fatalf("task result=%v, options evaluated %d times", err, calls)
			}
		})
	}
}

func TestTaskCooperatesWithParentCancellation(t *testing.T) {
	for _, tc := range taskCases() {
		t.Run(tc.name, func(t *testing.T) {
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			called := false
			err := tc.run(ctx, func(taskCtx context.Context) error {
				called = true
				cancel()
				<-taskCtx.Done()
				// A canceled context alone cannot release this task's token.
				token, err := tc.lock.LockContext(context.Background(), WithAcquireTimeout(time.Millisecond))
				if err != ErrLockTimeout {
					tc.lock.Unlock(token)
					t.Errorf("parent cancellation released a running task: %v", err)
				}
				return taskCtx.Err()
			})
			if !called || !errors.Is(err, context.Canceled) {
				t.Fatalf("parent cancellation result: called=%v err=%v", called, err)
			}
			token := tc.lock.Lock()
			tc.lock.Unlock(token)
		})
	}
}

func TestDoReadSharesWithOtherReaders(t *testing.T) {
	l := NewRWLocker()
	err := l.DoRead(context.Background(), func(context.Context) error {
		token, err := l.RLockContext(context.Background(), WithAcquireTimeout(time.Second))
		if err != nil {
			return err
		}
		defer l.RUnlock(token)
		writer, err := l.LockContext(context.Background(), WithAcquireTimeout(time.Millisecond))
		if err != ErrLockTimeout {
			l.Unlock(writer)
			t.Errorf("writer entered during DoRead: %v", err)
		}
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
	token := l.Lock()
	l.Unlock(token)
}

func TestLegacyHoldTimeoutOnlyCancelsManagedTasks(t *testing.T) {
	l := NewLocker()
	err := l.Do(context.Background(), func(ctx context.Context) error {
		<-ctx.Done()
		return ctx.Err()
	}, WithLockHoldTimeout(time.Millisecond))
	if !errors.Is(err, context.DeadlineExceeded) {
		t.Fatalf("legacy option did not map to task timeout: %v", err)
	}
}

func ExampleLocker_Do() {
	l := NewLocker()
	err := l.Do(context.Background(), func(ctx context.Context) error {
		fmt.Println("done")
		return nil
	}, WithTaskTimeout(time.Second))
	fmt.Println(err)
	// Output:
	// done
	// <nil>
}
