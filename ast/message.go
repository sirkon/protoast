package ast

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

var _ Unique = &MessageField{}
var _ Field = &MessageField{}

// MessageField представление поля message-а
type MessageField struct {
	unique

	Name     string
	Sequence int
	Type     Type
	Options  []*Option
}

func (*MessageField) isField() {}
