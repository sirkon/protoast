package ast

var _ Node = &Service{}

// Service представление для сервисов
type Service struct {
	unique

	File *File

	Name    string
	Methods []*Method
}

func (s *Service) node() {}
