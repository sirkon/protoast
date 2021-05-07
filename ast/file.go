package ast

import (
	"io"

	"github.com/sirkon/protoast/ast/internal/liner"
	"github.com/sirkon/protoast/internal/errors"
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

// Message поиск структуры по имени.
// Возвращает ошибку ErrorTypeNotFound если такой тип с таким именем не найден.
func (f *File) Message(name string) (*Message, error) {
	typ := f.Type(name)
	if typ == nil {
		return nil, ErrorTypeNotFound(name)
	}

	switch v := typ.(type) {
	case *Message:
		return v, nil
	default:
		return nil, unexpectedType(typ, &Message{})
	}
}

// Enum поиск перечисления по имени
// Возвращает ошибку ErrorTypeNotFound если такой тип с таким именем не найден.
func (f *File) Enum(name string) (*Enum, error) {
	typ := f.Type(name)
	if typ == nil {
		return nil, ErrorTypeNotFound(name)
	}

	switch v := typ.(type) {
	case *Enum:
		return v, nil
	default:
		return nil, unexpectedType(typ, &Enum{})
	}
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

func unexpectedType(typ Type, expected Type) error {
	return errors.Newf("type is %T, not %T", typ, expected)
}
