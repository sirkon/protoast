package ast

var _ ScalarNode = &Float32{}

// Float32 представление типа float
type Float32 struct {
	unique
}

func (*Float32) genericType() {}
func (*Float32) node()        {}
func (*Float32) scalar()      {}
