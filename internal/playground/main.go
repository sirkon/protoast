package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/sirkon/protoast/v2/internal/core"
)

func main() {
	if err := os.Chdir(filepath.Join("internal", "testdata", "schema")); err != nil {
		panic(err)
	}

	resolvers, err := core.Resolvers().WithProtoc().WithRoot(".", "./vendor").Build()
	if err != nil {
		panic(err)
	}

	registry, err := core.NewRegistry(resolvers...)
	if err != nil {
		panic(err)
	}

	f, err := registry.Proto("service/utopia/v1/service_hash_download.proto")
	if err != nil {
		panic(err)
	}

	fmt.Println(f.Name(), f.Package())
	for option := range registry.Options(f) {
		fmt.Println(option.Name(), option.Value().Value())
	}
}
