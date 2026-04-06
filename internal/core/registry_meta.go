package core

import (
	"iter"
	"slices"
	"text/scanner"

	"github.com/emicklei/proto"
	"github.com/sirkon/protoast/v2/internal/errors"
)

func (r *Registry) NodeIndex(node Node) string {
	switch n := node.(type) {
	case *File:
		return r.scopes[n.proto]
	case *Message:
		return r.scopes[n.proto]
	case *Enum:
		return r.scopes[n.proto]
	case *Service:
		return r.scopes[n.proto]
	case *Method:
		return r.scopes[n.proto]
	case *MessageField:
		switch m := n.proto.(type) {
		case *proto.NormalField:
			return r.scopes[m]
		case *proto.Oneof:
			return r.scopes[m]
		case *proto.MapField:
			return r.scopes[m]
		default:
			return ""
		}
	case *EnumValue:
		return r.scopes[n.proto]
	case *OneOf:
		return r.scopes[n.proto]
	case *OneOfBranch:
		return r.scopes[n.proto]
	case *Map:
		return r.scopes[n.proto]
	default:
		return ""
	}
}

func (r *Registry) NodeDescription(node Node) string {
	switch n := node.(type) {
	case *File:
		return "file"
	case *Syntax:
		return "syntax"
	case *Package:
		return "package"
	case *Import:
		return "import"
	case *Option:
		return "option"
	case *Message:
		return "message"
	case *Enum:
		return "enum"
	case *Service:
		return "service"
	case *Method:
		return "method"
	case *MessageField:
		return "message field"
	case *EnumValue:
		return "enum value"
	case *OneOf:
		return "oneof"
	case *OneOfBranch:
		return "oneof branch"
	case *Map:
		return "map[" + r.NodeDescription(n.Key()) + ", " + r.NodeDescription(n.Value(r)) + "] field"
	case *Repeated:
		return "[]" + r.NodeDescription(n.Type)
	case *Reserved:
		return "reserved"
	default:
		panic(errors.Newf("unsupported node type: %T", node))
	}
}

func (r *Registry) Comment(node Node) []string {
	switch n := node.(type) {
	case *Message:
		if n.proto.Comment != nil {
			return n.proto.Comment.Lines
		}
	case *MessageField:
		switch p := n.proto.(type) {
		case *proto.NormalField:
			if p.Comment != nil {
				return p.Comment.Lines
			}
		case *proto.Oneof:
			if p.Comment != nil {
				return p.Comment.Lines
			}
		case *proto.MapField:
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
	case *Reserved:
		if n.proto.Comment != nil {
			return n.proto.Comment.Lines
		}
	case *Import:
		if n.proto.Comment != nil {
			return n.proto.Comment.Lines
		}
	case *Syntax:
		if n.proto.Comment != nil {
			return n.proto.Comment.Lines
		}
	case *Package:
		if n.proto.Comment != nil {
			return n.proto.Comment.Lines
		}
	}

	return nil
}

func (r *Registry) Pos(pos Positionable) (res scanner.Position) {
	defer func() {
		if res.Filename != "" {
			return
		}

		node, ok := pos.(Node)
		if !ok {
			return
		}

		nodes := slices.Collect(r.NodeHierarchy(node))
		last := nodes[len(nodes)-1]
		file := last.(*File)
		res.Filename = file.proto.Filename
	}()

	switch n := pos.(type) {
	case *File:
		return scanner.Position{
			Filename: n.proto.Filename,
			Line:     1,
		}
	case *Message:
		return n.proto.Position
	case *MessageField:
		switch p := n.proto.(type) {
		case *proto.NormalField:
			return p.Position
		case *proto.Oneof:
			return p.Position
		case *proto.MapField:
			return p.Position
		default:
			panic(errors.Newf("unsupported message field type: %T", p))
		}
	case *Option:
		return n.proto.Position
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
	case *Reserved:
		return n.proto.Position
	case *Import:
		return n.proto.Position
	case *Syntax:
		return n.proto.Position
	case *Package:
		return n.proto.Position
	default:
		panic(errors.Newf("unsupported node type: %T", n))
	}
}

func (r *Registry) Options(node NodeOptionable) iter.Seq[*Option] {
	switch n := node.(type) {
	case *File:
		return seqOptions(r, n.Package(), registryOptionsFile, n.proto.Elements)
	case *Message:
		scope := r.scopes[n.proto]
		return seqOptions(r, scope, registryOptionsMessage, n.proto.Elements)
	case *MessageField:
		switch p := n.proto.(type) {
		case *proto.NormalField:
			return seqOptions(r, r.scopes[p], registryOptionsMessageField, p.Options)
		case *proto.Oneof:
			return seqOptions(r, r.scopes[p], registryOptionsOneof, p.Elements)
		case *proto.MapField:
			return seqOptions(r, r.scopes[p], registryOptionsMessageField, p.Options)
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

func (r *Registry) OptionNamed(node NodeOptionable, name string) *Option {
	switch n := node.(type) {
	case *File:
		return namedOption(r, name, n.Package(), registryOptionsFile, n.proto.Elements)
	case *Message:
		scope := r.scopes[n.proto]
		return namedOption(r, name, scope, registryOptionsMessage, n.proto.Elements)
	case *MessageField:
		switch p := n.proto.(type) {
		case *proto.Message:
			return namedOption(r, name, r.scopes[p], registryOptionsMessageField, p.Elements)
		case *proto.Oneof:
			return namedOption(r, name, r.scopes[p], registryOptionsOneof, p.Elements)
		case *proto.MapField:
			return namedOption(r, name, r.scopes[p], registryOptionsMessageField, p.Options)
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
