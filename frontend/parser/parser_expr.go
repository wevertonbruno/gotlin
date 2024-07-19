package parser

import (
	"fmt"
	"strconv"

	"gotlin/frontend/ast"
	"gotlin/frontend/token"
)

func (p *Parser) parsePrimaryExpr() (ast.Expr, error) {
	switch p.currentTokenKind() {
	case token.IntLit:
		value, err := strconv.ParseInt(p.advance().Spelling, 10, 64)
		return &ast.IntLiteral{
			Value: value,
		}, err
	case token.StringLit:
		return &ast.StringLiteral{
			Value: p.advance().Spelling,
		}, nil
	case token.BooleanLit:
		return &ast.BoolLiteral{
			Value: p.advance().Spelling == "true",
		}, nil
	case token.IDENTIFIER:
		return &ast.IdentifierExpr{
			Value: p.advance(),
		}, nil
	default:
		return nil, NewError(fmt.Sprintf("Expected primary expression, got %s", p.currentTokenKind()))
	}
}

func (p *Parser) parseBinaryExpr(left ast.Expr, precedence BindingPower) (ast.Expr, error) {
	operator := p.advance()

	right, err := p.parseExpr(precedence)
	if err != nil {
		return nil, err
	}

	return &ast.BinaryExpr{
		Left:  left,
		Op:    operator,
		Right: right,
	}, nil
}

func (p *Parser) parseUnaryExpr() (ast.Expr, error) {
	operator := p.advance()
	right, err := p.parseExpr(Default)
	if err != nil {
		return nil, err
	}

	return &ast.UnaryExpr{
		Op:    operator,
		Right: right,
	}, nil
}

func (p *Parser) parseGroupingExpr() (ast.Expr, error) {
	p.advance()
	expr, err := p.parseExpr(Default)
	if err != nil {
		return nil, err
	}

	_, err = p.expected(token.RParen)
	if err != nil {
		return nil, err
	}

	return &ast.GroupingExpr{
		Expr: expr,
	}, nil
}

func (p *Parser) parseNotNullExpr(left ast.Expr, precedence BindingPower) (ast.Expr, error) {
	p.advance()

	//TODO Check all possibilities
	if _, ok := left.(*ast.IdentifierExpr); !ok {
		return nil, NewError(fmt.Sprintf("Identifier expected"))
	}

	return &ast.NonNullableExpr{
		Expr: left,
	}, nil
}

func (p *Parser) parseExpr(precedence BindingPower) (ast.Expr, error) {
	currKind := p.currentTokenKind()
	nudHandler, exists := p.lookupTable.GetNUDHandlerIfExists(currKind)
	if !exists {
		return nil, NewError(fmt.Sprintf("Unhandled token %s", currKind))
	}

	left, err := nudHandler()
	if err != nil {
		return nil, err
	}

	for p.lookupTable.GetBpHandler(p.currentTokenKind()) > precedence {
		currKind = p.currentTokenKind()
		ledHandler, existsLed := p.lookupTable.GetLedHandlerIfExists(currKind)
		if !existsLed {
			return nil, NewError(fmt.Sprintf("Unhandled token %s", currKind))
		}

		left, err = ledHandler(left, p.lookupTable.GetBpHandler(currKind))
		if err != nil {
			return nil, err
		}
	}

	return left, nil
}
