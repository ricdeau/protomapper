package types

import (
	"fmt"
)

var _ Type = (*Map)(nil)
var _ fmt.Stringer = (*Map)(nil)

type Map struct {
	Key Primitive
	Val Type
}

func (m *Map) IsStruct() bool {
	return false
}

func MapOf(key Primitive, val Type) *Map {
	return &Map{Key: key, Val: val}
}

func (m *Map) GetName() string {
	return "map[" + m.Key.GetName() + "]" + m.Val.GetName()
}

func (m *Map) GetComment() []string {
	return nil
}

func (m *Map) GetFields() []Field {
	return nil
}

func (m *Map) String() string {
	return m.GetName()
}
