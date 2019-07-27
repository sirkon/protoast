package ast

var _ Hashable = &Uint64{}

// Uint64 представление для типа uint64
type Uint64 struct {
	unique
}

func (*Uint64) genericType() {}
func (*Uint64) hashable()    {}
func (*Uint64) node()        {}
func (*Uint64) scalar()      {}
