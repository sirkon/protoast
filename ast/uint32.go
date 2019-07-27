package ast

var _ Hashable = &Uint32{}

// Uint32 представление для типа uint32
type Uint32 struct {
	unique
}

func (*Uint32) genericType() {}
func (*Uint32) hashable()    {}
func (*Uint32) node()        {}
func (*Uint32) scalar()      {}
