package files

import (
	"fmt"
	"os"

	"github.com/emicklei/proto"
	"github.com/pkg/errors"
)

// New конструктор
func New(importMapping map[string]string) *Files {
	return &Files{
		importMapping: importMapping,
		trees:         map[string]*proto.Proto{},
	}
}

// Files доступ к AST proto-файлов
type Files struct {
	importMapping map[string]string
	trees         map[string]*proto.Proto
}

// File представление готового для разбора файла с данным путём импорта
func (f *Files) File(importPath string) (res *proto.Proto, err error) {
	fileName, ok := f.importMapping[importPath]
	if !ok {
		return nil, UnexpectedImportPath(importPath)
	}

	file, err := os.Open(fileName)
	if err != nil {
		return nil, errors.WithMessage(err, "opening proto file")
	}
	defer func() {
		if cErr := file.Close(); cErr != nil {
			if err == nil {
				err = errors.WithMessage(err, "closing proto file "+file.Name())
			} else {
				panic(fmt.Errorf("closing proto file " + file.Name()))
			}
		}
	}()

	parser := proto.NewParser(file)
	parser.Filename(file.Name())

	ast, err := parser.Parse()
	if err != nil {
		return nil, err
	}

	f.trees[importPath] = ast
	return ast, nil
}
