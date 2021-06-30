package ast

// Equal проверка равенства нод Unique
func Equal(x, y Unique) bool {
	return GetUnique(x) == GetUnique(y)
}

// SamePrimitive проверка равенства примитивных типов
func SamePrimitive(x, y Type) bool {
	xv, ok := x.(ScalarNode)
	if !ok {
		return false
	}

	yv, ok := x.(ScalarNode)
	if !ok {
		return false
	}

	return xv.equivalent(yv)
}
