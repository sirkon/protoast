package ast

import "strconv"

var _ Unique = &Option{}

// Option опция
type Option struct {
	unique

	Name      string
	Value     OptionValue
	Extension *Extension
}

// OptionValue значение опции
type OptionValue interface {
	Unique

	isOptionValue()
}

// EmbeddedOption представление встроенной опции
type EmbeddedOption struct {
	unique
	Value string
}

func (*EmbeddedOption) isOptionValue() {}

// EnumOption представление опций типа Enum
type EnumOption struct {
	unique
	Value *EnumValue
}

func (*EnumOption) isOptionValue() {}

// IntOption branch of OptionValue
type IntOption struct {
	unique
	Value int64
}

func (*IntOption) isOptionValue() {}

func (o *IntOption) String() string {
	return strconv.FormatInt(o.Value, 10)
}

// UintOption branch of OptionValue
type UintOption struct {
	unique
	Value uint64
}

func (*UintOption) isOptionValue() {}

func (o *UintOption) String() string {
	return strconv.FormatUint(o.Value, 10)
}

// FloatOption branch of OptionValue
type FloatOption struct {
	unique
	Value float64
}

func (*FloatOption) isOptionValue() {}

func (o *FloatOption) String() string {
	return strconv.FormatFloat(o.Value, ' ', 5, 64)
}

// StringOption branch of OptionValue
type StringOption struct {
	unique
	Value string
}

func (*StringOption) isOptionValue() {}

func (o *StringOption) String() string {
	return o.Value
}

// BoolOption branch of OptionValue
type BoolOption struct {
	unique
	Value bool
}

func (o *BoolOption) String() string {
	return strconv.FormatBool(o.Value)
}

func (*BoolOption) isOptionValue() {}

// ArrayOption branch of OptionValue
type ArrayOption struct {
	unique
	Value []OptionValue
}

func (*ArrayOption) isOptionValue() {}

// MapOption branch of OptionValue
type MapOption struct {
	unique
	Value map[string]OptionValue
}

func (*MapOption) isOptionValue() {}
