package ast

var _ Hashable = String{}

type String struct{}

func (String) genericType()	{}
func (String) hashable()	{}
func (String) node()		{}
func (String) scalar()		{}
