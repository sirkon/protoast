package ast

var _ Type = &Stream{}

// Stream представление для stream-аргументов и возвращаемых значений метода
type Stream struct {
	unique

	Type Type
}

func (s *Stream) node()        {}
func (s *Stream) genericType() {}
