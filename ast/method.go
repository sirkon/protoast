package ast

var _ Node = &Method{}

// Method представление для метода
type Method struct {
	unique

	File    *File
	Service *Service

	Name   string
	Input  Type
	Output Type

	Options []*Option
}

func (m *Method) node() {}

// InputMessage возвращает структуру запроса (минуя оборачивающий Stream, если нужно)
func (m *Method) InputMessage() *Message {
	return getMessage(m.Input)
}

// OutputMessage аналогично InputMessage, возвращает структуру ответа, при необходимости снимая stream
func (m *Method) OutputMessage() *Message {
	return getMessage(m.Output)
}

func getMessage(m Type) *Message {
	v, ok := m.(*Message)
	if ok {
		return v
	}

	return m.(*Stream).Type.(*Message)
}
