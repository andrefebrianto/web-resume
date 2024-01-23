package errorx

import (
	"reflect"
	"strings"
	"testing"
)

func TestNewStackFrame(t *testing.T) {
	type args struct {
		pc uintptr
	}
	tests := []struct {
		name      string
		args      args
		wantFrame StackFrame
	}{
		{
			name: "create new stackframe with program counter 0",
			args: args{
				pc: 0,
			},
			wantFrame: StackFrame{
				ProgramCounter: 0,
			},
		},
		{
			name: "create new stackframe with program counter 1",
			args: args{
				pc: 1,
			},
			wantFrame: StackFrame{
				ProgramCounter: 1,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotFrame := NewStackFrame(tt.args.pc); !reflect.DeepEqual(gotFrame, tt.wantFrame) {
				t.Errorf("NewStackFrame() = %v, want %v", gotFrame, tt.wantFrame)
			}
		})
	}
}

func TestStackFrame_String(t *testing.T) {
	err := firstCall()
	if !Is(err, ErrBadGateway) {
		t.Fatal(err)
	}

	e := err.(*Error)
	stack := string(e.Stack())
	if !strings.Contains(stack, "firstCall: return secondCall(5)") {
		t.Errorf("Stack trace does not contain source line: 'firstCall: return secondCall(5)'")
		t.Errorf(stack)
	}
	if !strings.Contains(stack, "stackframe_test.go:") {
		t.Errorf("Stack trace does not contain file name: 'stackframe_test.go:'")
		t.Errorf(stack)
	}
}

func firstCall() error {
	return secondCall(5)
}

func secondCall(i int) error {
	return thirdCall()
}

func thirdCall() error {
	return New(TypeBadGateway, TypeBadGateway.String(), nil)
}
