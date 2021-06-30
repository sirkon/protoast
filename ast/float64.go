package ast

var _ ScalarNode = &Float64{}

// Float64 представление типа double
type Float64 struct {
	unique
}

func (*Float64) equivalent(v ScalarNode) bool {
	_, ok := v.(*Float64)
	return ok
}

func (*Float64) genericType() {}
func (*Float64) node()        {}
func (*Float64) scalar()      {}
