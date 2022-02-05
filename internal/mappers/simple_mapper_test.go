package mappers

import (
	"testing"

	"github.com/ricdeau/protoast/ast"
	"github.com/ricdeau/protomapper/internal/dicts"
	"github.com/ricdeau/protomapper/internal/types"
	"github.com/stretchr/testify/require"
)

func TestFieldMappers_simpleMapperFor(t *testing.T) {
	fm := &FieldMappers{
		mappersDict: dicts.NewMappersDict(),
	}

	field := types.NewField("SomeField", nil, types.Int)

	tests := []struct {
		name          string
		protoField    ast.Field
		wantFromProto string
		wantToProto   string
	}{
		{
			name: "string",
			protoField: &ast.MessageField{
				Name: "Name",
				Type: &ast.String{},
			},
			wantFromProto: "result.SomeField = src.Name",
			wantToProto:   "result.Name = src.SomeField",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			gotFromProto := fm.simpleMapperFor(field, tt.protoField).GetFromProto("src", "result")
			gotToProto := fm.simpleMapperFor(field, tt.protoField).GetToProto("src", "result")

			require.Equal(t, tt.wantFromProto, gotFromProto)
			require.Equal(t, tt.wantToProto, gotToProto)
		})
	}
}
