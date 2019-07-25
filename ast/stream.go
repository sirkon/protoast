package ast

var _ Type = Stream{}

type Stream struct {
	Type Type
}

func (s Stream) node()		{}
func (s Stream) genericType()	{}
