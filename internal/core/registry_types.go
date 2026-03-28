package core

import (
	"fmt"
	"iter"

	"github.com/emicklei/proto"

	"github.com/sirkon/protoast/v2/internal/errors"
)

func (r *Registry) TypeName(typ Type) string {
	switch t := typ.(type) {
	case *Message:
		return r.scopes[t.proto]
	case *Enum:
		return r.scopes[t.proto]
	case *Repeated:
		return "repeated " + r.TypeName(t.Type)
	case *OneOf:
		return "oneof"
	case BuiltinType:
		return t.String()
	default:
		panic(fmt.Sprintf("unsupported type %T", t))
	}
}

func (r *Registry) NodeParent(node Node) Node {
	switch n := node.(type) {
	case *Repeated:
		return r.NodeParent(n.Type)
	case *Message:
		return r.wrap(n.proto.Parent)
	case *MessageField:
		switch p := n.proto.(type) {
		case *proto.NormalField:
			return r.wrap(p.Parent)
		case *proto.Oneof:
			return r.wrap(p.Parent)
		case *proto.MapField:
			return r.wrap(p.Parent)
		default:
			panic(errors.Newf("unsupported message field type: %T", p))
		}
	case *Enum:
		return r.wrap(n.proto.Parent)
	case *EnumValue:
		return r.wrap(n.proto.Parent)
	case *OneOf:
		return r.wrap(n.proto.Parent)
	case *OneOfBranch:
		return r.wrap(n.proto.Parent)
	case *Service:
		return r.wrap(n.proto.Parent)
	case *Method:
		return r.wrap(n.proto.Parent)
	case *Reserved:
		return r.wrap(n.proto.Parent)
	case *Import:
		return r.wrap(n.proto.Parent)
	case *Syntax:
		return r.wrap(n.proto.Parent)
	case *Package:
		return r.wrap(n.proto.Parent)
	case *Option:
		return r.wrap(n.protoOption.Parent)
	case *OptionValue:
		switch w := n.option.protoOption.Parent.(type) {
		case *proto.Option:
			return r.wrapOption(w, n.option.protoOptionClass)
		default:
			return r.wrap(n.option.protoOption.Parent)
		}
	case *File:
		return nil
	default:
		panic(fmt.Sprintf("unsupported node type: %T", n))
	}
}

func (r *Registry) NodeHierarchy(node Node) iter.Seq[Node] {
	return func(yield func(Node) bool) {
		for node != nil {
			if !yield(node) {
				return
			}

			node = r.NodeParent(node)
		}
	}
}

func (r *Registry) getTypeByName(scopeOf proto.Visitee, name string) (res Type) {
	defer func() {
		if v, ok := scopeOf.(*proto.NormalField); ok {
			if v.Repeated {
				res = &Repeated{
					Type: res,
				}
			}
		}
	}()

	if res := builtinType(name); res != nil {
		return res
	}

	scope := r.scopes[scopeOf]
	resolveName, ok := r.resolveNameRaw(scope, name)
	if !ok {
		return nil
	}

	obj := r.registry[resolveName]
	if obj == nil {
		return nil
	}

	if v, ok := r.wrap(obj).(Type); ok {
		return v
	}

	panic(errors.Newf("node %T does not represent a type", obj))
}

func (r *Registry) wrap(t proto.Visitee) (res Node) {
	if v, ok := r.cache[t]; ok {
		return v
	}
	defer func() {
		r.cache[t] = res
	}()

	switch n := t.(type) {
	case *proto.Message:
		return &Message{
			proto: n,
		}
	case *proto.NormalField:
		return &MessageField{
			proto: n,
		}
	case *proto.Enum:
		return &Enum{
			proto: n,
		}
	case *proto.EnumField:
		return &EnumValue{
			proto: n,
		}
	case *proto.Oneof:
		return &MessageField{
			proto: n,
		}
	case *proto.OneOfField:
		return &OneOfBranch{
			proto: n,
		}
	case *proto.MapField:
		return &MessageField{
			proto: n,
		}
	case *proto.Proto:
		return &File{
			proto: n,
		}
	case *proto.Service:
		return &Service{
			proto: n,
		}
	case *proto.RPC:
		return &Method{
			proto: n,
		}
	case *proto.Import:
		return &Import{
			proto: n,
		}
	case *proto.Reserved:
		return &Reserved{
			proto: n,
		}
	case *proto.Syntax:
		return &Syntax{
			proto: n,
		}
	case *proto.Package:
		return &Package{
			proto: n,
		}
	default:
		panic(errors.Newf("unsupported type %T", t))
	}
}

func (r *Registry) wrapOption(option *proto.Option, where *proto.Message) Node {
	return newOption(r, r.scopes[option], where, option)
}

func builtinType(name string) BuiltinType {
	res := builtinComparableType(name)
	if res != nil {
		return res
	}

	if name == "bytes" {
		return bytesTypeValue
	}

	return nil
}

func builtinComparableType(name string) ComparableType {
	switch name {
	case "float":
		return floatTypeValue
	case "double":
		return doubleTypeValue
	case "int32":
		return int32TypeValue
	case "int64":
		return int64TypeValue
	case "uint32":
		return uint32TypeValue
	case "uint64":
		return uint64TypeValue
	case "sint32":
		return sint32TypeValue
	case "sint64":
		return sint64TypeValue
	case "fixed32":
		return fixed32TypeValue
	case "fixed64":
		return fixed64TypeValue
	case "sfixed32":
		return sfixed32TypeValue
	case "sfixed64":
		return sfixed64TypeValue
	case "bool":
		return boolTypeValue
	case "string":
		return stringTypeValue
	}

	return nil
}

var (
	boolTypeValue     = &Bool{}
	floatTypeValue    = &Float{}
	doubleTypeValue   = &Double{}
	int32TypeValue    = &Int32{}
	int64TypeValue    = &Int64{}
	uint32TypeValue   = &Uint32{}
	uint64TypeValue   = &Uint64{}
	sint32TypeValue   = &Sint32{}
	sint64TypeValue   = &Sint64{}
	fixed32TypeValue  = &Fixed32{}
	fixed64TypeValue  = &Fixed64{}
	sfixed32TypeValue = &Sfixed32{}
	sfixed64TypeValue = &Sfixed64{}
	stringTypeValue   = &String{}
	bytesTypeValue    = &Bytes{}
)
