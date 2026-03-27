package core

import (
	"iter"

	"github.com/emicklei/proto"
)

// Enum represents enum named type.
type Enum struct {
	isNamedType

	proto *proto.Enum
}

func (e *Enum) Name() string {
	return e.proto.Name
}

func (e *Enum) Values() iter.Seq[*EnumValue] {
	return func(yield func(*EnumValue) bool) {
		for _, value := range e.proto.Elements {
			vv, ok := value.(*proto.EnumField)
			if !ok {
				continue
			}

			if !yield(&EnumValue{
				proto: vv,
			}) {
				return
			}
		}
	}
}

func (e *Enum) Value(name string) *EnumValue {
	for _, e := range e.proto.Elements {
		v, ok := e.(*proto.EnumField)
		if !ok {
			continue
		}

		if name != v.Name {
			continue
		}

		return &EnumValue{
			proto: v,
		}
	}

	return nil
}

type EnumValue struct {
	isFieldNode

	proto *proto.EnumField
}

func (e *EnumValue) Name() string {
	return e.proto.Name
}

func (e *EnumValue) Value() int {
	return e.proto.Integer
}
