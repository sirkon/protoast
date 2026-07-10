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
| Lazy iterators (Go 1.23) | ✅       | ❌               | ❌     |
| No descriptor generation | ✅       | ❌               | ❌     |
| Cyclic import resolution | ✅       | ⚠️ (complicated) | ✅     |
| Validation               | ❌       | ✅               | ✅     |


## Installation

```shell
go get github.com/sirkon/protoast/v2
```

## Quick Start / Usage

Here is a simple example showing how to load a `.proto` file, resolve its imports using the local environment, and iterate over its elements using Go 1.23 iterators.

```go
package main

import (
    "fmt"
    "os"
    "path/filepath"

    "github.com/sirkon/protoast/v2"
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

### Advanced Example: Enforcing API Contracts (Governance Linter)

In microservice architectures, you often want to enforce global standards — for instance, ensuring that every non-streaming RPC method returns a structure containing a unified `ResponseStatus` message for error handling and metadata.

The following example reads service definitions and validates that:

1. Every unary RPC response of certain services contains a status field.
2. This field is strictly typed as `ResponseStatus`.
3. The field is explicitly named `response_status` (enforcing naming consistency).
4. No other field is allowed to use that `ResponseStatus` type.
5. It respects a custom method option defined in `meta.v1` package:
   ```protobuf
   package meta.v1;
   
   extend google.protobuf.MethodOptions {
       google.protobuf.Any non_standard_method = 12345;
   }
   ```
   If a method is annotated with `(meta.v1.non_standard_method) = {}`, the linter will bypass all the rules above for that specific method.

```go
package main

import (
    "fmt"
    "os"

    "github.com/sirkon/protoast/v2"
)

func main() {
    schemaRoot := os.Getenv("SCHEMA_ROOT")

    // WithProtoc does not execute or invoke the protoc binary. Instead, it tells the
    // resolver to find a protoc binary to locate where the "stdlib" (google/protobuf/*)
    // proto-files are stored, as they are typically strictly coupled with the binary.
    // Use `WithRoot(googleProtobufFiles)` manually if your environment stores
    // these standard files in a non-standard or separate location.
    resolvers, err := protoast.Resolvers().WithProtoc().WithRoot(schemaRoot).Build()
    if err != nil {
        panic(fmt.Errorf("create resolvers for schema files: %w", err))
    }

    registry, err := protoast.NewRegistry(resolvers)
    if err != nil {
        panic(fmt.Errorf("create registry: %w", err))
    }

    metaStatus, err := registry.Proto("meta/v1/status.proto")
    if err != nil {
        panic(fmt.Errorf("look for status response definition file: %w", err))
    }

    responseStatusType := metaStatus.Message(registry, "ResponseStatus")
    if responseStatusType == nil {
        panic(fmt.Errorf("missing definition of ResponseStatus message"))
    }

    servicesPackages := [][2]string{
        {"service/storage/v1/service.proto", "Storage"},
        {"service/meta/v1/meta.proto", "Meta"},
    }
    noErrors := true
    for _, servicesPackage := range servicesPackages {
        serviceRoot, serviceName := servicesPackage[0], servicesPackage[1]

        serviceFile, err := registry.Proto(serviceRoot)
        if err != nil {
            panic(fmt.Errorf("look for service definition file %s: %w", serviceName, err))
        }

        service := serviceFile.Service(registry, serviceName)
        if service == nil {
            panic(fmt.Errorf("service %q not found in %q file", serviceName, serviceFile))
        }

    methodLoop:
        for method := range service.Methods(registry) {
            isStream, responseType := method.Output(registry)
            if isStream {
                continue
            }

            for option := range method.Options(registry) {
                if option.Is(registry, ".meta.v1.non_standard_method") {
                    continue methodLoop
                }
            }

            var responseStatusDetected bool
            for field := range responseType.Fields(registry) {
                if field.Type(registry) == responseStatusType {
                    if field.Name() != "response_status" {
                        fmt.Printf(
                            "%s only response_status field can use %s type\n",
                            registry.Pos(field),
                            registry.TypeName(responseStatusType),
                        )
                        noErrors = false
                        continue
                    }

                    responseStatusDetected = true
                    continue
                }
            }

            if responseStatusDetected {
                continue
            }

            fmt.Printf("%s missing `%s response_status` field in response message type of method %s\n",
                registry.Pos(method),
                registry.TypeName(responseStatusType),
                registry.NodeIndex(method), // Outputs FQN of a method, like .service.storage.v1.Service.MethodName.
            )
            noErrors = false
        }
    }

    if !noErrors {
        fmt.Println("schema does not satisfy criteria")
        os.Exit(1)
    }
}
```

### More

To see how `protoast` behaves in real-world scenarios at scale (including custom nested options validation, array tags parsing, and cross-file type lookups), check out our comprehensive integration test suite:
👉 [protoast_test.go](./protoast_test.go)

## License

Distributed under the **MIT License**. See `LICENSE` for more information.

