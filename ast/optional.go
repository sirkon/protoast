package ast

var _ Type = Optional{}

type Optional struct {
	Type Type
}

func (Optional) genericType()	{}
func (Optional) node()		{}
