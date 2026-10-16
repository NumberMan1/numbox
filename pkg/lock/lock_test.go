package lock

import (
	"context"
	"errors"
	"runtime"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

var lockerCases = []struct {
	name string
	new  func() Locker
}{
	{"mutex", NewLocker},
	{"rw-writer", func() Locker { return NewRWLocker() }},
}

func TestLockerDoesNotExpire(t *testing.T) {
	for _, tc := range lockerCases {
		t.Run(tc.name, func(t *testing.T) {
			l := tc.new()
			if l.Acquired(0) || l.Unlock(0) {
				t.Fatal("zero token must be invalid")
			}
			token := l.Lock(WithLockHoldTimeout(time.Nanosecond))
			defer l.Unlock(token)
			_, err := l.LockContext(context.Background(), WithAcquireTimeout(10*time.Millisecond))
			if !errors.Is(err, ErrLockTimeout) || !l.Acquired(token) {
				t.Fatalf("held lock expired or was stolen: acquired=%v err=%v", l.Acquired(token), err)
			}
			if l.Unlock(token + 1) {
				t.Fatal("wrong token released the lock")
			}
			calls := 0
			cb := WithUnlockCallback(func() { calls++ })
			if !l.Unlock(token, cb) || l.Acquired(token) || l.Unlock(token, cb) || calls != 1 {
				t.Fatal("unlock must consume the token and invoke its callback exactly once")
			}
			next := l.Lock()
			defer l.Unlock(next)
			if next == token || l.Unlock(token) {
				t.Fatal("stale token affected the new owner")
			}
		})
	}
}

func TestLockContextCancellation(t *testing.T) {
	for _, tc := range lockerCases {
		t.Run(tc.name, func(t *testing.T) {
			l := tc.new()
			ctx, cancel := context.WithCancel(context.Background())
			cancel()
			called := false
			token, err := l.LockContext(ctx, WithLockCallback(func() { called = true }))
			if token != 0 || !errors.Is(err, context.Canceled) || called {
				t.Fatalf("canceled acquisition succeeded: token=%v err=%v callback=%v", token, err, called)
			}

			owner := l.Lock()
			defer l.Unlock(owner)
			ctx, cancel = context.WithTimeout(context.Background(), 10*time.Millisecond)
			defer cancel()
			token, err = l.LockContext(ctx, WithAcquireTimeout(0))
			if token != 0 || !errors.Is(err, context.DeadlineExceeded) || !l.Acquired(owner) {
				t.Fatalf("wait cancellation changed ownership: token=%v err=%v", token, err)
			}
			// Cancellation of a successfully acquired token must not release it.
			ctx, cancel = context.WithCancel(context.Background())
			l.Unlock(owner)
			owner, err = l.LockContext(ctx)
			if err != nil {
				t.Fatal(err)
			}
			cancel()
			if !l.Acquired(owner) {
				t.Fatal("context cancellation released a manually managed token")
			}
			l.Unlock(owner)
		})
	}
}

func TestLockStillPanicsOnAcquireTimeout(t *testing.T) {
	for _, tc := range lockerCases {
		t.Run(tc.name, func(t *testing.T) {
			l := tc.new()
			token := l.Lock()
			defer l.Unlock(token)
			func() {
				defer func() {
					if got := recover(); got != ErrLockTimeout {
						t.Errorf("panic = %v, want ErrLockTimeout", got)
					}
				}()
				l.Lock(WithAcquireTimeout(time.Millisecond))
			}()
			if !l.Acquired(token) {
				t.Fatal("wait timeout released the owner")
			}
		})
	}
}

func TestRWReadersDoNotExpire(t *testing.T) {
	l := NewRWLocker()
	first := l.RLock(WithLockHoldTimeout(time.Nanosecond))
	defer l.RUnlock(first)
	second := l.RLock()
	defer l.RUnlock(second)
	if first == second || !l.Acquired(first) || !l.Acquired(second) {
		t.Fatal("readers must have independent valid tokens")
	}
	_, err := l.LockContext(context.Background(), WithAcquireTimeout(10*time.Millisecond))
	if !errors.Is(err, ErrLockTimeout) || !l.Acquired(first) {
		t.Fatalf("writer bypassed a reader: %v", err)
	}
	if l.Unlock(first) {
		t.Fatal("read token released a write lock")
	}
	if !l.RUnlock(first) || l.RUnlock(first) || !l.Acquired(second) {
		t.Fatal("releasing one reader affected another")
	}
	l.RUnlock(second)
	writer := l.Lock()
	defer l.Unlock(writer)
	if l.RUnlock(writer) {
		t.Fatal("write token released as a reader")
	}
	_, err = l.RLockContext(context.Background(), WithAcquireTimeout(time.Millisecond))
	if !errors.Is(err, ErrLockTimeout) {
		t.Fatalf("reader bypassed a writer: %v", err)
	}
}

func waitForWriter(t *testing.T, l *rwLocker) {
	t.Helper()
	timeout := time.After(time.Second)
	for atomic.LoadInt32(&l.writeIntention) == 0 {
		select {
		case <-timeout:
			t.Fatal("writer did not begin waiting")
		default:
			runtime.Gosched()
		}
	}
}

func TestRWWriterPriorityAndCancellation(t *testing.T) {
	l := NewRWLocker().(*rwLocker)
	reader := l.RLock()
	defer l.RUnlock(reader)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	writer := make(chan error, 1)
	go func() {
		token, err := l.LockContext(ctx, WithAcquireTimeout(0))
		if err == nil {
			l.Unlock(token)
		}
		writer <- err
	}()
	waitForWriter(t, l)
	if _, err := l.RLockContext(context.Background(), WithAcquireTimeout(time.Millisecond)); err != ErrLockTimeout {
		t.Fatalf("new reader bypassed a waiting writer: %v", err)
	}
	cancel()
	select {
	case err := <-writer:
		if !errors.Is(err, context.Canceled) {
			t.Fatalf("canceled writer returned %v", err)
		}
	case <-time.After(time.Second):
		t.Fatal("writer did not observe cancellation")
	}
	second, err := l.RLockContext(context.Background(), WithAcquireTimeout(time.Second))
	if err != nil {
		t.Fatalf("canceled writer retained write intention: %v", err)
	}
	l.RUnlock(second)
	l.RUnlock(reader)
	token := l.Lock()
	l.Unlock(token)
}

func TestConcurrentCriticalSections(t *testing.T) {
	for _, tc := range lockerCases {
		t.Run(tc.name, func(t *testing.T) {
			l := tc.new()
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			var wg sync.WaitGroup
			count := 0
			for range 8 {
				wg.Add(1)
				go func() {
					defer wg.Done()
					for range 100 {
						token, err := l.LockContext(ctx, WithAcquireTimeout(0))
						if err != nil {
							t.Error(err)
							return
						}
						count++
						l.Unlock(token)
					}
				}()
			}
			wg.Wait()
			if count != 800 {
				t.Fatalf("count = %d, want 800", count)
			}
		})
	}
}

func TestAcquisitionCallbackPanicReleasesToken(t *testing.T) {
	type acquireFunc func(...LockOption) Token
	for _, tc := range []struct {
		name string
		new  func() (Locker, acquireFunc)
	}{
		{"mutex", func() (Locker, acquireFunc) { l := NewLocker(); return l, l.Lock }},
		{"rw-writer", func() (Locker, acquireFunc) { l := NewRWLocker(); return l, l.Lock }},
		{"rw-reader", func() (Locker, acquireFunc) { l := NewRWLocker(); return l, l.RLock }},
	} {
		t.Run(tc.name, func(t *testing.T) {
			l, acquire := tc.new()
			func() {
				defer func() {
					if recover() != "callback" {
						t.Error("callback panic was not propagated")
					}
				}()
				acquire(WithLockCallback(func() { panic("callback") }))
			}()
			token, err := l.LockContext(context.Background(), WithAcquireTimeout(time.Second))
			if err != nil {
				t.Fatalf("panic leaked a token: %v", err)
			}
			l.Unlock(token)
		})
	}
}

func TestTimeoutMetrics(t *testing.T) {
	before := GetMetrics()
	l := NewLocker()
	token := l.Lock()
	l.LockContext(context.Background(), WithAcquireTimeout(time.Millisecond))
	l.Unlock(token)
	rw := NewRWLocker()
	token = rw.RLock()
	rw.LockContext(context.Background(), WithAcquireTimeout(time.Millisecond))
	rw.RUnlock(token)
	token = rw.Lock()
	rw.RLockContext(context.Background(), WithAcquireTimeout(time.Millisecond))
	rw.Unlock(token)
	after := GetMetrics()
	if after.LTimeOutTimes-before.LTimeOutTimes != 1 || after.RWTimeOutTimes-before.RWTimeOutTimes != 2 {
		t.Fatalf("timeout metrics: before=%+v after=%+v", before, after)
	}
}

func BenchmarkSyncLock(b *testing.B) {
	for range b.N {
		var l sync.Mutex
		l.Lock()
		l.Unlock()
	}
}

func BenchmarkLock(b *testing.B) {
	for range b.N {
		l := NewLocker()
		token := l.Lock()
		l.Unlock(token)
	}
}

func BenchmarkSyncRWLock(b *testing.B) {
	for range b.N {
		var l sync.RWMutex
		l.RLock()
		l.RUnlock()
		l.Lock()
		l.Unlock()
	}
}

func BenchmarkRWLock(b *testing.B) {
	for range b.N {
		l := NewRWLocker()
		token := l.RLock()
		l.RUnlock(token)
		token = l.Lock()
		l.Unlock(token)
	}
}
