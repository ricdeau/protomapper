package types

import (
	"fmt"
)

var _ Type = (*Array)(nil)
var _ fmt.Stringer = (*Array)(nil)

type Array struct {
	Elem Type
}

func (a *Array) IsStruct() bool {
	return false
}

func ArrayOf(elem Type) *Array {
	return &Array{Elem: elem}
}

func (a *Array) GetName() string {
	return "[]" + a.Elem.GetName()
}

func (a *Array) GetComment() []string {
	return nil
}

func (a *Array) GetFields() []Field {
	return nil
}

func (a *Array) String() string {
	return a.GetName()
}
