package ast

import "testing"

func TestEnum_String(t *testing.T) {
	file := &File{
		Package: "e.d.f",
	}
	tests := []struct {
		name string
		enum *Enum
		want string
	}{
		{
			name: "no parents",
			enum: &Enum{
				File: file,
				Name: "Enum",
			},
			want: "e.d.f.Enum",
		},
		{
			name: "there are parents",
			enum: &Enum{
				File: file,
				ParentMsg: &Message{
					File: file,
					ParentMsg: &Message{
						File:      file,
						ParentMsg: nil,
						Name:      "GranParent",
					},
					Name: "Parent",
				},
				Name: "Enum",
			},
			want: "e.d.f.GranParent.Parent.Enum",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.enum.String(); got != tt.want {
				t.Errorf("String() = %v, want %v", got, tt.want)
			}
		})
	}
}
