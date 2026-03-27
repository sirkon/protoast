package core

import (
	"unsafe"

	"github.com/emicklei/proto"
)

type messageFieldTypeVariant interface {
	isMessageFieldTypeVariant()
}

type isEmickleiNormalField struct {
	proto.NormalField
}

type isEmickleiOneOf struct {
	proto.Oneof
}

type isEmickleiMapField struct {
	proto.MapField
}

func (m *isEmickleiNormalField) isMessageFieldTypeVariant() {}
func (m *isEmickleiOneOf) isMessageFieldTypeVariant()       {}
func (m *isEmickleiMapField) isMessageFieldTypeVariant()    {}

func (m *isEmickleiNormalField) asProto() *proto.NormalField {
	return (*proto.NormalField)(unsafe.Pointer(m))
}

func (m *isEmickleiOneOf) asProto() *proto.Oneof {
	return (*proto.Oneof)(unsafe.Pointer(m))
}

func (m *isEmickleiMapField) asProto() *proto.MapField {
	return (*proto.MapField)(unsafe.Pointer(m))
}

func asEmickleiNormalField(p *proto.NormalField) *isEmickleiNormalField {
	return (*isEmickleiNormalField)(unsafe.Pointer(p))
}

func asEmickleiOneOf(p *proto.Oneof) *isEmickleiOneOf {
	return (*isEmickleiOneOf)(unsafe.Pointer(p))
}

func asEmickleiMapField(p *proto.MapField) *isEmickleiMapField {
	return (*isEmickleiMapField)(unsafe.Pointer(p))
}
