package registry

import (
	"sync"

	"github.com/ricdeau/protomapper/internal/types"
)

type MappersDict struct {
	inner map[string]types.FieldMapper
	mx    sync.RWMutex
}

func NewMappersDict() *MappersDict {
	return &MappersDict{
		inner: map[string]types.FieldMapper{},
	}
}
func (d *MappersDict) Get(pbTypeName string) types.FieldMapper {
	d.mx.RLock()
	defer d.mx.RUnlock()

	m, ok := d.inner[pbTypeName]
	if !ok {
		return nil
	}

	return m
}

func (d *MappersDict) Put(pbTypeName string, mapper types.FieldMapper) {
	d.mx.Lock()
	defer d.mx.Unlock()

	_, ok := d.inner[pbTypeName]
	if ok {
		return
	}

	d.inner[pbTypeName] = mapper
}
