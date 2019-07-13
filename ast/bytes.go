package ast

var _ ScalarNode = Bytes{}

// Bytes представление типа Bytes
type Bytes struct{}

func (Bytes) genericType() {}
func (Bytes) node()        {}
func (Bytes) scalar()      {}
