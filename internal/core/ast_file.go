package core

import (
	"iter"

	"github.com/emicklei/proto"
)

type File struct {
	isNode
	isNodeOptionable

	proto *proto.Proto
}

// Name returns file name.
func (f *File) Name() string {
	return f.proto.Filename
}

// Package returns package value.
func (f *File) Package() string {
	for _, element := range f.proto.Elements {
		v, ok := element.(*proto.Package)
		if ok {
			return v.Name
		}
	}

	return ""
}

// Messages defined at the top level.
func (f *File) Messages(r *Registry) iter.Seq[*Message] {
	return func(yield func(*Message) bool) {
		for _, element := range f.proto.Elements {
			v, ok := element.(*proto.Message)
			if ok && !v.IsExtend {
				if !yield(r.wrap(v).(*Message)) {
					return
				}
			}
		}
	}
}

// Message returns a message with the given name.
func (f *File) Message(r *Registry, name string) *Message {
	for _, element := range f.proto.Elements {
		v, ok := element.(*proto.Message)
		if ok && v.Name == name {
			return r.wrap(v).(*Message)
		}
	}

	return nil
}

// Enums defined at the top level.
func (f *File) Enums(r *Registry) iter.Seq[*Enum] {
	return func(yield func(*Enum) bool) {
		for _, element := range f.proto.Elements {
			v, ok := element.(*proto.Enum)
			if ok {
				if !yield(r.wrap(v).(*Enum)) {
					return
				}
			}
		}
	}
}

// Enum returns an enum with the given name.
func (f *File) Enum(r *Registry, name string) *Enum {
	for _, element := range f.proto.Elements {
		v, ok := element.(*proto.Enum)
		if ok {
			if v.Name == name {
				return r.wrap(v).(*Enum)
			}
		}
	}

	return nil
}

// Types returns named types defined at the top level.
func (f *File) Types(r *Registry) iter.Seq[NamedType] {
	return func(yield func(levelType NamedType) bool) {
		for _, element := range f.proto.Elements {
			var value NamedType
			switch v := element.(type) {
			case *proto.Message:
				if v.IsExtend {
					continue
				}
				value = r.wrap(v).(NamedType)
			case *proto.Enum:
				value = r.wrap(v).(NamedType)
			}
			if value != nil {
				if !yield(value) {
					return
				}
			}
		}
	}
}

// Type returns a named type with the given name.
func (f *File) Type(r *Registry, typename string) NamedType {
	for _, element := range f.proto.Elements {
		switch v := element.(type) {
		case *proto.Message:
			if v.Name != typename {
				continue
			}

			return r.wrap(v).(*Message)

		case *proto.Enum:
			if v.Name != typename {
				continue
			}

			return r.wrap(v).(*Enum)
		}
	}

	return nil
}

func (f *File) Everything(r *Registry) iter.Seq[Node] {
	return func(yield func(Node) bool) {
		for _, e := range f.proto.Elements {
			if v, ok := e.(*proto.Option); ok {
				if !yield(r.wrapOption(v, r.optionContextFile())) {
					return
				}
				continue
			}
			if !yield(r.wrap(e)) {
				return
			}
		}
	}
}
