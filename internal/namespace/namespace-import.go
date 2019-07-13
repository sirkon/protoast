package namespace

import (
	"strings"
	"text/scanner"

	"github.com/pkg/errors"

	"github.com/sirkon/prototypes/ast"
)

func newImport(main, importNs Namespace, builder *Builder) Namespace {
	return &nsImport{main: main, imports: []Namespace{importNs}, builder: builder}
}

type nsImport struct {
	main    Namespace
	imports []Namespace
	builder *Builder
}

func (ns *nsImport) getType(name string) (ast.Type, scanner.Position) {
	res, pos := ns.main.getType(name)
	if res != nil {
		return res, pos
	}

	for _, imp := range ns.imports {
		if imp.PkgName() == ns.PkgName() {
			res, pos = imp.getType(name)
			if res != nil {
				return res, pos
			}
		}
		if strings.HasPrefix(name, imp.PkgName()+".") {
			nm := name[len(imp.PkgName())+1:]
			res, pos = imp.getType(nm)
			if res != nil {
				return res, pos
			}
		}
	}

	return nil, pos
}

func (ns *nsImport) String() string              { return ns.main.String() }
func (ns *nsImport) SetPkgName(pkg string) error { return ns.main.SetPkgName(pkg) }
func (ns *nsImport) PkgName() string             { return ns.main.PkgName() }

func (ns *nsImport) WithImport(pkgNamespace Namespace) (Namespace, error) {
	ns.imports = append(ns.imports, pkgNamespace)
	return ns, nil
}

func (ns *nsImport) WithScope(name string) Namespace { return newScope(name, ns, ns.builder) }

func (ns *nsImport) GetType(name string) ast.Type {
	res, _ := ns.getType(name)
	return res
}

func (ns *nsImport) SetType(name string, def ast.Type, defPos scanner.Position) error {
	for _, imp := range ns.imports {
		if imp.PkgName() == ns.PkgName() {
			res, pos := imp.getType(name)
			if res != nil {
				return errors.Errorf("%s duplicate definition of %s which has been previously defined here %s", defPos, name, pos)
			}
		}
	}
	return ns.main.SetType(name, def, defPos)
}

func (ns *nsImport) Finalized() bool { return ns.main.Finalized() }
func (ns *nsImport) Finalize()       { ns.main.Finalize() }
