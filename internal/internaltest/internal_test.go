package internaltest

import (
	"fmt"
	"testing"

	"github.com/sirkon/protoast"
)

func TestGeneric(t *testing.T) {
	files := protoast.NewFilesViaResolver(func(f string) (string, error) {
		return fmt.Sprintf("../../testdata/%s", f), nil
	})

	builder := protoast.NewBuilder(files, func(err error) {
		t.Errorf("\r" + err.Error())
	})

	f, err := builder.AST("opts.proto")
	if err != nil {
		t.Fatal(err)
	}

	t.Log(f.Package)
}
