package ast

import (
	"path"
	"strings"

	"github.com/sirkon/protoast/internal/errors"
)

// NewPackage конструктор пакета прото-файлов. Если список files пуст возвращается ошибка
// ErrorPackageMissingFiles
func NewPackage(files []*File) (*Package, error) {
	if len(files) == 0 {
		return nil, ErrorPackageMissingFiles{}
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

	res := &Package{
		pkg:   pkg,
		gopkg: gopkg,
		files: files,
	}
	return res, nil
}

// Package представление пакета прото-файлов. Данная структура введена
// скорее для удобства, т.к. понятия пакет в смысле группа файлов с одинаковым
// именем package в protobuf нет — файлы не попадают на трансляцию
// автоматически, а добавляются туда вручную, при трансляции они могут попасть
// в различные целевые каталоги и т.д.
type Package struct {
	pkg   string
	gopkg string
	files []*File
}

// File отдать файл с данным именем. Функция path.Split вызванная на имени файла
// должна возвращать только пустое значение dir.
//
// Правильные названия файлов при вызове: error_codes.proto, marker.proto, atlas
// Неправильные названия файлов: atlas/atlas.proto, /atlas.proto
func (s *Package) File(name string) (*File, error) {
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
func (s *Package) Service(name string) *Service {
	for _, file := range s.files {
		service := file.Service(name)
		if service != nil {
			return service
		}
	}

	return nil
}

// Type поиск типа с данным именем в корневом пространстве имён пакета
func (s *Package) Type(name string) Type {
	for _, file := range s.files {
		typ := file.Type(name)
		if typ != nil {
			return typ
		}
	}

	return nil
}

// Message поиск структуры по имени.
// Возвращает ошибку ast.ErrorTypeNotFound если такой тип с таким именем не найден.
func (s *Package) Message(name string) (*Message, error) {
	typ := s.Type(name)
	if typ == nil {
		return nil, ErrorTypeNotFound(name)
	}

	switch v := typ.(type) {
	case *Message:
		return v, nil
	default:
		return nil, unexpectedType(typ, &Message{})
	}
}

// Enum поиск перечисления по имени
// Возвращает ошибку ast.ErrorTypeNotFound если такой тип с таким именем не найден.
func (s *Package) Enum(name string) (*Enum, error) {
	typ := s.Type(name)
	if typ == nil {
		return nil, ErrorTypeNotFound(name)
	}

	switch v := typ.(type) {
	case *Enum:
		return v, nil
	default:
		return nil, unexpectedType(typ, &Enum{})
	}
}

// ScanTypes пробежка по всем типам пакета, включая и вложенные
func (s *Package) ScanTypes(inspector func(p Type) bool) {
	for _, file := range s.files {
		file.ScanTypes(inspector)
	}
}

// Files получить все файлы текущего пакета
func (s *Package) Files() []*File {
	var res []*File
	for _, file := range s.files {
		res = append(res, file)
	}

	return res
}

// Services получить все сервисы текущего пакета
func (s *Package) Services() []*Service {
	var res []*Service
	for _, file := range s.files {
		res = append(res, file.Services...)
	}

	return res
}

// Types получить все типы текущего пакета лежащие в корне пакета
func (s *Package) Types() []Type {
	var res []Type
	for _, file := range s.files {
		res = append(res, file.Types...)
	}

	return res
}

// Pkg имя proto-пакета
func (s *Package) Pkg() string {
	return s.pkg
}

// GoPkg go-пакета
func (s *Package) GoPkg() string {
	return s.gopkg
}
