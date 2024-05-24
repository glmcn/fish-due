package panictool

import (
	"testing"
	"time"
)

func TestGetString(t *testing.T) {
	type args[T interface {
		string | int8 | int32 | int64 | *time.Time | *string | *int64
	}] struct {
		in T
	}
	type testCase[T interface {
		string | int8 | int32 | int64 | *time.Time | *string | *int64
	}] struct {
		name string
		args args[T]
		want string
	}
	tests := []testCase[ /* TODO: Insert concrete types here */ *string]{
		// TODO: Add test cases.
		{
			name: "",
			args: args[*string]{
				in: nil,
			},
			want: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetString(tt.args.in); got != tt.want {
				t.Errorf("GetString() = %v, want %v", got, tt.want)
			}
		})
	}
}
