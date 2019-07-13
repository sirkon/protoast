package ast

/*
Данный покет содержит описание определяемых структур данных

На данный момент не планируется AST-представление для сервисов, т.к. с ними по сути всех задач проще работать в
предлагаемом github.com/emicklei/proto стиле
*/

// Node обобщённое представление элемента – возможно понадобится добавить AST для всех элементов, не только для типов
type Node interface {
	node()
}

// Type обобщённое представление типа
type Type interface {
	Node
	genericType()
}

// ScalarNode к этому интерефейсу относятся все скалярные типы
type ScalarNode interface {
	Type
	scalar()
}

// Hashable только эти типы могут быть ключами map-ов
type Hashable interface {
	ScalarNode
	hashable()
}

// Options опции полей перечислений и сообщений
type Options map[string]string
