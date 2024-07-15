package namespace

import (
	"github.com/sirkon/protoast/ast"
)

func NewBuilderNaming(naming func(string, string) string) *Builder {
	return &Builder{
		mapping:     map[string]Namespace{},
		scopeNaming: naming,
	}
}

type Builder struct {
	mapping     map[string]Namespace
	scopeNaming func(first, last string) string
	files       map[string]*ast.File
}

func (nb *Builder) get(fileName string, ns Namespace) Namespace {
	res, ok := nb.mapping[fileName]
	if ok {
		return res
	}

	if ns == nil {
		ns = newPlain(fileName, nb)
	}

	nb.mapping[fileName] = ns
	return ns
}

func (nb *Builder) Get(fileName string) Namespace {
	return nb.get(fileName, nil)
}
