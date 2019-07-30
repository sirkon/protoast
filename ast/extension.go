package ast

var _ Type = &Extension{}

// Extension представление extension
type Extension struct {
	unique

	File      *File
	ParentMsg *Message

	Name   string
	Fields []*ExtensionField
	Types  []Type
}

func (*Extension) genericType() {}
func (*Extension) node()        {}

var _ Unique = &ExtensionField{}

// ExtensionField представление поля extension-а
type ExtensionField struct {
	unique

	Name     string
	Sequence int
	Type     Type
	Options  []*Option
}

// MessageToExtension копирует Message в Extension с сохранением всей информации
func MessageToExtension(msg *Message) *Extension {
	ext := &Extension{
		unique:    msg.unique,
		File:      msg.File,
		ParentMsg: msg.ParentMsg,
		Name:      msg.Name,
		Types:     msg.Types,
	}
	for _, f := range msg.Fields {
		ext.Fields = append(ext.Fields, &ExtensionField{
			unique:   f.unique,
			Name:     f.Name,
			Sequence: f.Sequence,
			Type:     f.Type,
			Options:  f.Options,
		})
	}
	return ext
}
