package ast

var _ Hashable = Sint32{}

type Sint32 struct{}

func (Sint32) genericType()	{}
func (Sint32) hashable()	{}
func (Sint32) node()		{}
func (Sint32) scalar()		{}
