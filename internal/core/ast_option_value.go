package core

import (
	"fmt"
	"strconv"
	"strings"
	"text/scanner"

	"github.com/emicklei/proto"

	"github.com/sirkon/protoast/v2/internal/errors"
)

func buildFromLiteral(r *Registry, option *proto.Option, field *proto.NormalField, literal *proto.Literal, ignoreRepeat bool) OptionValueVariant {
	if field.Repeated && !ignoreRepeat {
		var res []OptionValueVariant
		for _, l := range literal.Array {
			res = append(res, buildFromLiteral(r, nil, field, l, true))
		}
		return &OptionValueArray{
			isOptionValueVariant: isOptionValueVariant{
				option: option,
			},
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
			isOptionValueVariant: isOptionValueVariant{
				option: option,
			},
			proto: literal,
			Value: val,
		}
	case "int32", "sint32", "sfixed32", "int64", "sint64", "sfixed64":
		val, err := strconv.Atoi(literal.Source)
		if err != nil {
			panic(errors.Wrap(err, "convert literal to int"))
		}
		return &OptionValueInt{
			isOptionValueVariant: isOptionValueVariant{
				option: option,
			},
			proto: literal,
			Value: val,
		}
	case "uint32", "fixed32", "uin64", "fixed64":
		val, err := strconv.ParseUint(literal.Source, 10, 64)
		if err != nil {
			panic(errors.Wrap(err, "convert literal to uint"))
		}
		return &OptionValueUint{
			isOptionValueVariant: isOptionValueVariant{
				option: option,
			},
			proto: literal,
			Value: uint(val),
		}
	case "float", "double":
		val, err := strconv.ParseFloat(literal.Source, 64)
		if err != nil {
			panic(errors.Wrap(err, "convert literal to float"))
		}
		return &OptionValueFloat{
			isOptionValueVariant: isOptionValueVariant{
				option: option,
			},
			proto: literal,
			Value: val,
		}
	case "string":
		return &OptionValueString{
			isOptionValueVariant: isOptionValueVariant{
				option: option,
			},
			proto: literal,
			Value: literal.Source,
		}
	case "bytes":
		return &OptionValueBytes{
			isOptionValueVariant: isOptionValueVariant{
				option: option,
			},
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
			isOptionValueVariant: isOptionValueVariant{
				option: option,
			},
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
			rr := buildFromLiteral(r, nil, f.proto.(*proto.NormalField), lit.Literal, false)
			res = append(res, OptionValueMapItem{
				Key:   lit.Name,
				Value: rr,
			})
		}
		return &OptionValueMap{
			isOptionValueVariant: isOptionValueVariant{
				option: option,
			},
			proto: literal,
			Value: res,
		}
	default:
		panic(errors.Newf("unknown field type: %s", field.Type))
	}
}

type OptionValueVariant interface {
	Positionable
	fmt.Stringer
	isOptionValueVariantType()
}

type isOptionValueVariant struct {
	option *proto.Option
}

func (v *isOptionValueVariant) nodeProto() proto.Visitee { return v.option }
func (v *isOptionValueVariant) pos() scanner.Position    { return v.option.Position }
func (*isOptionValueVariant) isOptionValueVariantType()  {}

type OptionValueBool struct {
	isOptionValueVariant
	proto *proto.Literal

	Value bool
}

type OptionValueInt struct {
	isOptionValueVariant
	proto *proto.Literal

	Value int
}

type OptionValueUint struct {
	isOptionValueVariant
	proto *proto.Literal

	Value uint
}

type OptionValueFloat struct {
	isOptionValueVariant
	proto *proto.Literal

	Value float64
}

type OptionValueString struct {
	isOptionValueVariant
	proto *proto.Literal

	Value string
}

type OptionValueBytes struct {
	isOptionValueVariant
	proto *proto.Literal

	Value []byte
}

type OptionValueEnum struct {
	isOptionValueVariant
	proto *proto.Literal

	Value *EnumValue
}

type OptionValueArray struct {
	isOptionValueVariant
	proto *proto.Literal

	Value []OptionValueVariant
}

type OptionValueMap struct {
	isOptionValueVariant
	proto *proto.Literal

	Value []OptionValueMapItem
}

type OptionValueMapItem struct {
	isOptionValueVariant
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

var (
	_ Positionable = new(OptionValueBool)
	_ Positionable = new(OptionValueInt)
	_ Positionable = new(OptionValueUint)
	_ Positionable = new(OptionValueString)
	_ Positionable = new(OptionValueBytes)
	_ Positionable = new(OptionValueFloat)
	_ Positionable = new(OptionValueEnum)
	_ Positionable = new(OptionValueArray)
	_ Positionable = new(OptionValueMap)
)
