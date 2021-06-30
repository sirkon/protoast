package ast

var _ Hashable = &Sint64{}

// Sint64 представление типа sint64
type Sint64 struct {
	unique
}

func (*Sint64) equivalent(v ScalarNode) bool {
	_, ok := v.(*Sint64)
	return ok
}

func (*Sint64) genericType() {}
func (*Sint64) hashable()    {}
func (*Sint64) node()        {}
func (*Sint64) scalar()      {}
