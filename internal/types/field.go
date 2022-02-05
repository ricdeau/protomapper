package types

import (
	"fmt"
)

var _ Field = (*field)(nil)
var _ fmt.Stringer = (*field)(nil)

type field struct {
	name    string
	comment []string
	t       Type
}

func NewField(name string, comment []string, t Type) *field {
	return &field{
		name:    name,
		comment: comment,
		t:       t,
	}
}

func (f *field) GetName() string {
	return f.name
}

func (f *field) GetType() Type {
	return f.t
}

func (f *field) GetComment() []string {
	return f.comment
}

func (f *field) String() string {
	return f.GetName()
}
