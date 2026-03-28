package core

import (
	"fmt"
	"text/scanner"

	"github.com/emicklei/proto"
)

type Positionable interface {
	pos() scanner.Position
}

// Node represents any ast item.
type Node interface {
	Positionable
	nodeProto() proto.Visitee
}

type NodeOptionable interface {
	Node
	isNodeOptionableType()
}

// FieldNode represents message fields, oneof branches and enum values.
type FieldNode interface {
	Node
	isFieldNodeType()
}

// Type represents all types.
type Type interface {
	Node
	isTypeType()
}

// BuiltinType represents builtin types.
type BuiltinType interface {
	fmt.Stringer
	ComposableType
	isBuiltinTypeType()
}

// ComparableType represents types with comparable value, meaning
// every builtin type except bytes.
type ComparableType interface {
	BuiltinType
	isComparableTypeType()
}

// ComposableType represents builtins and named types.
type ComposableType interface {
	Type
	isComposableTypeType()
}

// NamedType represents named types. Meaning messages and enums.
type NamedType interface {
	ComposableType
	isNamedTypeType()
}

type isNodeOptionable struct{}

type isFieldNode struct{}

type isType struct{}

type isBuiltinType struct {
	isComposableType
}

type isNamedType struct {
	isComposableType
}

type isComposableType struct {
	isType
}

type isComparableType struct {
	isBuiltinType
}

func (*isNodeOptionable) isNodeOptionableType() {}
func (*isFieldNode) isFieldNodeType()           {}
func (*isType) isTypeType()                     {}
func (*isBuiltinType) isBuiltinTypeType()       {}
func (*isNamedType) isNamedTypeType()           {}
func (*isComparableType) isComparableTypeType() {}
func (*isComposableType) isComposableTypeType() {}

func (*isBuiltinType) nodeProto() proto.Visitee { return nil }
func (*isBuiltinType) pos() scanner.Position    { return scanner.Position{} }

var _ Node = new(isBuiltinType)
