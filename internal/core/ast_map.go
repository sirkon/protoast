package core

import (
	"strings"

	"github.com/emicklei/proto"

	"github.com/sirkon/protoast/v2/internal/errors"
)

type Map struct {
	isType

	proto *proto.MapField
}

func (m *Map) Key() ComparableType {
	keyType := builtinComparableType(m.proto.KeyType)
	if keyType == nil {
		panic(errors.Newf("key type %s is not supported in maps", m.proto.KeyType))
	}

	return keyType
}

func (m *Map) Value(r *Registry) ComposableType {
	valueType := builtinType(m.proto.Type)
	if valueType != nil {
		return valueType
	}

	var qualifiedTypeName string
	if !strings.HasPrefix(m.proto.Type, ".") {
		qualifiedTypeName = "." + m.proto.Type
	}

	visitee := r.registry[qualifiedTypeName]
	switch t := visitee.(type) {
	case *proto.Message:
		return &Message{
			proto: t,
		}
	case *proto.Enum:
		return &Enum{
			proto: t,
		}
	default:
		panic(errors.Newf("value type %s is not supported in maps", m.proto.Type))
	}
}
