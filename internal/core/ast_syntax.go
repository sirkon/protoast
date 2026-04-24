package core

import (
	"text/scanner"

	"github.com/emicklei/proto"
)

type Syntax struct {
	proto *proto.Syntax
}

var _ Node = new(Syntax)

func (s *Syntax) Value() string {
	return s.proto.Value
}

func (s *Syntax) nodeProto() proto.Visitee { return s.proto }
func (s *Syntax) pos() scanner.Position    { return s.proto.Position }
