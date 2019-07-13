package ast

var _ Type = Optional{}

// Optional представление опционального поля
type Optional struct {
	Type Type
}

func (Optional) genericType() {}
func (Optional) node()        {}
