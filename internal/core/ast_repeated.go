package core

import (
	"text/scanner"

	"github.com/emicklei/proto"
)

type Repeated struct {
	isType
	Type Type
}

var _ Node = new(Repeated)

func (r *Repeated) nodeProto() proto.Visitee { return r.Type.nodeProto() }
func (r *Repeated) pos() scanner.Position    { return r.Type.pos() }
