package mappers

import (
	"github.com/ricdeau/protoast/ast"
	"github.com/ricdeau/protomapper/internal/types"
)

func (m *FieldMappers) castMapperFor(field types.Field, protoField ast.Field) (result types.FieldMapper) {
	if result = m.mappersDict.Get(field, protoField); result != nil {
		return
	}

	res := &FieldMapper{
		Parent: m,
		Field:  field,
		Proto:  protoField,
	}

	res.FromProto = func(src, dest string) string {
		return dest + "." + res.Field.GetName() + " = " + scalarCast("int")(src+"."+goName(protoField))
	}

	res.ToProto = func(src, dest string) string {
		scalarType := ast.FieldType(protoField).(ast.ScalarNode)
		return dest + "." + goName(protoField) + " = " + cast(scalarType)(src+"."+res.Field.GetName())
	}

	m.mappersDict.Put(field, protoField, res)

	return res
}
