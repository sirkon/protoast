package namespace

// NewBuilder создание построителя пространства имён
func NewBuilder() *Builder {
	return &Builder{
		mapping: map[string]Namespace{},
	}
}

// Builder построение пространств имён
type Builder struct {
	mapping map[string]Namespace
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