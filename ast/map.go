package ast

var _ Type = Map{}

// Map представление типа map<K, V>
type Map struct {
	KeyType   Hashable
	ValueType Type
}

func (Map) genericType() {}
func (Map) node()        {}
