package errors

import (
	"errors"
	"fmt"
)

type Code int32

func (c Code) Int() int {
	return int(c)
}

func (c Code) Int32() int32 {
	return int32(c)
}

func (c Code) Int64() int64 {
	return int64(c)
}

type CodeError interface {
	SetStack(stack []byte)
	Stack() []byte
	SetInnerError(innerErr error)
	InnerError() error
	error
}

func newCodeError(code Code, message string) CodeError {
	return &codeError{
		Code: code,
		Info: message,
	}
}

type codeError struct {
	Code     Code   `json:"code"`
	Info     string `json:"info"`
	innerErr error
	stack    []byte
}

func (c *codeError) SetInnerError(innerErr error) {
	c.innerErr = innerErr
}

func (c *codeError) SetStack(stack []byte) {
	c.stack = stack
}

func (c *codeError) InnerError() error {
	if c.innerErr == nil {
		c.innerErr = errors.New(c.Info)
	}
	return c.innerErr
}

func (c *codeError) Stack() []byte {
	return c.stack
}

func (c *codeError) Error() string {
	return c.Info
}

// ToError 将提供的data转换为CodeError
// 如果data本身就是CodeError，则返回data本身
// 如果data不是，则返回CodeError(code, data)
func ToError(code Code, data any) CodeError {
	errInf, ok := data.(error)
	var codeError CodeError
	if ok && errors.As(errInf, &codeError) {
		return codeError
	}
	return newCodeError(code, fmt.Sprint(data))
}

var (
	Is = errors.Is
	As = errors.As
)
