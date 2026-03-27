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
		return wrapEmicklei(n.proto.Parent)
	case *MessageField:
		switch p := n.payload.(type) {
		case *isEmickleiNormalField:
			return wrapEmicklei(p.asProto().Parent)
		case *isEmickleiOneOf:
			return wrapEmicklei(p.asProto().Parent)
		case *isEmickleiMapField:
			return wrapEmicklei(p.asProto().Parent)
		default:
			panic(errors.Newf("unsupported message field type: %T", p))
		}
	case *Enum:
		return wrapEmicklei(n.proto.Parent)
	case *EnumValue:
		return wrapEmicklei(n.proto.Parent)
	case *OneOf:
		return wrapEmicklei(n.proto.Parent)
	case *OneOfBranch:
		return wrapEmicklei(n.proto.Parent)
	case *Service:
		return wrapEmicklei(n.proto.Parent)
	case *Method:
		return wrapEmicklei(n.proto.Parent)
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

func (r *Registry) TypeIsDefined(typ Type, ref string) bool {
	switch t := typ.(type) {
	case *Message:
		return t.proto == r.registry[ref]
	case *Enum:
		return t.proto == r.registry[ref]
	}

	return false
}

func (r *Registry) TypeIsGoogleProtobufAny(typ Type) bool {
	return r.TypeIsDefined(typ, ".google.protobuf.Any")
}

func (r *Registry) TypeIsGoogleProtobufEmpty(typ Type) bool {
	return r.TypeIsDefined(typ, ".google.protobuf.Empty")
}

func (r *Registry) TypeIsGoogleProtobufTimestamp(typ Type) bool {
	return r.TypeIsDefined(typ, ".google.protobuf.Timestamp")
}

func (r *Registry) TypeIsGoogleProtobufDuration(typ Type) bool {
	return r.TypeIsDefined(typ, ".google.protobuf.Duration")
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
	resolveName, ok := r.resolveName(scope, name)
	if !ok {
		return nil
	}

	obj := r.registry[resolveName]
	if obj == nil {
		return nil
	}

	if v, ok := wrapEmicklei(obj).(Type); ok {
		return v
	}

	panic(errors.Newf("node %T does not represent a type", obj))
}

func wrapEmicklei(t proto.Visitee) Node {
	switch n := t.(type) {
	case *proto.Message:
		return &Message{
			proto: n,
		}
	case *proto.Enum:
		return &Enum{
			proto: n,
		}
	case *proto.Oneof:
		return &OneOf{
			proto: n,
		}
	case *proto.MapField:
		return &Map{
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
	default:
		panic(errors.Newf("unsupported type %T", t))
	}
}

func builtinType(name string) BuiltinType {
	res := builtinComparableType(name)
	if res != nil {
		return res
	}

	if name == "bytes" {
		return &Bytes{}
	}

	return nil
}

func builtinComparableType(name string) ComparableType {
	switch name {
	case "double":
		return &Double{}
	case "float":
		return &Float{}
	case "int32":
		return &Int32{}
	case "int64":
		return &Int64{}
	case "uint32":
		return &Uint32{}
	case "uint64":
		return &Uint64{}
	case "sint32":
		return &Sint32{}
	case "sint64":
		return &Sint64{}
	case "fixed32":
		return &Fixed32{}
	case "fixed64":
		return &Fixed64{}
	case "sfixed32":
		return &Sfixed32{}
	case "sfixed64":
		return &Sfixed64{}
	case "bool":
		return &Bool{}
	case "string":
		return &String{}
	}

	return nil
}
