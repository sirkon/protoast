package ast

// Field представление сущности являющейся полем сообщения. Это может быть как непосредственно поле, так и ветвь oneof-а
type Field interface {
	isField() (string, Type, []*Option)
}

// FieldName возвращает название поля
func FieldName(f Field) string {
	n, _, _ := f.isField()
	return n
}

// FieldType возвращает тип поля
func FieldType(f Field) Type {
	_, t, _ := f.isField()
	return t
}

// FieldOptions возвращает опции поля
func FieldOptions(f Field) []*Option {
	_, _, opts := f.isField()
	return opts
}
