package types

import (
	"github.com/ricdeau/protoast/ast"
)

type FieldMapperFunc = func(dest string) string

type Type interface {
	GetName() string
	GetComment() []string
	GetFields() []Field
	IsStruct() bool
}

type Field interface {
	GetName() string
	GetComment() []string
	GetType() Type
}

type FieldMapper interface {
	FromPb(fieldName string) FieldMapperFunc
	ToPb(fieldName string) FieldMapperFunc
}

type TypeResolver interface {
	Resolve(typ Type, pbType ast.Type) string
}

type FieldResolver interface {
	Resolve(field Field, pbType ast.Type) string
}

type ImportsResolver interface {
	Resolve(typ Type, pbType ast.Type) []string
}
