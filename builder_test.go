package protoast

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/sirkon/protoast/ast"
)

func TestNamespaces_Get(t *testing.T) {
	mapping := map[string]string{
		"errors.proto":             "testdata/errors.proto",
		"sample.proto":             "testdata/sample.proto",
		"service.proto":            "testdata/service.proto",
		"users.proto":              "testdata/users.proto",
		"google/protobuf/any.poto": "testdata/google/protobuf/any.proto",
	}

	files := NewFiles(mapping)
	nss := NewBuilder(files, func(err error) {
		t.Errorf("\r%s", err)
	})
	ns, err := nss.Namespace("sample.proto")
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, ns.GetType("dos not exist"), nil)
	simpleType := ns.GetType("Simple")
	require.Equal(t, &ast.Message{
		File: &ast.File{
			Name:    "sample.proto",
			Package: "sample",
		},
		Name: "Simple",
		Fields: []ast.MessageField{
			{
				Name:     "anyField",
				Sequence: 1,
				Type:     ast.Any{},
			},
			{
				Name:     "boolField",
				Sequence: 2,
				Type:     ast.Bool{},
			},
			{
				Name:     "bytesField",
				Sequence: 3,
				Type:     ast.Bytes{},
			},
			{
				Name:     "fixed32Field",
				Sequence: 4,
				Type:     ast.Fixed32{},
			},
			{
				Name:     "fixed64Field",
				Sequence: 5,
				Type:     ast.Fixed64{},
			},
			{
				Name:     "floatField",
				Sequence: 6,
				Type:     ast.Float32{},
			},
			{
				Name:     "doubleField",
				Sequence: 7,
				Type:     ast.Float64{},
			},
			{
				Name:     "int32Field",
				Sequence: 8,
				Type:     ast.Int32{},
			},
			{
				Name:     "int64Field",
				Sequence: 9,
				Type:     ast.Int64{},
			},
			{
				Name:     "sfixed32Field",
				Sequence: 10,
				Type:     ast.Sfixed32{},
			},
			{
				Name:     "sfixed64Field",
				Sequence: 11,
				Type:     ast.Sfixed64{},
			},
			{
				Name:     "sint32Field",
				Sequence: 12,
				Type:     ast.Sint32{},
			},
			{
				Name:     "sint64Field",
				Sequence: 13,
				Type:     ast.Sint64{},
			},
			{
				Name:     "uint32Field",
				Sequence: 14,
				Type:     ast.Uint32{},
			},
			{
				Name:     "uint64Field",
				Sequence: 15,
				Type:     ast.Uint64{},
			},
		},
	}, simpleType)
	require.Equal(t, ns.GetType("Easy"), &ast.Enum{
		File: &ast.File{
			Name:    "sample.proto",
			Package: "sample",
		},
		Name: "Easy",
		Values: []ast.EnumValue{
			{
				Name:    "RESERVED",
				Integer: 0,
			},
			{
				Name:    "VALUE",
				Integer: 1,
			},
		},
	})
	userType := ns.GetType("User")
	require.Equal(t, userType, &ast.Message{
		File: userType.(*ast.Message).File,
		Name: "User",
		Fields: []ast.MessageField{
			{
				Name:     "id",
				Sequence: 1,
				Type:     ast.String{},
			},
			{
				Name:     "name",
				Sequence: 2,
				Type:     ast.String{},
			},
		},
	})

	sample := &ast.Message{
		File: &ast.File{
			Name:    "sample.proto",
			Package: "sample",
		},
		Name: "Response",
		Fields: []ast.MessageField{
			{
				Name:     "code",
				Sequence: 1,
				Type: &ast.Enum{
					File: &ast.File{
						Name:    "errors.proto",
						Package: "sample",
					},
					Name: "Error",
					Values: []ast.EnumValue{
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
					File: &ast.File{
						Name:    "users.proto",
						Package: "sample",
					},
					Name: "User",
					Fields: []ast.MessageField{
						{
							Name:     "id",
							Sequence: 1,
							Type:     ast.String{},
						},
						{
							Name:     "name",
							Sequence: 2,
							Type:     ast.String{},
						},
					},
				},
			},
			{
				Name:     "options",
				Sequence: 3,
				Type: ast.Map{
					KeyType:   ast.String{},
					ValueType: ast.String{},
				},
			},
			{
				Name:     "oo",
				Sequence: -1,
				Type: &ast.OneOf{
					Name: "oo",
					Branches: []ast.OneOfBranch{
						{
							Name:     "field1",
							Type:     ast.String{},
							Sequence: 4,
						},
						{
							Name:     "field2",
							Type:     ast.Int32{},
							Sequence: 5,
						},
					},
				},
			},
		},
	}
	sample.Fields[3].Type.(*ast.OneOf).ParentMsg = sample
	responseType := ns.GetType("Response")
	require.Equal(t, sample, responseType)

	ns, err = nss.Namespace("service.proto")
	if err != nil {
		t.Fatal(err)
	}

	srv := ns.GetService("Service")
	require.NotNil(t, srv)

	require.Equal(t, &ast.Service{
		File: srv.File,
		Name: "Service",
		Methods: []*ast.Method{
			{
				File:    srv.File,
				Service: srv,
				Name:    "Method1",
				Input:   simpleType,
				Output:  responseType,
			},
			{
				File:    srv.File,
				Service: srv,
				Name:    "Method2",
				Input:   ast.Stream{Type: simpleType},
				Output:  responseType,
				Options: []ast.MethodOption{
					{
						Name: "(common.option)",
						Values: []ast.OptionValue{
							{
								Name:  "status",
								Value: "200",
							},
							{
								Name:  "message",
								Value: "OK",
							},
						},
					},
					{
						Name: "(common.another_option)",
						Values: []ast.OptionValue{
							{
								Name:  "option",
								Value: "option",
							},
						},
					},
				},
			},
			{
				File:    srv.File,
				Service: srv,
				Name:    "Method3",
				Input:   simpleType,
				Output:  ast.Stream{Type: responseType},
			},
			{
				File:    srv.File,
				Service: srv,
				Name:    "Method4",
				Input:   ast.Stream{Type: simpleType},
				Output:  ast.Stream{Type: responseType},
			}},
	}, srv)
}
