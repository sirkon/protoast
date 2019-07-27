package ast

var _ Type = &Map{}

// Map представление типа map<KeyType, ValueType>
type Map struct {
	unique

	KeyType   Hashable
	ValueType Type
}

func (*Map) genericType() {}
func (*Map) node()        {}
