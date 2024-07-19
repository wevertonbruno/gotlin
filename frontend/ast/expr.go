package ast

import (
	"gotlin/frontend/token"
)

type Expr interface {
	expr()
}

type BinaryExpr struct {
	Left  Expr
	Right Expr
	Op    token.Token
}

func (e *BinaryExpr) expr() {}

type UnaryExpr struct {
	Op    token.Token
	Right Expr
}

func (e *UnaryExpr) expr() {}

type IntLiteral struct {
	Value int64
}

func (e *IntLiteral) expr() {}

type DoubleLiteral struct {
	Value float64
}

func (e *DoubleLiteral) expr() {}

type BoolLiteral struct {
	Value bool
}

func (e *BoolLiteral) expr() {}

type StringLiteral struct {
	Value string
}

func (e *StringLiteral) expr() {}

type GroupingExpr struct {
	Expr Expr
}

func (e *GroupingExpr) expr() {}

type IdentifierExpr struct {
	Value token.Token
}

func (e *IdentifierExpr) expr() {}

type NonNullableExpr struct {
	Expr Expr
}

func (e *NonNullableExpr) expr() {}
