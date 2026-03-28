package core

import (
	"iter"

	"github.com/emicklei/proto"
)

type Service struct {
	isNode
	isNodeOptionable

	proto *proto.Service
}

type Method struct {
	isNode
	isNodeOptionable

	proto *proto.RPC
}

// Name returns service name.
func (s *Service) Name() string {
	return s.proto.Name
}

// Methods returns all service methods.
func (s *Service) Methods(r *Registry) iter.Seq[*Method] {
	return func(yield func(*Method) bool) {
		for _, e := range s.proto.Elements {
			m, ok := e.(*proto.RPC)
			if !ok {
				continue
			}

			if !yield(r.wrap(m).(*Method)) {
				return
			}
		}
	}
}

// Method returns a method with the given name.
func (s *Service) Method(r *Registry, name string) *Method {
	for _, e := range s.proto.Elements {
		m, ok := e.(*proto.RPC)
		if !ok {
			continue
		}

		if m.Name != name {
			continue
		}

		return r.wrap(m).(*Method)
	}

	return nil
}

// Everything returns everything defined in the service.
func (s *Service) Everything(r *Registry) iter.Seq[Node] {
	return func(yield func(Node) bool) {
		for _, e := range s.proto.Elements {
			if v, ok := e.(*proto.Option); ok {
				if !yield(r.wrapOption(v, r.optionContextService())) {
					return
				}
				continue
			}
			if !yield(r.wrap(e)) {
				return
			}
		}
	}
}

// Name returns method name.
func (m *Method) Name() string {
	return m.proto.Name
}

// Input returns method input type.
func (m *Method) Input(r *Registry) (stream bool, typ *Message) {
	v := r.getTypeByName(m.proto, m.proto.RequestType).(*Message)
	return m.proto.StreamsRequest, v
}

// Output returns method output type.
func (m *Method) Output(r *Registry) (stream bool, typ *Message) {
	v := r.getTypeByName(m.proto, m.proto.ReturnsType).(*Message)
	return m.proto.StreamsReturns, v
}

func (m *Method) Options(r *Registry) iter.Seq[*Option] {
	return func(yield func(*Option) bool) {
		for _, e := range m.proto.Elements {
			v, ok := e.(*proto.Option)
			if !ok {
				continue
			}

			if !yield(r.wrapOption(v, r.optionContextMethod()).(*Option)) {
				return
			}
		}
	}
}
