package ast

/**
import (
	"gotlin/frontend/object"
	"gotlin/frontend/token"
)

type Interpreter struct{}

func NewInterpreter() *Interpreter {
	return &Interpreter{}
}

func (i *Interpreter) VisitBinaryExpr(expr *BinaryExpr) object.Object {
	left := i.evaluate(expr.Left)
	right := i.evaluate(expr.Right)

	switch {
	case checkType(object.IntType, left, right):
		return evaluateIntBinaryExpr(expr.Op, left, right)
	default:
		// TODO CHECK
		return object.NULL
	}
}

func (i *Interpreter) VisitUnaryExpr(expr *UnaryExpr) object.Object {
	right := i.evaluate(expr.Right)
	switch expr.Op.Kind {
	case token.NOT:
		return evaluateNegation(right)
	case token.OpMinus:
		value := right.(*object.Int).Value
		return &object.Int{Value: -value}
	default:
		return object.NULL
	}
}

func (i *Interpreter) VisitIntLiteral(expr *IntLiteral) object.Object {
	return &object.Int{
		Value: expr.Value,
	}
}

func (i *Interpreter) VisitBooleanLiteral(expr *BoolLiteral) object.Object {
	return &object.Boolean{
		Value: expr.Value,
	}
}

func (i *Interpreter) VisitGroupingExpr(expr *GroupingExpr) object.Object {
	return i.evaluate(expr.Expr)
}

func (i *Interpreter) evaluate(expr Expr) object.Object {
	return expr.Accept(i)
}

func evaluateNegation(right object.Object) object.Object {
	switch right {
	case object.TRUE:
		return object.FALSE
	case object.FALSE:
		return object.TRUE
	default:
		panic("Type Checking")
	}
}

func evaluateIntBinaryExpr(op token.Token, left object.Object, right object.Object) object.Object {
	leftVal := left.(*object.Int).Value
	rightVal := right.(*object.Int).Value

	switch op.Kind {
	case token.OpPlus:
		return &object.Int{Value: leftVal + rightVal}
	case token.OpMinus:
		return &object.Int{Value: leftVal - rightVal}
	case token.OpMulti:
		return &object.Int{Value: leftVal * rightVal}
	case token.OpDivide:
		return &object.Int{Value: leftVal / rightVal}
	case token.OpGt:
		return &object.Boolean{Value: leftVal > rightVal}
	case token.OpGte:
		return &object.Boolean{Value: leftVal >= rightVal}
	case token.OpLt:
		return &object.Boolean{Value: leftVal < rightVal}
	case token.OpLte:
		return &object.Boolean{Value: leftVal <= rightVal}
	default:
		//TODO CHECK
		return object.NULL
	}
}

func checkType(t object.Type, objs ...object.Object) bool {
	for _, obj := range objs {
		if obj.Type() != t {
			return false
		}
	}
	return true
}

*/
