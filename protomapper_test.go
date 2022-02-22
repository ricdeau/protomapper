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

	t1, t2 := m.Types().GetByName("TestType")
	require.NotNil(t, t1)
	require.NotNil(t, t2)
}

func TestTypeRenderer(t *testing.T) {
	m := NewProtoMapper(cfg)
	err := m.ResolveTypes(FilepathResolver("testdata/protos"), "types.proto")
	require.NoError(t, err)

	_, pbType := m.Types().GetByName("CompoundType")

	err = m.TypeRenderer().Render(pbType)
	require.NoError(t, err)
}

func TestConverterRenderer(t *testing.T) {
	m := NewProtoMapper(cfg)
	err := m.ResolveTypes(FilepathResolver("testdata/protos"), "types.proto")
	require.NoError(t, err)

	_, st := m.Types().GetByName("SimpleType")
	err = m.TypeRenderer().Render(st)
	require.NoError(t, err)
	err = m.ConvertersRenderer().Render(st)
	require.NoError(t, err)

	_, dat := m.Types().GetByName("Data")
	err = m.TypeRenderer().Render(dat)
	require.NoError(t, err)
	err = m.ConvertersRenderer().Render(dat)
	require.NoError(t, err)
}
