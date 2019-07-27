package ast

var _ Unique = &Import{}

// Import представление для импортов
type Import struct {
	unique

	Path string
}
