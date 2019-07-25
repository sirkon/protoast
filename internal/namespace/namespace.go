package namespace

import (
	"text/scanner"

	"github.com/sirkon/protoast/ast"
)

type Namespace interface {
	String() string

	SetPkgName(pkg string) error
	PkgName() string

	WithImport(pkgNamespace Namespace) (Namespace, error)
	WithScope(name string) Namespace
	GetType(name string) ast.Type
	GetService(name string) *ast.Service
	SetNode(name string, def ast.Node, defPos scanner.Position) error

	Finalized() bool
	Finalize()

	getNode(name string) (ast.Node, scanner.Position)
}
