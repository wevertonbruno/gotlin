package ast

import (
	"gotlin/frontend/token"
)

type Stmt interface {
	stmt()
}

type BlockStmt struct {
	Statements []Stmt
}

func (s *BlockStmt) stmt() {}

type ExprStmt struct {
	Expr Expr
}

func (s *ExprStmt) stmt() {}

type VariableDecl struct {
	Name     token.Token
	Type     Type
	Value    Expr
	ReadOnly bool
}

func (s *VariableDecl) stmt() {}

type AssignStmt struct {
	Assigne Expr
	Value   Expr
}

func (t *AssignStmt) stmt() {}

type ClassPrimaryConstructor struct {
	Parameters []ClassParam
}

type ClassParam struct {
	Name         string
	Type         Type
	DefaultValue Expr
	ReadOnly     bool
}

type ClassDeclStmt struct {
	Name               token.Token
	PrimaryConstructor *ClassPrimaryConstructor
	Params             []ClassParam
}

func (t *ClassDeclStmt) stmt() {}
