package ast

var _ Node = &File{}

type File struct {
	Name	string
	Package	string

	Services	[]*Service
}

func (*File) node()	{}
