package ast

var _ ScalarNode = &Bool{}

// Bool представление булевского типа
type Bool struct {
	unique
}

func (*Bool) genericType() {}
func (*Bool) node()        {}
func (*Bool) scalar()      {}
