package protomapper

import (
	"github.com/ricdeau/protomapper/internal/mappers"
)

// ScalarGoTypeName get corresponding go type for protobuf scalar type.
func ScalarGoTypeName(t ProtoScalar) string {
	return mappers.ScalarMapper{}.GoTypeName(t)
}
