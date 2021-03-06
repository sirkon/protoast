package protoast

import (
	"testing"

	"github.com/sirkon/protoast/internal/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/sirkon/protoast/ast"
)

type copier map[ast.Unique]ast.Unique

func (c copier) copyOptions(opts []*ast.Option) []*ast.Option {
	if opts == nil {
		return nil
	}
	var res []*ast.Option
	for _, opt := range opts {
		res = append(res, c.copyCat(opt).(*ast.Option))
	}
	return res
}

func (c copier) copyType(t ast.Type) ast.Type {
	return c.copyCat(t).(ast.Type)
}

func (c copier) copyOptionValue(vv ast.OptionValue) ast.OptionValue {
	switch v := vv.(type) {
	case *ast.EmbeddedOption:
		return &ast.EmbeddedOption{Value: v.Value}
	case *ast.BoolOption:
		return &ast.BoolOption{Value: v.Value}
	case *ast.StringOption:
		return &ast.StringOption{Value: v.Value}
	case *ast.IntOption:
		return &ast.IntOption{Value: v.Value}
	case *ast.UintOption:
		return &ast.UintOption{Value: v.Value}
	case *ast.FloatOption:
		return &ast.FloatOption{Value: v.Value}
	case *ast.ArrayOption:
		var opts []ast.OptionValue
		for _, opt := range v.Value {
			opts = append(opts, c.copyOptionValue(opt))
		}
		return &ast.ArrayOption{
			Value: opts,
		}
	case *ast.MapOption:
		opts := map[string]ast.OptionValue{}
		for k, vvv := range v.Value {
			opts[k] = c.copyOptionValue(vvv)
		}
		return &ast.MapOption{Value: opts}
	default:
		panic(errors.Newf("unsupported option type %T", vv))
	}
}

func (c copier) copyMsg(m *ast.Message) *ast.Message {
	if m == nil {
		return nil
	}
	return c.copyCat(m).(*ast.Message)
}

func (c copier) copyExt(e *ast.Extension) *ast.Extension {
	if e == nil {
		return nil
	}
	return c.copyCat(e).(*ast.Extension)
}

func (c copier) copyFile(f *ast.File) *ast.File {
	if f == nil {
		return f
	}
	return c.copyCat(f).(*ast.File)
}

func (c copier) copyCat(k ast.Unique) ast.Unique {
	pk, alreadyHere := c[k]
	if alreadyHere {
		return pk
	}
	switch v := k.(type) {
	case *ast.Any:
		return &ast.Any{
			File: c.copyFile(v.File),
		}
	case *ast.Bool:
		return &ast.Bool{}
	case *ast.Bytes:
		return &ast.Bytes{}
	case *ast.Comment:
		return &ast.Comment{
			Value: v.Value,
		}
	case *ast.Enum:
		var res ast.Enum
		c[k] = &res
		var values []*ast.EnumValue
		for _, val := range v.Values {
			values = append(values, c.copyCat(val).(*ast.EnumValue))
		}
		res.File = c.copyFile(v.File)
		res.ParentMsg = c.copyMsg(v.ParentMsg)
		res.Name = v.Name
		res.Values = values
		return c[k]
	case *ast.EnumValue:
		return &ast.EnumValue{
			Name:    v.Name,
			Integer: v.Integer,
			Options: c.copyOptions(v.Options),
		}
	case *ast.File:
		var res ast.File
		c[k] = &res
		var imports []*ast.Import
		for _, imp := range v.Imports {
			imports = append(imports, c.copyCat(imp).(*ast.Import))
		}
		var types []ast.Type
		for _, t := range v.Types {
			types = append(types, c.copyType(t))
		}
		var services []*ast.Service
		for _, s := range v.Services {
			services = append(services, c.copyCat(s).(*ast.Service))
		}
		res.Name = v.Name
		res.Package = v.Package
		res.Imports = imports
		res.Types = types
		res.Services = services
		res.Options = c.copyOptions(v.Options)
		return c[k]
	case *ast.Fixed32:
		return &ast.Fixed32{}
	case *ast.Fixed64:
		return &ast.Fixed64{}
	case *ast.Float32:
		return &ast.Float32{}
	case *ast.Float64:
		return &ast.Float64{}
	case *ast.Import:
		return &ast.Import{
			Path: v.Path,
			File: c.copyFile(v.File),
		}
	case *ast.Int32:
		return &ast.Int32{}
	case *ast.Int64:
		return &ast.Int64{}
	case *ast.Map:
		return &ast.Map{
			KeyType:   c.copyCat(v.KeyType).(ast.Hashable),
			ValueType: c.copyType(v.ValueType),
		}
	case *ast.Message:
		var res ast.Message
		c[k] = &res
		var fields []*ast.MessageField
		for _, f := range v.Fields {
			fields = append(fields, c.copyCat(f).(*ast.MessageField))
		}
		var types []ast.Type
		for _, t := range v.Types {
			types = append(types, c.copyType(t))
		}
		res.File = c.copyFile(v.File)
		res.ParentMsg = c.copyMsg(v.ParentMsg)
		res.Name = v.Name
		res.Fields = fields
		res.Types = types
		return c[k]
	case *ast.Extension:
		var res ast.Extension
		c[k] = &res
		var fields []*ast.ExtensionField
		for _, f := range v.Fields {
			fields = append(fields, c.copyCat(f).(*ast.ExtensionField))
		}
		var types []ast.Type
		for _, t := range v.Types {
			types = append(types, c.copyType(t))
		}
		res.File = c.copyFile(v.File)
		res.ParentMsg = c.copyMsg(v.ParentMsg)
		res.Name = v.Name
		res.Fields = fields
		res.Types = types
		return c[k]
	case *ast.ExtensionField:
		var res ast.ExtensionField
		c[k] = &res
		res.Name = v.Name
		res.Sequence = v.Sequence
		res.Type = c.copyType(v.Type)
		res.Options = c.copyOptions(v.Options)
		return c[k]
	case *ast.MessageField:
		var res ast.MessageField
		c[k] = &res
		res.Name = v.Name
		res.Sequence = v.Sequence
		res.Type = c.copyType(v.Type)
		res.Options = c.copyOptions(v.Options)
		return c[k]
	case *ast.Method:
		var res ast.Method
		c[k] = &res

		var methodOptions []*ast.MethodOption
		for _, mo := range v.Options {
			methodOptions = append(methodOptions, c.copyCat(mo).(*ast.MethodOption))
		}
		res.File = c.copyFile(v.File)
		res.Service = c.copyCat(v.Service).(*ast.Service)
		res.Name = v.Name
		res.Input = c.copyType(v.Input)
		res.Output = c.copyType(v.Output)
		res.Options = methodOptions
		return c[k]
	case *ast.MethodOption:
		var res ast.MethodOption
		c[k] = &res

		var movs []*ast.MethodOptionValue
		for _, mov := range v.Values {
			movs = append(movs, c.copyCat(mov).(*ast.MethodOptionValue))
		}
		res.Name = v.Name
		res.Values = movs
		res.Extension = c.copyExt(v.Extension)
		return c[k]
	case *ast.MethodOptionValue:
		return &ast.MethodOptionValue{
			Name:  v.Name,
			Value: v.Value,
		}
	case *ast.OneOf:
		var res ast.OneOf
		c[k] = &res

		var branches []*ast.OneOfBranch
		for _, b := range v.Branches {
			branches = append(branches, c.copyCat(b).(*ast.OneOfBranch))
		}
		res.ParentMsg = c.copyMsg(v.ParentMsg)
		res.Name = v.Name
		res.Branches = branches
		return c[k]
	case *ast.OneOfBranch:
		var res ast.OneOfBranch
		c[k] = &res
		res.Name = v.Name
		res.Type = c.copyType(v.Type)
		res.Sequence = v.Sequence
		res.Options = c.copyOptions(v.Options)
		return c[k]
	case *ast.Option:
		return &ast.Option{
			Name:      v.Name,
			Extension: c.copyExt(v.Extension),
			Value:     c.copyOptionValue(v.Value),
		}
	case *ast.Optional:
		return &ast.Optional{
			Type: c.copyType(v.Type),
		}
	case *ast.Repeated:
		return &ast.Repeated{
			Type: c.copyType(v.Type),
		}
	case *ast.Service:
		var res ast.Service
		c[k] = &res

		var rpcs []*ast.Method
		for _, r := range v.Methods {
			rpcs = append(rpcs, c.copyCat(r).(*ast.Method))
		}
		res.File = c.copyFile(v.File)
		res.Name = v.Name
		res.Methods = rpcs
		return c[k]
	case *ast.Sfixed32:
		return &ast.Sfixed32{}
	case *ast.Sfixed64:
		return &ast.Sfixed64{}
	case *ast.Sint32:
		return &ast.Sint32{}
	case *ast.Sint64:
		return &ast.Sint64{}
	case *ast.Stream:
		return &ast.Stream{
			Type: c.copyType(v.Type),
		}
	case *ast.String:
		return &ast.String{}
	case *ast.Uint32:
		return &ast.Uint32{}
	case *ast.Uint64:
		return &ast.Uint64{}
	default:
		panic(errors.Newf("unsupported type %T", k))
	}
}

func TestNamespaces_Get(t *testing.T) {
	mapping := map[string]string{
		"errors.proto":                     "testdata/errors.proto",
		"sample.proto":                     "testdata/sample.proto",
		"service.proto":                    "testdata/service.proto",
		"users.proto":                      "testdata/users.proto",
		"google/protobuf/any.proto":        "testdata/google/protobuf/any.proto",
		"google/protobuf/descriptor.proto": "testdata/google/protobuf/descriptor.proto",
		"common/options.proto":             "testdata/common/options.proto",
	}
	c := copier{}

	files := NewFiles(mapping)
	nss := NewBuilder(files, func(err error) {
		t.Errorf("\r%s", err)
	})
	ns, err := nss.Namespace("sample.proto")
	if err != nil {
		t.Fatal(err)
	}
	require.Equal(t, ns.GetType("dos not exist"), nil)

	simpleType := ns.GetType("Simple").(*ast.Message)
	anyFile, err := nss.AST("google/protobuf/any.proto")
	if err != nil {
		t.Error(err)
	}
	sampleSimpleMessage := &ast.Message{
		File: c.copyFile(simpleType.File),
		Name: "Simple",
		Fields: []*ast.MessageField{
			{
				Name:     "anyField",
				Sequence: 1,
				Type: &ast.Any{
					File: c.copyFile(anyFile),
				},
			},
			{
				Name:     "boolField",
				Sequence: 2,
				Type:     &ast.Bool{},
			},
			{
				Name:     "bytesField",
				Sequence: 3,
				Type:     &ast.Bytes{},
			},
			{
				Name:     "fixed32Field",
				Sequence: 4,
				Type:     &ast.Fixed32{},
			},
			{
				Name:     "fixed64Field",
				Sequence: 5,
				Type:     &ast.Fixed64{},
			},
			{
				Name:     "floatField",
				Sequence: 6,
				Type:     &ast.Float32{},
			},
			{
				Name:     "doubleField",
				Sequence: 7,
				Type:     &ast.Float64{},
			},
			{
				Name:     "int32Field",
				Sequence: 8,
				Type:     &ast.Int32{},
			},
			{
				Name:     "int64Field",
				Sequence: 9,
				Type:     &ast.Int64{},
			},
			{
				Name:     "sfixed32Field",
				Sequence: 10,
				Type:     &ast.Sfixed32{},
			},
			{
				Name:     "sfixed64Field",
				Sequence: 11,
				Type:     &ast.Sfixed64{},
			},
			{
				Name:     "sint32Field",
				Sequence: 12,
				Type:     &ast.Sint32{},
			},
			{
				Name:     "sint64Field",
				Sequence: 13,
				Type:     &ast.Sint64{},
			},
			{
				Name:     "uint32Field",
				Sequence: 14,
				Type:     &ast.Uint32{},
			},
			{
				Name:     "uint64Field",
				Sequence: 15,
				Type:     &ast.Uint64{},
			},
		},
	}
	require.Equal(t, sampleSimpleMessage, c.copyCat(simpleType))

	sampleEnum := &ast.Enum{
		File: c.copyFile(simpleType.File),
		Name: "Easy",
		Values: []*ast.EnumValue{
			{
				Name:    "RESERVED",
				Integer: 0,
			},
			{
				Name:    "VALUE",
				Integer: 1,
			},
		},
	}
	require.Equal(t, sampleEnum, c.copyCat(ns.GetType("Easy")))

	errorsFile, err := nss.AST("errors.proto")
	if err != nil {
		t.Fatal(err)
	}

	userType := ns.GetType("User").(*ast.Message)
	sampleUserMessage := &ast.Message{
		File: c.copyFile(userType.File),
		Name: "User",
		Fields: []*ast.MessageField{
			{
				Name:     "id",
				Sequence: 1,
				Type:     &ast.String{},
			},
			{
				Name:     "name",
				Sequence: 2,
				Type:     &ast.String{},
			},
		},
	}
	require.Equal(t, sampleUserMessage, c.copyCat(userType))

	sampleFile := c.copyFile(simpleType.File)
	sampleResponse := &ast.Message{
		File: sampleFile,
		Name: "Response",
		Fields: []*ast.MessageField{
			{
				Name:     "code",
				Sequence: 1,
				Type: &ast.Enum{
					File: c.copyFile(errorsFile),
					Name: "Error",
					Values: []*ast.EnumValue{
						{
							Name:    "RESERVED",
							Integer: 0,
						},
						{
							Name:    "OK",
							Integer: 200,
						},
						{
							Name:    "ERROR",
							Integer: 404,
						},
					},
				},
			},
			{
				Name:     "user",
				Sequence: 2,
				Type: &ast.Message{
					File: c.copyFile(userType.File),
					Name: "User",
					Fields: []*ast.MessageField{
						{
							Name:     "id",
							Sequence: 1,
							Type:     &ast.String{},
						},
						{
							Name:     "name",
							Sequence: 2,
							Type:     &ast.String{},
						},
					},
				},
			},
			{
				Name:     "options",
				Sequence: 3,
				Type: &ast.Map{
					KeyType:   &ast.String{},
					ValueType: &ast.String{},
				},
			},
			{
				Name:     "oo",
				Sequence: -1,
				Type: &ast.OneOf{
					Name: "oo",
					Branches: []*ast.OneOfBranch{
						{
							Name:     "field1",
							Type:     &ast.String{},
							Sequence: 4,
						},
						{
							Name:     "field2",
							Type:     &ast.Int32{},
							Sequence: 5,
						},
					},
				},
			},
		},
	}
	sampleResponse.Fields[3].Type.(*ast.OneOf).ParentMsg = sampleResponse
	responseType := ns.GetType("Response")
	require.Equal(t, sampleResponse, c.copyCat(responseType))

	// тестируется файл с типами
	file := &ast.File{
		Name:    "sample.proto",
		Package: "sample",
		Imports: []*ast.Import{
			{
				Path: "errors.proto",
				File: c.copyFile(errorsFile),
			},
			{
				Path: "users.proto",
				File: c.copyFile(func() *ast.File {
					res, err := nss.AST("users.proto")
					if err != nil {
						t.Fatal(err)
					}
					return c.copyFile(res)
				}()),
			},
			{
				Path: "google/protobuf/any.proto",
				File: c.copyFile(anyFile),
			},
		},
		Types: []ast.Type{sampleSimpleMessage, sampleEnum, sampleResponse},
	}
	fileToTest, err := nss.AST("sample.proto")
	if err != nil {
		t.Fatal(err)
	}
	require.Equal(t, file, c.copyCat(fileToTest))

	ns, err = nss.Namespace("service.proto")
	if err != nil {
		t.Fatal(err)
	}

	optionsFile, err := nss.AST("common/options.proto")
	if err != nil {
		t.Fatal(err)
	}
	ext := optionsFile.Extensions[0]

	srv := ns.GetService("Service")
	require.NotNil(t, srv)
	srvFile := c.copyFile(srv.File)
	srvSample := c.copyCat(srv).(*ast.Service)

	sampleService := &ast.Service{
		File: srvFile,
		Name: "Service",
		Methods: []*ast.Method{
			{
				File:    srvFile,
				Service: srvSample,
				Name:    "Method1",
				Input:   c.copyType(simpleType),
				Output:  c.copyType(responseType),
			},
			{
				File:    srvFile,
				Service: srvSample,
				Name:    "Method2",
				Input:   &ast.Stream{Type: c.copyType(simpleType)},
				Output:  c.copyType(responseType),
				Options: []*ast.MethodOption{
					{
						Name: "(common.option)",
						Values: []*ast.MethodOptionValue{
							{
								Name:  "status",
								Value: "200",
							},
							{
								Name:  "message",
								Value: "OK",
							},
						},
						Extension: c.copyExt(ext),
					},
					{
						Name: "(common.another_option)",
						Values: []*ast.MethodOptionValue{
							{
								Name:  "option",
								Value: "option",
							},
						},
						Extension: c.copyExt(ext),
					},
				},
			},
			{
				File:    srvFile,
				Service: srvSample,
				Name:    "Method3",
				Input:   c.copyType(simpleType),
				Output:  &ast.Stream{Type: c.copyType(responseType)},
			},
			{
				File:    srvFile,
				Service: srvSample,
				Name:    "Method4",
				Input:   &ast.Stream{Type: c.copyType(simpleType)},
				Output:  &ast.Stream{Type: c.copyType(responseType)},
			},
		},
	}
	require.Equal(t, sampleService, srvSample)

	serviceFile := &ast.File{
		Name:    "service.proto",
		Package: "sample",
		Imports: []*ast.Import{
			{
				Path: "sample.proto",
				File: sampleFile,
			},
			{
				Path: "common/options.proto",
				File: c.copyFile(optionsFile),
			},
		},
		Services: []*ast.Service{sampleService},
	}
	serviceAST, err := nss.AST("service.proto")
	if err != nil {
		t.Fatal(err)
	}
	require.Equal(t, serviceFile, c.copyFile(serviceAST))

	require.Equal(t, "testdata/service.proto:8:1", nss.Position(serviceAST.Services[0]).String())
	require.Equal(t, "testdata/service.proto:12:21", nss.PositionField(serviceAST.Services[0].Methods[1].Options[0].Values[0], &serviceAST.Services[0].Methods[1].Options[0].Values[0].Name).String())
}

func TestSubsample(t *testing.T) {
	mapping := map[string]string{
		"subsample.proto": "testdata/subsample.proto",
	}
	c := copier{}

	files := NewFiles(mapping)
	nss := NewBuilder(files, func(err error) {
		t.Errorf("\r%s", err)
	})
	file, err := nss.AST("subsample.proto")
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, "proto3", file.Syntax)

	subMessage := &ast.Message{
		Name: "SubMessage",
		Fields: []*ast.MessageField{
			{
				Name:     "field",
				Sequence: 1,
				Type:     &ast.String{},
			},
		},
	}
	subEnum := &ast.Enum{
		Name: "SubEnum",
		Values: []*ast.EnumValue{
			{
				Name:    "RESERVED",
				Integer: 0,
			},
		},
	}
	sampleMessage := &ast.Message{
		Name: "Message",
		Fields: []*ast.MessageField{
			{
				Name:     "subMsg",
				Sequence: 1,
				Type:     subMessage,
			},
			{
				Name:     "subEnum",
				Sequence: 2,
				Type:     subEnum,
			},
		},
		Types: []ast.Type{
			subMessage,
			subEnum,
		},
	}
	sampleFile := &ast.File{
		Name:    "subsample.proto",
		Package: "sample",
		Types:   []ast.Type{sampleMessage},
	}
	sampleMessage.File = sampleFile
	subMessage.File = sampleFile
	subMessage.ParentMsg = sampleMessage
	subEnum.ParentMsg = sampleMessage
	subEnum.File = sampleFile
	require.Equal(t, sampleFile, c.copyCat(file))
}

func TestSubsample2(t *testing.T) {
	mapping := map[string]string{
		"subsample.proto": "testdata/subsample2.proto",
	}
	c := copier{}

	files := NewFiles(mapping)
	nss := NewBuilder(files, func(err error) {
		t.Errorf("\r%s", err)
	})
	file, err := nss.AST("subsample.proto")
	if err != nil {
		t.Fatal(err)
	}

	subMessage := &ast.Message{
		Name: "SubMessage",
		Fields: []*ast.MessageField{
			{
				Name:     "field",
				Sequence: 1,
				Type:     &ast.String{},
			},
		},
	}
	subEnum := &ast.Enum{
		Name: "SubEnum",
		Values: []*ast.EnumValue{
			{
				Name:    "RESERVED",
				Integer: 0,
			},
		},
	}
	sampleMessage := &ast.Message{
		Name: "Message",
		Fields: []*ast.MessageField{
			{
				Name:     "subMsg",
				Sequence: 1,
				Type:     subMessage,
			},
			{
				Name:     "subEnum",
				Sequence: 2,
				Type:     subEnum,
			},
		},
		Types: []ast.Type{
			subEnum,
		},
	}
	sampleFile := &ast.File{
		Name:    "subsample.proto",
		Package: "sample",
		Types:   []ast.Type{subMessage, sampleMessage},
	}
	sampleMessage.File = sampleFile
	subMessage.File = sampleFile
	subEnum.ParentMsg = sampleMessage
	subEnum.File = sampleFile
	require.Equal(t, sampleFile, c.copyCat(file))
}

func TestMapValueType(t *testing.T) {
	const mapValueTypeName = "map-value-type.proto"
	mapping := map[string]string{
		mapValueTypeName: "testdata/map-value-type.proto",
	}
	c := copier{}
	files := NewFiles(mapping)
	nss := NewBuilder(files, func(err error) {
		t.Errorf("\r%s", err)
	})

	file, err := nss.AST(mapValueTypeName)
	if err != nil {
		t.Fatal(err)
	}

	subType := &ast.Message{
		Name: "Value",
		Fields: []*ast.MessageField{
			{
				Name:     "msg",
				Sequence: 1,
				Type:     &ast.String{},
			},
			{
				Name:     "code",
				Sequence: 2,
				Type:     &ast.Int32{},
			},
		},
	}
	expected := &ast.Message{
		Name: "Struct",
		Fields: []*ast.MessageField{
			{
				Name:     "m",
				Sequence: 1,
				Type: &ast.Map{
					KeyType:   &ast.String{},
					ValueType: subType,
				},
			},
			{
				Name:     "total",
				Sequence: 2,
				Type:     &ast.Int32{},
			},
		},
	}
	expectedFile := &ast.File{
		Name:    mapValueTypeName,
		Package: "map_value_type",
		Types:   []ast.Type{subType, expected},
	}
	subType.File = expectedFile
	expected.File = expectedFile

	require.Equal(t, expectedFile, c.copyCat(file))
}

func TestQuestionable(t *testing.T) {
	mapping := map[string]string{
		"questionable.proto": "testdata/questionable.proto",
		"subsample.proto":    "testdata/subsample.proto",
	}
	c := copier{}
	files := NewFiles(mapping)
	nss := NewBuilder(files, func(err error) {
		t.Errorf("\r%s", err)
	})

	questionable, err := nss.AST("questionable.proto")
	if err != nil {
		t.Fatal(err)
	}
	subsample, err := nss.AST("subsample.proto")
	if err != nil {
		t.Fatal(err)
	}
	thisType := c.copyType(questionable.Types[0]).(*ast.Message)
	thisType.File = nil
	assert.Equal(t, &ast.Message{
		Name: "Questionable",
		Fields: []*ast.MessageField{
			{
				Name:     "f",
				Sequence: 1,
				Type:     c.copyType(subsample.Types[0].(*ast.Message).Types[0]),
			},
		},
	}, thisType)
}

func TestCommentRegression(t *testing.T) {
	mapping := map[string]string{
		"comment-regression.proto": "testdata/comment-regression.proto",
	}
	files := NewFiles(mapping)
	nss := NewBuilder(files, func(err error) {
		t.Errorf("\r%s", err)
	})

	file, err := nss.AST("comment-regression.proto")
	if err != nil {
		t.Fatal(err)
	}

	t.Log(file)
}

func TestHarderOptions(t *testing.T) {
	mapping := map[string]string{
		"comment-regression.proto": "testdata/comment-regression.proto",
	}
	files := NewFiles(mapping)
	nss := NewBuilder(files, func(err error) {
		t.Errorf("\r%s", err)
	})

	file, err := nss.AST("comment-regression.proto")
	if err != nil {
		t.Fatal(err)
	}

	t.Log(file)
}

func TestOpts(t *testing.T) {
	mapping := map[string]string{
		"google/protobuf/any.proto":        "testdata/google/protobuf/any.proto",
		"google/protobuf/descriptor.proto": "testdata/google/protobuf/descriptor.proto",
		"common/options.proto":             "testdata/common/options.proto",
		"opts.proto":                       "testdata/opts.proto",
	}
	files := NewFiles(mapping)
	nss := NewBuilder(files, func(err error) {
		t.Errorf("\r%s", err)
	})

	file, err := nss.AST("opts.proto")
	if err != nil {
		t.Fatal(err)
	}

	t.Log(file)
}

func TestPackageName(t *testing.T) {
	mapping := map[string]string{
		"go_package_pkg.proto":     "testdata/go_package_pkg.proto",
		"go_package_full.proto":    "testdata/go_package_full.proto",
		"go_package_missing.proto": "testdata/go_package_missing.proto",
		"go_package_invalid.proto": "testdata/go_package_invalid.proto",
	}
	files := NewFiles(mapping)
	nss := NewBuilder(files, func(err error) {
		t.Logf("\r%s", err)
	})

	file, err := nss.AST("go_package_pkg.proto")
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, "name", file.GoPkg)

	file, err = nss.AST("go_package_full.proto")
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, "github.com/sirkon/testdata", file.GoPath)
	assert.Equal(t, "name", file.GoPkg)

	_, err = nss.AST("go_package_missing.proto")
	_, err = nss.AST("go_package_invalid.proto")
}
