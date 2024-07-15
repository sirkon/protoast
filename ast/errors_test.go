package ast

import (
	"io"
	"testing"
)

func TestIsErrorTypeNotFound(t *testing.T) {
	type args struct{}
	tests := []struct {
		name string
		err  error
		want bool
	}{
		{
			name: "ok",
			err:  ErrorTypeNotFound("123"),
			want: true,
		},
		{
			name: "ok as well in fact",
			err:  io.EOF,
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsErrorTypeNotFound(tt.err); got != tt.want {
				t.Errorf("IsErrorTypeNotFound() = %v, want %v", got, tt.want)
			}
		})
	}
}
