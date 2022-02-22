package mappers

import (
	"github.com/ricdeau/protomapper/internal/types"
)

type SimpleMapper struct{}

func (s SimpleMapper) FromPb(fieldName string) types.FieldMapperFunc {
	return func(src string) string {
		if src != "" {
			src += "."
		}
		return src + fieldName
	}
}

func (s SimpleMapper) ToPb(fieldName string) types.FieldMapperFunc {
	return func(src string) string {
		if src != "" {
			src += "."
		}
		return src + fieldName
	}
}
