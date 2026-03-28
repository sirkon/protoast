package core

import (
	"fmt"
	"strconv"
	"strings"
	"text/scanner"

	"github.com/emicklei/proto"

	"github.com/sirkon/protoast/v2/internal/errors"
)

type OptionValue struct {
	isNode

	option *Option
}

func (o OptionValue) Value() OptionValueVariant {
	return buildFromLiteral(o.option.registry, o.option.protoOptionField, &o.option.protoOption.Constant, false)
}

func buildFromLiteral(
	r *Registry,
	field *proto.NormalField,
	literal *proto.Literal,
	ignoreRepeat bool,
) OptionValueVariant {
	if field.Repeated && !ignoreRepeat {
		var res []OptionValueVariant
		for _, l := range literal.Array {
			res = append(res, buildFromLiteral(r, field, l, true))
		}
		return &OptionValueArray{
			Value: res,
		}
	}

	switch field.Type {
	case "bool":
		val, err := strconv.ParseBool(literal.Source)
		if err != nil {
			panic(errors.Wrap(err, "convert literal to bool"))
		}
		return &OptionValueBool{
			proto: literal,
			Value: val,
		}
	case "int32", "sint32", "sfixed32", "int64", "sint64", "sfixed64":
		val, err := strconv.Atoi(literal.Source)
		if err != nil {
			panic(errors.Wrap(err, "convert literal to int"))
		}
		return &OptionValueInt{
			proto: literal,
			Value: val,
		}
	case "uin32", "fixed32", "uin64", "fixed64":
		val, err := strconv.ParseUint(literal.Source, 10, 64)
		if err != nil {
			panic(errors.Wrap(err, "convert literal to uint"))
		}
		return &OptionValueUint{
			proto: literal,
			Value: uint(val),
		}
	case "float", "double":
		val, err := strconv.ParseFloat(literal.Source, 64)
		if err != nil {
			panic(errors.Wrap(err, "convert literal to float"))
		}
		return &OptionValueFloat{
			proto: literal,
			Value: val,
		}
	case "string":
		return &OptionValueString{
			proto: literal,
			Value: literal.Source,
		}
	case "bytes":
		return &OptionValueBytes{
			proto: literal,
			Value: []byte(literal.Source),
		}
	}

	switch t := r.getTypeByName(field, field.Type).(type) {
	case *Enum:
		val := t.Value(r, literal.Source)
		if val == nil {
			panic(errors.Newf("unknown enum %s value %s", field.Type, literal.Source))
		}
		return &OptionValueEnum{
			proto: literal,
			Value: val,
		}
	case *Message:
		var res []OptionValueMapItem
		for _, lit := range literal.OrderedMap {
			f := t.Field(r, lit.Name)
			if f == nil {
				panic(errors.Newf("unknown message %s field %s", t.Name(), lit.Name))
			}
			rr := buildFromLiteral(r, f.proto.(*proto.NormalField), lit.Literal, false)
			res = append(res, OptionValueMapItem{
				Key:   lit.Name,
				Value: rr,
			})
		}
		return &OptionValueMap{
			proto: literal,
			Value: res,
		}
	default:
		panic(errors.Newf("unknown field type: %s", field.Type))
	}
}

type OptionValueVariant interface {
	Node

	fmt.Stringer
	isOptionValueVariant()
}

type OptionValueBool struct {
	isNode
	proto *proto.Literal

	Value bool
}

type OptionValueInt struct {
	isNode
	proto *proto.Literal

	Value int
}

type OptionValueUint struct {
	isNode
	proto *proto.Literal

	Value uint
}

type OptionValueFloat struct {
	isNode
	proto *proto.Literal

	Value float64
}

type OptionValueString struct {
	isNode
	proto *proto.Literal

	Value string
}

type OptionValueBytes struct {
	isNode
	proto *proto.Literal

	Value []byte
}

type OptionValueEnum struct {
	isNode
	proto *proto.Literal

	Value *EnumValue
}

type OptionValueArray struct {
	isNode
	proto *proto.Literal

	Value []OptionValueVariant
}

type OptionValueMap struct {
	isNode
	proto *proto.Literal

	Value []OptionValueMapItem
	pos   scanner.Position
}

type OptionValueMapItem struct {
	isNode
	proto *proto.Literal

	Key   string
	Value OptionValueVariant
}

func (o *OptionValueBool) String() string {
	return strconv.FormatBool(o.Value)
}

func (o *OptionValueInt) String() string {
	return strconv.Itoa(o.Value)
}

func (o *OptionValueUint) String() string {
	return strconv.FormatUint(uint64(o.Value), 10)
}

func (o *OptionValueFloat) String() string {
	return strconv.FormatFloat(o.Value, 'f', -1, 64)
}

func (o *OptionValueString) String() string {
	return o.Value
}

func (o *OptionValueBytes) String() string {
	return fmt.Sprint(o.Value)
}

func (o *OptionValueEnum) String() string {
	return o.Value.proto.Parent.(*proto.Enum).Name + "." + o.Value.Name()
}

func (o *OptionValueArray) String() string {
	var res strings.Builder
	res.WriteByte('[')
	for i, value := range o.Value {
		if i > 0 {
			res.WriteString(", ")
		}
		res.WriteString(value.String())
	}
	res.WriteByte(']')

	return res.String()
}

func (o *OptionValueMap) String() string {
	var res strings.Builder
	res.WriteByte('{')
	for i, value := range o.Value {
		if i > 0 {
			res.WriteString(", ")
		}

		res.WriteString(value.Key)
		res.WriteString(": ")
		res.WriteString(value.Value.String())
	}
	res.WriteByte('}')

	return res.String()
}

func (o *OptionValueBool) isOptionValueVariant()   {}
func (o *OptionValueInt) isOptionValueVariant()    {}
func (o *OptionValueUint) isOptionValueVariant()   {}
func (o *OptionValueString) isOptionValueVariant() {}
func (o *OptionValueBytes) isOptionValueVariant()  {}
func (o *OptionValueFloat) isOptionValueVariant()  {}
func (o *OptionValueEnum) isOptionValueVariant()   {}
func (o *OptionValueArray) isOptionValueVariant()  {}
func (o *OptionValueMap) isOptionValueVariant()    {}
