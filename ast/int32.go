package ast

var _ Hashable = &Int32{}

// Int32 представление типа int32
type Int32 struct {
	unique
}

func (*Int32) genericType() {}
func (*Int32) hashable()    {}
func (*Int32) node()        {}
func (*Int32) scalar()      {}
