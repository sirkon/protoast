package ast

var _ Type = &Enum{}

// Enum представление типа enum
type Enum struct {
	unique

	File      *File
	ParentMsg *Message

	Name   string
	Values []*EnumValue
}

func (*Enum) genericType() {}
func (*Enum) node()        {}

var _ Unique = &EnumValue{}

// EnumValue представление поля
type EnumValue struct {
	unique

	Name    string
	Integer int
	Options []*Option
}
