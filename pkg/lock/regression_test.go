package lock

import (
	"sync"
	"testing"
	"time"
)

func TestAcquiredConcurrentWithLock(t *testing.T) {
	l := NewLocker()
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		for range 1000 {
			l.Acquired(1)
		}
	}()
	for range 1000 {
		token := l.Lock()
		l.Unlock(token)
	}
	wg.Wait()
}
func TestUnlockCallbackCanQueryAndReacquire(t *testing.T) {
	for _, l := range []Locker{NewLocker(), NewRWLocker()} {
		token := l.Lock()
		done := make(chan bool, 1)
		go func() {
			done <- l.Unlock(token, WithUnlockCallback(func() {
				if l.Acquired(token) {
					t.Error("released token still acquired")
				}
				next := l.Lock()
				l.Unlock(next)
			}))
		}()
		select {
		case ok := <-done:
			if !ok {
				t.Fatal("unlock failed")
			}
		case <-time.After(time.Second):
			t.Fatal("callback deadlocked")
		}
		if l.Unlock(token) {
			t.Fatal("duplicate unlock succeeded")
		}
	}
	rw := NewRWLocker()
	token := rw.RLock()
	done := make(chan bool, 1)
	go func() { done <- rw.RUnlock(token, WithUnlockCallback(func() { next := rw.Lock(); rw.Unlock(next) })) }()
	select {
	case ok := <-done:
		if !ok {
			t.Fatal("read unlock failed")
		}
	case <-time.After(time.Second):
		t.Fatal("read callback deadlocked")
	}
}
