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
log.Printf("%#v\n", file.Services[0])
```
