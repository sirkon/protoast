package ast

var _ Type = &OneOf{}

// OneOf представление oneof-поля в message-ах
type OneOf struct {
	Name     string
	Branches []OneOfBranch
}

func (*OneOf) genericType() {}
func (*OneOf) node()        {}

// OneOfBranch представление ветвления в oneof
type OneOfBranch struct {
	Name     string
	Type     Type
	Sequence int
	Options  Options
}
