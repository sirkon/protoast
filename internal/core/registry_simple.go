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

	if v.protoOptionField != sample {
		return nil
	}

	source := v.protoOption.Constant.Source
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

type GoPackageOption struct {
	Path string
	Name string
}

func (o *GoPackageOption) IsValid() bool {
	return o.Path != "" && o.Name != ""
}
