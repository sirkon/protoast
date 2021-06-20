package ast

import (
	"path"
	"strings"
)

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

// InputMessage возвращает структуру запроса (минуя оборачивающий Stream, если нужно)
func (m *Method) InputMessage() *Message {
	return getMessage(m.Input)
}

// OutputMessage аналогично InputMessage, возвращает структуру ответа, при необходимости снимая stream
func (m *Method) OutputMessage() *Message {
	return getMessage(m.Output)
}

func getMessage(m Type) *Message {
	v, ok := m.(*Message)
	if ok {
		return v
	}

	return m.(*Stream).Type.(*Message)
}

// URI метода как принято в gRPC
func (m *Method) URI() string {
	dir, _ := path.Split(m.Service.File.Name)
	dir, _ = path.Split(dir)
	dir = strings.TrimRight(dir, "/")

	var serviceFullName string
	if dir != "" {
		serviceFullName = strings.ReplaceAll(dir, "/", ".") + "."
	}
	serviceFullName += m.Service.File.Package

	return "/" + serviceFullName + "/" + m.Name
}
