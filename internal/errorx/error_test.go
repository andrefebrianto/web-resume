package errorx

import (
	"errors"
	"fmt"
	"reflect"
	"testing"
)

func TestNew(t *testing.T) {
	type args struct {
		errType Type
		msg     string
		e       interface{}
	}
	tests := []struct {
		name string
		args args
		want *Error
	}{
		{
			name: "Create new error with wrong type",
			args: args{
				errType: Type("INVALID_TYPE"),
				msg:     "Wrong type message",
			},
			want: &Error{
				Code:    TypeToCode[TypeInternal],
				Message: "Wrong type message",
				Type:    Type("INVALID_TYPE"),
			},
		},
		{
			name: "Create new invalid parameter error",
			args: args{
				errType: TypeInvalidParameter,
				msg:     "Username or email or password is invalid",
			},
			want: &Error{
				Code:    TypeToCode[TypeInvalidParameter],
				Message: "Username or email or password is invalid",
				Type:    TypeInvalidParameter,
			},
		},
		{
			name: "Create new error with existing predefined error",
			args: args{
				errType: TypeInternal,
				msg:     TypeInternal.String(),
				e:       ErrForbidden,
			},
			want: ErrForbidden,
		},
		{
			name: "Wrap error",
			args: args{
				errType: TypeForbidden,
				msg:     "You are not authorized to access this resource",
				e:       errors.New("forbidden"),
			},
			want: &Error{
				Code:    TypeToCode[TypeForbidden],
				Message: "You are not authorized to access this resource",
				Type:    TypeForbidden,
				Err:     errors.New("forbidden"),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := New(tt.args.errType, tt.args.msg, tt.args.e); !Is(got, tt.want) {
				t.Errorf("New() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWrap(t *testing.T) {
	type args struct {
		e interface{}
	}
	tests := []struct {
		name string
		args args
		want *Error
	}{
		{
			name: "Wrap error from nil",
			args: args{
				e: nil,
			},
			want: nil,
		},
		{
			name: "Wrap error from string",
			args: args{
				e: "invalid parameter",
			},
			want: &Error{
				Err: fmt.Errorf("%v", "invalid parameter"),
			},
		},
		{
			name: "Wrap error from integer",
			args: args{
				e: "123",
			},
			want: &Error{
				Err: fmt.Errorf("%v", 123),
			},
		},
		{
			name: "Wrap error from errorx package",
			args: args{
				e: ErrForbidden,
			},
			want: ErrForbidden,
		},
		{
			name: "Wrap standard error",
			args: args{
				e: errors.New("standard error"),
			},
			want: &Error{
				Err: errors.New("standard error"),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Wrap(tt.args.e); !Is(got, tt.want) {
				t.Errorf("Wrap() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestErrorf(t *testing.T) {
	type args struct {
		format string
		a      []interface{}
	}
	tests := []struct {
		name string
		args args
		want *Error
	}{
		{
			name: "create error",
			args: args{
				format: "error can't find user id %s in order %d",
				a:      []interface{}{"devoratio", 168},
			},
			want: &Error{
				Err: fmt.Errorf("error can't find user id %s in order %d", "devoratio", 168),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Errorf(tt.args.format, tt.args.a...); !Is(got, tt.want) {
				t.Errorf("Errorf() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIs(t *testing.T) {
	firstError := errors.New("standard error")
	secondError := errors.New("different error")

	type args struct {
		e        error
		original error
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "compare standard error with standard error",
			args: args{
				e:        firstError,
				original: secondError,
			},
			want: false,
		},
		{
			name: "compare standard error with wrapped error",
			args: args{
				e:        New(TypeForbidden, TypeForbidden.String(), firstError),
				original: firstError,
			},
			want: true,
		},
		{
			name: "compare wrapped error with standard error",
			args: args{
				e:        firstError,
				original: New(TypeForbidden, TypeForbidden.String(), firstError),
			},
			want: true,
		},
		{
			name: "compare wrapped error with wrapped error",
			args: args{
				e:        ErrForbidden,
				original: New(TypeForbidden, TypeForbidden.String(), firstError),
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Is(tt.args.e, tt.args.original); got != tt.want {
				t.Errorf("Is() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestError_Is(t *testing.T) {
	standardError := errors.New("standard error")

	type fields struct {
		Code    int
		Message string
		Type    Type
		Details map[string]interface{}
		Err     error
	}
	type args struct {
		err error
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		{
			name: "compare with standard error",
			fields: fields{
				Code:    TypeToCode[TypeBadGateway],
				Message: TypeBadGateway.String(),
				Type:    TypeBadGateway,
				Err:     standardError,
			},
			args: args{
				err: standardError,
			},
			want: false,
		},
		{
			name: "compare with equal pointer wrapped error",
			fields: fields{
				Code:    TypeToCode[TypeBadGateway],
				Message: TypeBadGateway.String(),
				Type:    TypeBadGateway,
				Err:     standardError,
			},
			args: args{
				err: ErrBadGateway,
			},
			want: true,
		},
		{
			name: "compare with equal wrapped error",
			fields: fields{
				Code:    TypeToCode[TypeBadGateway],
				Message: TypeBadGateway.String(),
				Type:    TypeBadGateway,
				Err:     standardError,
			},
			args: args{
				err: Error{
					Code:    TypeToCode[TypeBadGateway],
					Message: TypeBadGateway.String(),
					Type:    TypeBadGateway,
					Err:     standardError,
				},
			},
			want: true,
		},
		{
			name: "compare with different wrapped error",
			fields: fields{
				Code:    TypeToCode[TypeNotImplemented],
				Message: TypeNotImplemented.String(),
				Type:    TypeNotImplemented,
				Err:     standardError,
			},
			args: args{
				err: Error{
					Code:    TypeToCode[TypeBadGateway],
					Message: TypeBadGateway.String(),
					Type:    TypeBadGateway,
					Err:     standardError,
				},
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := Error{
				Code:    tt.fields.Code,
				Message: tt.fields.Message,
				Type:    tt.fields.Type,
				Details: tt.fields.Details,
				Err:     tt.fields.Err,
			}
			if got := e.Is(tt.args.err); got != tt.want {
				t.Errorf("Error.Is() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestError_Error(t *testing.T) {
	rootErr := errors.New("root error")
	type fields struct {
		Code    int
		Message string
		Type    Type
		Details map[string]interface{}
		Err     error
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "without root error",
			fields: fields{
				Code:    TypeToCode[TypeBadGateway],
				Message: TypeBadGateway.String(),
				Type:    TypeBadGateway,
			},
			want: fmt.Sprintf("error %d: %s", TypeToCode[TypeBadGateway], TypeBadGateway.String()),
		},
		{
			name: "with root error",
			fields: fields{
				Code:    TypeToCode[TypeBadGateway],
				Message: TypeBadGateway.String(),
				Type:    TypeBadGateway,
				Err:     rootErr,
			},
			want: fmt.Sprintf("error %d: %s: %s", TypeToCode[TypeBadGateway], TypeBadGateway.String(), rootErr),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := Error{
				Code:    tt.fields.Code,
				Message: tt.fields.Message,
				Type:    tt.fields.Type,
				Details: tt.fields.Details,
				Err:     tt.fields.Err,
			}
			if got := e.Error(); got != tt.want {
				t.Errorf("Error.Error() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestError_Stack(t *testing.T) {
	stack := []uintptr{0}
	stackFrames := []StackFrame{NewStackFrame(stack[0])}

	type fields struct {
		Code    int
		Message string
		Type    Type
		Details map[string]interface{}
		Err     error
		stack   []uintptr
		frames  []StackFrame
	}
	tests := []struct {
		name   string
		fields fields
		want   []byte
	}{
		{
			name: "current stack with created stack frame",
			fields: fields{
				Code:    TypeToCode[TypeServiceUnavailable],
				Message: TypeServiceUnavailable.String(),
				Type:    TypeServiceUnavailable,
				stack:   stack,
				frames:  stackFrames,
			},
			want: []byte{58, 48, 32, 40, 48, 120, 48, 41, 10, 9, 58, 32, 63, 63, 63, 10},
		},
		{
			name: "current stack with nil stack frame",
			fields: fields{
				Code:    TypeToCode[TypeServiceUnavailable],
				Message: TypeServiceUnavailable.String(),
				Type:    TypeServiceUnavailable,
				stack:   stack,
				frames:  nil,
			},
			want: []byte{58, 48, 32, 40, 48, 120, 48, 41, 10, 9, 58, 32, 63, 63, 63, 10},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := &Error{
				Code:    tt.fields.Code,
				Message: tt.fields.Message,
				Type:    tt.fields.Type,
				Details: tt.fields.Details,
				Err:     tt.fields.Err,
				stack:   tt.fields.stack,
				frames:  tt.fields.frames,
			}
			if got := err.Stack(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Error.Stack() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestError_Callers(t *testing.T) {
	type fields struct {
		Code    int
		Message string
		Type    Type
		Details map[string]interface{}
		Err     error
		stack   []uintptr
		frames  []StackFrame
	}
	tests := []struct {
		name   string
		fields fields
		want   []uintptr
	}{
		{
			name: "stack with MaxStackDepth length",
			fields: fields{
				stack: make([]uintptr, MaxStackDepth),
			},
			want: make([]uintptr, MaxStackDepth),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := &Error{
				Code:    tt.fields.Code,
				Message: tt.fields.Message,
				Type:    tt.fields.Type,
				Details: tt.fields.Details,
				Err:     tt.fields.Err,
				stack:   tt.fields.stack,
				frames:  tt.fields.frames,
			}
			if got := err.Callers(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Error.Callers() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestError_ErrorStack(t *testing.T) {
	type fields struct {
		Code    int
		Message string
		Type    Type
		Details map[string]interface{}
		Err     error
		stack   []uintptr
		frames  []StackFrame
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "",
			fields: fields{
				Code:    TypeToCode[TypeForbidden],
				Message: TypeForbidden.String(),
				Type:    TypeForbidden,
				Err:     errors.New("unauthorized"),
				stack:   make([]uintptr, MaxStackDepth),
				frames:  make([]StackFrame, 0),
			},
			want: "*errors.errorString error 403: FORBIDDEN: unauthorized\n",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := &Error{
				Code:    tt.fields.Code,
				Message: tt.fields.Message,
				Type:    tt.fields.Type,
				Details: tt.fields.Details,
				Err:     tt.fields.Err,
				stack:   tt.fields.stack,
				frames:  tt.fields.frames,
			}
			if got := err.ErrorStack(); got != tt.want {
				t.Errorf("Error.ErrorStack() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestError_TypeName(t *testing.T) {
	type fields struct {
		Code    int
		Message string
		Type    Type
		Details map[string]interface{}
		Err     error
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "type panic",
			fields: fields{
				Err: uncaughtPanic{message: "paniced"},
			},
			want: "panic",
		},
		{
			name: "type *errors.errorString",
			fields: fields{
				Err: errors.New("error"),
			},
			want: "*errors.errorString",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := &Error{
				Code:    tt.fields.Code,
				Message: tt.fields.Message,
				Type:    tt.fields.Type,
				Details: tt.fields.Details,
				Err:     tt.fields.Err,
			}
			if got := err.TypeName(); got != tt.want {
				t.Errorf("Error.TypeName() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestError_Unwrap(t *testing.T) {
	rootErr := errors.New("root error")

	tests := []struct {
		name    string
		errors  Error
		wantErr error
	}{
		{
			name:    "unwrap error type errors.errorString",
			errors:  *New(TypeBadGateway, TypeBadGateway.String(), rootErr),
			wantErr: rootErr,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.errors.Unwrap(); err != tt.wantErr {
				t.Errorf("Error.Unwrap() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
