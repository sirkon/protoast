# ProtoAST — Lightweight Protobuf AST & Type Resolver for Go

[![Go Reference](https://go.dev)](https://go.dev)
[![License: MIT](https://shields.io)](https://opensource.org)



`protoast` is a fast, dependency-free, and developer-friendly library designed to build a strictly typed **Abstract Syntax Tree (AST)** and perform **full type resolution** for Protocol Buffers (v2 and v3) in pure Go. 

Built on top of the excellent lexical parser [emicklei/proto](https://github.com/emicklei/proto), it provides a powerful, high-level API optimized for writing **custom Protobuf linters, static analyzers, and code generators**.

## Why protoast? (The Problem with protocompile)

Most Go developers building custom tools face a dilemma: use raw text parsers without type resolution, or pull in heavy, complex compilers like `bufbuild/protocompile`. 

`protocompile` meticulously clones the entire C++ `protoc` descriptor logic to maintain 100% bug-for-bug compatibility with legacy features. If you control your schemas and don't need 15 years of legacy workarounds, `protoast` provides a clean, native alternative.

### Key Advantages:
- **No Protobuf Descriptors Needed:** You don't have to deal with `google.protobuf.FileDescriptorProto`. The `protoast` API *is* the descriptor itself, containing even more context (like exact source code positions for IDEs).
- **Go 1.23 Iterators Support:** Leverages native `iter.Seq` (`for field := range msg.Fields(r)`) for zero-allocation, lazy, and clean tree traversal.
- **Lazy Parsing & Cyclic Dependency Resolution:** Safely handles complex dependency graphs and recursive imports (`A.proto` imports `B.proto` and vice-versa) out of the box. Single-pass evaluation on demand.
- **Deep Custom Options Inspection:** Parses complex, nested custom options, arrays, and extension values into structured Go interfaces, matching them with their actual definition types.
- **Zero External C/C++ Dependencies:** Pure Go. Compile it into a single static binary easily.

> **Note on Validation:** `protoast` is **not a validating compiler**. It assumes your `.proto` files are already structurally valid (e.g., compiled successfully via `protoc` or checked via `buf lint`).
> 
| Feature                  | protoast | protocompile     | protoc |
| ------------------------ | -------- | ---------------- | ------ |
| Pure Go, no C++          | ✅       | ❌               | ❌     |
| Lazy iterators (Go 1.23) | ✅       | ❌               | ❌     |
| No descriptor generation | ✅       | ❌               | ❌     |
| Cyclic import resolution | ✅       | ⚠️ (complicated) | ✅     |
| Validation               | ❌       | ✅               | ✅     |


## Installation

```shell
go get github.com/sirkon/protoast
```

## Quick Start / Usage

Here is a simple example showing how to load a `.proto` file, resolve its imports using the local environment, and iterate over its elements using Go 1.23 iterators.

```go
package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/sirkon/protoast"
)

func main() {
	// 1. Configure the path resolver. 
	// WithProtoc() automatically locates your local protoc installation to resolve standard imports like google/protobuf/any.proto
	resolvers, err := protoast.Resolvers().
		WithProtoc().
		WithRoot(".", "./vendor").
		Build()
	if err != nil {
		panic(err)
	}

	// 2. Initialize the type registry
	registry, err := protoast.NewRegistry(resolvers)
	if err != nil {
		panic(err)
	}

	// 3. Parse and resolve a target file
	f, err := registry.Proto("service/utopia/v1/service_hash_download.proto")
	if err != nil {
		panic(err)
	}

	fmt.Printf("Package: %s (%s)\n", f.Name(), f.Package())

	// 4. Iterate over file options using lazy iterators
	for option := range registry.Options(f) {
		fmt.Printf("File Option -> %s: %v\n", option.Name(), option.Value())
	}
    
	// 5. Deep traversal over Messages and Fields
	for msg := range f.Messages(registry) {
		fmt.Println("Message:", msg.Name())
		
		for field := range msg.Fields(registry) {
			fmt.Printf("  Field: %s.%s\n", msg.Name(), field.Name())
		}
	}
}
```

## Advanced Examples

To see how `protoast` behaves in real-world scenarios at scale (including custom nested options validation, array tags parsing, and cross-file type lookups), check out our comprehensive integration test suite:
👉 [protoast_test.go](./protoast_test.go)

## License

Distributed under the **MIT License**. See `LICENSE` for more information.
