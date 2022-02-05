package dicts

import (
	"sync"

	"github.com/ricdeau/protoast/ast"
	"github.com/ricdeau/protomapper/internal/types"
)

type mappersDictKey struct {
	f types.Field
	p ast.Field
}

type MappersDict struct {
	inner map[mappersDictKey]types.FieldMapper
	mx    sync.RWMutex
}

func NewMappersDict() *MappersDict {
	return &MappersDict{
		inner: map[mappersDictKey]types.FieldMapper{},
	}
}
func (d *MappersDict) Get(f types.Field, p ast.Field) types.FieldMapper {
	d.mx.RLock()
	defer d.mx.RUnlock()

	return d.inner[mappersDictKey{f, p}]
}

func (d *MappersDict) Put(f types.Field, p ast.Field, v types.FieldMapper) {
	d.mx.Lock()
	defer d.mx.Unlock()

	key := mappersDictKey{f, p}
	_, ok := d.inner[key]
	if ok {
		return
	}

	d.inner[key] = v
}
