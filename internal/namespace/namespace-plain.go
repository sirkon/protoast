package namespace

import (
	"text/scanner"

	"github.com/pkg/errors"

	"github.com/sirkon/prototypes/ast"
)

type nodeTuple struct {
	item ast.Node
	pos  scanner.Position
}

func newPlain(name string, builder *Builder) Namespace {
	return &plain{
		name:    name,
		builder: builder,
		ns:      map[string]nodeTuple{},
	}
}

// plain описание текущего пространства имён
type plain struct {
	name string
	pkg  string

	builder *Builder
	ns      map[string]nodeTuple

	final bool
}

func (n *plain) getType(name string) (ast.Type, scanner.Position) {
	var pos scanner.Position
	res, ok := n.ns[name]
	if !ok {
		return nil, pos
	}

	item := n.GetType(name)
	if item == nil {
		return nil, pos
	}

	return item, res.pos
}

func (n *plain) WithImport(pkgNamespace Namespace) (Namespace, error) {
	return newImport(n, pkgNamespace, n.builder), nil
}

func (n *plain) WithScope(name string) Namespace {
	return newScope(name, n, n.builder)
}

func (n *plain) GetType(name string) ast.Type {
	res, ok := n.ns[name]
	if !ok {
		return nil
	}

	v, ok := res.item.(ast.Type)
	if !ok {
		return nil
	}

	return v
}

func (n *plain) SetType(name string, def ast.Type, defPos scanner.Position) error {
	prev, ok := n.ns[name]
	if ok {
		return errors.Errorf("%s duplicate type %s declartion, the previous one was %s", defPos, name, prev.pos)
	}

	n.ns[name] = nodeTuple{
		item: def,
		pos:  defPos,
	}
	return nil
}

func (n *plain) Finalized() bool { return n.final }
func (n *plain) Finalize()       { n.final = true }
func (n *plain) String() string  { return n.name }

func (n *plain) SetPkgName(pkg string) error {
	n.pkg = pkg
	return nil
}

func (n *plain) PkgName() string { return n.pkg }
