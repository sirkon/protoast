package core

import (
	"strings"

	"github.com/sirkon/protoast/v2/internal/errors"
)

func (r *Registry) TypeIsDefined(typ Type, ref string) bool {
	switch t := typ.(type) {
	case *Message:
		return t.proto == r.registry[ref]
	case *Enum:
		return t.proto == r.registry[ref]
	}

	return false
}

func (r *Registry) TypeIsGoogleProtobufAny(typ Type) bool {
	return r.TypeIsDefined(typ, ".google.protobuf.Any")
}

func (r *Registry) TypeIsGoogleProtobufEmpty(typ Type) bool {
	return r.TypeIsDefined(typ, ".google.protobuf.Empty")
}

func (r *Registry) TypeIsGoogleProtobufTimestamp(typ Type) bool {
	return r.TypeIsDefined(typ, ".google.protobuf.Timestamp")
}

func (r *Registry) TypeIsGoogleProtobufDuration(typ Type) bool {
	return r.TypeIsDefined(typ, ".google.protobuf.Duration")
}

func (r *Registry) GoPackageOption(node Node) *GoPackageOption {
	v, ok := node.(*Option)
	if !ok {
		return nil
	}

	sample, ok := r.registry[".google.protobuf.FileOptions.go_package"]
	if !ok {
		panic(errors.New("no go_package option detected in registry"))
	}

	if v.optionField != sample {
		return nil
	}

	source := v.proto.Constant.Source
	pos := strings.IndexByte(source, ';')
	if pos < 0 {
		return &GoPackageOption{
			Path: source,
		}
	}

	return &GoPackageOption{
		Path: source[:pos],
		Name: source[pos+1:],
	}
}

func (r *Registry) GoPackage(node Node) *GoPackageOption {
	file := r.NodeFile(node)
	if file == nil {
		return nil
	}

	for a := range file.Everything(r) {
		v, ok := a.(*Option)
		if !ok {
			continue
		}

		if v.Name() != "go_option" {
			continue
		}

		return r.GoPackageOption(v)
	}

	return nil
}

type GoPackageOption struct {
	Path string
	Name string
}

func (o *GoPackageOption) IsValid() bool {
	return o.Path != "" && o.Name != ""
}
