package mappers

import (
	"testing"

	"github.com/ricdeau/protomapper/internal/registry"
	"github.com/stretchr/testify/require"
)

func TestFieldMappers_simpleMapperFor(t *testing.T) {

	const fieldName = "SomeField"

	tests := []struct {
		name          string
		wantFromProto string
		wantToProto   string
	}{
		{
			name:          "string",
			wantFromProto: "src.SomeField",
			wantToProto:   "src.SomeField",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			mapper := registry.Mappers.Get(tt.name)
			gotFromProto := mapper.FromPb(fieldName)("src")
			gotToProto := mapper.ToPb(fieldName)("src")

			require.Equal(t, tt.wantFromProto, gotFromProto)
			require.Equal(t, tt.wantToProto, gotToProto)
		})
	}
}
