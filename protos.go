package protoast

import (
	"bytes"
	"fmt"

	"github.com/emicklei/proto"
	"github.com/pkg/errors"
)

type protos struct {
	files Files
	trees map[string]*proto.Proto
}

func (p *protos) fileProto(importPath string) (res *proto.Proto, err error) {
	fileData, err := p.files.File(importPath)
	if err != nil {
		return nil, errors.WithMessagef(err, "getting %s file data", importPath)
	}
	absPath, ok := p.files.Abs(importPath)
	if !ok {
		return nil, errors.Errorf("couldn't find absolute path for import %s", importPath)
	}
	file := bytes.NewBuffer(fileData)

	parser := proto.NewParser(file)
	parser.Filename(absPath)

	ast, err := parser.Parse()
	if err != nil {
		return nil, err
	}

	p.trees[importPath] = ast
	return ast, nil
}

var _ error = unexpectedImportPath("")

type unexpectedImportPath string

func (f unexpectedImportPath) Error() string {
	return fmt.Sprintf("unexpected import %s", string(f))
}
