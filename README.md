# ProtoTypes
A library to represent protobuf types definitions shaped into ASTes.

## Usage

```go
mapping := map[string]string{
	"file.proto": "/var/lib/protofiles/file.proto",
}
ns := prototypes.NewNamespaces(mapping, func(err error) {
	log.Println(err)
})

// retrieves namespace object to get ASTes of types defined in the file
fileTypes, err := ns.Get("file.proto")
if err != nil {
	log.Fatal(err)
}

// output AST of type Type defined in file.proto
log.Printf("%#v\n", fileTypes.Get("Type"))
```
