package ast

var _ ScalarNode = &Bytes{}

// Bytes представление типа bytes
type Bytes struct {
	unique
}

func (*Bytes) genericType() {}
func (*Bytes) node()        {}
func (*Bytes) scalar()      {}
