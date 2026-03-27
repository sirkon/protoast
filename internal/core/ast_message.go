package core

import (
	"iter"

	"github.com/emicklei/proto"

	"github.com/sirkon/protoast/v2/internal/errors"
)

type Message struct {
	isNamedType

	proto *proto.Message
}

type MessageField struct {
	isFieldNode

	payload messageFieldTypeVariant
}

func (m *Message) Name() string {
	return m.proto.Name
}

func (m *Message) Fields() iter.Seq[*MessageField] {
	return func(yield func(*MessageField) bool) {
		for _, element := range m.proto.Elements {
			var field MessageField
			switch t := element.(type) {
			case *proto.NormalField:
				field.payload = asEmickleiNormalField(t)
			case *proto.Oneof:
				field.payload = asEmickleiOneOf(t)
			case *proto.MapField:
				field.payload = asEmickleiMapField(t)
			default:
				continue
			}

			if !yield(&field) {
				return
			}
		}
	}
}

func (m *Message) Field(name string) *MessageField {
	var field MessageField

	for _, element := range m.proto.Elements {
		switch t := element.(type) {
		case *proto.NormalField:
			if t.Name != name {
				continue
			}
			field.payload = asEmickleiNormalField(t)
		case *proto.Oneof:
			if t.Name != name {
				continue
			}
			field.payload = asEmickleiOneOf(t)
		case *proto.MapField:
			if t.Name != name {
				continue
			}
			field.payload = asEmickleiMapField(t)
		default:
			continue
		}

		return &field
	}

	return nil
}

func (m *Message) Types() iter.Seq[NamedType] {
	return func(yield func(NamedType) bool) {
		for _, element := range m.proto.Elements {
			var value NamedType
			switch e := element.(type) {
			case *proto.Message:
				if e.IsExtend {
					continue
				}

				value = &Message{
					proto: e,
				}
			case *proto.Enum:
				value = &Enum{
					proto: e,
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

func (m *Message) Type(typename string) NamedType {
	for _, element := range m.proto.Elements {
		switch v := element.(type) {
		case *proto.Message:
			if v.IsExtend {
				continue
			}
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

func (m *MessageField) Name() string {
	switch p := m.payload.(type) {
	case *isEmickleiNormalField:
		return p.Name
	case *isEmickleiOneOf:
		return p.Name
	case *isEmickleiMapField:
		return p.Name
	default:
		panic(errors.Newf("message came with invalid payload %T", m.payload))
	}
}

func (m *MessageField) Type(r *Registry) Type {
	switch p := m.payload.(type) {
	case *isEmickleiNormalField:
		normalField := p.asProto()
		return r.getTypeByName(normalField, normalField.Type)
	case *isEmickleiOneOf:
		return &OneOf{
			proto: p.asProto(),
		}
	case *isEmickleiMapField:
		return &Map{
			proto: p.asProto(),
		}
	default:
		panic(errors.Newf("message came with invalid payload %T", m.payload))
	}
}
