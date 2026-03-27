# ProtoAST
A library to represent protobuf services definitions shaped into ASTes.

## Usage

`````go
prjResolver, err := protoast.NewPathResolver(".", "./vendor")
if err != nil {
    returm fmt.Errorf("resolve schema paths")
}

protocResolver, err := protoast.NewProtocResolver()
if err != nil {
    return fmt.Errorf("resolve protoc distribution protos")
}

resolver := protoc.WithResolvers(prjResolver, protocResolver)
ns := protoast.New(resolver, func(err error) {
    log.Println(err)
})

// retrieves AST for file.proto
file, err := ns.AST("file.proto")
if err != nil {
    log.Fatal(err)
}

// output AST of the first service defined in file.proto
log.Printf("%#v", file.Services[0])

// returns comment for the first service in a file
log.Printf("%#v", ns.Comment(file.Services[0])) 

// returns position of the second service's name
log.Printf("%#v", ns.Position(file.Service[1], &file.Service[1].Name))
```

You can also use `protoast.NewFilesViaResolver` constructor to use a callback function instead of `map[string]string`
to resolve import names.
