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

// Method поиск метода по имени
func (s *Service) Method(name string) *Method {
	for _, method := range s.Methods {
		if method.Name == name {
			return method
		}
	}

	return nil
}

func (s *Service) node() {}
