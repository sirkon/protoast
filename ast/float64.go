package ast

var _ ScalarNode = Float64{}

type Float64 struct{}

func (Float64) genericType()	{}
func (Float64) node()		{}
func (Float64) scalar()		{}
