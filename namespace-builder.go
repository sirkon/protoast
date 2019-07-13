package prototypes

import (
	"github.com/emicklei/proto"
	"github.com/pkg/errors"

	"github.com/sirkon/prototypes/internal/files"
	"github.com/sirkon/prototypes/internal/namespace"
)

// NewNamespaces конструктор вычислителя пространств имён для данного проекта proto-файлов.
// Входные параметры:
//    imports:        <путь импортирования> => <абсолютный путь импортируемого файла>
//    errProcessing:  функция обработки ошибок разбора
func NewNamespaces(imports map[string]string, errProcessing func(error)) *Namespaces {
	return &Namespaces{
		files:         files.New(imports),
		nsBuilder:     namespace.NewBuilder(),
		errProcessing: errProcessing,
	}
}

// Namespaces вычисление пространств имён для файлов
type Namespaces struct {
	files         *files.Files
	nsBuilder     *namespace.Builder
	errProcessing func(error)
}

// Get получение пространства имён для данного файла
func (s *Namespaces) Get(importPath string) (Namespace, error) {
	return s.get(importPath)
}

// Proto получение готового к обходу представления файла предоставляемого библиотекой gitub.com/emicklei/proto
func (s *Namespaces) Proto(importPath string) (*proto.Proto, error) {
	return s.files.File(importPath)
}

func (s *Namespaces) get(importPath string) (namespace.Namespace, error) {
	ns := s.nsBuilder.Get(importPath)
	if ns.Finalized() {
		return ns, nil
	}

	protoItem, err := s.files.File(importPath)
	if err != nil {
		return nil, errors.WithMessagef(err, "getting file %s", importPath)
	}

	if err = s.processFile(ns, protoItem); err != nil {
		return nil, errors.WithMessagef(err, "processing file %s", importPath)
	}
	ns.Finalize()

	return ns, nil
}

func (s *Namespaces) processFile(ns namespace.Namespace, p *proto.Proto) error {
	errChan := make(chan error)
	pf := prefetcher{
		ns:     ns,
		nss:    s,
		errors: errChan,
	}

	v := typesVisitor{
		ns:     ns,
		nss:    s,
		errors: errChan,
	}

	var errCount int
	go func() {
		for err := range errChan {
			s.errProcessing(err)
			errCount++
		}
	}()

	p.Accept(&pf)
	p.Accept(&v)
	close(errChan)

	switch errCount {
	case 0:
		return nil
	case 1:
		return errors.New("an error occured")
	default:
		return errors.Errorf("%d errors occured", errCount)
	}
}
