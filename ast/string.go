package ast

var _ Hashable = &String{}

// String представление для стрового типа
type String struct {
	unique
}

func (*String) equivalent(v ScalarNode) bool {
	_, ok := v.(*String)
	return ok
}

func (*String) String() string {
	return "string"
}

func (*String) genericType() {}
func (*String) hashable()    {}
func (*String) node()        {}
func (*String) scalar()      {}
