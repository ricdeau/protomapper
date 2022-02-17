package mappers

import (
	"github.com/ricdeau/protoast/ast"
	"github.com/ricdeau/protomapper/internal/types"
)

type EnumMapper struct {
	typ *ast.Enum
}

func NewEnumMapper(typ *ast.Enum) *EnumMapper {
	return &EnumMapper{typ: typ}
}

func (e *EnumMapper) FromProto(fieldName string) types.FieldMapperFunc {
	return func(src string) string {
		if src != "" {
			src += "."
		}
		return src + fieldName + ".String()"
	}
}

func (e *EnumMapper) ToProto(fieldName string) types.FieldMapperFunc {
	return func(src string) string {
		if src != "" {
			src += "."
		}
		pbN := GoTypeName(e.typ)
		val := pbN + "_value[" + src + fieldName + "]"
		return call(pbN)(val)
	}
}
