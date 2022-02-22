package mappers

import (
	"fmt"

	"github.com/ricdeau/protoast/ast"
	"github.com/ricdeau/protomapper/internal/helpers"
	"github.com/ricdeau/protomapper/internal/types"
)

type CastMapper struct {
	typ        ast.Named
	targetType types.Type
}

func NewCastMapper(typ ast.Named, targetType types.Type) *CastMapper {
	return &CastMapper{typ: typ, targetType: targetType}
}

func (c *CastMapper) FromPb(fieldName string) types.FieldMapperFunc {
	return func(src string) string {
		if src != "" {
			src += "."
		}
		return call(c.targetType.GetName())(src + fieldName)
	}
}

func (c *CastMapper) ToPb(fieldName string) types.FieldMapperFunc {
	return func(src string) string {
		if src != "" {
			src += "."
		}
		scalarType := c.typ.(ast.ScalarNode)
		return cast(scalarType)(src + fieldName)
	}
}

func (c *CastMapper) CastFuncFromProto() string {
	scalarType := c.typ.(ast.ScalarNode)
	return fmt.Sprintf("func (x %s) { return %s }",
		helpers.GoTypeName(scalarType),
		call(c.targetType.GetName())("x"),
	)
}

func (c *CastMapper) CastFuncToProto() string {
	scalarType := c.typ.(ast.ScalarNode)
	return fmt.Sprintf("func (x %s) { return %s }",
		c.targetType.GetName(),
		call(helpers.GoTypeName(scalarType))("x"),
	)
}

func cast(protoType ast.ScalarNode) func(s string) string {
	return call(helpers.GoTypeName(protoType))
}

func call(callName string) func(s string) string {
	return func(s string) string {
		return callName + "(" + s + ")"
	}
}
