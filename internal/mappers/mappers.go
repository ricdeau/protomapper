package mappers

import (
	"fmt"

	"github.com/iancoleman/strcase"
	"github.com/ricdeau/protoast/ast"
	"github.com/ricdeau/protomapper/internal/registry"
	"github.com/ricdeau/protomapper/internal/types"
)

var (
	simpleMapper = new(SimpleMapper)
)

func init() {
	registry.Mappers.Put(new(ast.String).GetFullName(), simpleMapper)
	registry.Mappers.Put(new(ast.Bool).GetFullName(), simpleMapper)
	registry.Mappers.Put(new(ast.Bytes).GetFullName(), simpleMapper)
	registry.Mappers.Put(repeatedKey(new(ast.String)), simpleMapper)
	registry.Mappers.Put(repeatedKey(new(ast.Bool)), simpleMapper)
	registry.Mappers.Put(repeatedKey(new(ast.Bytes)), simpleMapper)
	registry.Mappers.Put(new(ast.Int32).GetFullName(), NewCastMapper(new(ast.Int32), types.Int))
	registry.Mappers.Put(new(ast.Int64).GetFullName(), NewCastMapper(new(ast.Int64), types.Int))
	registry.Mappers.Put(new(ast.Uint32).GetFullName(), NewCastMapper(new(ast.Uint32), types.Int))
	registry.Mappers.Put(new(ast.Uint64).GetFullName(), NewCastMapper(new(ast.Uint64), types.Int))
	registry.Mappers.Put(new(ast.Sint32).GetFullName(), NewCastMapper(new(ast.Sint32), types.Int))
	registry.Mappers.Put(new(ast.Sint64).GetFullName(), NewCastMapper(new(ast.Sint64), types.Int))
	registry.Mappers.Put(new(ast.Fixed32).GetFullName(), NewCastMapper(new(ast.Fixed32), types.Int))
	registry.Mappers.Put(new(ast.Fixed64).GetFullName(), NewCastMapper(new(ast.Fixed64), types.Int))
	registry.Mappers.Put(new(ast.Sfixed32).GetFullName(), NewCastMapper(new(ast.Sfixed32), types.Int))
	registry.Mappers.Put(new(ast.Sfixed64).GetFullName(), NewCastMapper(new(ast.Sfixed64), types.Int))
	registry.Mappers.Put(new(ast.Float32).GetFullName(), NewCastMapper(new(ast.Float32), types.Float64))
	registry.Mappers.Put(new(ast.Float64).GetFullName(), NewCastMapper(new(ast.Float64), types.Float64))

}

func AddMapper(pbType ast.Type) {
	switch v := pbType.(type) {
	case *ast.Enum:
		registry.Mappers.Put(v.GetFullName(), NewEnumMapper(v))
	case *ast.Repeated:
		elType := v.Type.(ast.Named)
		registry.Mappers.Put(repeatedKey(elType), NewSliceMapper(v.Type))
	}
}

func GetMapper(protoField ast.Field) (types.FieldMapper, error) {
	fieldType := ast.FieldType(protoField)
	var key string
	switch v := fieldType.(type) {
	case ast.Named:
		key = v.GetFullName()
	case *ast.Repeated:
		key = repeatedKey(v.Type.(ast.Named))
	}
	mapper := registry.Mappers.Get(key)
	if mapper == nil {
		return nil, fmt.Errorf("mapper for type %T not registered", fieldType)
	}

	return mapper, nil
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
		return enumGoName(v)
	default:
		panic(fmt.Sprintf("%T unsupported", v))
	}
}

func repeatedKey(typ ast.Named) string {
	return "repeated " + typ.GetFullName()
}

func enumGoName(e *ast.Enum) string {
	name := strcase.ToCamel(e.Name)
	msg := e.ParentMsg
	for msg != nil {
		name = strcase.ToCamel(e.ParentMsg.Name) + "_" + name
		msg = msg.ParentMsg
	}
	return "pb." + name
}
