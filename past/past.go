package past

import (
	"github.com/sirkon/protoast/v2/internal/core"
)

type (
	Node           = core.Node
	NodeOptionable = core.NodeOptionable
	FieldNode      = core.FieldNode
	Type           = core.Type
	BuiltinType    = core.BuiltinType
	ComparableType = core.ComparableType
	ComposableType = core.ComposableType
	NamedType      = core.NamedType

	File               = core.File
	Service            = core.Service
	Method             = core.Method
	Message            = core.Message
	MessageField       = core.MessageField
	Enum               = core.Enum
	EnumValue          = core.EnumValue
	Map                = core.Map
	Repeated           = core.Repeated
	OneOf              = core.OneOf
	OneOfBranch        = core.OneOfBranch
	Option             = core.Option
	OptionValue        = core.Option
	OptionValueVariant = core.OptionValueVariant
	Reserved           = core.Reserved
	Import             = core.Import
	Syntax             = core.Syntax
	Package            = core.Package

	OptionValueBool    = core.OptionValueBool
	OptionValueInt     = core.OptionValueInt
	OptionValueUint    = core.OptionValueUint
	OptionValueFloat   = core.OptionValueFloat
	OptionValueString  = core.OptionValueString
	OptionValueBytes   = core.OptionValueBytes
	OptionValueEnum    = core.OptionValueEnum
	OptionValueArray   = core.OptionValueArray
	OptionValueMap     = core.OptionValueMap
	OptionValueMapItem = core.OptionValueMapItem
)
