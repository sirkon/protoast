package ast

// Field представление сущности являющейся полем сообщения. Это может быть как непосредственно поле, так и ветвь oneof-а
type Field interface {
	Unique
	isField() (string, Type, []*Option, int)
}

// FieldName возвращает название поля
func FieldName(f Field) string {
	n, _, _, _ := f.isField()
	return n
}

// FieldType возвращает тип поля
func FieldType(f Field) Type {
	_, t, _, _ := f.isField()
	return t
}

// FieldOptions возвращает опции поля
func FieldOptions(f Field) []*Option {
	_, _, opts, _ := f.isField()
	return opts
}

// FieldIndex возвращает индекс поля
func FieldIndex(f Field) int {
	_, _, _, index := f.isField()
	return index
}
