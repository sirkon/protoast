package ast

import "fmt"

var _ Type = &Map{}

// Map представление типа map<KeyType, ValueType>
type Map struct {
	unique

	KeyType   Hashable
	ValueType Type
}

func (m *Map) String() string {
	return fmt.Sprintf("map[%s]%s", m.KeyType, m.ValueType)
}

func (*Map) genericType() {}
func (*Map) node()        {}
