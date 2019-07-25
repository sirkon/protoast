package protoast

import (
	"github.com/emicklei/proto"

	"github.com/sirkon/protoast/ast"
	"github.com/sirkon/protoast/internal/namespace"
)

var _ proto.Visitor = &prefetcher{}

type prefetcher struct {
	file	*ast.File
	ns	namespace.Namespace
	nss	*Builder
	curMsg	*ast.Message

	errors	chan<- error
}

func (p *prefetcher) VisitMessage(m *proto.Message) {
	message := &ast.Message{
		ParentMsg:	p.curMsg,
		File:		p.file,
		Name:		m.Name,
	}
	v := &prefetcher{
		ns:	p.ns.WithScope(m.Name),
		nss:	p.nss,
		errors:	p.errors,
		curMsg:	message,
	}

	for _, e := range m.Elements {
		e.Accept(v)
	}

	if err := p.ns.SetNode(m.Name, message, m.Position); err != nil {
		v.errors <- err
	}
	p.curMsg = nil
}

func (p *prefetcher) VisitService(v *proto.Service)	{}
func (p *prefetcher) VisitSyntax(s *proto.Syntax)	{}

func (p *prefetcher) VisitPackage(pkg *proto.Package) {
	p.file.Package = pkg.Name
	if err := p.ns.SetPkgName(pkg.Name); err != nil {
		p.errors <- err
	}
}

func (p *prefetcher) VisitOption(o *proto.Option)	{}

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

func (p *prefetcher) VisitNormalField(i *proto.NormalField)	{}
func (p *prefetcher) VisitEnumField(i *proto.EnumField)		{}

func (p *prefetcher) VisitEnum(e *proto.Enum) {
	enum := &ast.Enum{
		ParentMsg:	p.curMsg,
		File:		p.file,
		Name:		e.Name,
		Values:		nil,
	}
	if err := p.ns.SetNode(e.Name, enum, e.Position); err != nil {
		p.errors <- err
	}
}

func (p *prefetcher) VisitComment(e *proto.Comment)	{}

func (p *prefetcher) VisitOneof(o *proto.Oneof)			{}
func (p *prefetcher) VisitOneofField(o *proto.OneOfField)	{}
func (p *prefetcher) VisitReserved(r *proto.Reserved)		{}
func (p *prefetcher) VisitRPC(r *proto.RPC)			{}
func (p *prefetcher) VisitMapField(f *proto.MapField)		{}
func (p *prefetcher) VisitGroup(g *proto.Group)			{}
func (p *prefetcher) VisitExtensions(e *proto.Extensions)	{}
