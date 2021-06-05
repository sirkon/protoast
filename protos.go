package protoast

import (
	"bytes"
	"fmt"

	"github.com/emicklei/proto"
	"github.com/sirkon/protoast/internal/errors"
)

type protos struct {
	files Files
	trees map[string]*proto.Proto
}

func (p *protos) fileProto(importPath string) (res *proto.Proto, err error) {
	fileData, err := p.files.File(importPath)
	if err != nil {
		return nil, errors.Wrap(err, "read proto file data")
	}
	absPath, err := p.files.Abs(importPath)
	if err != nil {
		return nil, errors.Wrap(err, "compute absolute path for the file")
	}
	file := bytes.NewBuffer(fileData)

	parser := proto.NewParser(file)
	parser.Filename(absPath)

	ast, err := parser.Parse()
	if err != nil {
		return nil, errors.Wrap(err, "parse file")
	}

	p.trees[importPath] = ast
	return ast, nil
}

var _ error = unexpectedImportPath("")

type unexpectedImportPath string

func (f unexpectedImportPath) Error() string {
	return fmt.Sprintf("unexpected import %s", string(f))
}
