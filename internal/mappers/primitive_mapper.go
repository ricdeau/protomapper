package mappers

import (
	"fmt"

	"github.com/ricdeau/protoast/ast"
)

type ScalarMapper struct {
}

func (ScalarMapper) GoTypeName(t ast.ScalarNode) string {
	switch v := t.(type) {
	case *ast.Float64:
		return "float64"
	case *ast.Float32:
		return "float32"
	case *ast.Int32, *ast.Sint32, *ast.Sfixed32:
		return "int32"
	case *ast.Int64, *ast.Sint64, *ast.Sfixed64:
		return "int64"
	case *ast.Uint32, *ast.Fixed32:
		return "uint32"
	case *ast.Uint64, *ast.Fixed64:
		return "uint64"
	case *ast.Bool:
		return "bool"
	case *ast.String:
		return "string"
	case *ast.Bytes:
		return "[]byte"
	default:
		panic(fmt.Sprintf("%T not a scalar type", v))
	}
}
