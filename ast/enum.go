package ast

var _ Type = &Enum{}

// Enum описание enum-типов
type Enum struct {
	File   *File
	Name   string
	Values []EnumValue
}

func (*Enum) genericType() {}
func (*Enum) node()        {}

// EnumValue описание значения перечисления
type EnumValue struct {
	Name    string
	Integer int
	Options Options
}
