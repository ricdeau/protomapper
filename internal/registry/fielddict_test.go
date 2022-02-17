package registry

import (
	"testing"

	"github.com/ricdeau/protoast/ast"
	"github.com/ricdeau/protomapper/internal/types"
	"github.com/stretchr/testify/require"
)

func TestNewFieldDict(t *testing.T) {
	d := NewFieldDict()
	require.NotNil(t, d)
	require.NotNil(t, d.inner)
}

func TestFieldDict_Get(t *testing.T) {
	key := &ast.MessageField{
		Name: "field",
	}
	val := types.NewField("field", nil, types.String)
	d := NewFieldDict()
	d.inner[key] = val

	got := d.Get(key)
	require.Equal(t, val, got)
}

func TestFieldDict_Put(t *testing.T) {
	key := &ast.MessageField{
		Name: "field",
	}
	val := types.NewField("field", nil, types.String)
	d := NewFieldDict()

	require.Len(t, d.inner, 0)

	d.Put(key, val)

	require.Len(t, d.inner, 1)
}

func TestFieldDict_PutIfNotExist(t *testing.T) {
	key := &ast.MessageField{
		Name: "field",
	}
	val := types.NewField("field", nil, types.String)
	d := NewFieldDict()
	d.inner[key] = val

	val2 := types.NewField("field2", nil, types.Int)

	d.PutIfNotExist(key, val2)

	got := d.Get(key)

	require.Equal(t, val, got)
}
