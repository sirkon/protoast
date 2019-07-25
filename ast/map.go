package ast

var _ Type = Map{}

type Map struct {
	KeyType		Hashable
	ValueType	Type
}

func (Map) genericType()	{}
func (Map) node()		{}
