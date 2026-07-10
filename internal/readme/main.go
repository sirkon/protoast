package main

import (
	"fmt"
	"os"

	"github.com/sirkon/protoast/v2"
)

func main() {
	schemaRoot := os.Getenv("SCHEMA_ROOT")

	// WithProtoc() does not execute the protoc binary. Instead, it locates the
	// protoc binary path via `which protoc` to find where the "stdlib" (google/protobuf/*)
	// proto-files are stored, as they are typically strictly coupled with the binary.
	// Use `WithRoot(googleProtobufFiles)` manually if your environment stores
	// these standard files in a non-standard or separate location.
	//
	// Note: Because this relies strictly on `which protoc`, the WithProtoc()
	// method will not work on Windows environments.
	resolvers, err := protoast.Resolvers().WithProtoc().WithRoot(schemaRoot).Build()
	if err != nil {
		panic(fmt.Errorf("create resolvers for schema files: %w", err))
	}

	registry, err := protoast.NewRegistry(resolvers)
	if err != nil {
		panic(fmt.Errorf("create registry: %w", err))
	}

	meta, err := registry.Proto("meta/v1/status.proto")
	if err != nil {
		panic(fmt.Errorf("look for status response definition file: %w", err))
	}

	responseStatusType := meta.Message(registry, "ResponseStatus")
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

		for method := range service.Methods(registry) {
			isStream, responseType := method.Output(registry)
			if isStream {
				continue
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
