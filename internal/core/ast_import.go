package core

import (
	"text/scanner"

	"github.com/emicklei/proto"
)

type Import struct {
	proto *proto.Import
}

// Path returns a path imported.
func (i *Import) Path() string {
	return i.proto.Filename
}

var _ Node = new(Import)

func (i *Import) nodeProto() proto.Visitee { return i.proto }
func (i *Import) pos() scanner.Position    { return i.proto.Position }
