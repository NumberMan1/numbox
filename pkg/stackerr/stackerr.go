package stackerr

import (
	"errors"
	"fmt"

	pkgerrors "github.com/pkg/errors"
)

// Callers returns a pkgerrors.StackTrace with the given skip depth.
func Callers(skip int) pkgerrors.StackTrace {
	err := pkgerrors.New("")
	var tracer interface {
		StackTrace() pkgerrors.StackTrace
	}
	if errors.As(err, &tracer) {
		st := tracer.StackTrace()
		if len(st) > skip {
			return st[skip:]
		}
		return st
	}
	return nil
}

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
		e.innerErr = errors.New(e.msg)
	}
	return e.innerErr
}

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
