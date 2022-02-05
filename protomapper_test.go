package protomapper

import (
	"testing"

	"github.com/stretchr/testify/require"
)

var cfg = &Config{
	AppName:         "tests",
	TypesDir:        "testdata/types",
	TypesGoPackage:  "types",
	ConvertersDir:   "testdata/converters",
	ConvertersGoPkg: "converters",
	TypesImport:     "github.com/ricdeau/protomapper/testdata/types",
	PbImport:        "github.com/ricdeau/protomapper/testdata/gen/protos",
}

func TestTypeMapper(t *testing.T) {
	m := NewProtoMapper(cfg)
	err := m.ResolveTypes(FilepathResolver("testdata/protos"), "types.proto")
	require.NoError(t, err)

	protoType, _ := m.Types().GetByName("TestType")
	goType, err := m.TypeMapper().FromProtoType(protoType)
	require.NoError(t, err)
	require.NotNil(t, goType)

	t1, t2 := m.Types().GetByName("TestType")
	require.NotNil(t, t1)
	require.NotNil(t, t2)
}

func TestTypeRenderer(t *testing.T) {
	m := NewProtoMapper(cfg)
	err := m.ResolveTypes(FilepathResolver("testdata/protos"), "types.proto")
	require.NoError(t, err)

	protoType, _ := m.Types().GetByName("TestType")
	tp, err := m.TypeMapper().FromProtoType(protoType)
	require.NoError(t, err)

	err = m.TypeRenderer().Render(tp)
	require.NoError(t, err)
}

func TestConverterRenderer(t *testing.T) {
	m := NewProtoMapper(cfg)
	err := m.ResolveTypes(FilepathResolver("testdata/protos"), "types.proto")
	require.NoError(t, err)

	protoType, _ := m.Types().GetByName("SimpleType")
	tp, err := m.TypeMapper().FromProtoType(protoType)
	require.NoError(t, err)

	err = m.TypeRenderer().Render(tp)
	require.NoError(t, err)

	err = m.ConvertersRenderer().Render(tp)
	require.NoError(t, err)
}
