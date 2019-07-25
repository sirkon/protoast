package ast

var _ Hashable = Fixed32{}

type Fixed32 struct{}

func (Fixed32) genericType()	{}
func (Fixed32) hashable()	{}
func (Fixed32) node()		{}
func (Fixed32) scalar()		{}
