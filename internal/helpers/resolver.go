package helpers

import (
	"fmt"

	"github.com/iancoleman/strcase"
	"github.com/ricdeau/protoast/ast"
	"github.com/ricdeau/protomapper/internal/types"
)

var CamelCaseName = TypeResolverFunc(func(typ types.Type, pbType ast.Type) string {
	if typ != nil {
		return strcase.ToCamel(typ.GetName())
	}
	named, ok := pbType.(ast.Named)
	if ok {
		return strcase.ToCamel(named.GetName())
	}
	panic(fmt.Sprintf("neither go type nor pb named type provided: %s, %s", typ, pbType))
})

var SnakeCaseGoTypeFile = TypeResolverFunc(func(typ types.Type, pbType ast.Type) string {
	return strcase.ToSnake(typ.GetName()) + ".go"
})

var DefaultImportsResolver = ImportsResolverFunc(func(typ types.Type, pbType ast.Type) []string {
	return nil
})

var StandardFieldResolver = FieldResolverFunc(func(field types.Field, pbType ast.Type) string {
	return field.GetName() + " " + field.GetType().GetName()
})

type TypeResolverFunc func(typ types.Type, pbType ast.Type) string

func (f TypeResolverFunc) Resolve(typ types.Type, pbType ast.Type) string {
	return f(typ, pbType)
}

type FieldResolverFunc func(field types.Field, pbType ast.Type) string

func (f FieldResolverFunc) Resolve(field types.Field, pbType ast.Type) string {
	return f(field, pbType)
}

type ImportsResolverFunc func(typ types.Type, pbType ast.Type) []string

func (f ImportsResolverFunc) Resolve(typ types.Type, pbType ast.Type) []string {
	return f(typ, pbType)
}
