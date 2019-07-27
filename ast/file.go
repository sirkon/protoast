package ast

var _ Node = &File{}

// File представление для файла
type File struct {
	unique

	Name    string
	Package string

	Imports  []*Import
	Types    []Type
	Services []*Service
	Options  []*Option
}

func (*File) node() {}
