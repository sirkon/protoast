package ast

var _ Type = Repeated{}

// Repeated в данном случае мы отошли от неудачной практики аттрибутов "repeated" и ввели специальый тип – Repeated,
// что позволит унифицировать обработку типов
type Repeated struct {
	Type Type
}

func (Repeated) genericType() {}
func (Repeated) node()        {}
