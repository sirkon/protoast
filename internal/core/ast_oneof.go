package core

import (
	"iter"

	"github.com/emicklei/proto"
)

type OneOf struct {
	isType

	proto *proto.Oneof
}

type OneOfBranch struct {
	isFieldNode

	proto *proto.OneOfField
}

func (o *OneOf) Branches() iter.Seq[*OneOfBranch] {
	return func(yield func(*OneOfBranch) bool) {
		for _, e := range o.proto.Elements {
			v, ok := e.(*proto.OneOfField)
			if !ok {
				continue
			}

			if !yield(&OneOfBranch{
				proto: v,
			}) {
				return
			}
		}
	}
}

func (o *OneOf) Branch(name string) *OneOfBranch {
	for _, e := range o.proto.Elements {
		v, ok := e.(*proto.OneOfField)
		if !ok {
			continue
		}

		if name != v.Name {
			continue
		}

		return &OneOfBranch{
			proto: v,
		}
	}

	return nil
}
