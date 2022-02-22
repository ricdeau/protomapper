package registry

import (
	"sync"

	"github.com/ricdeau/protoast/ast"
	"github.com/ricdeau/protomapper/internal/types"
)

type TypeDict struct {
	types   map[ast.Type]types.Type
	pbTypes map[types.Type]ast.Type
	names   map[string]types.Type
	mx      sync.RWMutex
}

func NewTypeDict() *TypeDict {
	return &TypeDict{
		types:   map[ast.Type]types.Type{},
		pbTypes: map[types.Type]ast.Type{},
		names:   map[string]types.Type{},
	}
}

func (d *TypeDict) GetPbType(typ types.Type) ast.Type {
	d.mx.RLock()
	defer d.mx.RUnlock()

	return d.pbTypes[typ]
}

func (d *TypeDict) GetType(pbType ast.Type) types.Type {
	d.mx.RLock()
	defer d.mx.RUnlock()

	return d.types[pbType]
}

func (d *TypeDict) GetByName(name string) (types.Type, ast.Type) {
	d.mx.RLock()
	defer d.mx.RUnlock()

	typ, ok := d.names[name]
	if !ok {
		return nil, nil
	}

	return typ, d.pbTypes[typ]
}

func (d *TypeDict) Put(pbType ast.Type, typ types.Type) {
	d.mx.Lock()
	defer d.mx.Unlock()

	if t := d.types[pbType]; t == nil {
		d.types[pbType] = typ
		d.names[typ.GetName()] = typ
	}

	if t := d.pbTypes[typ]; t == nil {
		d.pbTypes[typ] = pbType
	}
}
