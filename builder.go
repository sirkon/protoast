package protoast

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"reflect"
	"strings"
	"text/scanner"

	"github.com/emicklei/proto"
	"github.com/pkg/errors"

	"github.com/sirkon/protoast/ast"
	"github.com/sirkon/protoast/internal/namespace"
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
	errCount      int
}

// SameDirProtos Отдать список все файлы из директории содержащей указанный файл, включая данный.
func (s *Builder) SameDirProtos(file *ast.File) ([]*ast.File, error) {
	abs, err := s.protos.files.Abs(file.Name)
	if err != nil {
		return nil, errors.WithMessagef(err, "get absolute path name for %s", file.Name)
	}
	dir, _ := filepath.Split(abs)
	lst, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, errors.WithMessagef(err, "read directory %s of file %s", dir, file.Name)
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
				panic(errors.Errorf("no position was set for %T.%s", k, val.Type().Field(i).Name))
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
		return nil, errors.WithMessagef(err, "getting namespace for %s", importPath)
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
		return nil, nil, errors.WithMessagef(err, "getting file %s", importPath)
	}

	if err = s.processFile(ns, protoItem, importPath); err != nil {
		return nil, nil, errors.WithMessagef(err, "processing file %s", importPath)
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
		return errors.Errorf("%d errors occured", s.errCount)
	}
}

func (s *Builder) processError(err error) {
	if err != nil {
		s.errCount++
	}
	s.errProcessing(err)
}
