package types

type FieldMapperFunc func(src, dest string) string

type FieldMapper interface {
	GetFromProto(src, dest string) string
	GetToProto(src, dest string) string
}
