package mappers

import (
	"testing"

	"github.com/ricdeau/protoast/ast"
	"github.com/ricdeau/protomapper/internal/dicts"
	"github.com/ricdeau/protomapper/internal/types"
	"github.com/stretchr/testify/require"
)

func TestFieldMappers_castMapperFor(t *testing.T) {
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
			name: "int32",
			protoField: &ast.MessageField{
				Name: "Integer",
				Type: &ast.Int32{},
			},
			wantFromProto: "result.SomeField = int(src.Integer)",
			wantToProto:   "result.Integer = int32(src.SomeField)",
		},
		{
			name: "sfixed64",
			protoField: &ast.MessageField{
				Name: "Fixed",
				Type: &ast.Sfixed64{},
			},
			wantFromProto: "result.SomeField = int(src.Fixed)",
			wantToProto:   "result.Fixed = int64(src.SomeField)",
		},
		{
			name: "uint32",
			protoField: &ast.MessageField{
				Name: "Integer",
				Type: &ast.Uint32{},
			},
			wantFromProto: "result.SomeField = int(src.Integer)",
			wantToProto:   "result.Integer = uint32(src.SomeField)",
		},
		{
			name: "float32",
			protoField: &ast.MessageField{
				Name: "Float",
				Type: &ast.Float32{},
			},
			wantFromProto: "result.SomeField = int(src.Float)",
			wantToProto:   "result.Float = float32(src.SomeField)",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			gotFromProto := fm.castMapperFor(field, tt.protoField).GetFromProto("src", "result")
			gotToProto := fm.castMapperFor(field, tt.protoField).GetToProto("src", "result")

			require.Equal(t, tt.wantFromProto, gotFromProto)
			require.Equal(t, tt.wantToProto, gotToProto)
		})
	}
}
