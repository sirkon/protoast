package protoast

import (
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/pkg/errors"
)

// Files абстракция для работы с файлами
type Files interface {
	File(name string) ([]byte, error)
	Abs(name string) (string, error)
}

// NewFiles отдаёт реализацию Files построенную на готовом соответствии
func NewFiles(mapping map[string]string) Files {
	return &files{
		mapping: mapping,
	}
}

type files struct {
	mapping map[string]string
}

func (f *files) Abs(name string) (string, error) {
	res, ok := f.mapping[name]
	if !ok {
		return "", errors.Errorf("cannot resolve %s", name)
	}
	return res, nil
}

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

// NewFilesViaResolver отдаёт реализацию Files опирающуюся на функцию-резолвер
func NewFilesViaResolver(resolver func(string) (string, error)) Files {
	return backResolver(resolver)
}

type backResolver func(name string) (string, error)

func (b backResolver) File(path string) ([]byte, error) {
	absPath, err := b.Abs(path)
	if err != nil {
		return nil, err
	}
	data, err := ioutil.ReadFile(absPath)
	if err != nil {
		return nil, errors.WithMessagef(err, "reading %s", path)
	}
	return data, nil
}

func (b backResolver) Abs(path string) (string, error) {
	res, err := b(path)
	if err != nil {
		return "", errors.WithMessagef(err, "resolving %s", path)
	}
	absPath, err := filepath.Abs(res)
	if err != nil {
		return "", errors.WithMessagef(err, "expanding %s - resolved %s", res, path)
	}
	return absPath, nil
}
