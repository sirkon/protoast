package core

import (
	"strings"

	"github.com/emicklei/proto"
	"github.com/sirkon/protoast/v2/internal/errors"
)

type visitorDemark struct {
	r *Registry

	file     *proto.Proto
	scope    string
	isExtend bool
}

func (v *visitorDemark) scopedName(n string) string {
	return v.scope + "." + n
}

func (v *visitorDemark) VisitMessage(m *proto.Message) {
	prevScope := v.scope
	if m.IsExtend {
		v.isExtend = true
	} else {
		v.scope = v.scopedName(m.Name)
		v.r.registry[v.scope] = m
		v.r.scopes[m] = v.scope
	}

	for _, e := range m.Elements {
		e.Accept(v)
	}

	v.isExtend = false
	v.scope = prevScope
}

func (v *visitorDemark) VisitService(s *proto.Service) {
	scopedName := v.scopedName(s.Name)
	v.r.registry[scopedName] = s
	prevScope := v.scope
	v.scope = scopedName
	v.r.registry[v.scope] = s
	v.r.scopes[s] = v.scope
	for _, e := range s.Elements {
		e.Accept(v)
	}
	v.scope = prevScope
}

func (v *visitorDemark) VisitSyntax(s *proto.Syntax) {}

func (v *visitorDemark) VisitPackage(p *proto.Package) {
	if strings.HasPrefix(p.Name, ".") {
		v.scope = p.Name
	} else {
		if p.Name != "" {
			v.scope = "." + p.Name
		}
	}
	v.r.scopes[v.file] = v.scope
}

func (v *visitorDemark) VisitOption(o *proto.Option) {}

func (v *visitorDemark) VisitImport(i *proto.Import) {
	if _, ok := v.r.protos[i.Filename]; ok {
		return
	}

	file, err := v.r.protoFile(i.Filename)
	if err != nil {
		panic(errors.Wrap(err, "get proto file for import "+i.Filename))
	}

	vv := &visitorDemark{
		r: v.r,
	}
	file.Accept(vv)
}

func (v *visitorDemark) VisitNormalField(f *proto.NormalField) {
	scopedName := v.scopedName(f.Name)
	v.r.registry[scopedName] = f
	v.r.scopes[f] = scopedName
}

func (v *visitorDemark) VisitEnumField(f *proto.EnumField) {
	scopedName := v.scopedName(f.Name)
	v.r.registry[scopedName] = f
}

func (v *visitorDemark) VisitEnum(e *proto.Enum) {
	scopedName := v.scopedName(e.Name)
	v.r.registry[scopedName] = e
	v.r.scopes[e] = scopedName
	for _, e := range e.Elements {
		e.Accept(v)
	}
}

func (v *visitorDemark) VisitComment(c *proto.Comment) {}

func (v *visitorDemark) VisitOneof(o *proto.Oneof) {
	scopedName := v.scopedName(o.Name)
	prev := v.scope
	v.r.registry[scopedName] = o
	v.r.scopes[o] = scopedName
	for _, e := range o.Elements {
		e.Accept(v)
	}
	v.scope = prev
}
func (v *visitorDemark) VisitOneofField(f *proto.OneOfField) {
	scopedName := v.scopedName(f.Name)
	v.r.registry[scopedName] = f
	v.r.scopes[f] = scopedName
}

func (v *visitorDemark) VisitReserved(r *proto.Reserved) {}

func (v *visitorDemark) VisitRPC(r *proto.RPC) {
	scopedName := v.scopedName(r.Name)
	v.r.registry[scopedName] = r
	v.r.scopes[r] = scopedName
}

func (v *visitorDemark) VisitMapField(f *proto.MapField) {
	scopedName := v.scopedName(f.Name)
	v.r.registry[scopedName] = f
	v.r.scopes[f] = scopedName
}

func (v *visitorDemark) VisitGroup(g *proto.Group)           {}
func (v *visitorDemark) VisitExtensions(e *proto.Extensions) {}
