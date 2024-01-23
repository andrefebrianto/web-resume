package errorx

import (
	"bytes"
	"errors"
	"fmt"
	"net/http"
	"reflect"
	"runtime"
)

// This wrapper is inspired by https://github.com/go-errors/errors

const MaxStackDepth = 50

type Error struct {
	Code    int
	Message string
	Type    Type
	Details map[string]interface{}
	Err     error

	stack  []uintptr
	frames []StackFrame
}

func New(errType Type, msg string, e interface{}) *Error {
	c, cFound := TypeToCode[errType]
	if !cFound {
		c = http.StatusInternalServerError
	}

	var err error

	switch e := e.(type) {
	case *Error:
		return e
	case error:
		err = e
	default:
		err = fmt.Errorf("%v", e)
	}

	stack := make([]uintptr, MaxStackDepth)
	length := runtime.Callers(2, stack[:])

	return &Error{
		Code:    c,
		Message: msg,
		Type:    errType,
		Err:     err,
		stack:   stack[:length],
	}
}

func Wrap(e interface{}) *Error {
	if e == nil {
		return nil
	}

	var err error

	switch e := e.(type) {
	case *Error:
		return e
	case error:
		err = e
	default:
		err = fmt.Errorf("%v", e)
	}

	stack := make([]uintptr, MaxStackDepth)

	// The skip parameter indicates how far up the stack to start the stacktrace. 0 is from
	// the current call, 1 from its caller, etc.
	length := runtime.Callers(3, stack[:])
	return &Error{
		Err:   err,
		stack: stack[:length],
	}
}

func Errorf(format string, a ...interface{}) *Error {
	return Wrap(fmt.Errorf(format, a...))
}

func Is(e error, original error) bool {
	if errors.Is(e, original) {
		return true
	}

	if original, ok := original.(*Error); ok {
		return Is(e, original.Err)
	}

	return false
}

func (e Error) Is(err error) bool {
	var pointer *Error
	var instance Error
	if errors.As(err, &pointer) {
		return cmp(&e, pointer)
	}
	if errors.As(err, &instance) {
		return cmp(&e, &instance)
	}
	return false
}

func cmp(e1 *Error, e2 *Error) bool {
	if e1.Code != e2.Code &&
		e1.Message != e2.Message &&
		e1.Type != e2.Type {
		return false
	}
	return true
}

func (e Error) Error() string {
	if e.Err == nil {
		return fmt.Sprintf("error %d: %s", e.Code, e.Message)
	}
	return fmt.Sprintf("error %d: %s: %s", e.Code, e.Message, e.Err.Error())
}

func (err *Error) Stack() []byte {
	buf := bytes.Buffer{}

	for _, frame := range err.StackFrames() {
		buf.WriteString(frame.String())
	}

	return buf.Bytes()
}

func (err *Error) Callers() []uintptr {
	return err.stack
}

func (err *Error) ErrorStack() string {
	return fmt.Sprintf("%s %s\n%s", err.TypeName(), err.Error(), string(err.Stack()))
}

func (err *Error) StackFrames() []StackFrame {
	if err.frames == nil {
		err.frames = make([]StackFrame, len(err.stack))

		for i, pc := range err.stack {
			err.frames[i] = NewStackFrame(pc)
		}
	}

	return err.frames
}

func (err *Error) TypeName() string {
	if _, ok := err.Err.(uncaughtPanic); ok {
		return "panic"
	}
	return reflect.TypeOf(err.Err).String()
}

func (e Error) Unwrap() error {
	return e.Err
}
