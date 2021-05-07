package protoast

import (
	"io/ioutil"
	"path/filepath"

	"github.com/sirkon/protoast/internal/errors"
)

// Files абстракция для работы с файлами
type Files interface {
	File(name string) ([]byte, error)
	Abs(name string) (string, error)
}

// NewFiles отдаёт реализацию Files построенную на готовом соответствии
//
// Deprecated: сущность имеет малый смысл, только хорошее название забрала
func NewFiles(mapping map[string]string) Files {
	return &files{
		mapping: mapping,
	}
}

type files struct {
	mapping map[string]string
}

// Abs для реализации Files
func (f *files) Abs(name string) (string, error) {
	res, ok := f.mapping[name]
	if !ok {
		return "", errors.New("no such fie")
	}
	return res, nil
}

// File для реализации Files
func (f *files) File(path string) ([]byte, error) {
	absPath, ok := f.mapping[path]
	if !ok {
		return nil, errors.New("no such file")
	}

	data, err := ioutil.ReadFile(absPath)
	if err != nil {
		return nil, errors.Wrap(err, "read file")
	}

	return data, nil
}

// NewFilesViaResolver отдаёт реализацию Files опирающуюся на функцию-резолвер
func NewFilesViaResolver(resolver func(string) (string, error)) Files {
	return backResolver(resolver)
}

type backResolver func(name string) (string, error)

// File для реализации Files
func (b backResolver) File(path string) ([]byte, error) {
	absPath, err := b.Abs(path)
	if err != nil {
		return nil, errors.Wrap(err, "compute absolute path")
	}
	data, err := ioutil.ReadFile(absPath)
	if err != nil {
		return nil, errors.Wrapf(err, "read file")
	}
	return data, nil
}

// Abs для реализации Files
func (b backResolver) Abs(path string) (string, error) {
	res, err := b(path)
	if err != nil {
		return "", errors.Wrap(err, "resolve path")
	}
	absPath, err := filepath.Abs(res)
	if err != nil {
		return "", errors.Wrap(err, "compute absolute path")
	}
	return absPath, nil
}
