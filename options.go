package protomapper

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
