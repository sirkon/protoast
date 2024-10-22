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

var (
	_ Unique = &OneOfBranch{}
	_ Field  = &OneOfBranch{}
)

func (*OneOf) String() string {
	return "oneof"
}

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

func (o *OneOfBranch) isField() (string, Type, []*Option, int) {
	return o.Name, o.Type, o.Options, o.Sequence
}
