package mappers

import (
	"testing"

	"github.com/ricdeau/protomapper/internal/registry"
	"github.com/stretchr/testify/require"
)

func TestFieldMappers_castMapperFor(t *testing.T) {
	const fieldName = "SomeField"
	tests := []struct {
		name          string
		wantFromProto string
		wantToProto   string
	}{
		{
			name:          "int32",
			wantFromProto: "int(src.SomeField)",
			wantToProto:   "int32(src.SomeField)",
		},
		{
			name:          "sfixed64",
			wantFromProto: "int(src.SomeField)",
			wantToProto:   "int64(src.SomeField)",
		},
		{
			name:          "uint32",
			wantFromProto: "int(src.SomeField)",
			wantToProto:   "uint32(src.SomeField)",
		},
		{
			name:          "float",
			wantFromProto: "float64(src.SomeField)",
			wantToProto:   "float32(src.SomeField)",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			mapper := registry.Mappers.Get(tt.name)
			gotFromProto := mapper.FromProto(fieldName)("src")
			gotToProto := mapper.ToProto(fieldName)("src")

			require.Equal(t, tt.wantFromProto, gotFromProto)
			require.Equal(t, tt.wantToProto, gotToProto)
		})
	}
}
