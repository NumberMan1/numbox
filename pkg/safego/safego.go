// Package safego starts fire-and-forget goroutines with panic reporting.
// Safe means panic-isolated only; it does not imply race safety, retries,
// cancellation, waiting, or shutdown draining.
package safego

import (
	"fmt"
	"log/slog"

	"github.com/NumberMan1/numbox/pkg/stackerr"
)

type panicReporter func(error)

// Go starts a fire-and-forget goroutine.
// A panic is reported with its stack and recovered.
func Go(work func()) {
	goWithReporter(work, reportPanic)
}

// GoWithReporter uses a caller-supplied panic reporter; nil uses slog.
// The reporter must not panic. Go does not wait for work or provide cancellation.
func GoWithReporter(work func(), reporter func(error)) {
	if reporter == nil {
		reporter = reportPanic
	}
	goWithReporter(work, reporter)
}

func goWithReporter(work func(), reporter panicReporter) {
	go run(work, reporter)
}

func run(work func(), reporter panicReporter) {
	defer func() {
		if value := recover(); value != nil {
			reporter(stackerr.FromPanic(value))
		}
	}()
	work()
}

func reportPanic(err error) {
	slog.Error("goroutine panic", "error", fmt.Sprintf("%+v", err))
}
