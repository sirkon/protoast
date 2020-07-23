package ast

import (
	"io"

	"github.com/sirkon/protoast/ast/internal/liner"
)

var _ Node = &File{}

// File представление для файла
type File struct {
	unique

	Name    string
	Package string
	Syntax  string

	Imports    []*Import
	Types      []Type
	Extensions []*Extension
	Services   []*Service
	Options    []*Option

	GoPath string
	GoPkg  string
}

func (*File) node() {}

func (f *File) print(dest io.Writer, printer *Printer) error {
	l := liner.New(dest)
	l.Line(`syntax = "$"`, f.Syntax)
	l.Newl()
	l.Line("package $", f.Package)
	l.Newl()

	for _, imp := range f.Imports {
		printer.Plan(imp.File)
		l.Line(`import "$";`, imp.Path)
	}

	l.Newl()

	return nil
}
