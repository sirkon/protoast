package namespace

import (
	"fmt"
	"text/scanner"

	"github.com/pkg/errors"

	"github.com/sirkon/prototypes/ast"
)

func newScope(name string, outer Namespace, builder *Builder) Namespace {
	scopeName := outer.String() + ":scope=" + name
	ns := &scope{
		name:    scopeName,
		outer:   outer,
		current: newPlain(scopeName, builder),
		builder: builder,
	}

	return builder.get(scopeName, ns)
}

type scope struct {
	name    string
	outer   Namespace
	current Namespace
	builder *Builder

	final bool
}

func (s *scope) getType(name string) (ast.Type, scanner.Position) {
	res, pos := s.outer.getType(name)
	if res != nil {
		return res, pos
	}

	return s.current.getType(name)
}

func (s *scope) WithImport(pkgNamespace Namespace) (Namespace, error) {
	return nil, fmt.Errorf("cannot import in inner scopes")
}

func (s *scope) WithScope(name string) Namespace {
	ns := newScope(
		name,
		s,
		s.builder,
	)

	return s.builder.get(s.String(), ns)
}

func (s *scope) GetType(name string) ast.Type {
	res := s.outer.GetType(name)
	if res != nil {
		return res
	}

	return s.current.GetType(name)
}

func (s *scope) SetType(name string, def ast.Type, defPos scanner.Position) error {
	return s.current.SetType(name, def, defPos)
}

func (s *scope) Finalized() bool { return s.final }
func (s *scope) Finalize()       { s.outer.Finalize() }
func (s *scope) PkgName() string { return s.outer.PkgName() }
func (s *scope) String() string  { return s.name }
func (s *scope) SetPkgName(pkg string) error {
	return errors.New("package directive is not allowed in inner scopes")
}
