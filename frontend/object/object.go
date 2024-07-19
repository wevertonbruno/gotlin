package object

import "fmt"

type Type string

const (
	IntType     Type = "Int"
	BooleanType Type = "Boolean"
	NullType    Type = "null"
	UnitType    Type = "Unit"
)

var (
	NULL  = &Null{}
	TRUE  = &Boolean{true}
	FALSE = &Boolean{false}
	UNIT  = &Unit{}
)

type Object interface {
	Type() Type
	Inspect() string
}

type Int struct {
	Value int64
}

func (i *Int) Inspect() string { return fmt.Sprintf("%d", i.Value) }
func (i *Int) Type() Type      { return IntType }

type Boolean struct {
	Value bool
}

func (b *Boolean) Inspect() string { return fmt.Sprintf("%t", b.Value) }
func (b *Boolean) Type() Type      { return BooleanType }

type Null struct{}

func (n *Null) Inspect() string { return "null" }
func (n *Null) Type() Type      { return NullType }

type Unit struct{}

func (u *Unit) Inspect() string { return "Unit" }
func (n *Unit) Type() Type      { return UnitType }
