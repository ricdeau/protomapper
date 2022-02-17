package registry

import (
	"testing"

	"github.com/ricdeau/protoast/ast"
	"github.com/ricdeau/protomapper/internal/types"
	"github.com/stretchr/testify/require"
)

func TestNewTypeDict(t *testing.T) {
	d := NewTypeDict()
	require.NotNil(t, d)
	require.NotNil(t, d.inner)
	require.NotNil(t, d.names)
}

func TestTypeDict_Get(t *testing.T) {
	key := &ast.Message{
		Name: "SomeMessage",
	}
	val := &types.Struct{TypeName: "SomeMessage"}
	d := NewTypeDict()
	d.inner[key] = val

	got := d.Get(key)
	require.Equal(t, val, got)
}

func TestTypeDict_GetByName(t *testing.T) {
	key := &ast.Message{
		Name: "SomeMessage",
	}
	val := &types.Struct{TypeName: "SomeMessage"}
	d := NewTypeDict()
	d.inner[key] = val
	d.names[key.Name] = key

	k, v := d.GetByName("SomeMessage")
	require.Equal(t, key, k)
	require.Equal(t, val, v)
}

func TestTypeDict_Put(t *testing.T) {
	key := &ast.Message{
		Name: "SomeMessage",
	}
	val := &types.Struct{TypeName: "SomeMessage"}
	d := NewTypeDict()

	require.Len(t, d.inner, 0)
	require.Len(t, d.names, 0)

	d.Put(key, val)

	require.Len(t, d.inner, 1)
	require.Len(t, d.names, 1)
}

func TestTypeDict_PutIfNotExist(t *testing.T) {
	key := &ast.Message{
		Name: "SomeMessage",
	}
	val := &types.Struct{TypeName: "SomeMessage"}
	d := NewTypeDict()
	d.inner[key] = val
	d.names[key.Name] = key

	val2 := &types.Struct{TypeName: "SomeMessage2"}

	d.PutIfNotExist(key, val2)

	got := d.Get(key)

	require.Equal(t, val, got)
}
