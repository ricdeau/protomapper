package mappers

import (
	"github.com/ricdeau/protoast/ast"
	"github.com/ricdeau/protomapper/internal/registry"
	"github.com/ricdeau/protomapper/internal/types"
)

type MessageMapper struct {
	msg *ast.Message
}

func NewMessageMapper(msg *ast.Message) *MessageMapper {
	return &MessageMapper{msg: msg}
}

func (m *MessageMapper) FromPb(fieldName string) types.FieldMapperFunc {
	return func(src string) string {
		if src != "" {
			src += "."
		}
		typ := registry.Types.GetType(m.msg)
		methName := typ.GetName() + "FromPb"
		return call(methName)(src + fieldName)
	}
}

func (m *MessageMapper) ToPb(fieldName string) types.FieldMapperFunc {
	return func(src string) string {
		if src != "" {
			src += "."
		}
		typ := registry.Types.GetType(m.msg)
		methName := typ.GetName() + "ToPb"
		return call(methName)(src + fieldName)
	}
}
