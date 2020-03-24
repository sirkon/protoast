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

var _ Unique = &MessageField{}

// MessageField представление поля message-а
type MessageField struct {
	unique

	Name     string
	Sequence int
	Type     Type
	Options  []*Option
}
