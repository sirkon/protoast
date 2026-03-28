package core

import (
	"fmt"
)

// Node represents any ast item.
type Node interface {
	isNodeType()
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

type isNode struct{}

type isNodeOptionable struct{}

type isFieldNode struct {
	isNode
}

type isType struct {
	isNode
}

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

func (*isNode) isNodeType()                     {}
func (*isNodeOptionable) isNodeOptionableType() {}
func (*isFieldNode) isFieldNodeType()           {}
func (*isType) isTypeType()                     {}
func (*isBuiltinType) isBuiltinTypeType()       {}
func (*isNamedType) isNamedTypeType()           {}
func (*isComparableType) isComparableTypeType() {}
func (*isComposableType) isComposableTypeType() {}
