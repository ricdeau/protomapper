package mappers

import (
	"fmt"

	"github.com/ricdeau/protoast/ast"
	"github.com/ricdeau/protomapper/internal/dicts"
	"github.com/ricdeau/protomapper/internal/types"
)

var _ types.FieldMapper = (*FieldMapper)(nil)

type FieldMappers struct {
	mappersDict *dicts.MappersDict
	fieldsDict  *dicts.FieldDict
}

func NewFieldMappers(mappersDict *dicts.MappersDict, fieldsDict *dicts.FieldDict) *FieldMappers {
	return &FieldMappers{
		mappersDict: mappersDict,
		fieldsDict:  fieldsDict,
	}
}

type FieldMapper struct {
	Parent    *FieldMappers
	Proto     ast.Field
	Field     types.Field
	FromProto types.FieldMapperFunc
	ToProto   types.FieldMapperFunc
}

func (m *FieldMapper) GetFromProto(src, dest string) string {
	return m.FromProto(src, dest)
}

func (m *FieldMapper) GetToProto(src, dest string) string {
	return m.ToProto(src, dest)
}

func (m *FieldMappers) GetMapper(protoField ast.Field) types.FieldMapper {
	field := m.fieldsDict.Get(protoField)
	if field == nil {
		panic(fmt.Sprintf("field for %T not registered", protoField))
	}

	switch field.GetType() {
	case types.String, types.Bool:
		return m.simpleMapperFor(field, protoField)
	case types.Int:
		return m.castMapperFor(field, protoField)
	default:
		panic(fmt.Sprintf("mapper for %T not supported", field.GetType()))
	}
}
