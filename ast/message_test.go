package ast

import "testing"

func TestMessage_String(t *testing.T) {
	file := &File{
		Package: "a.b.c",
	}
	tests := []struct {
		name string
		msg  *Message
		want string
	}{
		{
			name: "no parents",
			msg: &Message{
				File:      file,
				ParentMsg: nil,
				Name:      "Msg",
			},
			want: "a.b.c.Msg",
		},
		{
			name: "there are parents",
			msg: &Message{
				File: file,
				ParentMsg: &Message{
					File: file,
					ParentMsg: &Message{
						File: file,
						Name: "GranParent",
					},
					Name: "Parent",
				},
				Name: "Msg",
			},
			want: "a.b.c.GranParent.Parent.Msg",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.msg.String(); got != tt.want {
				t.Errorf("String() = %v, want %v", got, tt.want)
			}
		})
	}
}
