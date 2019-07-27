package ast

var _ Type = &Optional{}

// Optional представление для опциональных полей
type Optional struct {
	unique

	Type Type
}

func (*Optional) genericType() {}
func (*Optional) node()        {}
