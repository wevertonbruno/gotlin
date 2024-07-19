package ast

type Type interface {
	tp()
}

type NullableType struct {
	Type Type
}

func (nt *NullableType) tp() {}

type TypeName struct {
	Name string
}

func (id *TypeName) tp() {}

type ArrayType struct {
	Underlying Type
}

func (id *ArrayType) tp() {}
