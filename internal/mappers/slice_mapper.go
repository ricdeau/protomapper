package mappers

import (
	"fmt"

	"github.com/ricdeau/protoast/ast"
	"github.com/ricdeau/protomapper/internal/registry"
	"github.com/ricdeau/protomapper/internal/types"
)

const mapSlicePattern = "Map(src.%[1]s, func(x %[2]s) %[3]s { return %[4]s} )"

type SliceMapper struct {
	typ ast.Type
}

func NewSliceMapper(typ ast.Type) *SliceMapper {
	return &SliceMapper{typ: typ}
}

func (s *SliceMapper) FromProto(fieldName string) types.FieldMapperFunc {
	return func(src string) string {
		m := registry.Mappers.Get(s.typ.(ast.Named).GetFullName())
		if m == nil {
			AddMapper(s.typ.(ast.Type))
			m = registry.Mappers.Get(s.typ.(ast.Named).GetFullName())
		}

		targetType := registry.Types.Get(s.typ)
		mapFunc := m.FromProto("x")("")
		return fmt.Sprintf(mapSlicePattern,
			fieldName, GoTypeName(s.typ), targetType.GetName(), mapFunc)
	}
}

func (s *SliceMapper) ToProto(fieldName string) types.FieldMapperFunc {
	return func(src string) string {
		m := registry.Mappers.Get(s.typ.(ast.Named).GetFullName())
		if m == nil {
			AddMapper(s.typ.(ast.Type))
			m = registry.Mappers.Get(s.typ.(ast.Named).GetFullName())
		}

		srcType := registry.Types.Get(s.typ)
		mapFunc := m.ToProto("x")("")
		return fmt.Sprintf(mapSlicePattern,
			fieldName, srcType.GetName(), GoTypeName(s.typ), mapFunc)
	}
}
