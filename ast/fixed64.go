package ast

var _ Hashable = &Fixed64{}

// Fixed64 представление типа fixed64
type Fixed64 struct {
	unique
}

func (*Fixed64) equivalent(v ScalarNode) bool {
	_, ok := v.(*Fixed64)
	return ok
}

func (*Fixed64) genericType() {}
func (*Fixed64) hashable()    {}
func (*Fixed64) node()        {}
func (*Fixed64) scalar()      {}
