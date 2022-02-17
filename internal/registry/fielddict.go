package registry

import (
	"sync"

	"github.com/ricdeau/protoast/ast"
	"github.com/ricdeau/protomapper/internal/types"
)

type FieldDict struct {
	inner map[ast.Field]types.Field
	mx    sync.RWMutex
}

func NewFieldDict() *FieldDict {
	return &FieldDict{
		inner: map[ast.Field]types.Field{},
	}
}

func (d *FieldDict) Get(k ast.Field) types.Field {
	d.mx.RLock()
	defer d.mx.RUnlock()

	return d.inner[k]
}

func (d *FieldDict) Put(k ast.Field, v types.Field) {
	d.mx.Lock()
	defer d.mx.Unlock()

	d.inner[k] = v
}

func (d *FieldDict) PutIfNotExist(k ast.Field, v types.Field) {
	d.mx.Lock()
	defer d.mx.Unlock()

	if val := d.inner[k]; val != nil {
		return
	}

	d.inner[k] = v
}
