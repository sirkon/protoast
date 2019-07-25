package ast

var _ ScalarNode = Float32{}

type Float32 struct{}

func (Float32) genericType()	{}
func (Float32) node()		{}
func (Float32) scalar()		{}
