# ProtoAST
Lazy parser for protobuf built over excellent [emicklei/proto](https://github.com/emicklei/proto).

## Usage

```go
package main

import (
    "fmt"
    "os"
    "path/filepath"

    "github.com/sirkon/protoast/v2"
)

func main() {
    home, err := os.UserHomeDir()
    if err != nil {
        panic(err)
    }

    if err := os.Chdir(home); err != nil {
        panic(err)
    }

    if err := os.Chdir(filepath.Join("Sources", "work", "utopia", "internal", "schema")); err != nil {
        panic(err)
    }

    resolvers, err := protoast.Resolvers().WithProtoc().WithRoot(".", "./vendor").Build()
    if err != nil {
        panic(err)
    }

    registry, err := protoast.NewRegistry(resolvers)
    if err != nil {
        panic(err)
    }

    f, err := registry.Proto("service/utopia/v1/service_hash_download.proto")
    if err != nil {
        panic(err)
    }

    fmt.Println("package data:", f.Name(), f.Package())
    for option := range registry.Options(f) {
        fmt.Println("file option:", option.Name(), option.Value())
    }
    
    for msg := range f.Messages(registry) {
        fmt.Println("message:", msg.Name())
        
        for field := range msg.Fields(registry) {
            fmt.Println("field:", msg.Name() + "." + field.Name())
        }
    }
}
```


