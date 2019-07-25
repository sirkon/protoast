package ast

var _ Node = &Service{}

type Service struct {
	File	*File

	Name	string
	Methods	[]*Method
}

func (s *Service) node()	{}
