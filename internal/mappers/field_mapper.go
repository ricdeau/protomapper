package mappers

import (
	"github.com/ricdeau/protomapper/internal/registry"
)

type FieldMappers struct {
	mappersDict *registry.MappersDict
	fieldsDict  *registry.FieldDict
}

func NewFieldMappers(mappersDict *registry.MappersDict, fieldsDict *registry.FieldDict) *FieldMappers {
	return &FieldMappers{
		mappersDict: mappersDict,
		fieldsDict:  fieldsDict,
	}
}
