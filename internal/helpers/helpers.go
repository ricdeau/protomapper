package helpers

import (
	"fmt"
	"path/filepath"

	"github.com/iancoleman/strcase"
	"github.com/ricdeau/protoast/ast"
)

func PbGoFieldName(pbField ast.Field) string {
	return strcase.ToCamel(pbField.GetName())
}

func PbGoStructName(msg *ast.Message) string {
	if msg.GetParentMsg() == nil {
		return strcase.ToCamel(msg.GetName())
	}

	return PbGoStructName(msg.GetParentMsg()) + "_" + msg.GetName()
}

func PbGoEnumName(e *ast.Enum) string {
	name := strcase.ToCamel(e.Name)
	msg := e.ParentMsg
	for msg != nil {
		name = strcase.ToCamel(e.ParentMsg.Name) + "_" + name
		msg = msg.ParentMsg
	}
	return name
}

func PkgFromDir(dir string) string {
	return filepath.Base(dir)
}

func GoTypeName(t ast.Type) string {
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
	case *ast.Enum:
		return "pb." + PbGoEnumName(v)
	case *ast.Message:
		return "*pb." + PbGoStructName(v)
	case *ast.Repeated:
		return "[]" + GoTypeName(v.Type)
	default:
		panic(fmt.Sprintf("%T unsupported", v))
	}
}
