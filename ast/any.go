package ast

var _ Type = Any{}

// Any представление типа google.protobuf.any
type Any struct{}

func (Any) genericType() {}
func (Any) node()        {}
