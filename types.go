package protomapper

import (
	"path/filepath"

	"github.com/ricdeau/protoast/ast"
	"github.com/ricdeau/protomapper/internal/types"
)

// FilepathResolver resolver for files in directory path.
var FilepathResolver = func(path ...string) FileResolver {
	return func(s string) (string, error) {
		path := append(path, s)
		return filepath.Join(path...), nil
	}
}

type ProtoType = ast.Type
type ProtoField = ast.Field
type ProtoOption = ast.Option
type ProtoScalar = ast.ScalarNode

type Type = types.Type
type Field = types.Field

type FileResolver func(string) (string, error)
type ProtoMapperFunc func(protoField ProtoField) func(src, dest string) string

// TypeMapper mapper for protobuf types.
type TypeMapper interface {
	FromProtoType(t ProtoType) (Type, error)
	FromProtoField(f ProtoField) (Field, error)
}

// Renderer renderer for types.
type Renderer interface {
	Render(t Type) error
	SetTypeNameResolver(func(t Type) string)
}
