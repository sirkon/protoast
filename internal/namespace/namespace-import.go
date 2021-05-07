package namespace

import (
	"strings"
	"text/scanner"

	"github.com/sirkon/protoast/internal/errors"

	"github.com/sirkon/protoast/ast"
)

func newImport(main, importNs Namespace, builder *Builder) Namespace {
	return &nsImport{main: main, imports: []Namespace{importNs}, builder: builder}
}

type nsImport struct {
	main    Namespace
	imports []Namespace
	builder *Builder
}

func (ns *nsImport) getNode(name string) (ast.Node, scanner.Position) {
	var res ast.Node
	var pos scanner.Position

	for _, imp := range ns.imports {
		if imp.PkgName() == ns.PkgName() {
			res, pos = imp.getNode(name)
			if res != nil {
				return res, pos
			}
		}
		if strings.HasPrefix(name, imp.PkgName()+".") {
			nm := name[len(imp.PkgName())+1:]
			res, pos = imp.getNode(nm)
			if res != nil {
				return res, pos
			}
		}
	}

	res, pos = ns.main.getNode(name)
	if res != nil {
		return res, pos
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
	res, _ := ns.getNode(name)
	typeItem, _ := res.(ast.Type)
	return typeItem
}

func (ns *nsImport) GetService(name string) *ast.Service {
	res, _ := ns.getNode(name)
	srv, _ := res.(*ast.Service)
	return srv
}

func (ns *nsImport) SetNode(name string, def ast.Node, defPos scanner.Position) error {
	for _, imp := range ns.imports {
		if imp.PkgName() == ns.PkgName() {
			res, pos := imp.getNode(name)
			if res != nil {
				return errors.Newf("%s duplicate definition of %s which has been previously defined here %s", defPos, name, pos)
			}
		}
	}
	return ns.main.SetNode(name, def, defPos)
}

func (ns *nsImport) Finalized() bool { return ns.main.Finalized() }
func (ns *nsImport) Finalize()       { ns.main.Finalize() }
