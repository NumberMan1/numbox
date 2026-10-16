package stackerr

import (
	"errors"
	"fmt"

	pkgerrors "github.com/pkg/errors"
)

// Trace is the stack representation exposed by this package.
type Trace = pkgerrors.StackTrace

// Capture records the current call stack after skipping skipFrames helper
// frames. A non-positive skipFrames value keeps the complete captured stack.
func Capture(skipFrames int) Trace {
	err := pkgerrors.New("")
	var tracer Tracer
	if errors.As(err, &tracer) {
		stack := tracer.StackTrace()
		skipFrames = max(0, min(skipFrames, len(stack)))
		return stack[skipFrames:]
	}
	return nil
}

// Tracer is implemented by errors that expose a stack through StackTrace.
// The method name intentionally matches github.com/pkg/errors.
type Tracer interface {
	StackTrace() Trace
}

// TraceFrom returns the first Trace exposed by err's standard error chain.
// It returns false when err is nil or no error in the chain implements Tracer.
func TraceFrom(err error) (Trace, bool) {
	if err == nil {
		return nil, false
	}
	var tracer Tracer
	if !errors.As(err, &tracer) {
		return nil, false
	}
	return tracer.StackTrace(), true
}

// FromPanic converts a recovered value to an error and adds a stack when
// needed. When value is already an error, its identity and error chain are
// preserved.
func FromPanic(value any) error {
	panicErr, ok := value.(error)
	if !ok {
		panicErr = errors.New(fmt.Sprint(value))
	}
	return WithStack(panicErr)
}

// WithStack returns err with a stack. It preserves nil and returns an error
// that already implements Tracer unchanged, so a chain is not stacked twice.
func WithStack(err error) error {
	if err == nil {
		return nil
	}
	if _, ok := TraceFrom(err); ok {
		return err
	}
	return pkgerrors.WithStack(err)
}

// Callers is the compatibility name for Capture.
func Callers(skip int) Trace { return Capture(skip + 1) }

type StackError interface {
	SetInnerError(innerErr error)
	InnerError() error
	StackTrace() pkgerrors.StackTrace
	error
}

func FormatStackError(arg any, errTipArg ...string) error {
	errTip := "business panic recover"

	if len(errTipArg) > 0 {
		errTip = errTipArg[0]
	}

	errPanic, ok := arg.(error)
	if !ok {
		errPanic = errors.New(fmt.Sprint(arg))
	}

	errStack := NewStackError(errPanic.Error())
	errStack.SetInnerError(pkgerrors.Wrap(errPanic, errTip))
	return errStack
}

func NewStackError(message string) StackError {
	return &stackError{
		msg:   message,
		stack: Callers(2),
	}
}

type stackError struct {
	msg      string
	innerErr error
	stack    pkgerrors.StackTrace
}

func (e *stackError) SetInnerError(innerErr error) {
	e.innerErr = innerErr
}

func (e *stackError) InnerError() error {
	if e.innerErr == nil {
		return errors.New(e.msg)
	}
	return e.innerErr
}

// Unwrap exposes the original cause to errors.Is and errors.As.
func (e *stackError) Unwrap() error { return e.innerErr }

func (e *stackError) StackTrace() pkgerrors.StackTrace {
	return e.stack
}

func (e *stackError) Error() string {
	return e.msg
}

func (e *stackError) Format(s fmt.State, verb rune) {
	switch verb {
	case 'v':
		if s.Flag('+') {
			fmt.Fprintf(s, "%s", e.msg)
			e.stack.Format(s, verb)
			return
		}
		fallthrough
	case 's':
		fmt.Fprintf(s, "%s", e.msg)
	case 'q':
		fmt.Fprintf(s, "%q", e.msg)
	}
}
