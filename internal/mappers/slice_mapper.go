package mappers

import (
	"fmt"

	"github.com/ricdeau/protoast/ast"
	"github.com/ricdeau/protomapper/internal/helpers"
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

func (s *SliceMapper) FromPb(fieldName string) types.FieldMapperFunc {
	return func(src string) string {
		m := registry.Mappers.Get(s.typ.(ast.Named).GetFullName())
		if m == nil {
			AddMapper(s.typ.(ast.Type))
			m = registry.Mappers.Get(s.typ.(ast.Named).GetFullName())
		}

		targetType := registry.Types.GetType(s.typ)
		targetTypeName := targetType.GetName()
		if _, ok := targetType.(*types.Struct); ok {
			targetTypeName = "*types." + targetTypeName
		}
		mapFunc := m.FromPb("x")("")
		return fmt.Sprintf(mapSlicePattern,
			fieldName, helpers.GoTypeName(s.typ), targetTypeName, mapFunc)
	}
}

func (s *SliceMapper) ToPb(fieldName string) types.FieldMapperFunc {
	return func(src string) string {
		m := registry.Mappers.Get(s.typ.(ast.Named).GetFullName())
		if m == nil {
			AddMapper(s.typ.(ast.Type))
			m = registry.Mappers.Get(s.typ.(ast.Named).GetFullName())
		}

		srcType := registry.Types.GetType(s.typ)
		srcTypeName := srcType.GetName()
		if _, ok := srcType.(*types.Struct); ok {
			srcTypeName = "*types." + srcTypeName
		}
		mapFunc := m.ToPb("x")("")
		return fmt.Sprintf(mapSlicePattern,
			fieldName, srcTypeName, helpers.GoTypeName(s.typ), mapFunc)
	}
}
