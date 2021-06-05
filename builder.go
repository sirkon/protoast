package protoast

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"reflect"
	"strings"
	"text/scanner"

	"github.com/emicklei/proto"
	"github.com/sirkon/protoast/internal/errors"
	"github.com/sirkon/protoast/internal/namespace"

	"github.com/sirkon/protoast/ast"
)

// NewBuilder конструктор построителя AST-представления
func NewBuilder(imports Files, errProcessing func(error)) *Builder {
	return newBuilderWithCustomNaming(imports, errProcessing, defaultNaming)
}

func defaultNaming(first, last string) string {
	return first + "." + last
}

func newBuilderWithCustomNaming(imports Files, errProcessing func(error), scopeNaming func(string, string) string) *Builder {
	return &Builder{
		protos: &protos{
			files: imports,
			trees: map[string]*proto.Proto{},
		},
		nsBuilder:     namespace.NewBuilderNaming(scopeNaming),
		asts:          map[string]*ast.File{},
		finalNses:     map[string]Namespace{},
		errProcessing: errProcessing,
		comments:      map[string]*ast.Comment{},
		positions:     map[string]scanner.Position{},
		uniqueContext: ast.UniqueContext{},
		imports:       imports,
	}
}

// StringRef создаёт ссылку на строку
func StringRef(value string) *string {
	return &value
}

// Builder построитель структурированной информации
type Builder struct {
	protos        *protos
	nsBuilder     *namespace.Builder
	asts          map[string]*ast.File
	finalNses     map[string]Namespace
	errProcessing func(error)

	comments      map[string]*ast.Comment
	positions     map[string]scanner.Position
	uniqueContext ast.UniqueContext
	imports       Files
	errCount      int
}

// SameDirProtos Отдать список все файлы из директории содержащей указанный файл, включая данный.
//
// Deprecated: лучше использовать SamePackage
func (s *Builder) SameDirProtos(file *ast.File) ([]*ast.File, error) {
	abs, err := s.protos.files.Abs(file.Name)
	if err != nil {
		return nil, errors.Wrapf(err, "get absolute path name for the file")
	}
	dir, _ := filepath.Split(abs)
	lst, err := os.ReadDir(dir)
	if err != nil {
		return nil, errors.Wrap(err, "read file's directory "+dir)
	}

	var res []*ast.File
	relDir, _ := filepath.Split(file.Name)
	for _, info := range lst {
		if info.IsDir() {
			continue
		}
		if !strings.HasSuffix(info.Name(), ".proto") {
			continue
		}
		_, base := filepath.Split(info.Name())
		newFile, err := s.AST(filepath.Join(relDir, base))
		if err != nil {
			return nil, err
		}
		res = append(res, newFile)
	}

	return res, nil
}

// SamePackage отдать все файлы пакета для данного файла
func (s *Builder) SamePackage(file *ast.File) (*ast.Package, error) {
	files, err := s.SameDirProtos(file)
	if err != nil {
		return nil, errors.Wrap(err, "get list of proto files from the same directory")
	}

	res, err := ast.NewPackage(files)
	if err != nil {
		return nil, errors.Wrap(err, "form package of the given file")
	}

	return res, nil
}

// Package отдать все файлы пакета для proto-файлов из данной директории.
// Будет работать только для резолвера-функции. Если в директории нет файлов
// то будет возвращена ошибка ast.ErrorPackageMissingFiles
func (s *Builder) Package(dir string) (*ast.Package, error) {
	abs, err := s.imports.Abs(dir)
	if err != nil {
		return nil, errors.Wrap(err, "get absolute path for the directory")
	}

	readDir, err := os.ReadDir(abs)
	if err != nil {
		return nil, errors.Wrap(err, "list files in the directory")
	}

	var files []*ast.File
	for _, entry := range readDir {
		if entry.IsDir() {
			continue
		}

		if !strings.HasSuffix(entry.Name(), ".proto") {
			continue
		}

		file, err := s.AST(path.Join(dir, entry.Name()))
		if err != nil {
			return nil, errors.Wrap(err, "get AST for "+entry.Name())
		}

		files = append(files, file)
	}

	res, err := ast.NewPackage(files)
	if err != nil {
		return nil, errors.Wrap(err, "form proto-package for directory "+dir)
	}

	return res, nil
}

// Comment возвращает комментарий для сущности реализующей Unique
func (s *Builder) Comment(k ast.Unique) *ast.Comment {
	return s.comments[ast.GetUnique(k)]
}

// CommentField возвращает комментарий для поля сущности реализующей Unique
func (s *Builder) CommentField(k ast.Unique, fieldAddr interface{}) *ast.Comment {
	return s.comments[ast.GetFieldKey(k, fieldAddr)]
}

// Position возвращает позицию данного Unique
func (s *Builder) Position(k ast.Unique) scanner.Position {
	res, ok := s.positions[ast.GetUnique(k)]
	var name string
	switch v := k.(type) {
	case *ast.Message:
		name = fmt.Sprintf("message %s::%s.%s %s", v.File.Name, v.File.Package, v.Name, ast.GetUnique(k))
	case *ast.Enum:
		name = fmt.Sprintf("enum %s.%s", v.File.Package, v.Name)
	case *ast.Service:
		name = fmt.Sprintf("service %s.%s", v.File.Package, v.Name)
	case *ast.Extension:
		name = fmt.Sprintf("extension %s.%s", v.File.Package, v.Name)
	default:
		name = fmt.Sprintf("%T", v)
	}
	if !ok {
		panic(errors.New(name))
	}
	return res
}

// PositionField возвращает позицию данного для поля данного Uniq
func (s *Builder) PositionField(k ast.Unique, fieldAddr interface{}) scanner.Position {
	res, ok := s.positions[ast.GetFieldKey(k, fieldAddr)]
	if !ok {
		val := reflect.ValueOf(k).Elem()
		addr := reflect.ValueOf(fieldAddr).Pointer()
		for i := 0; i < val.NumField(); i++ {
			if val.Field(i).Addr().Pointer() == addr {
				panic(errors.Newf("no position was set for %T.%s", k, val.Type().Field(i).Name))
			}
		}
		panic("must not be here")
	}
	return res
}

// AST представление для файла с данным относительным путём
func (s *Builder) AST(importPath string) (*ast.File, error) {
	_, res, err := s.get(importPath)
	return res, err
}

// Namespace получение пространства имён для данного файла
func (s *Builder) Namespace(importPath string) (Namespace, error) {
	_, _, err := s.get(importPath)
	if err != nil {
		return nil, errors.Wrap(err, "get namespace for the file")
	}
	return s.finalNses[importPath], nil
}

// Scope пространство имён для данной области видимости
func (s *Builder) Scope(access string) Namespace {
	return s.nsBuilder.Get(access)
}

// Proto представление для прохода визиторами из github.com/emicklei/proto
func (s *Builder) Proto(importPath string) (*proto.Proto, error) {
	return s.protos.fileProto(importPath)
}

func (s *Builder) get(importPath string) (namespace.Namespace, *ast.File, error) {
	ns := s.nsBuilder.Get(importPath)
	if ns.Finalized() {
		return ns, s.asts[importPath], nil
	}
	protoItem, err := s.protos.fileProto(importPath)
	if err != nil {
		return nil, nil, errors.Wrap(err, "get file "+importPath)
	}

	if err = s.processFile(ns, protoItem, importPath); err != nil {
		return nil, nil, errors.Wrap(err, "process file "+importPath)
	}
	ns.Finalize()

	return ns, s.asts[importPath], nil
}

func (s *Builder) processFile(ns namespace.Namespace, p *proto.Proto, importPath string) error {
	errChan := make(chan error)
	pf := prefetcher{
		file: &ast.File{
			Name: importPath,
		},
		ns:     ns,
		nss:    s,
		errors: s.processError,
	}

	v := typesVisitor{
		file:   pf.file,
		ns:     ns,
		nss:    s,
		errors: s.errProcessing,
		enumCtx: struct {
			item        *ast.Enum
			prevField   map[string]scanner.Position
			prevInteger map[int]scanner.Position
		}{},
		msgCtx: struct {
			onMsg       bool
			item        *ast.Message
			prevField   map[string]scanner.Position
			prevInteger map[int]scanner.Position
		}{},
		oneOf: nil,
	}

	p.Accept(&pf)
	p.Accept(&v)
	close(errChan)

	s.asts[importPath] = v.file
	s.finalNses[importPath] = v.ns

	switch s.errCount {
	case 0:
		return nil
	case 1:
		return errors.New("an error occured")
	default:
		return errors.Newf("%d errors occured", s.errCount)
	}
}

func (s *Builder) processError(err error) {
	if err != nil {
		s.errCount++
	}
	s.errProcessing(err)
}
