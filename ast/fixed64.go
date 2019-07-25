package ast

var _ Hashable = Fixed64{}

type Fixed64 struct{}

func (Fixed64) genericType()	{}
func (Fixed64) hashable()	{}
func (Fixed64) node()		{}
func (Fixed64) scalar()		{}
