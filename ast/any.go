package ast

var _ Type = Any{}

type Any struct{}

func (Any) genericType()	{}
func (Any) node()		{}
