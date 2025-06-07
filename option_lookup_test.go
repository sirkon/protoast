package protoast

import (
	"fmt"
	"text/scanner"

	"github.com/sirkon/protoast/ast"
)

func Example_typesVisitor_optionLookup() {
	v := new(typesVisitor)
	v.file = new(ast.File)
	v.errors = func(err error) {
		fmt.Print(err)
	}

	v.optionLookup("(deprecated)", scanner.Position{}, messageOptions)
	v.optionLookup("(google.api.message_visibility).restriction", scanner.Position{}, messageOptions)
	v.optionLookup("(common.v1.log)", scanner.Position{}, fieldOptions)
	// Output:
	//
	//
	// <input> unknown option (common.v1.log, belong to google.protobuf.FieldOptions)
}
