package ast

var _ ScalarNode = &Bytes{}

// Bytes представление типа bytes
type Bytes struct {
	unique
}

func (*Bytes) equivalent(v ScalarNode) bool {
	_, ok := v.(*Bytes)
	return ok
}

func (*Bytes) genericType() {}
func (*Bytes) node()        {}
func (*Bytes) scalar()      {}
