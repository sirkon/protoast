package ast

type Node interface {
	node()
}

type Type interface {
	Node
	genericType()
}

type ScalarNode interface {
	Type
	scalar()
}

type Hashable interface {
	ScalarNode
	hashable()
}

type Options map[string]string
