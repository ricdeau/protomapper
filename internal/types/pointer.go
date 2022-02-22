package types

import (
	"fmt"
)

var _ Type = (*Pointer)(nil)
var _ fmt.Stringer = (*Pointer)(nil)

type Pointer struct {
	Type Type
}

func (p *Pointer) IsStruct() bool {
	return false
}

func PointerOf(t Type) *Pointer {
	return &Pointer{Type: t}
}

func (p *Pointer) GetName() string {
	return "*" + p.Type.GetName()
}

func (p *Pointer) GetComment() []string {
	return p.Type.GetComment()
}

func (p Pointer) GetFields() []Field {
	return p.Type.GetFields()
}

func (p *Pointer) String() string {
	return p.GetName()
}
