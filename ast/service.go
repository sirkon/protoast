package ast

var _ Node = &Service{}

// Service представление для сервисов
type Service struct {
	unique

	File *File

	Name    string
	Methods []*Method
	Options []*Option
}

func (s *Service) node() {}
