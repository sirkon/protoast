package fileset

import (
	"path"
	"strings"

	"github.com/sirkon/protoast/internal/errors"

	"github.com/sirkon/protoast/ast"
)

// NewFileSet конструктор пакета прото-файлов
func NewFileSet(files []*ast.File) (*FileSet, error) {
	if len(files) == 0 {
		return nil, errors.New("at least one file is required to form a set of proto files")
	}

	var pkg string
	var gopkg string
	for _, file := range files {
		if pkg == "" || gopkg == "" {
			pkg = file.Package
			gopkg = file.GoPkg
			continue
		}

		if pkg != file.Package {
			return nil, errors.Newf(
				"cannot form a package of proto files related to different proto packages %s and %s",
				pkg,
				file.Package,
			)
		}

		if gopkg != file.GoPkg {
			return nil, errors.Newf(
				"cannot form a package of proto files related to different go packages %s and %s",
				gopkg,
				file.GoPkg,
			)
		}
	}

	res := &FileSet{
		pkg:   pkg,
		gopkg: gopkg,
		files: files,
	}
	return res, nil
}

// FileSet представление пакета прото-файлов
type FileSet struct {
	pkg   string
	gopkg string
	files []*ast.File
}

// File отдать файл с данным именем. Функция path.Split вызванная на имени файла
// должна возвращать только пустое значение dir.
//
// Правильные названия файлов при вызове: error_codes.proto, marker.proto, atlas
// Неправильные названия файлов: atlas/atlas.proto, /atlas.proto
func (s *FileSet) File(name string) (*ast.File, error) {
	if name == "" {
		return nil, errors.New("empty name")
	}

	dir, base := path.Split(name)
	if dir != "" {
		return nil, errors.Newf("%s: package file name must not have directory parts", name)
	}

	if base == "" {
		return nil, errors.Newf("%s: invalid package file name", name)
	}

	if !strings.HasSuffix(name, ".proto") {
		return nil, errors.Newf("%s: package file name missing .proto suffix", name)
	}

	for _, file := range s.files {
		_, base := path.Split(file.Name)
		if base == name {
			return file, nil
		}
	}

	return nil, errors.Newf("%s: file not found in %s package", name, s.pkg)
}

// Service поиск сервиса с данным именем в пространстве имён пакета
func (s *FileSet) Service(name string) *ast.Service {
	for _, file := range s.files {
		service := file.Service(name)
		if service != nil {
			return service
		}
	}

	return nil
}

// Type поиск типа с данным именем в корневом пространстве имён пакета
func (s *FileSet) Type(name string) ast.Type {
	for _, file := range s.files {
		typ := file.Type(name)
		if typ != nil {
			return typ
		}
	}

	return nil
}

// ScanTypes пробежка по всем типам пакета, включая и вложенные
func (s *FileSet) ScanTypes(inspector func(p ast.Type) bool) {
	for _, file := range s.files {
		file.ScanTypes(inspector)
	}
}

// Files получить все файлы текущего пакета
func (s *FileSet) Files() []*ast.File {
	var res []*ast.File
	for _, file := range s.files {
		res = append(res, file)
	}

	return res
}

// Services получить все сервисы текущего пакета
func (s *FileSet) Services() []*ast.Service {
	var res []*ast.Service
	for _, file := range s.files {
		res = append(res, file.Services...)
	}

	return res
}

// Types получить все типы текущего пакета лежащие в корне пакета
func (s *FileSet) Types() []ast.Type {
	var res []ast.Type
	for _, file := range s.files {
		res = append(res, file.Types...)
	}

	return res
}
