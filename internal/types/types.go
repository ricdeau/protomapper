package types

type FieldMapperFunc = func(dest string) string

type Type interface {
	GetName() string
	GetComment() []string
	GetFields() []Field
}

type Field interface {
	GetName() string
	GetComment() []string
	GetType() Type
}

type FieldMapper interface {
	FromProto(fieldName string) FieldMapperFunc
	ToProto(fieldName string) FieldMapperFunc
}
