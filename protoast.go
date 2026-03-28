package protoast

import (
	"github.com/sirkon/protoast/v2/internal/core"
)

// Registry a central repository for lazy traversal over protobuf structures.
type Registry = core.Registry

// PathResolversBuilder is a builder powered by fluent API to construct proto files paths resolutions.
type PathResolversBuilder = core.PathResolversBuilder

// PathResolver is typically a single root proto file path resolution entity.
// You can compose them, but why bother? [Resolvers] builder covers everything.
type PathResolver = core.PathResolver

// Resolvers a builder for file path resolvers.
func Resolvers() *PathResolversBuilder {
	return &core.PathResolversBuilder{}
}

// NewRegistry constructs a new registry with given resolvers.
func NewRegistry(resolvers []PathResolver) (*Registry, error) {
	return core.NewRegistry(resolvers...)
}
