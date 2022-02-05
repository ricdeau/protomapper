package mappers

import (
	"github.com/iancoleman/strcase"
	"github.com/ricdeau/protoast/ast"
)

func goName(protoField ast.Field) string {
	return strcase.ToCamel(protoField.GetName())
}

func cast(protoType ast.ScalarNode) func(s string) string {
	return scalarCast(ScalarMapper{}.GoTypeName(protoType))
}

func scalarCast(typeName string) func(s string) string {
	return func(s string) string {
		return typeName + "(" + s + ")"
	}
}
