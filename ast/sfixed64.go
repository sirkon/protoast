package ast

var _ Hashable = &Sfixed64{}

// Sfixed64 представление для типа sfixed64
type Sfixed64 struct {
	unique
}

func (*Sfixed64) equivalent(v ScalarNode) bool {
	_, ok := v.(*Sfixed64)
	return ok
}

func (*Sfixed64) genericType() {}
func (*Sfixed64) hashable()    {}
func (*Sfixed64) node()        {}
func (*Sfixed64) scalar()      {}
