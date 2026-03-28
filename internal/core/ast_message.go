package core

import (
	"iter"

	"github.com/emicklei/proto"

	"github.com/sirkon/protoast/v2/internal/errors"
)

type Message struct {
	isNamedType
	isNodeOptionable

	proto *proto.Message
}

type MessageField struct {
	isFieldNode
	isNodeOptionable

	proto proto.Visitee
}

// Name returns message name.
func (m *Message) Name() string {
	return m.proto.Name
}

func (m *Message) IsExtension() bool {
	return m.proto.IsExtend
}

// Fields returns top level fields of the message.
func (m *Message) Fields(r *Registry) iter.Seq[*MessageField] {
	return func(yield func(*MessageField) bool) {
		for _, element := range m.proto.Elements {
			var field Node
			switch e := element.(type) {
			case *proto.NormalField:
				field = r.wrap(e)
			case *proto.OneOfField:
				field = r.wrap(e)
			case *proto.MapField:
				field = r.wrap(e)
			default:
				continue
			}
			if !yield(field.(*MessageField)) {
				return
			}
		}
	}
}

// Field returns top level field with the given name.
func (m *Message) Field(r *Registry, name string) *MessageField {
	for _, element := range m.proto.Elements {
		var p proto.Visitee
		switch t := element.(type) {
		case *proto.NormalField:
			if t.Name != name {
				continue
			}
			p = t
		case *proto.Oneof:
			if t.Name != name {
				continue
			}
			p = t
		case *proto.MapField:
			if t.Name != name {
				continue
			}
			p = t
		default:
			continue
		}

		return r.wrap(p).(*MessageField)
	}

	return nil
}

// Types returns named types defined at the top level of the message.
func (m *Message) Types(r *Registry) iter.Seq[NamedType] {
	return func(yield func(NamedType) bool) {
		for _, element := range m.proto.Elements {
			var value NamedType
			switch e := element.(type) {
			case *proto.Message:
				if e.IsExtend {
					continue
				}

				value = r.wrap(e).(*Message)
			case *proto.Enum:
				value = r.wrap(e).(*Enum)
			}
			if value != nil {
				if !yield(value) {
					return
				}
			}
		}
	}
}

// Type returns top level named type with given name defined at the top level of the message.
func (m *Message) Type(r *Registry, typename string) NamedType {
	for _, element := range m.proto.Elements {
		switch v := element.(type) {
		case *proto.Message:
			if v.IsExtend {
				continue
			}
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

// Everything returns everything defined at the top level of the message.
func (m *Message) Everything(r *Registry) iter.Seq[Node] {
	return func(yield func(Node) bool) {
		for _, e := range m.proto.Elements {
			if v, ok := e.(*proto.Option); ok {
				if !yield(r.wrapOption(v, r.optionContextMessage())) {
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

// Name returns field name.
func (m *MessageField) Name() string {
	switch p := m.proto.(type) {
	case *proto.NormalField:
		return p.Name
	case *proto.Oneof:
		return p.Name
	case *proto.MapField:
		return p.Name
	default:
		panic(errors.Newf("message came with invalid payload %T", m.proto))
	}
}

// Type returns field type.
func (m *MessageField) Type(r *Registry) (res Type) {
	if v, ok := r.ftcache[m]; ok {
		return v
	}
	defer func() {
		r.ftcache[m] = res
	}()
	switch p := m.proto.(type) {
	case *proto.NormalField:
		normalField := p
		return r.getTypeByName(normalField, normalField.Type)
	case *proto.Oneof:
		return &OneOf{
			proto: p,
		}
	case *proto.MapField:
		return &Map{
			proto: p,
		}
	default:
		panic(errors.Newf("message came with invalid payload %T", m.proto))
	}
}
