package protoast

import (
	"fmt"
	"text/scanner"
)

func Example_typesVisitor_optionLookup() {
	v := new(typesVisitor)

	v.errors = func(err error) {
		fmt.Print(err)
	}

	v.optionLookup("(deprecated)", scanner.Position{}, messageOptions)
	v.optionLookup("(google.api.message_visibility).restriction", scanner.Position{}, messageOptions)
	// Output:
	//
	//
}
