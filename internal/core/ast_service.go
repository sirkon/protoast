package core

import (
	"iter"
	"text/scanner"

	"github.com/emicklei/proto"
)

type Service struct {
	isNodeOptionable

	proto *proto.Service
}

type Method struct {
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

var (
	_ Node = new(Service)
	_ Node = new(Method)
)

func (s *Service) nodeProto() proto.Visitee { return s.proto }
func (s *Service) pos() scanner.Position    { return s.proto.Position }
func (m *Method) nodeProto() proto.Visitee  { return m.proto }
func (m *Method) pos() scanner.Position     { return m.proto.Position }
