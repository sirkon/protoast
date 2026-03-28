package core

import (
	"text/scanner"

	"github.com/emicklei/proto"
)

type Package struct {
	proto *proto.Package
}

func (p *Package) Package() string {
	return p.proto.Name
}

var _ Node = new(Package)

func (p *Package) nodeProto() proto.Visitee { return p.proto }
func (p *Package) pos() scanner.Position    { return p.proto.Position }
