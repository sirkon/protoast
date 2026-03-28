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

	Bool     = core.Bool
	Int32    = core.Int32
	Int64    = core.Int64
	Sint32   = core.Sint32
	Sint64   = core.Sint64
	Sfixed32 = core.Sfixed32
	Sfixed64 = core.Sfixed64
	Uint32   = core.Uint32
	Uint64   = core.Uint64
	Fixed32  = core.Fixed32
	Fixed64  = core.Fixed64
	Float    = core.Float
	Double   = core.Double
	String   = core.String
	Bytes    = core.Bytes

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
