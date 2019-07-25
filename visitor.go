package protoast

import (
	"text/scanner"

	"github.com/emicklei/proto"

	"github.com/sirkon/protoast/ast"
	"github.com/sirkon/protoast/internal/namespace"
)

var _ proto.Visitor = &typesVisitor{}

type typesVisitor struct {
	file	*ast.File
	ns	namespace.Namespace
	nss	*Builder

	errors	chan<- error

	enumCtx	struct {
		item		*ast.Enum
		prevField	map[string]scanner.Position
		prevInteger	map[int]scanner.Position
	}

	msgCtx	struct {
		onMsg		bool
		item		*ast.Message
		prevField	map[string]scanner.Position
		prevInteger	map[int]scanner.Position
	}

	oneOf	*ast.OneOf

	service	*ast.Service
}

func (tv *typesVisitor) VisitMessage(m *proto.Message) {
	v := &typesVisitor{
		ns:	tv.ns.WithScope(m.Name),
		nss:	tv.nss,
		errors:	tv.errors,
	}

	msg := tv.ns.GetType(m.Name)
	if msg == nil {
		panic("internal error: message must be predeclared on the prefetch phase")
	}
	v.msgCtx.item = msg.(*ast.Message)
	v.msgCtx.prevField = map[string]scanner.Position{}
	v.msgCtx.prevInteger = map[int]scanner.Position{}

	for _, e := range m.Elements {
		e.Accept(v)
	}
	v.msgCtx.item = nil
}

func (tv *typesVisitor) VisitService(v *proto.Service) {
	tv.service = &ast.Service{
		File:	tv.file,
		Name:	v.Name,
	}

	for _, e := range v.Elements {
		e.Accept(tv)
	}
	if err := tv.ns.SetNode(v.Name, tv.service, v.Position); err != nil {
		tv.errors <- err
	}
	tv.file.Services = append(tv.file.Services, tv.service)
}

func (tv *typesVisitor) VisitSyntax(s *proto.Syntax)	{}

func (tv *typesVisitor) VisitPackage(p *proto.Package) {
	if err := tv.ns.SetPkgName(p.Name); err != nil {
		tv.errors <- err
	}
}
func (tv *typesVisitor) VisitOption(o *proto.Option)	{}

func (tv *typesVisitor) VisitImport(i *proto.Import) {
	importNs, _, err := tv.nss.get(i.Filename)
	if err != nil {
		tv.errors <- errPosf(i.Position, "reading import %s: %s", i.Filename, err)
		return
	}

	tv.ns, err = tv.ns.WithImport(importNs)
	if err != nil {
		tv.errors <- errPos(i.Position, err)
	}
}

func (tv *typesVisitor) VisitNormalField(i *proto.NormalField) {
	if prev, ok := tv.msgCtx.prevField[i.Name]; ok {
		tv.errors <- errPosf(i.Position, "duplicate field %s, the previous definition was in %s", i.Name, prev)
	}
	if prev, ok := tv.msgCtx.prevInteger[i.Sequence]; ok {
		tv.errors <- errPosf(i.Position, "duplicate field sequence %d, the previous valuy was in %s", i.Sequence, prev)
	}
	tv.msgCtx.prevField[i.Name] = i.Position
	tv.msgCtx.prevInteger[i.Sequence] = i.Position

	var options ast.Options
	for _, o := range i.Options {
		if options == nil {
			options = ast.Options{}
		}
		options[o.Name] = o.Constant.Source
	}

	t := standardType(i.Type)
	if t == nil {
		t = tv.ns.GetType(i.Type)
	}
	if t == nil {
		tv.errors <- errPosf(i.Position, "unknown type %s", i.Type)
		return
	}
	if i.Optional {
		t = ast.Optional{
			Type: t,
		}
	}
	if i.Repeated {
		t = ast.Repeated{
			Type: t,
		}
	}
	tv.msgCtx.item.Fields = append(tv.msgCtx.item.Fields, ast.MessageField{
		Name:		i.Name,
		Sequence:	i.Sequence,
		Type:		t,
		Options:	options,
	})
}

func standardType(typeName string) ast.Type {
	switch typeName {
	case "bool":
		return ast.Bool{}
	case "google.protobuf.Any":
		return ast.Any{}
	case "bytes":
		return ast.Bytes{}
	case "fixed32":
		return ast.Fixed32{}
	case "fixed64":
		return ast.Fixed64{}
	case "float":
		return ast.Float32{}
	case "double":
		return ast.Float64{}
	case "int32":
		return ast.Int32{}
	case "int64":
		return ast.Int64{}
	case "sfixed32":
		return ast.Sfixed32{}
	case "sfixed64":
		return ast.Sfixed64{}
	case "sint32":
		return ast.Sint32{}
	case "sint64":
		return ast.Sint64{}
	case "string":
		return ast.String{}
	case "uint32":
		return ast.Uint32{}
	case "uint64":
		return ast.Uint64{}
	}
	return nil
}

func (tv *typesVisitor) VisitEnumField(i *proto.EnumField) {
	if prev, ok := tv.enumCtx.prevField[i.Name]; ok {
		tv.errors <- errPosf(i.Position, "duplicate enum field %s, the previous definition was in %s", i.Name, prev)
	}
	if prev, ok := tv.enumCtx.prevInteger[i.Integer]; ok {
		tv.errors <- errPosf(i.Position, "duplicate enum field id %d, the previous field with the same id was in %s", i.Integer, prev)
	}
	tv.enumCtx.prevField[i.Name] = i.Position
	tv.enumCtx.prevInteger[i.Integer] = i.Position
	var options ast.Options
	if i.ValueOption != nil {
		options = ast.Options{}
		options[i.ValueOption.Name] = i.ValueOption.Constant.Source
	}
	tv.enumCtx.item.Values = append(tv.enumCtx.item.Values, ast.EnumValue{
		Name:		i.Name,
		Integer:	i.Integer,
		Options:	options,
	})
}

func (tv *typesVisitor) VisitEnum(e *proto.Enum) {
	enum := tv.ns.GetType(e.Name)
	if enum == nil {
		panic("internal error: enum must be predeclared on the prefetch phase")
	}
	tv.enumCtx.item = enum.(*ast.Enum)
	tv.enumCtx.prevField = map[string]scanner.Position{}
	tv.enumCtx.prevInteger = map[int]scanner.Position{}
	for _, e := range e.Elements {
		e.Accept(tv)
	}
}

func (tv *typesVisitor) VisitComment(e *proto.Comment)	{}

func (tv *typesVisitor) VisitOneof(o *proto.Oneof) {
	if prev, ok := tv.msgCtx.prevField[o.Name]; ok {
		tv.errors <- errPosf(o.Position, "duplicate field %s, the previous definition was in %s", o.Name, prev)
	}
	tv.msgCtx.prevField[o.Name] = o.Position

	tv.oneOf = &ast.OneOf{
		ParentMsg:	tv.msgCtx.item,
		Name:		o.Name,
	}
	tv.msgCtx.item.Fields = append(tv.msgCtx.item.Fields, ast.MessageField{
		Name:		o.Name,
		Sequence:	-1,
		Type:		tv.oneOf,
	})

	for _, e := range o.Elements {
		e.Accept(tv)
	}
}

func (tv *typesVisitor) VisitOneofField(o *proto.OneOfField) {
	if prev, ok := tv.msgCtx.prevField[o.Name]; ok {
		tv.errors <- errPosf(o.Position, "duplicate field %s, the previous definition was in %s", o.Name, prev)
	}
	if prev, ok := tv.msgCtx.prevInteger[o.Sequence]; ok {
		tv.errors <- errPosf(o.Position, "duplicate field sequence %d, the previous valuy was in %s", o.Sequence, prev)
	}
	tv.msgCtx.prevField[o.Name] = o.Position
	tv.msgCtx.prevInteger[o.Sequence] = o.Position

	var options ast.Options
	for _, o := range o.Options {
		if options == nil {
			options = ast.Options{}
		}
		options[o.Name] = o.Constant.Source
	}

	t := standardType(o.Type)
	if t == nil {
		t = tv.ns.GetType(o.Type)
	}
	if t == nil {
		tv.errors <- errPosf(o.Position, "unknown type %s", o.Type)
		return
	}
	tv.oneOf.Branches = append(tv.oneOf.Branches, ast.OneOfBranch{
		Name:		o.Name,
		Type:		t,
		Sequence:	o.Sequence,
		Options:	options,
	})
}

func (tv *typesVisitor) VisitReserved(r *proto.Reserved)	{}

func (tv *typesVisitor) VisitRPC(r *proto.RPC) {

	req := tv.ns.GetType(r.RequestType)
	if req == nil {
		tv.errors <- errPosf(r.Position, "unknown type %s", r.RequestType)
		return
	}
	if r.StreamsRequest {
		req = ast.Stream{
			Type: req,
		}
	}

	resp := tv.ns.GetType(r.ReturnsType)
	if resp == nil {
		tv.errors <- errPosf(r.Position, "unknown type %s", r.RequestType)
		return
	}
	if r.StreamsReturns {
		resp = ast.Stream{
			Type: resp,
		}
	}

	var mos []ast.MethodOption
	for _, o := range r.Options {
		mo := ast.MethodOption{
			Name: o.Name,
		}
		for _, v := range o.Constant.OrderedMap {
			mo.Values = append(mo.Values, ast.OptionValue{
				Name:	v.Name,
				Value:	v.Source,
			})
		}
		mos = append(mos, mo)
	}

	rpc := &ast.Method{
		File:		tv.file,
		Service:	tv.service,
		Name:		r.Name,
		Input:		req,
		Output:		resp,
		Options:	mos,
	}
	tv.service.Methods = append(tv.service.Methods, rpc)

	if err := tv.ns.SetNode(tv.service.Name+"::"+r.Name, rpc, r.Position); err != nil {
		tv.errors <- errPosf(r.Position, "duplicate method %s in service %s", r.Name, tv.service.Name)
		return
	}
}

func (tv *typesVisitor) VisitMapField(f *proto.MapField) {
	if prev, ok := tv.msgCtx.prevField[f.Name]; ok {
		tv.errors <- errPosf(f.Position, "duplicate field %s, the previous definition was in %s", f.Name, prev)
	}
	if prev, ok := tv.msgCtx.prevInteger[f.Sequence]; ok {
		tv.errors <- errPosf(f.Position, "duplicate field sequence %d, the previous valuy was in %s", f.Sequence, prev)
	}
	tv.msgCtx.prevField[f.Name] = f.Position
	tv.msgCtx.prevInteger[f.Sequence] = f.Position

	var options ast.Options
	for _, o := range f.Options {
		if options == nil {
			options = ast.Options{}
		}
		options[o.Name] = o.Constant.Source
	}

	keyRawType := standardType(f.KeyType)
	if keyRawType == nil {
		tv.errors <- errPosf(f.Position, "invalid map key type %s", f.Type)
		return
	}
	keyType, isHashable := keyRawType.(ast.Hashable)
	if !isHashable {
		tv.errors <- errPosf(f.Position, "invalid map key type %s", f.Type)
	}

	valType := standardType(f.KeyType)
	if valType == nil {
		valType = tv.ns.GetType(f.Type)
		if valType == nil {
			tv.errors <- errPosf(f.Position, "unknown value type %s", f.Type)
			return
		}
	}
	tv.msgCtx.item.Fields = append(tv.msgCtx.item.Fields, ast.MessageField{
		Name:		f.Name,
		Sequence:	f.Sequence,
		Type: ast.Map{
			KeyType:	keyType,
			ValueType:	valType,
		},
		Options:	options,
	})
}

func (tv *typesVisitor) VisitGroup(g *proto.Group)		{}
func (tv *typesVisitor) VisitExtensions(e *proto.Extensions)	{}
