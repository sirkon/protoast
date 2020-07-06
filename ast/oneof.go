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
var _ Field = &OneOfBranch{}

// OneOfBranch представление для ветви
type OneOfBranch struct {
	unique

	Name     string
	Type     Type
	ParentOO *OneOf
	Sequence int
	Options  []*Option
}

func (o *OneOfBranch) isField() (Type, []*Option) {
	return o.Type, o.Options
}
