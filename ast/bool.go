package ast

var _ ScalarNode = Bool{}

type Bool struct{}

func (Bool) genericType()	{}
func (Bool) node()		{}
func (Bool) scalar()		{}
