package ast

var _ ScalarNode = Bool{}

// Bool представление типа bool
type Bool struct{}

func (Bool) genericType() {}
func (Bool) node()        {}
func (Bool) scalar()      {}
