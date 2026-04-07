package core

import (
	"iter"
	"strings"
	"text/scanner"

	"github.com/emicklei/proto"
)

type Option struct {
	proto       *proto.Option
	registry    *Registry
	optionClass *proto.Message
	optionField *proto.NormalField
}

func newOption(r *Registry, scope string, class *proto.Message, option *proto.Option) *Option {
	name := parenthesisReplacer.Replace(option.Name)
	name, ok := r.resolveNameRaw(scope, name)
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

	res := &Option{
		registry:    r,
		optionClass: class,
		optionField: field,
		proto:       option,
	}
	return res
}

func (o *Option) Name() string {
	// An option can be either predefined or custom. A custom one
	// does not refer an "extend" message but option container like
	// google.protobuf.descriptor.FileOptions.
	if o.optionField.Parent == o.optionClass {
		scope := o.registry.scopes[o.optionClass]
		return "(" + scope + ")" + "." + o.optionField.Name
	} else {
		return o.proto.Name
	}
}

func (o *Option) Value() OptionValueVariant {
	return buildFromLiteral(o.registry, o.proto, o.optionField, &o.proto.Constant, false)
}

// Is checks if given option has this qualified name. Meaning .x.y.z, not x.y.z.
func (o *Option) Is(r *Registry, name string) bool {
	return r.NodeIndex(r.wrap(o.optionField)) == name
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
	registryOptionsFile         = ".google.protobuf.FileOptions"
	registryOptionsMessage      = ".google.protobuf.MessageOptions"
	registryOptionsMessageField = ".google.protobuf.FieldOptions"
	registryOptionsEnum         = ".google.protobuf.EnumOptions"
	registryOptionsEnumValue    = ".google.protobuf.EnumValueOptions"
	registryOptionsOneof        = ".google.protobuf.OneofOptions"
	registryOptionsService      = ".google.protobuf.ServiceOptions"
	registryOptionsMethod       = ".google.protobuf.MethodOptions"
)

var parenthesisReplacer = strings.NewReplacer("(", "", ")", "")

var _ Node = new(Option)

func (o *Option) nodeProto() proto.Visitee { return o.proto }
func (o *Option) pos() scanner.Position    { return o.proto.Position }
