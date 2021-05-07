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

// Type поиск типа по имени.
func (f *File) Type(name string) Type {
	for _, typ := range f.Types {
		switch v := typ.(type) {
		case *Message:
			if v.Name == name {
				return v
			}
		case *Enum:
			if v.Name == name {
				return v
			}
		}
	}

	return nil
}

// ScanTypes пробежка по типам данного пакета
func (f *File) ScanTypes(inspector func(typ Type) bool) {
	for _, typ := range f.Types {
		switch v := typ.(type) {
		case *Message:
			v.ScanTypes(inspector)
		default:
			inspector(typ)
		}
	}
}

// Service поиск сервиса по имени
func (f *File) Service(name string) *Service {
	for _, service := range f.Services {
		if service.Name == name {
			return service
		}
	}

	return nil
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
