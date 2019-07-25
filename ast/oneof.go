package ast

var _ Type = &OneOf{}

type OneOf struct {
	ParentMsg	*Message
	Name		string
	Branches	[]OneOfBranch
}

func (*OneOf) genericType()	{}
func (*OneOf) node()		{}

type OneOfBranch struct {
	Name		string
	Type		Type
	Sequence	int
	Options		Options
}
