package prototypes

import (
	"text/scanner"

	"github.com/emicklei/proto"

	"github.com/sirkon/prototypes/ast"
	"github.com/sirkon/prototypes/internal/namespace"
)

var _ proto.Visitor = &typesVisitor{}

type typesVisitor struct {
	ns  namespace.Namespace
	nss *Namespaces

	errors chan<- error

	enumCtx struct {
		item        *ast.Enum
		prevField   map[string]scanner.Position
		prevInteger map[int]scanner.Position
	}

	msgCtx struct {
		onMsg       bool
		item        *ast.Message
		prevField   map[string]scanner.Position
		prevInteger map[int]scanner.Position
	}

	oneOf *ast.OneOf
}

func (tv *typesVisitor) VisitMessage(m *proto.Message) {
	v := &typesVisitor{
		ns:     tv.ns.WithScope(m.Name),
		nss:    tv.nss,
		errors: tv.errors,
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
}

func (tv *typesVisitor) VisitService(v *proto.Service) {}

func (tv *typesVisitor) VisitSyntax(s *proto.Syntax)   {}
func (tv *typesVisitor) VisitPackage(p *proto.Package) { tv.ns.SetPkgName(p.Name) }
func (tv *typesVisitor) VisitOption(o *proto.Option)   {}

func (tv *typesVisitor) VisitImport(i *proto.Import) {
	importNs, err := tv.nss.get(i.Filename)
	if err != nil {
		tv.errors <- ErrPosf(i.Position, "reading import %s: %s", i.Filename, err)
		return
	}

	tv.ns, err = tv.ns.WithImport(importNs)
	if err != nil {
		tv.errors <- ErrPos(i.Position, err)
	}
}

func (tv *typesVisitor) VisitNormalField(i *proto.NormalField) {
	if prev, ok := tv.msgCtx.prevField[i.Name]; ok {
		tv.errors <- ErrPosf(i.Position, "duplicate field %s, the previous definition was in %s", i.Name, prev)
	}
	if prev, ok := tv.msgCtx.prevInteger[i.Sequence]; ok {
		tv.errors <- ErrPosf(i.Position, "duplicate field sequence %d, the previous valuy was in %s", i.Sequence, prev)
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
		tv.errors <- ErrPosf(i.Position, "unknown type %s", i.Type)
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
		Name:     i.Name,
		Sequence: i.Sequence,
		Type:     t,
		Options:  options,
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
		tv.errors <- ErrPosf(i.Position, "duplicate enum field %s, the previous definition was in %s", i.Name, prev)
	}
	if prev, ok := tv.enumCtx.prevInteger[i.Integer]; ok {
		tv.errors <- ErrPosf(i.Position, "duplicate enum field id %d, the previous field with the same id was in %s", i.Integer, prev)
	}
	tv.enumCtx.prevField[i.Name] = i.Position
	tv.enumCtx.prevInteger[i.Integer] = i.Position
	var options ast.Options
	if i.ValueOption != nil {
		options = ast.Options{}
		options[i.ValueOption.Name] = i.ValueOption.Constant.Source
	}
	tv.enumCtx.item.Values = append(tv.enumCtx.item.Values, ast.EnumValue{
		Name:    i.Name,
		Integer: i.Integer,
		Options: options,
	})
}

func (tv *typesVisitor) VisitEnum(e *proto.Enum) {
	enum := tv.ns.GetType(e.Name)
	if enum == nil {
		panic("internal error: enum must be predeclared on prefetch phase")
	}
	tv.enumCtx.item = enum.(*ast.Enum)
	tv.enumCtx.prevField = map[string]scanner.Position{}
	tv.enumCtx.prevInteger = map[int]scanner.Position{}
	for _, e := range e.Elements {
		e.Accept(tv)
	}
}

func (tv *typesVisitor) VisitComment(e *proto.Comment) {}

func (tv *typesVisitor) VisitOneof(o *proto.Oneof) {
	if prev, ok := tv.msgCtx.prevField[o.Name]; ok {
		tv.errors <- ErrPosf(o.Position, "duplicate field %s, the previous definition was in %s", o.Name, prev)
	}
	tv.msgCtx.prevField[o.Name] = o.Position

	tv.oneOf = &ast.OneOf{}
	tv.msgCtx.item.Fields = append(tv.msgCtx.item.Fields, ast.MessageField{
		Name:     o.Name,
		Sequence: -1,
		Type:     tv.oneOf,
	})

	for _, e := range o.Elements {
		e.Accept(tv)
	}
}

func (tv *typesVisitor) VisitOneofField(o *proto.OneOfField) {
	if prev, ok := tv.msgCtx.prevField[o.Name]; ok {
		tv.errors <- ErrPosf(o.Position, "duplicate field %s, the previous definition was in %s", o.Name, prev)
	}
	if prev, ok := tv.msgCtx.prevInteger[o.Sequence]; ok {
		tv.errors <- ErrPosf(o.Position, "duplicate field sequence %d, the previous valuy was in %s", o.Sequence, prev)
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
		tv.errors <- ErrPosf(o.Position, "unknown type %s", o.Type)
		return
	}
	tv.oneOf.Branches = append(tv.oneOf.Branches, ast.OneOfBranch{
		Name:     o.Name,
		Type:     t,
		Sequence: o.Sequence,
		Options:  options,
	})
}

func (tv *typesVisitor) VisitReserved(r *proto.Reserved) {}
func (tv *typesVisitor) VisitRPC(r *proto.RPC)           {}

func (tv *typesVisitor) VisitMapField(f *proto.MapField) {
	if prev, ok := tv.msgCtx.prevField[f.Name]; ok {
		tv.errors <- ErrPosf(f.Position, "duplicate field %s, the previous definition was in %s", f.Name, prev)
	}
	if prev, ok := tv.msgCtx.prevInteger[f.Sequence]; ok {
		tv.errors <- ErrPosf(f.Position, "duplicate field sequence %d, the previous valuy was in %s", f.Sequence, prev)
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
		tv.errors <- ErrPosf(f.Position, "invalid map key type %s", f.Type)
		return
	}
	keyType, isHashable := keyRawType.(ast.Hashable)
	if !isHashable {
		tv.errors <- ErrPosf(f.Position, "invalid map key type %s", f.Type)
	}

	valType := standardType(f.KeyType)
	if valType == nil {
		valType = tv.ns.GetType(f.Type)
		if valType == nil {
			tv.errors <- ErrPosf(f.Position, "unknown value type %s", f.Type)
			return
		}
	}
	tv.msgCtx.item.Fields = append(tv.msgCtx.item.Fields, ast.MessageField{
		Name:     f.Name,
		Sequence: f.Sequence,
		Type: ast.Map{
			KeyType:   keyType,
			ValueType: valType,
		},
		Options: options,
	})
}

func (tv *typesVisitor) VisitGroup(g *proto.Group) {
}

func (tv *typesVisitor) VisitExtensions(e *proto.Extensions) {}
