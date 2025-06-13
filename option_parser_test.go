package protoast

import (
	"fmt"
	"reflect"
	"testing"
)

func Test_normalizeOptionName(t *testing.T) {
	tests := []struct {
		name    string
		opt     string
		want    string
		wantErr bool
	}{
		{
			name: "ok trivial",
			opt:  "ident",
			want: "ident",
		},
		{
			name: "ok dot seq",
			opt:  "common.v1.option",
			want: "common.v1.option",
		},
		{
			name: "ok full bracket",
			opt:  "(common.v1.option)",
			want: "common.v1.option",
		},
		{
			name: "ok partial bracket",
			opt:  "(common.v1.option).log.(parser.human)",
			want: "common.v1.option.log.parser.human",
		},
		{
			name:    "invalid empty",
			opt:     "",
			want:    "",
			wantErr: true,
		},
		{
			name:    "invalid free dot",
			opt:     "option.v1.",
			want:    "",
			wantErr: true,
		},
		{
			name:    "unexpected close bracket",
			opt:     "common.v1.option)",
			want:    "",
			wantErr: true,
		},
		{
			name:    "unpaired open bracket",
			opt:     "((common.v1.option)",
			want:    "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := normalizeOptionName(tt.opt)
			if err != nil && tt.wantErr {
				t.Log("expected error:", err)
			}
			if (err != nil) != tt.wantErr {
				t.Errorf("normalizeOptionName() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("normalizeOptionName() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getLexems(t *testing.T) {
	dot := lexemType{typ: lexemTypeCodeDot}
	open := lexemType{typ: lexemTypeCodeOpen}
	cls := lexemType{typ: lexemTypeCodeClose}
	ident := func(s string) lexemType {
		return lexemType{
			value: s,
			typ:   lexemTypeCodeIdent,
		}
	}

	tests := []struct {
		name    string
		s       string
		wantRes []lexemType
		wantErr bool
	}{
		{
			name:    "trivial",
			s:       "legacy",
			wantRes: []lexemType{ident("legacy")},
		},
		{
			name:    "dotted",
			s:       "common.v1.log",
			wantRes: []lexemType{ident("common"), dot, ident("v1"), dot, ident("log")},
		},
		{
			name:    "full commas",
			s:       "(common.v1.log)",
			wantRes: []lexemType{open, ident("common"), dot, ident("v1"), dot, ident("log"), cls},
		},
		{
			name:    "part commas",
			s:       "(common.v1).log",
			wantRes: []lexemType{open, ident("common"), dot, ident("v1"), cls, dot, ident("log")},
		},
		{
			name:    "invalid",
			s:       "&&&",
			wantRes: nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotRes, err := getLexems(tt.s)
			if (err != nil) != tt.wantErr {
				t.Errorf("getLexems() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("getLexems() gotRes = %v, want %v", gotRes, tt.wantRes)
			}
		})
	}
}

func exampleGetIdentifierStr(s string) {
	id, rest := getIdentifier(s)
	fmt.Println(id + "-" + rest + "-!")
}

func Example_getIdentifier() {
	exampleGetIdentifierStr("123")
	exampleGetIdentifierStr("_123")
	exampleGetIdentifierStr("abc")
	exampleGetIdentifierStr("abc.")

	// output:
	// -123-!
	// _123--!
	// abc--!
	// abc-.-!
}
