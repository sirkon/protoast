package ast

var _ Hashable = Sfixed32{}

// Sfixed32 представление типа sfixed32
type Sfixed32 struct{}

func (Sfixed32) genericType() {}
func (Sfixed32) hashable()    {}
func (Sfixed32) node()        {}
func (Sfixed32) scalar()      {}
