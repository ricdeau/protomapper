package types

type Type interface {
	GetName() string
	GetComment() []string
	GetFields() []Field
}

type Field interface {
	GetName() string
	GetComment() []string
	GetType() Type
}
