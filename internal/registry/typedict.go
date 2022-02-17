package registry

import (
	"sync"

	"github.com/ricdeau/protoast/ast"
	"github.com/ricdeau/protomapper/internal/types"
)

type TypeDict struct {
	inner map[ast.Type]types.Type
	names map[string]ast.Type
	mx    sync.RWMutex
}

func NewTypeDict() *TypeDict {
	return &TypeDict{
		inner: map[ast.Type]types.Type{},
		names: map[string]ast.Type{},
	}
}

func (d *TypeDict) Get(k ast.Type) types.Type {
	d.mx.RLock()
	defer d.mx.RUnlock()

	return d.inner[k]
}

func (d *TypeDict) GetByName(name string) (ast.Type, types.Type) {
	d.mx.RLock()
	defer d.mx.RUnlock()

	t, ok := d.names[name]
	if !ok {
		return nil, nil
	}

	return t, d.inner[t]
}

func (d *TypeDict) Put(k ast.Type, v types.Type) {
	d.mx.Lock()
	defer d.mx.Unlock()

	if n, ok := k.(ast.Named); ok {
		d.names[n.GetName()] = k
	}

	d.inner[k] = v
}

func (d *TypeDict) PutIfNotExist(k ast.Type, v types.Type) {
	d.mx.Lock()
	defer d.mx.Unlock()

	if val := d.inner[k]; val != nil {
		return
	}

	if n, ok := k.(ast.Named); ok {
		d.names[n.GetName()] = k
	}

	d.inner[k] = v
}
