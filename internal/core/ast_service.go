package core

import (
	"iter"

	"github.com/emicklei/proto"
)

type Service struct {
	isNode

	proto *proto.Service
}

func (s *Service) Name() string {
	return s.proto.Name
}

func (s *Service) Methods() iter.Seq[*Method] {
	return func(yield func(*Method) bool) {
		for _, e := range s.proto.Elements {
			m, ok := e.(*proto.RPC)
			if !ok {
				continue
			}

			method := &Method{
				proto: m,
			}
			if !yield(method) {
				return
			}
		}
	}
}

func (s *Service) Method(name string) *Method {
	for _, e := range s.proto.Elements {
		m, ok := e.(*proto.RPC)
		if !ok {
			continue
		}

		if m.Name != name {
			continue
		}

		return &Method{
			proto: m,
		}
	}

	return nil
}

type Method struct {
	isNode

	proto *proto.RPC
}

func (m *Method) Name() string {
	return m.proto.Name
}

func (m *Method) Input(r *Registry) (stream bool, typ *Message) {
	v := r.getTypeByName(m.proto, m.proto.RequestType).(*Message)
	return m.proto.StreamsRequest, v
}

func (m *Method) Output(r *Registry) (stream bool, typ *Message) {
	v := r.getTypeByName(m.proto, m.proto.ReturnsType).(*Message)
	return m.proto.StreamsReturns, v
}
