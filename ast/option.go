package ast

var _ Unique = &Option{}

// Option опция поля
type Option struct {
	unique

	Name  string
	Value string
}
