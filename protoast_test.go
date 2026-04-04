package protoast_test

import (
	"iter"
	"testing"

	"github.com/alecthomas/assert/v2"

	"github.com/sirkon/protoast/v2"
	"github.com/sirkon/protoast/v2/internal/errors"
	"github.com/sirkon/protoast/v2/past"
)

func TestProtoAST(t *testing.T) {
	resolvers, err := protoast.Resolvers().WithProtoc().WithRoot("./testdata").Build()
	if err != nil {
		t.Fatal(errors.Wrap(err, "build resolvers"))
	}

	r, err := protoast.NewRegistry(resolvers)
	if err != nil {
		t.Fatal(errors.Wrap(err, "create registry"))
	}

	data, err := r.Proto("data.proto")
	if err != nil {
		t.Fatal(errors.Wrap(err, "get data.proto"))
	}

	// Imports
	for i, imp := range enumerate(data.Imports(r)) {
		switch i {
		case 0:
			assert.Equal(t, "meta.proto", imp.Path())
		case 1:
			assert.Equal(t, "google/protobuf/any.proto", imp.Path())
		default:
			t.Fatalf("unexpected import %q", imp.Path())
		}
	}

	// File options
	for i, opt := range enumerate(r.Options(data)) {
		switch i {
		case 0:
			assert.Equal(t, "(.google.protobuf.FileOptions).go_package", opt.Name())
			v := opt.Value().(*past.OptionValueString)
			assert.Equal(t, "gopkg/pb;pb", v.String())
		case 1:
			assert.Equal(t, "(pb.file_tag)", opt.Name())
			v := opt.Value().(*past.OptionValueString)
			assert.Equal(t, "file.tag", v.String())

			pos := r.Pos(opt)
			assert.HasSuffix(t, pos.Filename, "data.proto")
			assert.Equal(t, pos.Line, 10)
		default:
			t.Fatalf("unexpected file option %s", opt.Name())
		}
	}

	// Messages
	assert.Equal(t, 2, iterLen(data.Messages(r)))

	// Must be EnumValuePayload with the single field.
	localMsg := data.Message(r, "EnumValuePayload")
	assert.NotEqual(t, nil, localMsg)
	assert.Equal(t, 1, iterLen(localMsg.Everything(r)), "no of anything in the message")
	assert.Equal(t, 1, iterLen(localMsg.Fields(r)), "no of fields in the message")

	// Must be Message
	msg := data.Message(r, "Message")
	assert.NotEqual(t, nil, msg)
	checkMessage(t, r, msg)

	// Check Enum.
	enum := data.Enum(r, "Enum")
	assert.NotEqual(t, nil, enum)
	var count int
	for value := range enum.Values(r) {
		count++
		switch value.Value() {
		case 0:
			assert.Equal(t, "ENUM_UNSPECIFIED", value.Name())
		case 1:
			assert.Equal(t, "ENUM_VALUE_1", value.Name())
		case 2:
			assert.Equal(t, "ENUM_VALUE_2", value.Name())
		default:
			t.Fatalf("unexpected enum value %q = %d", value.Name(), value.Value())
		}
	}
	if count != 3 {
		t.Errorf("expected 3 values in Enum, got %d", count)
	}

	// Check service.
	s := data.Service(r, "Service")
	assert.NotEqual(t, nil, s)
	serviceData := map[string][]string{}
	for method := range s.Methods(r) {
		var payload []string
		stream, msg := method.Input(r)
		payload = append(payload, streamType(stream), msg.Name())
		stream, msg = method.Output(r)
		payload = append(payload, streamType(stream), msg.Name())
		serviceData[method.Name()] = payload
	}
	wantServiceData := map[string][]string{
		"MethodUU": {"unary", "Message", "unary", "InnerMessage"},
		"MethodUS": {"unary", "EnumValuePayload", "stream", "Message"},
		"MethodSU": {"stream", "InnerMessage", "unary", "EnumValuePayload"},
		"MethodSS": {"stream", "InnerMessage", "stream", "InnerMessage"},
	}
	assert.Equal(t, wantServiceData, serviceData, "check service signature")
}

func streamType(isStream bool) string {
	if isStream {
		return "stream"
	}

	return "unary"
}

func checkMessage(t *testing.T, r *protoast.Registry, msg *past.Message) {
	names := map[string]bool{}
	for f := range msg.Fields(r) {
		names[f.Name()] = true
		switch f.Name() {
		case "i32":
			assertField[*past.Int32](t, r, f, 1)
		case "si32":
			assertField[*past.Sint32](t, r, f, 2)
		case "sf32":
			assertField[*past.Sfixed32](t, r, f, 3)
		case "i64":
			assertField[*past.Int64](t, r, f, 11)
		case "si64":
			assertField[*past.Sint64](t, r, f, 12)
		case "sf64":
			assertField[*past.Sfixed64](t, r, f, 13)
		case "u32":
			assertField[*past.Uint32](t, r, f, 21)
		case "uf32":
			assertField[*past.Fixed32](t, r, f, 22)
		case "u64":
			assertField[*past.Uint64](t, r, f, 31)
		case "uf64":
			assertField[*past.Fixed64](t, r, f, 32)
		case "f32":
			assertField[*past.Float](t, r, f, 41)
		case "f64":
			assertField[*past.Double](t, r, f, 42)
		case "b":
			assertField[*past.Bool](t, r, f, 51)
		case "str":
			assertField[*past.String](t, r, f, 61)
		case "raw":
			assertField[*past.Bytes](t, r, f, 62)
		case "bools":
			assertField[*past.Repeated](t, r, f, 63)
			rpt := f.Type(r).(*past.Repeated)
			if _, ok := rpt.Type.(*past.Bool); !ok {
				t.Errorf("field bools expected to have repeated %T type, got repeated %T", new(past.Bool), rpt.Type)
			}
		case "any":
			assertFieldExact(t, r, f, ".google.protobuf.Any", 101)
		case "local_msg":
			assertFieldExact(t, r, f, ".pb.EnumValuePayload", 102)
		case "inner_msg":
			assertFieldExact(t, r, f, ".pb.Message.InnerMessage", 103)
		case "local_enum":
			assertFieldExact(t, r, f, ".pb.Enum", 104)
		case "inner_enum":
			assertFieldExact(t, r, f, ".pb.Message.InnerEnum", 105)
		case "map_field":
			assertField[*past.Map](t, r, f, 111)
			m := f.Type(r).(*past.Map)
			assertType[*past.String](t, m.Key(), "check key type of map_field")
			assertType[*past.Bool](t, m.Value(r), "check value type of map_field")
		case "payload":
			oo, ok := f.Type(r).(*past.OneOf)
			if !ok {
				t.Errorf("field payload was expected to be %T, got %T", new(past.OneOf), f.Type(r))
			}
			for br := range oo.Branches(r) {
				names[f.Name()+"."+br.Name()] = true
				switch br.Name() {
				case "branch1":
					assertBranch[*past.String](t, r, br, 112)
				case "branch2":
					assertBranch[*past.Bool](t, r, br, 113)
				default:
					t.Errorf("unexpected branch %s %d", br.Name(), br.Value())
				}
			}
		default:
			t.Errorf("unexpected field %q %T = %d", f.Name(), f.Type(r), f.Value())
		}
	}

	// Check if all fields were iterated (including branches of oneof).
	requiredFields := []string{
		"i32", "si32", "sf32", "i64", "si64", "sf64", "u32", "uf32", "u64", "uf64", "f32", "f64",
		"b", "str", "raw", "bools", "any", "local_msg", "inner_msg", "local_enum", "inner_enum",
		"map_field", "payload", "payload.branch1", "payload.branch2",
	}
	for _, name := range requiredFields {
		if !names[name] {
			t.Errorf("missing required field %q", name)
		}
	}
}

func assertField[T past.Type](t *testing.T, r *protoast.Registry, field *past.MessageField, value int) {
	if v, ok := field.Type(r).(T); !ok {
		var zero T
		t.Errorf("field %s expected to have %T type, got %T", field.Name(), zero, v)
	}
	if field.Value() != value {
		t.Errorf("field %s expected to have %d sequence number, got %d", field.Name(), value, field.Value())
	}
}

func assertBranch[T past.Type](t *testing.T, r *protoast.Registry, field *past.OneOfBranch, value int) {
	if v, ok := field.Type(r).(T); !ok {
		var zero T
		t.Errorf("field %s expected to have %T type, got %T", field.Name(), zero, v)
	}
	if field.Value() != value {
		t.Errorf("field %s expected to have %d sequence number, got %d", field.Name(), value, field.Value())
	}
}

func assertFieldExact(t *testing.T, r *protoast.Registry, field *past.MessageField, typName string, value int) {
	typ := r.NodeByFullName(typName)
	if field.Type(r) != typ {
		t.Errorf("field %s expected to have type %v, got %v", field.Name(), r.NodeIndex(typ), r.NodeIndex(field.Type(r)))
	}

	if field.Value() != value {
		t.Errorf("field %s expected to have %d sequence number, got %d", field.Name(), value, field.Value())
	}
}

func assertType[T past.Type](t *testing.T, typ past.Type, what string) {
	if _, ok := typ.(T); !ok {
		var zero T
		t.Errorf("%s: expected type %T, got %T", what, zero, typ)
	}
}

func enumerate[T any](it iter.Seq[T]) iter.Seq2[int, T] {
	return func(yield func(int, T) bool) {
		var i int
		for t := range it {
			if !yield(i, t) {
				return
			}

			i++
		}
	}
}

func iterLen[T any](it iter.Seq[T]) int {
	var count int
	for range it {
		count++
	}

	return count
}
