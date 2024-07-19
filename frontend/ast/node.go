package ast

import (
	"gotlin/frontend/object"
)

type Node interface {
	Accept(visitor Visitor) object.Object
}

type Visitor interface {
	VisitProgram(p *Program) object.Object
	VisitBinaryExpr(expr *BinaryExpr) object.Object
	VisitUnaryExpr(expr *UnaryExpr) object.Object
	VisitIntLiteral(expr *IntLiteral) object.Object
	VisitDoubleLiteral(expr *DoubleLiteral) object.Object
	VisitBooleanLiteral(expr *BoolLiteral) object.Object
	VisitGroupingExpr(expr *GroupingExpr) object.Object
	VisitStringLiteral(expr *StringLiteral) object.Object

	VisitBlockStmt(expr *BlockStmt) object.Object
	VisitExprStmt(expr *ExprStmt) object.Object
	//VisitAssignStmt(expr *Ass) object.Object

	//VisitVarDecl(expr *VarDecl) object.Object
	VisitVariableDecl(expr *VariableDecl) object.Object
	VisitType(expr *Type) object.Object
	VisitVariable(identifier *IdentifierExpr) object.Object
}

type Program struct {
	Statements []Stmt
}

func (p *Program) Accept(visitor Visitor) object.Object {
	return visitor.VisitProgram(p)
}
