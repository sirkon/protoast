package ast

// Node представление базовой ноды
type Node interface {
	Unique
	node()
}

// Type представление типа
type Type interface {
	Node
	genericType()
}

// ScalarNode скалярные типы
type ScalarNode interface {
	Type
	scalar()
	equivalent(v ScalarNode) bool
}

// Hashable типы могущие являться ключами словарей
type Hashable interface {
	ScalarNode
	hashable()
}
