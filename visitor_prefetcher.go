package protoast

import (
	"github.com/emicklei/proto"
	"github.com/pkg/errors"

	"github.com/sirkon/protoast/ast"
	"github.com/sirkon/protoast/internal/namespace"
)

var _ proto.Visitor = &prefetcher{}

type prefetcher struct {
	file   *ast.File
	ns     namespace.Namespace
	nss    *Builder
	curMsg *ast.Message

	errors chan<- error
}

func (p *prefetcher) VisitMessage(m *proto.Message) {
	var message *ast.Message
	if m.IsExtend {
		msg := p.ns.GetType(m.Name)
		if msg == nil {
			p.errors <- errors.Errorf("%s failed to find type %s to extend", m.Position, m.Name)
			return
		}
		var ok bool
		message, ok = msg.(*ast.Message)
		if !ok {
			p.errors <- errors.Errorf("%s type %s turned to be not a message (%T)", m.Position, m.Name, msg)
		}
	} else {
		message = &ast.Message{
			ParentMsg: p.curMsg,
			File:      p.file,
			Name:      m.Name,
		}
	}
	if p.curMsg != nil {
		p.curMsg.Types = append(p.curMsg.Types, message)
	}
	v := &prefetcher{
		file:   p.file,
		ns:     p.ns.WithScope(m.Name),
		nss:    p.nss,
		errors: p.errors,
		curMsg: message,
	}

	for _, e := range m.Elements {
		e.Accept(v)
	}

	if m.IsExtend {
		return
	}
	if err := p.ns.SetNode(m.Name, message, m.Position); err != nil {
		v.errors <- err
	}
}

func (p *prefetcher) VisitService(v *proto.Service) {}
func (p *prefetcher) VisitSyntax(s *proto.Syntax)   {}

func (p *prefetcher) VisitPackage(pkg *proto.Package) {
	p.file.Package = pkg.Name
	if err := p.ns.SetPkgName(pkg.Name); err != nil {
		p.errors <- err
	}
}

func (p *prefetcher) VisitOption(o *proto.Option) {}

func (p *prefetcher) VisitImport(i *proto.Import) {
	ins, _, err := p.nss.get(i.Filename)
	if err != nil {
		p.errors <- errPosf(i.Position, "reading import %s: %s", i.Filename, err)
		return
	}

	p.ns, err = p.ns.WithImport(ins)
	if err != nil {
		p.errors <- errPos(i.Position, err)
	}
}

func (p *prefetcher) VisitNormalField(i *proto.NormalField) {}
func (p *prefetcher) VisitEnumField(i *proto.EnumField)     {}

func (p *prefetcher) VisitEnum(e *proto.Enum) {
	enum := &ast.Enum{
		ParentMsg: p.curMsg,
		File:      p.file,
		Name:      e.Name,
		Values:    nil,
	}
	if p.curMsg != nil {
		p.curMsg.Types = append(p.curMsg.Types, enum)
	}
	if err := p.ns.SetNode(e.Name, enum, e.Position); err != nil {
		p.errors <- err
	}
}

func (p *prefetcher) VisitComment(e *proto.Comment) {}

func (p *prefetcher) VisitOneof(o *proto.Oneof)           {}
func (p *prefetcher) VisitOneofField(o *proto.OneOfField) {}
func (p *prefetcher) VisitReserved(r *proto.Reserved)     {}
func (p *prefetcher) VisitRPC(r *proto.RPC)               {}
func (p *prefetcher) VisitMapField(f *proto.MapField)     {}
func (p *prefetcher) VisitGroup(g *proto.Group)           {}
func (p *prefetcher) VisitExtensions(e *proto.Extensions) {}
