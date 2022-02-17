package mappers

import (
	"github.com/iancoleman/strcase"
	"github.com/pkg/errors"
	"github.com/ricdeau/protoast/ast"
	"github.com/ricdeau/protomapper/internal/registry"
	"github.com/ricdeau/protomapper/internal/types"
)

type TypeMapper struct {
	excludeMessageField func(field ast.Field) bool
}

func NewTypeMapper(excludeMessageField func(field ast.Field) bool) *TypeMapper {
	return &TypeMapper{
		excludeMessageField: excludeMessageField,
	}
}

func (m *TypeMapper) FromProtoType(t ast.Type) (result types.Type, err error) {
	if val := registry.Types.Get(t); val != nil {
		return val, nil
	}

	switch v := t.(type) {
	case *ast.String, *ast.Enum:
		result = types.String
	case *ast.Bool:
		result = types.Bool
	case *ast.Int32, *ast.Int64, *ast.Uint32, *ast.Uint64, *ast.Sint32, *ast.Sint64, *ast.Fixed32, *ast.Fixed64, *ast.Sfixed32, *ast.Sfixed64:
		result = types.Int
	case *ast.Float32, *ast.Float64:
		result = types.Float64
	case *ast.Bytes:
		result = types.ArrayOf(types.Byte)
	case *ast.Repeated:
		elem, err := m.FromProtoType(v.Type)
		if err != nil {
			return nil, errors.Wrap(err, "get array elem")
		}
		result = types.ArrayOf(elem)
	case *ast.Map:
		key, err := m.FromProtoType(v.KeyType)
		if err != nil {
			return nil, errors.Wrap(err, "get map key")
		}

		val, err := m.FromProtoType(v.ValueType)
		if err != nil {
			return nil, errors.Wrap(err, "get map value")
		}

		if key, ok := key.(types.Primitive); !ok {
			return nil, errors.Errorf("invalid key type: %s", key)
		} else {
			result = types.MapOf(key, val)
		}
	case *ast.Message:
		name := strcase.ToCamel(v.GetName())
		s := types.NewStruct(name, v.GetComment().GetLines())
		err := m.fillFormMessage(v, s)
		if err != nil {
			return nil, errors.Wrap(err, "fill fields from message")
		}
		result = s
	default:
		return nil, errors.Errorf("unsupported type %T", t)
	}

	registry.Types.PutIfNotExist(t, result)
	AddMapper(t)

	return
}

func (m *TypeMapper) FromProtoField(f ast.Field) (types.Field, error) {
	if val := registry.Fields.Get(f); val != nil {
		return val, nil
	}

	fieldName := strcase.ToCamel(f.GetName())
	fieldProtoType := ast.FieldType(f)
	fieldType, err := m.FromProtoType(fieldProtoType)
	if err != nil {
		return nil, errors.Wrap(err, "get field type")
	}

	field := types.NewField(fieldName, f.GetComment().GetLines(), fieldType)

	registry.Fields.PutIfNotExist(f, field)
	registry.Types.PutIfNotExist(fieldProtoType, fieldType)

	return field, nil
}

func (m *TypeMapper) fillFormMessage(msg *ast.Message, s *types.Struct) error {
	for _, field := range msg.AllFields() {
		if m.excludeMessageField(field) {
			continue
		}

		f, err := m.FromProtoField(field)
		if err != nil {
			return err
		}
		s.Fields = append(s.Fields, f)
	}

	return nil
}
