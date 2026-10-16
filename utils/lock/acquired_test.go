package lock

import (
	"runtime"
	"sync"
	"testing"
	"time"
)

func TestLockerAcquiredConcurrentTransitions(t *testing.T) {
	lk := NewLocker()
	stale := lk.Lock()
	lk.Unlock(stale)

	start := make(chan struct{})
	done := make(chan struct{})
	var readers sync.WaitGroup
	for i := 0; i < 8; i++ {
		readers.Add(1)
		go func() {
			defer readers.Done()
			<-start
			for {
				select {
				case <-done:
					return
				default:
					if lk.Acquired(stale) {
						t.Error("a released token became acquired again")
						return
					}
					runtime.Gosched()
				}
			}
		}()
	}
	defer func() {
		close(done)
		readers.Wait()
	}()
	close(start)

	for i := 0; i < 1000; i++ {
		token := lk.Lock(WithLockHoldTimeout(time.Hour))
		if !lk.Acquired(token) {
			t.Fatal("new token should hold the lock")
		}
		runtime.Gosched()
		if !lk.Unlock(token) {
			t.Fatal("current token should unlock successfully")
		}
		if lk.Acquired(token) {
			t.Fatal("released token should not hold the lock")
		}
	}
}

func TestLockerUnlockCallbackCanQueryAcquired(t *testing.T) {
	lk := NewLocker()
	token := lk.Lock(WithLockHoldTimeout(time.Hour))
	result := make(chan bool, 1)
	go func() {
		lk.Unlock(token, WithUnlockCallback(func() {
			result <- lk.Acquired(token)
		}))
	}()

	select {
	case acquired := <-result:
		if acquired {
			t.Fatal("unlock callback should observe a released token")
		}
	case <-time.After(5 * time.Second):
		t.Fatal("unlock callback deadlocked while querying Acquired")
	}
}
