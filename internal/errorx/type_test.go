package errorx

import (
	"testing"
)

func TestType_String(t *testing.T) {
	tests := []struct {
		name string
		tr   Type
		want string
	}{
		{
			name: "Test Not Found Error",
			tr:   TypeNotFound,
			want: TypeNotFound.String(),
		},
		{
			name: "Test Internal Error",
			tr:   TypeInternal,
			want: TypeInternal.String(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.tr.String(); got != tt.want {
				t.Errorf("Type.String() = %v, want %v", got, tt.want)
			}
		})
	}
}
