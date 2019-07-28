package protoast

import (
	"text/scanner"

	"github.com/emicklei/proto"

	"github.com/sirkon/protoast/ast"
	"github.com/sirkon/protoast/internal/namespace"
)

var _ proto.Visitor = &typesVisitor{}

type typesVisitor struct {
	file *ast.File
	ns   namespace.Namespace
	nss  *Builder

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

	service *ast.Service
}

func (tv *typesVisitor) regInfo(k ast.Unique, comment *proto.Comment, pos scanner.Position) {
	ast.SetKey(k, tv.nss.uniqueContext)
	key := ast.GetKey(k)

	if comment != nil {
		cmt := &ast.Comment{
			Value: comment.Message(),
		}
		tv.regInfo(cmt, nil, comment.Position)
		tv.nss.comments[key] = cmt
	}

	tv.nss.positions[key] = pos
}

func (tv *typesVisitor) regFieldInfo(k ast.Unique, fieldAddr interface{}, comment *proto.Comment, pos scanner.Position) {
	ast.SetKey(k, tv.nss.uniqueContext)
	key := ast.GetFieldKey(k, fieldAddr)

	if comment != nil {
		cmt := &ast.Comment{
			Value: comment.Message(),
		}
		tv.regInfo(cmt, nil, comment.Position)
		tv.nss.comments[key] = cmt
	}

	tv.nss.positions[key] = pos
}

func (tv *typesVisitor) VisitMessage(m *proto.Message) {
	v := &typesVisitor{
		ns:     tv.ns.WithScope(m.Name),
		file:   tv.file,
		nss:    tv.nss,
		errors: tv.errors,
	}

	msg := tv.ns.GetType(m.Name).(*ast.Message)
	if msg == nil {
		panic("internal error: message must be predeclared on the prefetch phase")
	}
	v.msgCtx.item = msg
	v.msgCtx.prevField = map[string]scanner.Position{}
	v.msgCtx.item.File = tv.file
	v.msgCtx.prevInteger = map[int]scanner.Position{}

	if msg.ParentMsg == nil && !m.IsExtend {
		tv.file.Types = append(tv.file.Types, msg)
	}

	tv.regInfo(msg, m.Comment, m.Position)
	tv.regFieldInfo(msg, &msg.Name, m.Comment, m.Position)

	for _, e := range m.Elements {
		e.Accept(v)
	}
	v.msgCtx.item = nil
}

func (tv *typesVisitor) VisitService(v *proto.Service) {
	tv.service = &ast.Service{
		File: tv.file,
		Name: v.Name,
	}

	for _, e := range v.Elements {
		e.Accept(tv)
	}
	if err := tv.ns.SetNode(v.Name, tv.service, v.Position); err != nil {
		tv.errors <- err
	}
	tv.file.Services = append(tv.file.Services, tv.service)
	tv.regInfo(tv.service, v.Comment, v.Position)
	tv.regFieldInfo(tv.service, &tv.service.Name, v.Comment, v.Position)
}

func (tv *typesVisitor) VisitSyntax(s *proto.Syntax) {}

func (tv *typesVisitor) VisitPackage(p *proto.Package) {
	tv.file.Package = p.Name
	if err := tv.ns.SetPkgName(p.Name); err != nil {
		tv.errors <- err
	}
	tv.regFieldInfo(tv.file, &tv.file.Package, p.Comment, p.Position)
}
func (tv *typesVisitor) VisitOption(o *proto.Option) {
	option := &ast.Option{
		Name:  o.Name,
		Value: o.Constant.Source,
	}
	tv.file.Options = append(tv.file.Options, option)
	tv.regInfo(option, o.Comment, o.Position)
	tv.regFieldInfo(option, &option.Name, nil, o.Position)
	tv.regFieldInfo(option, &option.Value, nil, o.Constant.Position)
}

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
	imp := &ast.Import{
		Path: i.Filename,
	}
	tv.file.Imports = append(tv.file.Imports, imp)
	tv.regInfo(imp, i.Comment, i.Position)
	tv.regFieldInfo(imp, &imp.Path, nil, i.Position)
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

	var options []*ast.Option
	for _, o := range i.Options {
		option := &ast.Option{
			Name:  o.Name,
			Value: o.Constant.Source,
		}
		options = append(options, option)
		tv.regInfo(option, o.Comment, o.Position)
		tv.regFieldInfo(option, &option.Name, nil, o.Position)
		tv.regFieldInfo(option, &option.Value, nil, o.Constant.Position)
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
		t = &ast.Optional{
			Type: t,
		}
	}
	if i.Repeated {
		t = &ast.Repeated{
			Type: t,
		}
	}
	field := &ast.MessageField{
		Name:     i.Name,
		Sequence: i.Sequence,
		Type:     t,
		Options:  options,
	}
	tv.regInfo(field, i.Comment, i.Position)
	tv.regFieldInfo(field, &field.Name, nil, i.Position)
	tv.regFieldInfo(field, &field.Sequence, nil, i.Position)
	tv.regInfo(t, nil, i.Position)
	tv.msgCtx.item.Fields = append(tv.msgCtx.item.Fields, field)
}

func standardType(typeName string) ast.Type {
	switch typeName {
	case "bool":
		return &ast.Bool{}
	case "google.protobuf.Any":
		return &ast.Any{}
	case "bytes":
		return &ast.Bytes{}
	case "fixed32":
		return &ast.Fixed32{}
	case "fixed64":
		return &ast.Fixed64{}
	case "float":
		return &ast.Float32{}
	case "double":
		return &ast.Float64{}
	case "int32":
		return &ast.Int32{}
	case "int64":
		return &ast.Int64{}
	case "sfixed32":
		return &ast.Sfixed32{}
	case "sfixed64":
		return &ast.Sfixed64{}
	case "sint32":
		return &ast.Sint32{}
	case "sint64":
		return &ast.Sint64{}
	case "string":
		return &ast.String{}
	case "uint32":
		return &ast.Uint32{}
	case "uint64":
		return &ast.Uint64{}
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
	var options []*ast.Option
	if i.ValueOption != nil {
		option := &ast.Option{
			Name:  i.ValueOption.Name,
			Value: i.ValueOption.Constant.Source,
		}
		tv.regInfo(option, i.ValueOption.Comment, i.ValueOption.Position)
		tv.regFieldInfo(option, &option.Name, nil, i.ValueOption.Position)
		tv.regFieldInfo(option, &option.Value, nil, i.ValueOption.Constant.Position)
		options = append(options, option)
	}
	value := &ast.EnumValue{
		Name:    i.Name,
		Integer: i.Integer,
		Options: options,
	}
	tv.enumCtx.item.Values = append(tv.enumCtx.item.Values, value)
	tv.regInfo(value, i.Comment, i.Position)
	tv.regFieldInfo(value, &value.Name, nil, i.Position)
	tv.regFieldInfo(value, &value.Integer, nil, i.Position)
}

func (tv *typesVisitor) VisitEnum(e *proto.Enum) {
	enum := tv.ns.GetType(e.Name).(*ast.Enum)
	if enum == nil {
		panic("internal error: enum must be predeclared on the prefetch phase")
	}

	tv.enumCtx.item = enum
	tv.enumCtx.prevField = map[string]scanner.Position{}
	tv.enumCtx.prevInteger = map[int]scanner.Position{}
	for _, e := range e.Elements {
		e.Accept(tv)
	}

	if enum.ParentMsg == nil {
		tv.file.Types = append(tv.file.Types, enum)
	}

	tv.regInfo(enum, e.Comment, e.Position)
	tv.regFieldInfo(enum, &enum.Name, nil, e.Position)
}

func (tv *typesVisitor) VisitComment(e *proto.Comment) {}

func (tv *typesVisitor) VisitOneof(o *proto.Oneof) {
	if prev, ok := tv.msgCtx.prevField[o.Name]; ok {
		tv.errors <- errPosf(o.Position, "duplicate field %s, the previous definition was in %s", o.Name, prev)
	}
	tv.msgCtx.prevField[o.Name] = o.Position

	oo := &ast.OneOf{
		ParentMsg: tv.msgCtx.item,
		Name:      o.Name,
	}
	tv.oneOf = oo
	tv.regInfo(oo, o.Comment, o.Position)
	tv.regFieldInfo(oo, &oo.Name, o.Comment, o.Position)

	oob := &ast.MessageField{
		Name:     o.Name,
		Sequence: -1,
		Type:     tv.oneOf,
	}
	tv.msgCtx.item.Fields = append(tv.msgCtx.item.Fields, oob)
	tv.regInfo(oob, o.Comment, o.Position)
	tv.regFieldInfo(oob, &oob.Name, nil, o.Position)
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

	var options []*ast.Option
	for _, o := range o.Options {
		option := &ast.Option{
			Name:  o.Name,
			Value: o.Constant.Source,
		}
		options = append(options, option)
		tv.regInfo(option, o.Comment, o.Position)
		tv.regFieldInfo(option, &option.Name, nil, o.Position)
		tv.regFieldInfo(option, &option.Value, nil, o.Constant.Position)
	}

	t := standardType(o.Type)
	if t == nil {
		t = tv.ns.GetType(o.Type)
	}
	tv.regInfo(t, nil, o.Position)
	if t == nil {
		tv.errors <- errPosf(o.Position, "unknown type %s", o.Type)
		return
	}
	b := &ast.OneOfBranch{
		Name:     o.Name,
		Type:     t,
		Sequence: o.Sequence,
		Options:  options,
	}
	tv.regInfo(b, o.Comment, o.Position)
	tv.regFieldInfo(b, &b.Name, nil, o.Position)
	tv.regFieldInfo(b, &b.Sequence, nil, o.Position)

	tv.oneOf.Branches = append(tv.oneOf.Branches, b)
}

func (tv *typesVisitor) VisitReserved(r *proto.Reserved) {}

func (tv *typesVisitor) VisitRPC(r *proto.RPC) {

	req := tv.ns.GetType(r.RequestType)
	if req == nil {
		tv.errors <- errPosf(r.Position, "unknown type %s", r.RequestType)
		return
	}
	if r.StreamsRequest {
		req = &ast.Stream{
			Type: req,
		}
	}
	tv.regInfo(req, nil, r.Position)

	resp := tv.ns.GetType(r.ReturnsType)
	if resp == nil {
		tv.errors <- errPosf(r.Position, "unknown type %s", r.RequestType)
		return
	}
	if r.StreamsReturns {
		resp = &ast.Stream{
			Type: resp,
		}
	}
	tv.regInfo(resp, nil, r.Position)

	var mos []*ast.MethodOption
	for _, o := range r.Options {
		mo := &ast.MethodOption{
			Name: o.Name,
		}
		tv.regInfo(mo, o.Comment, o.Position)
		tv.regFieldInfo(mo, &mo.Name, nil, o.Position)
		for _, v := range o.Constant.OrderedMap {
			value := &ast.MethodOptionValue{
				Name:  v.Name,
				Value: v.Source,
			}
			tv.regInfo(value, nil, v.Position)
			tv.regFieldInfo(value, &value.Name, nil, v.Position)
			tv.regFieldInfo(value, &value.Value, nil, v.Literal.Position)
			mo.Values = append(mo.Values, value)
		}
		mos = append(mos, mo)
	}

	rpc := &ast.Method{
		File:    tv.file,
		Service: tv.service,
		Name:    r.Name,
		Input:   req,
		Output:  resp,
		Options: mos,
	}
	tv.service.Methods = append(tv.service.Methods, rpc)
	tv.regInfo(rpc, r.Comment, r.Position)
	tv.regFieldInfo(rpc, &rpc.Name, nil, r.Position)

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

	var options []*ast.Option
	for _, o := range f.Options {
		option := &ast.Option{
			Name:  o.Name,
			Value: o.Constant.Source,
		}
		tv.regInfo(option, o.Comment, o.Position)
		tv.regFieldInfo(option, &option.Name, nil, o.Position)
		tv.regFieldInfo(option, &option.Value, nil, o.Constant.Position)
		options = append(options, option)
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
	tv.regInfo(keyType, nil, f.Position)

	valType := standardType(f.KeyType)
	if valType == nil {
		valType = tv.ns.GetType(f.Type)
		if valType == nil {
			tv.errors <- errPosf(f.Position, "unknown value type %s", f.Type)
			return
		}
	}
	tv.regInfo(valType, nil, f.Position)

	fieldType := &ast.Map{
		KeyType:   keyType,
		ValueType: valType,
	}
	tv.regInfo(fieldType, f.Comment, f.Position)

	field := &ast.MessageField{
		Name:     f.Name,
		Sequence: f.Sequence,
		Type:     fieldType,
		Options:  options,
	}
	tv.msgCtx.item.Fields = append(tv.msgCtx.item.Fields, field)
	tv.regInfo(field, f.Comment, f.Position)
	tv.regFieldInfo(field, &field.Name, nil, f.Position)
	tv.regFieldInfo(field, &field.Sequence, nil, f.Position)
}

func (tv *typesVisitor) VisitGroup(g *proto.Group)           {}
func (tv *typesVisitor) VisitExtensions(e *proto.Extensions) {}
