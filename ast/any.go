package ast

var _ Type = &Any{}

// Any представление типа golang.protobuf.Any
type Any struct {
	unique
}

func (*Any) genericType() {}
func (*Any) node()        {}
