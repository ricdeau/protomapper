package protomapper

import (
	"github.com/ricdeau/protomapper/internal/mappers"
)

type Option interface {
	apply(mapper *ProtoMapper)
}

var _ Option = (*excludeMessageFields)(nil)

type excludeMessageFields func(field ProtoField) bool

func WithExcludeMessageFields(exclude func(field ProtoField) bool) Option {
	return excludeMessageFields(exclude)
}

func (o excludeMessageFields) apply(mapper *ProtoMapper) {
	mapper.excludeMessageField = o
}

var _ Option = (*fieldMapperOption)(nil)

type fieldMapperOption struct {
	*mappers.FieldMapper
}

func WithFieldMapper(field Field, fromProto, toProto ProtoMapperFunc) Option {
	return &fieldMapperOption{
		FieldMapper: &mappers.FieldMapper{
			Field: field,
		},
	}
}

func (o *fieldMapperOption) apply(mapper *ProtoMapper) {

}
