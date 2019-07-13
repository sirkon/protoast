package ast

var _ Type = &Message{}

// Message определение структуры типа message
type Message struct {
	Name   string
	Fields []MessageField
}

func (*Message) genericType() {}
func (*Message) node()        {}

// MessageField описание поля message-а
type MessageField struct {
	Name     string
	Sequence int
	Type     Type
	Options  Options
}
