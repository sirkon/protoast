package namespace

import (
	"sort"
)

// DefaultScopeNaming конструирование имени скоупа по-умолчанию
func DefaultScopeNaming(first, last string) string {
	return first + "::" + last
}

// NewBuilder создание построителя пространства имён
func NewBuilder() *Builder {
	return NewBuilderNaming(DefaultScopeNaming)
}

// NewBuilderNaming построитель пространства имён с настраиваемой функцией получения имени вложенной области видимости
func NewBuilderNaming(naming func(string, string) string) *Builder {
	return &Builder{
		mapping:     map[string]Namespace{},
		scopeNaming: naming,
	}
}

// Builder построение пространств имён
type Builder struct {
	mapping     map[string]Namespace
	scopeNaming func(first, last string) string
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

// Get получение пространства имён для proto-файла с данным путём.
func (nb *Builder) Get(fileName string) Namespace {
	return nb.get(fileName, nil)
}

// Scopes получение списка имён в пространстве.
func (nb *Builder) Scopes() []string {
	var res []string
	for name := range nb.mapping {
		res = append(res, name)
	}
	sort.Strings(res)
	return res
}
