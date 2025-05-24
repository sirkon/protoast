package protoast

import (
	"path"
	"strings"
	"text/scanner"

	"github.com/sirkon/protoast/internal/errors"

	"github.com/sirkon/protoast/ast"
)

type optionType string

const (
	fileOptions    = "google.protobuf.FileOptions"
	serviceOptions = "google.protobuf.ServiceOptions"
	methodOptions  = "google.protobuf.MethodOptions"
	messageOptions = "google.protobuf.MessageOptions"
	enumOptions    = "google.protobuf.EnumOptions"
	fieldOptions   = "google.protobuf.FieldOptions"

	enumValueOptions = "google.protobuf.EnumValueOptions"
	oneofOptions     = "google.protobuf.OneofOptions"

	// Теоретически есть и другие опции, но они не предоставляются парсером
)

func (tv *typesVisitor) optionLookup(name string, pos scanner.Position, ot optionType) *ast.Extension {
	name = optName(name)

	if d, ok := ignoreOpts[ot]; ok {
		if _, ok := d[name]; ok {
			return nil
		}
	}

	var fileFilter func(*ast.File) bool
	if !strings.ContainsRune(name, '.') {
		// ищем только среди файлов лежащих в данной директории и имеющих одинаковое имя пакета
		base, _ := path.Split(tv.file.Name)
		fileFilter = func(file *ast.File) bool {
			fBase, _ := path.Split(file.Name)
			return fBase == base && file.Package == tv.file.Package
		}
		name = tv.file.Package + "." + name
	} else {
		fileFilter = func(file *ast.File) bool {
			return strings.HasPrefix(name, file.Package)
		}
	}

	files := make([]*ast.File, 1, len(tv.file.Imports)+1)
	files[0] = tv.file
	for _, imp := range tv.file.Imports {
		files = append(files, imp.File)
	}
	for _, file := range files {
		if !fileFilter(file) {
			continue
		}
		optionName := name[len(file.Package)+1:]
		for _, e := range file.Extensions {
			if e.Name != string(ot) {
				continue
			}
			for _, f := range e.Fields {
				if optionName == f.Name {
					return e
				}
			}
		}
	}
	tv.errors(errors.Newf("%s unknown option (%s, belong to %s)", pos, name, ot))
	return nil
}

const (
	messageVisibility = "(google.api.message_visibility).restriction"
	fieldVisibility   = "(google.api.field_visibility).restriction"
)

// options эти игнорируем
var ignoreOpts = map[optionType]map[string]struct{}{
	fileOptions: {
		"optimize_for":         {},
		"go_package":           {},
		"java_package":         {},
		"java_outer_classname": {},
		"csharp_namespace":     {},
		"objc_class_prefix":    {},
		"cc_enable_arenas":     {},
		"java_multiple_files":  {},
	},
	fieldOptions: {
		"default":          {},
		"deprecated":       {},
		"packed":           {},
		"type_name":        {},
		"type_extendee":    {},
		"default_value":    {},
		"oneof_index":      {},
		"json_name":        {},
		"retention":        {},
		"targets":          {},
		"edition_defaults": {},
		"feature_support":  {},
		fieldVisibility:    {},
	},
	oneofOptions: {
		"deprecated": {},
	},
	enumValueOptions: {
		"deprecated": {},
	},
	methodOptions: {
		"deprecated": {},
	},
	messageOptions: {
		"deprecated":      {},
		messageVisibility: {},
	},
}

var visibilityOpts = map[string]struct{}{
	messageVisibility: {},
	fieldVisibility:   {},
}

func optName(name string) string {
	if _, ok := visibilityOpts[name]; ok {
		return name
	}

	name = strings.TrimPrefix(name, "(")
	name = strings.TrimSuffix(name, ")")

	return name
}
