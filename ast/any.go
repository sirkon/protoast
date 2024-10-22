package ast

var _ Type = &Any{}

// Any представление типа golang.protobuf.Any
type Any struct {
	unique

	File *File
}

func (*Any) String() string { return "any" }
func (*Any) genericType()   {}
func (*Any) node()          {}
