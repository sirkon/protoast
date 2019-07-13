package ast

var _ Hashable = Sfixed64{}

// Sfixed64 представление типа sfixed64
type Sfixed64 struct{}

func (Sfixed64) genericType() {}
func (Sfixed64) hashable()    {}
func (Sfixed64) node()        {}
func (Sfixed64) scalar()      {}
