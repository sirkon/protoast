package ast

var _ Hashable = &Sint32{}

// Sint32 представление для типа sint32
type Sint32 struct {
	unique
}

func (*Sint32) equivalent(v ScalarNode) bool {
	_, ok := v.(*Sint32)
	return ok
}

func (*Sint32) genericType() {}
func (*Sint32) hashable()    {}
func (*Sint32) node()        {}
func (*Sint32) scalar()      {}
