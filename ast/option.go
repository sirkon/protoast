package ast

var _ Unique = &Option{}

// Option опция
type Option struct {
	unique

	Name      string
	Value     string
	Extension *Extension
}
