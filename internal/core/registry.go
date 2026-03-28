package core

import (
	"bytes"
	"os"

	"github.com/emicklei/proto"

	"github.com/sirkon/protoast/v2/internal/errors"
)

type Registry struct {
	resolvers []PathResolver

	protos   map[string]*proto.Proto
	registry map[string]proto.Visitee
	scopes   map[proto.Visitee]string

	cache   map[proto.Visitee]Node
	ftcache map[*MessageField]Type
}

func NewRegistry(resolvers ...PathResolver) (*Registry, error) {
	res := &Registry{
		resolvers: resolvers,
		protos:    map[string]*proto.Proto{},
		registry:  map[string]proto.Visitee{},
		scopes:    map[proto.Visitee]string{},
		cache:     map[proto.Visitee]Node{},
		ftcache:   map[*MessageField]Type{},
	}
	if err := res.demarkFile("google/protobuf/descriptor.proto"); err != nil {
		return nil, errors.Wrap(err, "set up proto descriptor")
	}

	return res, nil
}

func (r *Registry) Proto(path string) (*File, error) {
	for i := range 2 {
		file, ok := r.protos[path]
		if ok {
			return &File{proto: file}, nil
		}

		if i == 1 {
			break
		}

		if err := r.demarkFile(path); err != nil {
			return nil, errors.Wrap(err, "resolve proto file "+path)
		}
	}

	return nil, errors.New("proto file not found")
}

func (r *Registry) demarkFile(path string) error {
	file, err := r.protoFile(path)
	if err != nil {
		return err
	}

	v := &visitorDemark{
		r: r,
	}
	file.Accept(v)

	return nil
}

func (r *Registry) protoFile(path string) (*proto.Proto, error) {
	if res, ok := r.protos[path]; ok {
		return res, nil
	}

	var protoName string
	for _, resolver := range r.resolvers {
		name, err := resolver.Resolve(path)
		if err != nil {
			if errors.Is(err, os.ErrNotExist) {
				continue
			}

			return nil, errors.Wrap(err, "resolve proto file path with "+resolver.String())
		}

		protoName = name
		break
	}

	if protoName == "" {
		return nil, errors.New("not found")
	}

	parsed, err := readProtoFile(protoName)
	if err != nil {
		return nil, errors.Wrap(err, "get proto definition from resolved file "+protoName)
	}

	r.protos[path] = parsed
	parsed.Filename = path
	return parsed, nil
}

func (r *Registry) optionContextFile() *proto.Message {
	return r.registry[registryOptionsFile].(*proto.Message)
}

func (r *Registry) optionContextMessage() *proto.Message {
	return r.registry[registryOptionsMessage].(*proto.Message)
}

func (r *Registry) optionContextMessageField() *proto.Message {
	return r.registry[registryOptionsMessageField].(*proto.Message)
}

func (r *Registry) optionContextEnum() *proto.Message {
	return r.registry[registryOptionsEnum].(*proto.Message)
}

func (r *Registry) optionContextEnumValue() *proto.Message {
	return r.registry[registryOptionsEnumValue].(*proto.Message)
}

func (r *Registry) optionContextOneof() *proto.Message {
	return r.registry[registryOptionsOneof].(*proto.Message)
}

func (r *Registry) optionContextService() *proto.Message {
	return r.registry[registryOptionsService].(*proto.Message)
}

func (r *Registry) optionContextMethod() *proto.Message {
	return r.registry[registryOptionsMethod].(*proto.Message)
}

func readProtoFile(protoName string) (*proto.Proto, error) {
	file, err := os.ReadFile(protoName)
	if err != nil {
		return nil, errors.Wrap(err, "read file")
	}

	parsed, err := proto.NewParser(bytes.NewReader(file)).Parse()
	if err != nil {
		return nil, errors.Wrap(err, "parse file")
	}

	return parsed, nil
}
