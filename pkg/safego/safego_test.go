package safego

import (
	"fmt"
	"strings"
	"sync/atomic"
	"testing"
	"time"
)

func TestGoReturnsBeforeWorkCompletes(t *testing.T) {
	workStarted := make(chan struct{})
	releaseWork := make(chan struct{})
	callReturned := make(chan struct{})
	defer close(releaseWork)

	go func() {
		Go(func() {
			close(workStarted)
			<-releaseWork
		})
		close(callReturned)
	}()

	select {
	case <-workStarted:
	case <-time.After(time.Second):
		t.Fatal("work did not start")
	}
	select {
	case <-callReturned:
	case <-time.After(time.Second):
		t.Fatal("Go blocked until work completed")
	}
}

func TestGoWithReporterRecoversPanicAfterWorkDefers(t *testing.T) {
	var deferred atomic.Bool
	reported := make(chan error, 1)

	goWithReporter(func() {
		defer deferred.Store(true)
		panic("boom")
	}, func(err error) {
		reported <- err
	})

	select {
	case err := <-reported:
		if !deferred.Load() {
			t.Fatal("work defer did not run before panic report")
		}
		if !strings.Contains(err.Error(), "boom") {
			t.Fatalf("error = %v, want panic payload", err)
		}
		if formatted := fmt.Sprintf("%+v", err); !strings.Contains(formatted, "safego.go") {
			t.Fatalf("formatted error = %q, want safego stack", formatted)
		}
	case <-time.After(time.Second):
		t.Fatal("panic was not reported")
	}
}
