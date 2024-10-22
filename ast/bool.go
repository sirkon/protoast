package ast

var _ ScalarNode = &Bool{}

// Bool представление булевского типа
type Bool struct {
	unique
}

func (*Bool) equivalent(v ScalarNode) bool {
	_, ok := v.(*Bool)
	return ok
}

func (*Bool) String() string {
	return "bool"
}

func (*Bool) genericType() {}
func (*Bool) node()        {}
func (*Bool) scalar()      {}
