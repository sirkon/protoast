# ProtoAST
A library to represent protobuf services definitions shaped into ASTes.

## Usage

```go
mapping := map[string]string{
	"file.proto": "/var/lib/protofiles/file.proto",
}
files := protoast.NewFiles(mapping)
ns := protoast.NewBuilder(files, func(err error) {
	log.Println(err)
})

// retrieves AST for file.proto
file, err := ns.AST("file.proto")
if err != nil {
	log.Fatal(err)
}

// output AST of type Type defined in file.proto
log.Printf("%#v", file.Services[0])

// returns comment for the first service in a file
log.Printf("%#v", ns.Comment(file.Services[0])) 

// returns position of the second service's name
log.Printf("%#v", ns.Position(file.Service[1], &file.Service[1].Name))
```

You can also use `protoast.NewFilesViaResolver` constructor to use a callback function instead of `map[string]string`
to resolve import names.