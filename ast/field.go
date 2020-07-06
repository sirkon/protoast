package ast

// Field представление сущности являющейся полем сообщения. Это может быть как непосредственно поле, так и ветвь oneof-а
type Field interface {
	isField() (Type, []*Option)
}

// FieldType возвращает тип поля
func FieldType(f Field) Type {
	t, _ := f.isField()
	return t
}

// FieldOptions возвращает опции поля
func FieldOptions(f Field) []*Option {
	_, opts := f.isField()
	return opts
}
