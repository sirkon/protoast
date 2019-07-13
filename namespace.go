package prototypes

import (
	"github.com/sirkon/prototypes/ast"
)

// Namespace представление пространства имён файла
type Namespace interface {
	GetType(name string) ast.Type
}
