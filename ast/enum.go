package ast

var _ Type = &Enum{}

type Enum struct {
	File		*File
	ParentMsg	*Message
	Name		string
	Values		[]EnumValue
}

func (*Enum) genericType()	{}
func (*Enum) node()		{}

type EnumValue struct {
	Name	string
	Integer	int
	Options	Options
}
