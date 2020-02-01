package ast

var _ Unique = &Comment{}

// Comment представление комментария
type Comment struct {
	unique

	Value string
	Lines []string
}
