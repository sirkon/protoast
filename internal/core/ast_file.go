package core

import (
	"iter"

	"github.com/emicklei/proto"
)

type File struct {
	isNode

	proto *proto.Proto
}

func (f *File) Name() string {
	return f.proto.Filename
}

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
func (f *File) Messages() iter.Seq[*Message] {
	return func(yield func(*Message) bool) {
		for _, element := range f.proto.Elements {
			v, ok := element.(*proto.Message)
			if ok && !v.IsExtend {
				if !yield(&Message{
					proto: v,
				}) {
					return
				}
			}
		}
	}
}

// Message returns a message with the given name.
func (f *File) Message(name string) *Message {
	for _, element := range f.proto.Elements {
		v, ok := element.(*proto.Message)
		if ok && v.Name == name {
			return &Message{
				proto: v,
			}
		}
	}

	return nil
}

func (f *File) Enums() iter.Seq[*Enum] {
	return func(yield func(*Enum) bool) {
		for _, element := range f.proto.Elements {
			v, ok := element.(*proto.Enum)
			if ok {
				if !yield(&Enum{
					proto: v,
				}) {
					return
				}
			}
		}
	}
}

func (f *File) Enum(name string) *Enum {
	for _, element := range f.proto.Elements {
		v, ok := element.(*proto.Enum)
		if ok {
			if v.Name == name {
				return &Enum{
					proto: v,
				}
			}
		}
	}

	return nil
}

func (f *File) Types() iter.Seq[NamedType] {
	return func(yield func(levelType NamedType) bool) {
		for _, element := range f.proto.Elements {
			var value NamedType
			switch v := element.(type) {
			case *proto.Message:
				if v.IsExtend {
					continue
				}
				value = &Message{
					proto: v,
				}
			case *proto.Enum:
				value = &Enum{
					proto: v,
				}
			}
			if value != nil {
				if !yield(value) {
					return
				}
			}
		}
	}
}

func (f *File) Type(typename string) NamedType {
	for _, element := range f.proto.Elements {
		switch v := element.(type) {
		case *proto.Message:
			if v.Name != typename {
				continue
			}

			return &Message{
				proto: v,
			}

		case *proto.Enum:
			if v.Name != typename {
				continue
			}

			return &Enum{
				proto: v,
			}
		}
	}

	return nil
}
