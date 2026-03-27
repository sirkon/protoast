package core

import (
	"iter"
	"text/scanner"

	"github.com/sirkon/protoast/v2/internal/errors"
)

func (r *Registry) Comment(node Node) []string {
	switch n := node.(type) {
	case *Message:
		if n.proto.Comment != nil {
			return n.proto.Comment.Lines
		}
	case *MessageField:
		switch p := n.payload.(type) {
		case *isEmickleiNormalField:
			if p.Comment != nil {
				return p.Comment.Lines
			}
		case *isEmickleiOneOf:
			if p.Comment != nil {
				return p.Comment.Lines
			}
		case *isEmickleiMapField:
			if p.Comment != nil {
				return p.Comment.Lines
			}
		default:
			panic(errors.Newf("unsupported field type: %T", p))
		}
	case *Enum:
		if n.proto.Comment != nil {
			return n.proto.Comment.Lines
		}
	case *EnumValue:
		if n.proto.Comment != nil {
			return n.proto.Comment.Lines
		}
	case *OneOf:
		if n.proto.Comment != nil {
			return n.proto.Comment.Lines
		}
	case *OneOfBranch:
		if n.proto.Comment != nil {
			return n.proto.Comment.Lines
		}
	case *Map:
		if n.proto.Comment != nil {
			return n.proto.Comment.Lines
		}
	case *Service:
		if n.proto.Comment != nil {
			return n.proto.Comment.Lines
		}
	case *Method:
		if n.proto.Comment != nil {
			return n.proto.Comment.Lines
		}
	}

	return nil
}

func (r *Registry) Pos(node Node) scanner.Position {
	switch n := node.(type) {
	case *File:
		return scanner.Position{
			Filename: n.proto.Filename,
			Line:     1,
		}
	case *Message:
		return n.proto.Position
	case *MessageField:
		switch p := n.payload.(type) {
		case *isEmickleiNormalField:
			return p.asProto().Position
		case *isEmickleiOneOf:
			return p.asProto().Position
		case *isEmickleiMapField:
			return p.asProto().Position
		default:
			panic(errors.Newf("unsupported message field type: %T", p))
		}
	case *Option:
		return n.protoOption.Position
	case *OptionValue:
		return n.option.protoOption.Constant.Position
	case *OptionValueBool:
		return n.proto.Position
	case *OptionValueInt:
		return n.proto.Position
	case *OptionValueUint:
		return n.proto.Position
	case *OptionValueFloat:
		return n.proto.Position
	case *OptionValueString:
		return n.proto.Position
	case *OptionValueBytes:
		return n.proto.Position
	case *OptionValueEnum:
		return n.proto.Position
	case *OptionValueArray:
		return n.proto.Position
	case *OptionValueMap:
		return n.proto.Position
	case *OptionValueMapItem:
		return n.proto.Position
	case *Enum:
		return n.proto.Position
	case *EnumValue:
		return n.proto.Position
	case *OneOf:
		return n.proto.Position
	case *OneOfBranch:
		return n.proto.Position
	case *Map:
		return n.proto.Position
	case *Service:
		return n.proto.Position
	case *Method:
		return n.proto.Position
	default:
		panic(errors.Newf("unsupported node type: %T", n))
	}
}

func (r *Registry) Options(node Node) iter.Seq[*Option] {
	switch n := node.(type) {
	case *File:
		return seqOptions(r, n.Package(), registryOptionsFile, n.proto.Elements)
	case *Message:
		scope := r.scopes[n.proto]
		return seqOptions(r, scope, registryOptionsMessage, n.proto.Elements)
	case *MessageField:
		switch p := n.payload.(type) {
		case *isEmickleiNormalField:
			return seqOptions(r, r.scopes[p.asProto()], registryOptionsMessageFields, p.Options)
		case *isEmickleiOneOf:
			return seqOptions(r, r.scopes[p.asProto()], registryOptionsOneof, p.Elements)
		case *isEmickleiMapField:
			return seqOptions(r, r.scopes[p.asProto()], registryOptionsMessageFields, p.Options)
		default:
			panic(errors.Newf("unsupported payload type: %T", n))
		}
	case *Enum:
		scope := r.scopes[n.proto]
		return seqOptions(r, scope, registryOptionsEnum, n.proto.Elements)
	case *EnumValue:
		scope := r.scopes[n.proto]
		return seqOptions(r, scope, registryOptionsEnumValue, n.proto.Elements)
	case *OneOf:
		scope := r.scopes[n.proto]
		return seqOptions(r, scope, registryOptionsOneof, n.proto.Elements)
	case *Service:
		scope := r.scopes[n.proto]
		return seqOptions(r, scope, registryOptionsService, n.proto.Elements)
	case *Method:
		scope := r.scopes[n.proto]
		return seqOptions(r, scope, registryOptionsMethod, n.proto.Elements)
	default:
		panic(errors.Newf("unsupported node type: %T", n))
	}
}

func (r *Registry) OptionNamed(node Node, name string) *Option {
	switch n := node.(type) {
	case *File:
		return namedOption(r, name, n.Package(), registryOptionsFile, n.proto.Elements)
	case *Message:
		scope := r.scopes[n.proto]
		return namedOption(r, name, scope, registryOptionsMessage, n.proto.Elements)
	case *MessageField:
		switch p := n.payload.(type) {
		case *isEmickleiNormalField:
			return namedOption(r, name, r.scopes[p.asProto()], registryOptionsMessageFields, p.Options)
		case *isEmickleiOneOf:
			return namedOption(r, name, r.scopes[p.asProto()], registryOptionsOneof, p.Elements)
		case *isEmickleiMapField:
			return namedOption(r, name, r.scopes[p.asProto()], registryOptionsMessageFields, p.Options)
		default:
			panic(errors.Newf("unsupported payload type: %T", n))
		}
	case *Enum:
		scope := r.scopes[n.proto]
		return namedOption(r, name, scope, registryOptionsEnum, n.proto.Elements)
	case *EnumValue:
		scope := r.scopes[n.proto]
		return namedOption(r, name, scope, registryOptionsEnumValue, n.proto.Elements)
	case *OneOf:
		scope := r.scopes[n.proto]
		return namedOption(r, name, scope, registryOptionsOneof, n.proto.Elements)
	case *Service:
		scope := r.scopes[n.proto]
		return namedOption(r, name, scope, registryOptionsService, n.proto.Elements)
	case *Method:
		scope := r.scopes[n.proto]
		return namedOption(r, name, scope, registryOptionsMethod, n.proto.Elements)
	default:
		panic(errors.Newf("unsupported node type: %T", n))
	}
}
