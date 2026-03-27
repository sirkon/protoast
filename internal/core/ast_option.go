package core

import (
	"iter"
	"strings"

	"github.com/emicklei/proto"
)

type Option struct {
	isNode

	registry         *Registry
	protoOptionClass *proto.Message
	protoOptionField *proto.NormalField
	protoOption      *proto.Option
}

func newOption(r *Registry, scope string, class *proto.Message, option *proto.Option) *Option {
	name := parenthesisReplacer.Replace(option.Name)
	name, ok := r.resolveName(scope, name)
	var field *proto.NormalField
	if ok {
		field = r.registry[name].(*proto.NormalField)
	} else {
		// There's only one choice for a valid proto - a builtin option.
		for _, fld := range class.Elements {
			vv, ok := fld.(*proto.NormalField)
			if !ok {
				continue
			}

			if vv.Name != option.Name {
				continue
			}

			field = vv
			break
		}
	}
	if field == nil {
		panic("option field not found")
	}

	return &Option{
		registry:         r,
		protoOptionClass: class,
		protoOptionField: field,
		protoOption:      option,
	}
}

func (o *Option) Name() string {
	// An option can be either predefined or custom. A custom one
	// does not refer an "extend" message but option container like
	// google.protobuf.descriptor.FileOptions.
	if o.protoOptionField.Parent == o.protoOptionClass {
		scope := o.registry.scopes[o.protoOptionClass]
		return "(" + scope + ")" + "." + o.protoOptionField.Name
	} else {
		return o.protoOption.Name
	}
}

func (o *Option) Value() OptionValue {
	return OptionValue{
		option: o,
	}
}

func seqOptions[T proto.Visitee](r *Registry, scope, className string, elements []T) iter.Seq[*Option] {
	class := r.registry[className].(*proto.Message)

	return func(yield func(*Option) bool) {
		for _, element := range elements {
			var vv proto.Visitee = element
			option, ok := vv.(*proto.Option)
			if !ok {
				continue
			}

			if !yield(newOption(r, scope, class, option)) {
				return
			}
		}
	}
}

func namedOption[T proto.Visitee](r *Registry, name string, scope, className string, elements []T) *Option {
	class := r.registry[className].(*proto.Message)

	for _, element := range elements {
		var vv proto.Visitee = element
		option, ok := vv.(*proto.Option)
		if !ok {
			continue
		}

		if option.Name != name {
			continue
		}

		return newOption(r, scope, class, option)
	}

	return nil
}

const (
	registryOptionsFile          = ".google.protobuf.FileOptions"
	registryOptionsMessage       = ".google.protobuf.MessageOptions"
	registryOptionsMessageFields = ".google.protobuf.FieldOptions"
	registryOptionsEnum          = ".google.protobuf.EnumOptions"
	registryOptionsEnumValue     = ".google.protobuf.EnumValueOptions"
	registryOptionsOneof         = ".google.protobuf.OneofOptions"
	registryOptionsService       = ".google.protobuf.ServiceOptions"
	registryOptionsMethod        = ".google.protobuf.MethodOptions"
)

var parenthesisReplacer = strings.NewReplacer("(", "", ")", "")
