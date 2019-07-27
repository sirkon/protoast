package ast

var _ Hashable = &String{}

// String представление для стрового типа
type String struct {
	unique
}

func (*String) genericType() {}
func (*String) hashable()    {}
func (*String) node()        {}
func (*String) scalar()      {}
