package types

import (
	"fmt"
)

var _ Type = Primitive("")
var _ fmt.Stringer = Primitive("")

const (
	String  = Primitive("string")
	Bool    = Primitive("bool")
	Int     = Primitive("int")
	Float64 = Primitive("float64")
	Byte    = Primitive("byte")
)

type Primitive string

func (t Primitive) IsStruct() bool {
	return false
}

func (t Primitive) GetName() string {
	return string(t)
}

func (t Primitive) GetComment() []string {
	return nil
}

func (t Primitive) GetFields() []Field {
	return nil
}

func (t Primitive) String() string {
	return t.GetName()
}
