package mappers

import (
	"github.com/ricdeau/protoast/ast"
	"github.com/ricdeau/protomapper/internal/types"
)

func (m *FieldMappers) simpleMapperFor(field types.Field, protoField ast.Field) (result types.FieldMapper) {
	if result = m.mappersDict.Get(field, protoField); result != nil {
		return
	}

	res := &FieldMapper{
		Parent: m,
		Field:  field,
		Proto:  protoField,
	}

	res.FromProto = func(src, dest string) string {
		fieldName := res.Field.GetName()
		protoFieldName := goName(res.Proto)
		return dest + "." + fieldName + " = " + src + "." + protoFieldName
	}

	res.ToProto = func(src, dest string) string {
		fieldName := res.Field.GetName()
		protoFieldName := goName(res.Proto)
		return dest + "." + protoFieldName + " = " + src + "." + fieldName
	}

	m.mappersDict.Put(field, protoField, res)

	return res
}
