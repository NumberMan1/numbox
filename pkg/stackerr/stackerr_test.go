package stackerr

import (
	"errors"
	"fmt"
	"strings"
	"testing"

	pkgerrors "github.com/pkg/errors"
)

func TestTraceFrom(t *testing.T) {
	if stack, ok := TraceFrom(nil); ok || stack != nil {
		t.Fatalf("TraceFrom(nil) = %#v, %v; want nil, false", stack, ok)
	}

	rawErr := errors.New("raw")
	if stack, ok := TraceFrom(rawErr); ok || stack != nil {
		t.Fatalf("TraceFrom(rawErr) = %#v, %v; want nil, false", stack, ok)
	}

	stackedErr := pkgerrors.WithStack(rawErr)
	stack, ok := TraceFrom(stackedErr)
	if !ok || len(stack) == 0 {
		t.Fatalf("TraceFrom(stackedErr) = %#v, %v; want stack, true", stack, ok)
	}
}

func TestWithStackAddsStackAndPreservesIdentity(t *testing.T) {
	sourceErr := errors.New("db failed")

	err := WithStack(sourceErr)

	if !errors.Is(err, sourceErr) {
		t.Fatalf("errors.Is(%v, sourceErr) = false, want true", err)
	}
	var tracer Tracer
	if !errors.As(err, &tracer) {
		t.Fatalf("error %#v does not expose StackTrace", err)
	}
}

func TestWithStackKeepsNilAndStackedError(t *testing.T) {
	if WithStack(nil) != nil {
		t.Fatal("WithStack(nil) should return nil")
	}

	sourceErr := pkgerrors.WithStack(errors.New("db failed"))
	if got := WithStack(sourceErr); got != sourceErr {
		t.Fatalf("WithStack returned %#v, want original stacked error", got)
	}
}

func TestFromPanicWrapsNonErrorWithStack(t *testing.T) {
	err := FromPanic("boom")

	if err.Error() != "boom" {
		t.Fatalf("FromPanic error = %q, want boom", err.Error())
	}
	var tracer Tracer
	if !errors.As(err, &tracer) {
		t.Fatalf("FromPanic = %#v, want Tracer", err)
	}
	if formatted := fmt.Sprintf("%+v", tracer.StackTrace()); !strings.Contains(formatted, "TestFromPanicWrapsNonErrorWithStack") {
		t.Fatalf("stack trace = %q, want caller stack", formatted)
	}
}

func TestFromPanicPreservesOriginalError(t *testing.T) {
	sourceErr := errors.New("cache refresh failed")

	err := FromPanic(sourceErr)

	if !errors.Is(err, sourceErr) {
		t.Fatalf("errors.Is(%v, sourceErr) = false, want true", err)
	}
}

func TestLegacyErrorChainAndCaptureBounds(t *testing.T) {
	cause := errors.New("original")
	wrapped := FormatStackError(cause, "context")
	if !errors.Is(wrapped, cause) {
		t.Fatal("legacy wrapper hides cause")
	}
	if len(Capture(-1)) == 0 || len(Capture(100000)) != 0 {
		t.Fatal("invalid stack skip handling")
	}
}
