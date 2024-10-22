package ast

import "fmt"

var _ Type = &Repeated{}

// Repeated представление для полей с repeeated
type Repeated struct {
	unique

	Type Type
}

func (r *Repeated) String() string {
	return fmt.Sprintf("[]%s", r.Type)
}

func (*Repeated) genericType() {}
func (*Repeated) node()        {}
