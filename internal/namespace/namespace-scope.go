package namespace

import (
	"text/scanner"

	"github.com/sirkon/protoast/internal/errors"

	"github.com/sirkon/protoast/ast"
)

func newScope(name string, outer Namespace, builder *Builder) Namespace {
	scopeName := builder.scopeNaming(outer.String(), name)
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

func (s *scope) GetService(name string) *ast.Service {
	res := s.outer.GetService(name)
	if res != nil {
		return res
	}

	return s.current.GetService(name)
}

func (s *scope) getNode(name string) (ast.Node, scanner.Position) {
	res, pos := s.outer.getNode(name)
	if res != nil {
		return res, pos
	}

	return s.current.getNode(name)
}

func (s *scope) WithImport(pkgNamespace Namespace) (Namespace, error) {
	return nil, errors.New("cannot import in inner scopes")
}

func (s *scope) WithScope(name string) Namespace {
	ns := newScope(
		name,
		s,
		s.builder,
	)

	return s.builder.get(ns.String(), ns)
}

func (s *scope) GetType(name string) ast.Type {
	res := s.current.GetType(name)
	if res != nil {
		return res
	}

	return s.outer.GetType(name)
}

func (s *scope) SetNode(name string, def ast.Node, defPos scanner.Position) error {
	return s.current.SetNode(name, def, defPos)
}

func (s *scope) Finalized() bool { return s.final }
func (s *scope) Finalize()       { s.outer.Finalize() }
func (s *scope) PkgName() string { return s.outer.PkgName() }
func (s *scope) String() string  { return s.name }
func (s *scope) SetPkgName(pkg string) error {
	return errors.New("package directive is not allowed in inner scopes")
}
