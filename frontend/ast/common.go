package ast

type ParameterWithOptionalType struct {
	Name string
	Type Type
}

type FunctionBody struct {
	Expr  Expr
	Block []Stmt
}
