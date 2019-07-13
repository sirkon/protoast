package namespace

import (
	"text/scanner"

	"github.com/sirkon/prototypes/ast"
)

// Namespace работа с
type Namespace interface {
	String() string

	SetPkgName(pkg string) error
	PkgName() string

	WithImport(pkgNamespace Namespace) (Namespace, error)
	WithScope(name string) Namespace
	GetType(name string) ast.Type
	SetType(name string, def ast.Type, defPos scanner.Position) error

	Finalized() bool
	Finalize()

	getType(name string) (ast.Type, scanner.Position)
}
