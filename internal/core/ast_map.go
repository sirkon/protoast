package core

import (
	"strings"
	"text/scanner"

	"github.com/emicklei/proto"

	"github.com/sirkon/protoast/v2/internal/errors"
)

type Map struct {
	isType

	proto *proto.MapField
}

// Key returns map key type.
func (m *Map) Key() ComparableType {
	keyType := builtinComparableType(m.proto.KeyType)
	if keyType == nil {
		panic(errors.Newf("key type %s is not supported in maps", m.proto.KeyType))
	}

	return keyType
}

// Value returns map value type.
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
		return r.wrap(t).(*Message)
	case *proto.Enum:
		return r.wrap(t).(*Enum)
	default:
		panic(errors.Newf("value type %s is not supported in maps", m.proto.Type))
	}
}

var _ Node = new(Map)

func (m *Map) nodeProto() proto.Visitee { return m.proto }
func (m *Map) pos() scanner.Position    { return m.proto.Position }
