package ast

var _ Hashable = &Fixed32{}

// Fixed32 представление типа fixed32
type Fixed32 struct {
	unique
}

func (*Fixed32) equivalent(v ScalarNode) bool {
	_, ok := v.(*Fixed32)
	return ok
}

func (*Fixed32) genericType() {}
func (*Fixed32) hashable()    {}
func (*Fixed32) node()        {}
func (*Fixed32) scalar()      {}
