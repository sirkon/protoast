package core

import (
	"github.com/emicklei/proto"
)

type Import struct {
	isNode

	proto *proto.Import
}

// Path returns a path imported.
func (i *Import) Path() string {
	return i.proto.Filename
}
