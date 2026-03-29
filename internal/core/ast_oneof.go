package core

import (
	"iter"
	"text/scanner"

	"github.com/emicklei/proto"
)

type OneOf struct {
	isType
	isNodeOptionable

	proto *proto.Oneof
}

type OneOfBranch struct {
	isFieldNode

	proto *proto.OneOfField
}

// Branches returns all branches.
func (o *OneOf) Branches(r *Registry) iter.Seq[*OneOfBranch] {
	return func(yield func(*OneOfBranch) bool) {
		for _, e := range o.proto.Elements {
			v, ok := e.(*proto.OneOfField)
			if !ok {
				continue
			}

			if !yield(r.wrap(v).(*OneOfBranch)) {
				return
			}
		}
	}
}

// Branch returns a branch with the given name.
func (o *OneOf) Branch(r *Registry, name string) *OneOfBranch {
	for _, e := range o.proto.Elements {
		v, ok := e.(*proto.OneOfField)
		if !ok {
			continue
		}

		if name != v.Name {
			continue
		}

		return r.wrap(v).(*OneOfBranch)
	}

	return nil
}

func (o *OneOfBranch) Name() string { return o.proto.Name }

func (o *OneOfBranch) Type(r *Registry) Type {
	return r.wrap(o.proto).(Type)
}

func (o *OneOfBranch) Value() int {
	return o.proto.Sequence
}

var (
	_ Node = new(OneOf)
	_ Node = new(OneOfBranch)
)

func (o *OneOf) nodeProto() proto.Visitee       { return o.proto }
func (o *OneOf) pos() scanner.Position          { return o.proto.Position }
func (o *OneOfBranch) nodeProto() proto.Visitee { return o.proto }
func (o *OneOfBranch) pos() scanner.Position    { return o.proto.Position }
