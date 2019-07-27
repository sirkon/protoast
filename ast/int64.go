package ast

var _ Hashable = &Int64{}

// Int64 представление типа int64
type Int64 struct {
	unique
}

func (*Int64) genericType() {}
func (*Int64) hashable()    {}
func (*Int64) node()        {}
func (*Int64) scalar()      {}
