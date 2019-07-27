package ast

var _ Type = &OneOf{}

// OneOf представление для oneof поля message-а
type OneOf struct {
	unique

	ParentMsg *Message

	Name     string
	Branches []*OneOfBranch
}

func (*OneOf) genericType() {}
func (*OneOf) node()        {}

var _ Unique = &OneOfBranch{}

// OneOfBranch представление для ветви
type OneOfBranch struct {
	unique

	Name     string
	Type     Type
	Sequence int
	Options  []*Option
}
