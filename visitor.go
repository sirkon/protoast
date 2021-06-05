package protoast

import (
	"fmt"
	"strconv"
	"strings"
	"text/scanner"

	"github.com/emicklei/proto"
	"github.com/sirkon/protoast/internal/errors"
	"github.com/sirkon/protoast/internal/namespace"

	"github.com/sirkon/protoast/ast"
)

var _ proto.Visitor = &typesVisitor{}

type typesVisitor struct {
	file *ast.File
	ns   namespace.Namespace
	nss  *Builder

	errors func(err error)

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
	serviceCtx *ast.Service

	// очень крайне ненадёжно, здесь должен быть стэк OneOf-ов
	oneOf *ast.OneOf

	service *ast.Service
}

func (tv *typesVisitor) regInfo(k ast.Unique, comment *proto.Comment, pos scanner.Position) {
	ast.SetUnique(k, tv.nss.uniqueContext)
	key := ast.GetUnique(k)

	if comment != nil {
		cmt := &ast.Comment{
			Value: comment.Message(),
			Lines: comment.Lines,
		}
		tv.regInfo(cmt, nil, comment.Position)
		tv.nss.comments[key] = cmt
	}

	tv.nss.positions[key] = pos
}

func (tv *typesVisitor) regFieldInfo(k ast.Unique, fieldAddr interface{}, comment *proto.Comment, pos scanner.Position) {
	ast.SetUnique(k, tv.nss.uniqueContext)
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
	msg := &ast.Message{
		File:      tv.file,
		Name:      m.Name,
		ParentMsg: tv.msgCtx.item,
	}
	prev := tv.msgCtx
	defer func() {
		tv.msgCtx = prev
	}()
	tv.processDirectMessage(m, msg)
	if m.IsExtend {
		ext := ast.MessageToExtension(msg)
		tv.regInfo(ext, m.Comment, m.Position)
		tv.regFieldInfo(ext, &ext.Name, m.Comment, m.Position)
		tv.file.Extensions = append(tv.file.Extensions, ext)
	}

	realMsg := tv.ns.GetType(m.Name).(*ast.Message)
	if realMsg == nil {
		panic("internal error: message must be predeclared on the prefetch phase")
	}
	realMsg.Fields = append(realMsg.Fields, msg.Fields...)
	realMsg.Options = msg.Options
	*msg = *realMsg

	if realMsg.ParentMsg == nil && !m.IsExtend {
		tv.file.Types = append(tv.file.Types, realMsg)
	}
	tv.regInfo(realMsg, m.Comment, m.Position)
	tv.regFieldInfo(realMsg, &realMsg.Name, m.Comment, m.Position)
}

func (tv *typesVisitor) processDirectMessage(m *proto.Message, msg *ast.Message) {
	v := &typesVisitor{
		ns:     tv.ns.WithScope(m.Name),
		file:   tv.file,
		nss:    tv.nss,
		errors: tv.errors,
	}
	v.msgCtx.item = msg
	v.msgCtx.prevField = map[string]scanner.Position{}
	v.msgCtx.item.File = tv.file
	v.msgCtx.prevInteger = map[int]scanner.Position{}
	for _, e := range m.Elements {
		e.Accept(v)
	}
}

func (tv *typesVisitor) VisitService(v *proto.Service) {
	tv.service = &ast.Service{
		File: tv.file,
		Name: v.Name,
	}

	tv.serviceCtx = tv.service
	for _, e := range v.Elements {
		e.Accept(tv)
	}
	tv.serviceCtx = nil

	if err := tv.ns.SetNode(v.Name, tv.service, v.Position); err != nil {
		tv.errors(err)
	}
	tv.file.Services = append(tv.file.Services, tv.service)
	tv.regInfo(tv.service, v.Comment, v.Position)
	tv.regFieldInfo(tv.service, &tv.service.Name, v.Comment, v.Position)

}

func (tv *typesVisitor) VisitSyntax(s *proto.Syntax) {
	tv.file.Syntax = s.Value
	tv.regFieldInfo(tv.file, &tv.file.Syntax, s.Comment, s.Position)
}

func (tv *typesVisitor) VisitPackage(p *proto.Package) {
	tv.file.Package = p.Name
	if tv.file.GoPkg == "" {
		tv.file.GoPkg = p.Name
	}
	if err := tv.ns.SetPkgName(p.Name); err != nil {
		tv.errors(err)
	}
	tv.regFieldInfo(tv.file, &tv.file.Package, p.Comment, p.Position)
}
func (tv *typesVisitor) VisitOption(o *proto.Option) {
	var option *ast.Option
	if tv.serviceCtx != nil {
		option = tv.feedOption(o, serviceOptions)
	} else if tv.msgCtx.item != nil {
		option = tv.feedOption(o, messageOptions)
	} else {
		option = tv.feedOption(o, fileOptions)

		if o.Name == "go_package" {
			if o.Constant.Source == "" {
				tv.errors(errPosf(o.Constant.Position, "missing go_package option value"))
			}
			parts := strings.Split(o.Constant.Source, ";")
			switch len(parts) {
			case 1:
				tv.file.GoPkg = parts[0]
			case 2:
				tv.file.GoPath = parts[0]
				tv.file.GoPkg = parts[1]
			default:
				tv.errors(errPosf(
					o.Constant.Position,
					`invalid go_package option value, can only be either "<path>;<pkg>" or "<pkg>"`,
				))
			}
		}
	}

	if tv.serviceCtx != nil {
		tv.serviceCtx.Options = append(tv.serviceCtx.Options, option)
	} else if tv.msgCtx.item != nil {
		tv.msgCtx.item.Options = append(tv.msgCtx.item.Options, option)
	} else {
		tv.file.Options = append(tv.file.Options, option)
	}
	tv.regInfo(option, o.Comment, o.Position)
	tv.regFieldInfo(option, &option.Name, nil, o.Position)
	tv.regFieldInfo(option, &option.Value, nil, o.Constant.Position)
}

func (tv *typesVisitor) feedOption(o *proto.Option, opts optionType) *ast.Option {
	res := &ast.Option{
		Name:      o.Name,
		Extension: tv.optionLookup(o.Name, o.Position, opts),
	}

	res.Value = tv.literalToOptionValueWithExt(o.Name, &o.Constant, res.Extension)

	return res
}

func (tv *typesVisitor) literalToOptionValueWithExt(name string, l *proto.Literal, ext *ast.Extension) (result ast.OptionValue) {
	defer func() {
		if result != nil {
			tv.regInfo(result, nil, l.Position)
		}
	}()
	switch {
	case l.IsString:
		return &ast.StringOption{Value: l.Source}
	case len(l.Array) > 0:
		var res ast.ArrayOption
		for _, item := range l.Array {
			res.Value = append(res.Value, tv.literalToOptionValueWithExt(name, item, ext))
		}
		return &res
	case l.Source != "":
		if ext == nil {
			return &ast.EmbeddedOption{Value: l.Source}
		}
		shortName := getShortName(name)
		// здесь вычисляем реальный тип основываясь на записи в соответствующем расширении
		for _, f := range ext.Fields {
			f := f
			if f.Name == shortName {
				value := tv.fromType(f.Type, name, l)
				if value != nil {
					return value
				}
			}
		}
		tv.errors(errPosf(l.Position, "invalid literal value %s for option %s", l.Source, name))
	default:
		res := &ast.MapOption{
			Value: map[string]ast.OptionValue{},
		}
	outerLoop:
		for _, item := range l.OrderedMap {
			item := item
			shortName := getShortName(name)
			for _, f := range ext.Fields {
				f := f
				if f.Name == shortName {
					switch v := f.Type.(type) {
					case *ast.Message:
						for _, f := range v.Fields {
							f := f
							if vv, ok := f.Type.(*ast.OneOf); ok {
								for _, b := range vv.Branches {
									b := b
									if item.Name == b.Name {
										res.Value[item.Name] = tv.literalToOptionValueWithOneof(item.Name, item.Literal, vv)
										continue outerLoop
									}
								}
							}
							if item.Name == f.Name {
								res.Value[item.Name] = tv.literalToOptionValueWithMsg(item.Name, item.Literal, v)
								continue outerLoop
							}
						}
						tv.errors(errPosf(item.Position, "unknown option %s", item.Name))
					case *ast.Optional:
						msg := v.Type.(*ast.Message)
						for _, f := range msg.Fields {
							f := f
							if vv, ok := f.Type.(*ast.OneOf); ok {
								for _, b := range vv.Branches {
									if item.Name == b.Name {
										res.Value[item.Name] = tv.literalToOptionValueWithOneof(item.Name, item.Literal, vv)
										continue outerLoop
									}
								}
							}
							if item.Name == f.Name {
								res.Value[item.Name] = tv.literalToOptionValueWithMsg(item.Name, item.Literal, msg)
								continue outerLoop
							}
						}
						tv.errors(errPosf(item.Position, "unknown option %s", item.Name))
					default:
						tv.errors(errPosf(l.Position, "invalid type %T for option %s", f.Type, name))
						return nil
					}

				}
				continue outerLoop
			}
			tv.errors(errPosf(l.Position, "invalid option %s", name))
		}
		return res
	}

	return nil
}

func (tv *typesVisitor) literalToOptionValueWithMsg(name string, l *proto.Literal, msg *ast.Message) (result ast.OptionValue) {
	defer func() {
		if result != nil {
			tv.regInfo(result, nil, l.Position)
		}
	}()

	switch {
	case l.IsString:
		return &ast.StringOption{Value: l.Source}
	case len(l.Array) > 0:
		var res ast.ArrayOption
		for _, item := range l.Array {
			res.Value = append(res.Value, tv.literalToOptionValueWithMsg(name, item, msg))
		}
		return &res
	case len(l.OrderedMap) > 0:
		res := &ast.MapOption{
			Value: map[string]ast.OptionValue{},
		}
	outerLoop:
		for _, item := range l.OrderedMap {
			shortName := getShortName(name)
			for _, f := range msg.Fields {
				if f.Name == shortName {
					res.Value[item.Name] = tv.fromType(f.Type, shortName, l)
					continue outerLoop
				}
			}
			// res.Value[item.Name] = tv.literalToOptionValueWithMsg(item.Name, item.Literal, ext)
			tv.errors(errPosf(l.Position, "invalid option %s", name))
		}
		return res
	case l.Source != "":
		if msg == nil {
			return &ast.EmbeddedOption{Value: l.Source}
		}
		shortName := getShortName(name)
		// здесь вычисляем реальный тип основываясь на записи в соответствующем расширении
		for _, f := range msg.Fields {
			if f.Name == shortName {
				value := tv.fromType(f.Type, name, l)
				if value != nil {
					return value
				}
			}
		}
		tv.errors(errPosf(l.Position, "invalid literal value %s for option %s", l.Source, name))
	default:
		return &ast.MapOption{}
	}

	return nil
}

func (tv *typesVisitor) literalToOptionValueWithOneof(name string, l *proto.Literal, oo *ast.OneOf) (result ast.OptionValue) {
	defer func() {
		if result != nil {
			tv.regInfo(result, nil, l.Position)
		}
	}()

	switch {
	case l.IsString:
		return &ast.StringOption{Value: l.Source}
	case len(l.Array) > 0:
		var res ast.ArrayOption
		for _, item := range l.Array {
			res.Value = append(res.Value, tv.literalToOptionValueWithOneof(name, item, oo))
		}
		return &res
	case l.Source != "":
		if oo == nil {
			return &ast.EmbeddedOption{Value: l.Source}
		}
		shortName := getShortName(name)
		// здесь вычисляем реальный тип основываясь на записи в соответствующем расширении
		for _, b := range oo.Branches {
			if b.Name == shortName {
				value := tv.fromType(b.Type, name, l)
				if value != nil {
					return value
				}
			}
		}
		tv.errors(errPosf(l.Position, "invalid literal value %s for option %s", l.Source, name))
	default:
		res := &ast.MapOption{
			Value: map[string]ast.OptionValue{},
		}
	outerLoop:
		for _, item := range l.OrderedMap {
			shortName := getShortName(name)
			for _, b := range oo.Branches {
				if b.Name == shortName {
					res.Value[item.Name] = tv.fromType(b.Type, shortName, l)
					continue outerLoop
				}
			}
			// res.Value[item.Name] = tv.literalToOptionValueWithMsg(item.Name, item.Literal, ext)
			tv.errors(errPosf(l.Position, "invalid option %s", name))
		}
		return res
	}

	return nil
}

func getShortName(name string) string {
	items := strings.Split(strings.Trim(name, `()`), ".")
	return items[len(items)-1]
}

func (tv *typesVisitor) fromType(fieldType ast.Type, name string, l *proto.Literal) ast.OptionValue {
	switch v := fieldType.(type) {
	case *ast.Int32:
		return tv.intValue(name, l)
	case *ast.Int64:
		return tv.intValue(name, l)
	case *ast.Uint32:
		return tv.uintValue(name, l)
	case *ast.Uint64:
		return tv.uintValue(name, l)
	case *ast.Fixed32:
		return tv.intValue(name, l)
	case *ast.Fixed64:
		return tv.intValue(name, l)
	case *ast.Sfixed32:
		return tv.intValue(name, l)
	case *ast.Sfixed64:
		return tv.intValue(name, l)
	case *ast.Sint32:
		return tv.intValue(name, l)
	case *ast.Sint64:
		return tv.intValue(name, l)
	case *ast.Float32:
		return tv.floatValue(name, l)
	case *ast.Float64:
		return tv.floatValue(name, l)
	case *ast.Bool:
		return tv.boolValue(name, l)
	case *ast.String:
		return &ast.StringOption{Value: l.Source}
	case *ast.Enum:
		for _, ev := range v.Values {
			ev := ev
			if ev.Name == l.Source {
				return &ast.EnumOption{Value: ev}
			}
		}
	case *ast.Optional:
		return tv.fromType(v.Type, name, l)
	case *ast.Repeated:
		return tv.fromType(v.Type, name, l)
	default:
		tv.errors(errPosf(l.Position, "type %T is not supported for option values", fieldType))
	}
	return nil
}

func (tv *typesVisitor) intValue(name string, l *proto.Literal) *ast.IntOption {
	res, err := strconv.ParseInt(l.Source, 10, 64)
	if err != nil {
		tv.errors(errPosf(l.Position, "invalid value for option %s: %s", name, err))
	}
	return &ast.IntOption{Value: res}
}

func (tv *typesVisitor) uintValue(name string, l *proto.Literal) *ast.UintOption {
	res, err := strconv.ParseUint(l.Source, 10, 64)
	if err != nil {
		tv.errors(errPosf(l.Position, "invalid value for option %s: %s", name, err))
	}
	return &ast.UintOption{Value: res}
}

func (tv *typesVisitor) floatValue(name string, l *proto.Literal) *ast.FloatOption {
	res, err := strconv.ParseFloat(l.Source, 64)
	if err != nil {
		tv.errors(errPosf(l.Position, "invalid value for option %s: %s", name, err))
	}
	return &ast.FloatOption{Value: res}
}

func (tv *typesVisitor) boolValue(name string, l *proto.Literal) *ast.BoolOption {
	switch l.Source {
	case "true":
		return &ast.BoolOption{Value: true}
	case "false":
		return &ast.BoolOption{Value: false}
	}
	tv.errors(errPosf(l.Position, "invalid value for option %s: only 'true' or 'false' are allowed, got %s", name, l.Source))

	return nil
}

func (tv *typesVisitor) VisitImport(i *proto.Import) {
	importNs, _, err := tv.nss.get(i.Filename)
	if err != nil {
		tv.errors(errPosf(i.Position, "reading import %s: %s", i.Filename, err))
		return
	}
	importFile, err := tv.nss.AST(i.Filename)
	if err != nil {
		tv.errors(errPosf(i.Position, "reading import %s: %s", i.Filename, err))
	}

	tv.ns, err = tv.ns.WithImport(importNs)
	if err != nil {
		tv.errors(errPos(i.Position, err))
	}
	imp := &ast.Import{
		Path: i.Filename,
		File: importFile,
	}
	tv.file.Imports = append(tv.file.Imports, imp)
	tv.regInfo(imp, i.Comment, i.Position)
	tv.regFieldInfo(imp, &imp.Path, nil, i.Position)
}

func (tv *typesVisitor) VisitNormalField(i *proto.NormalField) {
	if prev, ok := tv.msgCtx.prevField[i.Name]; ok {
		tv.errors(errPosf(i.Position, "duplicate field %s, the previous definition was in %s", i.Name, prev))
	}
	if prev, ok := tv.msgCtx.prevInteger[i.Sequence]; ok {
		tv.errors(errPosf(i.Position, "duplicate field sequence %d, the previous value was in %s", i.Sequence, prev))
	}
	tv.msgCtx.prevField[i.Name] = i.Position
	tv.msgCtx.prevInteger[i.Sequence] = i.Position

	var options []*ast.Option
	for _, o := range i.Options {
		o := o
		option := tv.feedOption(o, fieldOptions)
		options = append(options, option)
		tv.regInfo(option, o.Comment, o.Position)
		tv.regFieldInfo(option, &option.Name, nil, o.Position)
		tv.regFieldInfo(option, &option.Value, nil, o.Constant.Position)
	}

	t := tv.standardType(i.Type)
	if t == nil {
		t = tv.ns.GetType(i.Type)
	}
	if t == nil {
		if strings.HasPrefix(i.Type, tv.file.Package+".") {
			t = tv.ns.GetType(i.Type[len(tv.file.Package)+1:])
		}
	}
	if t == nil {
		lastIndex := strings.LastIndex(i.Type, ".")
		if lastIndex >= 0 {
			curPkgParts := strings.Split(tv.file.Package, ".")
			for len(curPkgParts) > 0 {
				fullName := strings.Join(curPkgParts, ".") + "." + i.Type
				t = tv.ns.GetType(fullName)
				if t != nil {
					break
				}
				curPkgParts = curPkgParts[:len(curPkgParts)-1]
			}
		}
	}
	if t == nil {
		tv.errors(errPosf(i.Position, "unknown type %s", i.Type))
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

func (tv *typesVisitor) standardType(typeName string) ast.Type {
	switch typeName {
	case "bool":
		return &ast.Bool{}
	case "google.protobuf.Any":
		file, err := tv.nss.AST("google/protobuf/any.proto")
		if err != nil {
			tv.errors(errors.Wrap(err, "google.protobuf.Any must have google/protobuf/any.proto import"))
		}
		return &ast.Any{
			File: file,
		}
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
		tv.errors(errPosf(i.Position, "duplicate enum field %s, the previous definition was in %s", i.Name, prev))
	}
	if prev, ok := tv.enumCtx.prevInteger[i.Integer]; ok {
		tv.errors(errPosf(i.Position, "duplicate enum field id %d, the previous field with the same id was in %s", i.Integer, prev))
	}
	tv.enumCtx.prevField[i.Name] = i.Position
	tv.enumCtx.prevInteger[i.Integer] = i.Position
	var options []*ast.Option
	if i.ValueOption != nil {
		option := tv.feedOption(i.ValueOption, enumValueOptions)
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
		tv.errors(errPosf(o.Position, "duplicate field %s, the previous definition was in %s", o.Name, prev))
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
		tv.errors(errPosf(o.Position, "duplicate field %s, the previous definition was in %s", o.Name, prev))
	}
	if prev, ok := tv.msgCtx.prevInteger[o.Sequence]; ok {
		tv.errors(errPosf(o.Position, "duplicate field sequence %d, the previous valuy was in %s", o.Sequence, prev))
	}
	tv.msgCtx.prevField[o.Name] = o.Position
	tv.msgCtx.prevInteger[o.Sequence] = o.Position

	var options []*ast.Option
	for _, o := range o.Options {
		option := tv.feedOption(o, fieldOptions)
		options = append(options, option)
		tv.regInfo(option, o.Comment, o.Position)
		tv.regFieldInfo(option, &option.Name, nil, o.Position)
		tv.regFieldInfo(option, &option.Value, nil, o.Constant.Position)
	}

	t := tv.standardType(o.Type)
	if t == nil {
		t = tv.ns.GetType(o.Type)
	}
	if t == nil {
		lastIndex := strings.LastIndex(o.Type, ".")
		if lastIndex >= 0 {
			curPkgParts := strings.Split(tv.file.Package, ".")
			for len(curPkgParts) > 0 {
				fullName := strings.Join(curPkgParts, ".") + "." + o.Type
				t = tv.ns.GetType(fullName)
				if t != nil {
					break
				}
				curPkgParts = curPkgParts[:len(curPkgParts)-1]
			}
		}
	}
	if t == nil {
		tv.errors(errPosf(o.Position, "unknown type %s", o.Type))
		return
	}
	tv.regInfo(t, nil, o.Position)
	b := &ast.OneOfBranch{
		Name:     o.Name,
		Type:     t,
		ParentOO: tv.oneOf,
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

	rpc := &ast.Method{
		File:    tv.file,
		Service: tv.service,
		Name:    r.Name,
	}

	req := tv.ns.GetType(r.RequestType)
	if req == nil {
		fmt.Println("lol kek")
		tv.errors(errPosf(r.Position, "unknown type %s", r.RequestType))
		return
	}
	if r.StreamsRequest {
		req = &ast.Stream{
			Type: req,
		}
	}

	resp := tv.ns.GetType(r.ReturnsType)
	if resp == nil {
		tv.errors(errPosf(r.Position, "unknown type %s", r.RequestType))
		return
	}
	if r.StreamsReturns {
		resp = &ast.Stream{
			Type: resp,
		}
	}

	var mos []*ast.MethodOption
	for _, o := range r.Options {
		mo := &ast.MethodOption{
			Name:      o.Name,
			Extension: tv.optionLookup(o.Name, o.Position, methodOptions),
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

	rpc.Input = req
	rpc.Output = resp
	rpc.Options = mos

	tv.service.Methods = append(tv.service.Methods, rpc)
	tv.regInfo(rpc, r.Comment, r.Position)
	tv.regFieldInfo(rpc, &rpc.Name, nil, r.Position)
	tv.regFieldInfo(rpc, &rpc.Input, nil, r.Position)
	tv.regFieldInfo(rpc, &rpc.Output, nil, r.Position)

	if err := tv.ns.SetNode(tv.service.Name+"::"+r.Name, rpc, r.Position); err != nil {
		tv.errors(errPosf(r.Position, "duplicate method %s in service %s", r.Name, tv.service.Name))
		return
	}
}

func (tv *typesVisitor) VisitMapField(f *proto.MapField) {
	if prev, ok := tv.msgCtx.prevField[f.Name]; ok {
		tv.errors(errPosf(f.Position, "duplicate field %s, the previous definition was in %s", f.Name, prev))
	}
	if prev, ok := tv.msgCtx.prevInteger[f.Sequence]; ok {
		tv.errors(errPosf(f.Position, "duplicate field sequence %d, the previous valuy was in %s", f.Sequence, prev))
	}
	tv.msgCtx.prevField[f.Name] = f.Position
	tv.msgCtx.prevInteger[f.Sequence] = f.Position

	var options []*ast.Option
	for _, o := range f.Options {
		option := tv.feedOption(o, fieldOptions)
		tv.regInfo(option, o.Comment, o.Position)
		tv.regFieldInfo(option, &option.Name, nil, o.Position)
		tv.regFieldInfo(option, &option.Value, nil, o.Constant.Position)
		options = append(options, option)
	}

	keyRawType := tv.standardType(f.KeyType)
	if keyRawType == nil {
		tv.errors(errPosf(f.Position, "invalid map key type %s", f.Type))
		return
	}
	keyType, isHashable := keyRawType.(ast.Hashable)
	if !isHashable {
		tv.errors(errPosf(f.Position, "invalid map key type %s", f.Type))
	}
	tv.regInfo(keyType, nil, f.Position)

	t := tv.standardType(f.Type)
	if t == nil {
		t = tv.ns.GetType(f.Type)

	}
	if t == nil {
		lastIndex := strings.LastIndex(f.Type, ".")
		if lastIndex >= 0 {
			curPkgParts := strings.Split(tv.file.Package, ".")
			for len(curPkgParts) > 0 {
				fullName := strings.Join(curPkgParts, ".") + "." + f.Type
				t = tv.ns.GetType(fullName)
				if t != nil {
					break
				}
				curPkgParts = curPkgParts[:len(curPkgParts)-1]
			}
		}
	}
	if t == nil {
		tv.errors(errPosf(f.Position, "unknown value type %s", f.Type))
		return
	}
	tv.regInfo(t, nil, f.Position)

	fieldType := &ast.Map{
		KeyType:   keyType,
		ValueType: t,
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
