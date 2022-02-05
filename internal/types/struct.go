package types

import (
	"fmt"
)

var _ Type = (*Struct)(nil)
var _ fmt.Stringer = (*Struct)(nil)

type Struct struct {
	TypeName string
	Comment  []string
	Fields   []Field
}

func NewStruct(name string, comment []string) *Struct {
	return &Struct{
		TypeName: name,
		Comment:  comment,
	}
}

func (s *Struct) GetName() string {
	return s.TypeName
}

func (s *Struct) GetComment() []string {
	return s.Comment
}

func (s *Struct) GetFields() []Field {
	return s.Fields
}

func (s *Struct) String() string {
	return s.GetName()
}
