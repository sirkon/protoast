package ast

import "strings"

var _ Type = &Message{}

// Message представление message
type Message struct {
	unique

	File      *File
	ParentMsg *Message

	Name    string
	Fields  []*MessageField
	Types   []Type
	Options []*Option
}

func (*Message) genericType() {}
func (*Message) node()        {}

// AllFields возвращает поля и ветви oneof данного сообщения единым списком
func (m *Message) AllFields() []Field {
	res := make([]Field, 0, len(m.Fields))

	for _, f := range m.Fields {
		switch v := f.Type.(type) {
		case *OneOf:
			for _, b := range v.Branches {
				res = append(res, b)
			}
		default:
			res = append(res, f)
		}
	}

	return res
}

// Field возвращает поле мессаджа по имени. Не производит нормализации по oneof-ам.
func (m *Message) Field(name string) *MessageField {
	for _, field := range m.Fields {
		if field.Name == name {
			return field
		}
	}

	return nil
}

// FieldOneof возвращает поле мессаджа по имени, при этом залезает, если нужно, внутрь oneof-а
func (m *Message) FieldOneof(name string) Field {
	for _, field := range m.Fields {
		if field.Name == name {
			return field
		}

		oo, ok := field.Type.(*OneOf)
		if !ok {
			continue
		}

		for _, branch := range oo.Branches {
			if branch.Name == name {
				return branch
			}
		}
	}

	return nil
}

// Type поиск подтипа по имени
func (m *Message) Type(name string) Type {
	for _, typ := range m.Types {
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

// String референс-имя сообщения, включает в себя название пакета,
// имена родительских сообщений, в пространстве имён которых оно определено.
func (m *Message) String() string {
	var buf strings.Builder
	if m.ParentMsg == nil {
		buf.WriteString(m.File.Package)
	} else {
		buf.WriteString(m.ParentMsg.String())
	}
	buf.WriteByte('.')
	buf.WriteString(m.Name)

	return buf.String()
}

// Message поиск подструктуры по имени.
// Возвращает ошибку ErrorTypeNotFound если такой тип с таким именем не найден.
func (m *Message) Message(name string) (*Message, error) {
	typ := m.Type(name)
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

// Enum поиск вложенного перечисления по имени.
// Возвращает ошибку ErrorTypeNotFound если такой тип с таким именем не найден.
func (m *Message) Enum(name string) (*Enum, error) {
	typ := m.Type(name)
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

// ScanTypes пробежка по внутренним типам данной структуры
func (m *Message) ScanTypes(inspector func(typ Type) bool) {
	if !inspector(m) {
		return
	}

	for _, typ := range m.Types {
		inspector(typ)
	}

	return
}

var (
	_ Unique = &MessageField{}
	_ Field  = &MessageField{}
)

// MessageField представление поля message-а
type MessageField struct {
	unique

	Name     string
	Sequence int
	Type     Type
	Options  []*Option
}

func (m *MessageField) isField() (string, Type, []*Option, int) {
	return m.Name, m.Type, m.Options, m.Sequence
}
