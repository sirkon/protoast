package namespace

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
