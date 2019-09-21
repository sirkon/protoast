package ast

var _ Node = &Method{}

// Method представление для метода
type Method struct {
	unique

	File    *File
	Service *Service

	Name   string
	Input  Type
	Output Type

	Options []*MethodOption
}

func (m *Method) node() {}

var _ Unique = &MethodOption{}

// MethodOption представление для опции метода
type MethodOption struct {
	unique

	Name      string
	Extension *Extension
	Values    []*MethodOptionValue
}

var _ Unique = &MethodOptionValue{}

// MethodOptionValue представление для значения опции метода
type MethodOptionValue struct {
	unique

	Name  string
	Value string
}
