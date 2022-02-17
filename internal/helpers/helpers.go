package helpers

import (
	"github.com/iancoleman/strcase"
	"github.com/ricdeau/protoast/ast"
)

func GoName(protoField ast.Field) string {
	return strcase.ToCamel(protoField.GetName())
}
