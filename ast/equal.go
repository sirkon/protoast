package ast

// Equal проверка равенства нод Unique
func Equal(x, y Unique) bool {
	return GetUnique(x) == GetUnique(y)
}
