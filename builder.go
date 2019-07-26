package protoast

import (
	"text/scanner"

	"github.com/emicklei/proto"
	"github.com/pkg/errors"

	"github.com/sirkon/protoast/ast"
	"github.com/sirkon/protoast/internal/namespace"
)

func NewBuilder(imports Files, errProcessing func(error)) *Builder {
	return NewBuilderWithCustomNaming(imports, errProcessing, DefaultNaming)
}

func DefaultNaming(first, last string) string {
	return first + "::" + last
}

func NewBuilderWithCustomNaming(imports Files, errProcessing func(error), scopeNaming func(string, string) string) *Builder {
	return &Builder{
		protos: &protos{
			files: imports,
			trees: map[string]*proto.Proto{},
		},
		nsBuilder:     namespace.NewBuilderNaming(scopeNaming),
		asts:          map[string]*ast.File{},
		finalNses:     map[string]Namespace{},
		errProcessing: errProcessing,
	}

}

type Builder struct {
	protos        *protos
	nsBuilder     *namespace.Builder
	asts          map[string]*ast.File
	finalNses     map[string]Namespace
	errProcessing func(error)
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
		errors: errChan,
	}

	v := typesVisitor{
		file:   pf.file,
		ns:     ns,
		nss:    s,
		errors: errChan,
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

	s.asts[importPath] = v.file
	s.finalNses[importPath] = v.ns

	switch errCount {
	case 0:
		return nil
	case 1:
		return errors.New("an error occured")
	default:
		return errors.Errorf("%d errors occured", errCount)
	}
}
