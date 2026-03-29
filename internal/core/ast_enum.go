package core

import (
	"iter"
	"text/scanner"

	"github.com/emicklei/proto"
)

// Enum represents enum named type.
type Enum struct {
	isNamedType
	isNodeOptionable

	proto *proto.Enum
}

type EnumValue struct {
	isFieldNode
	isNodeOptionable
	isType

	proto *proto.EnumField
}

// Name returns enum name.
func (e *Enum) Name() string {
	return e.proto.Name
}

// Values returns all enum values.
func (e *Enum) Values(r *Registry) iter.Seq[*EnumValue] {
	return func(yield func(*EnumValue) bool) {
		for _, value := range e.proto.Elements {
			vv, ok := value.(*proto.EnumField)
			if !ok {
				continue
			}

			v := r.wrap(vv)
			if !yield(v.(*EnumValue)) {
				return
			}
		}
	}
}

// Value returns enum value with the given name.
func (e *Enum) Value(r *Registry, name string) *EnumValue {
	for _, e := range e.proto.Elements {
		v, ok := e.(*proto.EnumField)
		if !ok {
			continue
		}

		if name != v.Name {
			continue
		}

		return r.wrap(v).(*EnumValue)
	}

	return nil
}

// Everything returns everything defined in this enum.
func (e *Enum) Everything(r *Registry) iter.Seq[Node] {
	return func(yield func(Node) bool) {
		for _, v := range e.proto.Elements {
			if vv, ok := v.(*proto.Option); ok {
				if !yield(r.wrapOption(vv, r.optionContextEnum())) {
					return
				}
				continue
			}

			if !yield(r.wrap(v)) {
				return
			}
		}
	}
}

func (e *EnumValue) Name() string {
	return e.proto.Name
}

func (e *EnumValue) Value() int {
	return e.proto.Integer
}

var _ Node = new(Enum)

var _ Node = new(EnumValue)

func (e *Enum) nodeProto() proto.Visitee      { return e.proto }
func (e *Enum) pos() scanner.Position         { return e.proto.Position }
func (e *EnumValue) nodeProto() proto.Visitee { return e.proto }
func (e *EnumValue) pos() scanner.Position    { return e.proto.Position }
