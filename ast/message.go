package ast

var _ Type = &Message{}

type Message struct {
	File		*File
	ParentMsg	*Message
	Name		string
	Fields		[]MessageField
}

func (*Message) genericType()	{}
func (*Message) node()		{}

type MessageField struct {
	Name		string
	Sequence	int
	Type		Type
	Options		Options
}
