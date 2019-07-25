package ast

var _ Node = &Method{}

type Method struct {
	File	*File

	Service	*Service
	Name	string
	Input	Type
	Output	Type

	Options	[]MethodOption
}

func (m *Method) node()	{}

type MethodOption struct {
	Name	string
	Values	[]OptionValue
}

type OptionValue struct {
	Name	string
	Value	string
}
