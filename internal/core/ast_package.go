package core

import (
	"github.com/emicklei/proto"
)

type Package struct {
	isNode

	proto *proto.Package
}

func (p *Package) Package() string {
	return p.proto.Name
}
