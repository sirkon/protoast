package ast

var _ Node = &File{}

// File представление файла (без вложенного AST, только имя файла и название пакета Go)
type File struct {
	Name  string
	GoPkg string
}

func (*File) node() {}
