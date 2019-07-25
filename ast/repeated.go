package ast

var _ Type = Repeated{}

type Repeated struct {
	Type Type
}

func (Repeated) genericType()	{}
func (Repeated) node()		{}
