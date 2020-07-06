package ast

// Field представление сущности являющейся полем сообщения. Это может быть как непосредственно поле, так и ветвь oneof-а
type Field interface {
	isField()
}
