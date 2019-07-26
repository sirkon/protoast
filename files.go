package protoast

import (
	"io/ioutil"
	"os"

	"github.com/pkg/errors"
)

// Files абстракция для работы с файлами
type Files interface {
	File(name string) ([]byte, error)
}

// NewFiles отдаёт реализацию Files
func NewFiles(mapping map[string]string) Files {
	return &files{
		mapping: mapping,
	}
}

type files struct {
	mapping map[string]string
}

// fileProto ...
func (f *files) File(path string) ([]byte, error) {
	absPath, ok := f.mapping[path]
	if !ok {
		return nil, errors.WithMessagef(os.ErrNotExist, "file %s", path)
	}

	data, err := ioutil.ReadFile(absPath)
	if err != nil {
		return nil, errors.WithMessagef(err, "reading %s", path)
	}

	return data, nil
}
