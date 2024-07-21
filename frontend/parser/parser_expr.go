package parser

import (
	"fmt"
	"strconv"

	"gotlin/frontend/ast"
	"gotlin/frontend/token"
)

func (p *Parser) parsePrimaryExpr() (ast.Expr, error) {
	switch p.currentTokenKind() {
	case token.INTLIT:
		value, err := strconv.ParseInt(p.advance().Spelling, 10, 64)
		return &ast.IntLiteral{
			Value: value,
		}, err
	case token.STRINGLIT:
		return &ast.StringLiteral{
			Value: p.advance().Spelling,
		}, nil
	case token.BOOLEANLIT:
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

	_, err = p.expected(token.CLOSE_PAREN)
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

func (p *Parser) parseCallExpr(left ast.Expr, precedence BindingPower) (ast.Expr, error) {
	switch left.(type) {
	case *ast.IdentifierExpr, *ast.CallExpr:
		break
	default:
		return nil, NewError(fmt.Sprintf("Expression '%s' cannot be invoked as a function.", left))
	}

	var args []ast.Expr
	p.advance()
	for p.hasTokens() && p.currentTokenKind() != token.CLOSE_PAREN {
		arg, err := p.parseExpr(Default)
		if err != nil {
			return nil, err
		}

		args = append(args, arg)
		if p.currentTokenKind() != token.CLOSE_PAREN {
			_, err = p.expected(token.COMMA)
			if err != nil {
				return nil, err
			}
		}
	}

	_, err := p.expected(token.CLOSE_PAREN)
	if err != nil {
		return nil, err
	}

	return &ast.CallExpr{
		Callee: left,
		Args:   args,
	}, nil
}

func (p *Parser) parseFunctionLiteral() (ast.Expr, error) {
	p.advance()
	_, err := p.expected(token.OPEN_PAREN)
	if err != nil {
		return nil, err
	}

	var funcParameters []*ast.ParameterWithOptionalType
	for p.hasTokens() && p.currentTokenKind() != token.CLOSE_PAREN {
		funcParameter, err2 := p.expected(token.IDENTIFIER)
		if err2 != nil {
			return nil, err2
		}

		_, err2 = p.expected(token.COLON)
		if err2 != nil {
			return nil, err2
		}

		funcParameterType, err2 := p.parseType(Default)
		if err2 != nil {
			return nil, err2
		}

		funcParameters = append(funcParameters, &ast.ParameterWithOptionalType{
			Name: funcParameter.Spelling,
			Type: funcParameterType,
		})

		if p.currentTokenKind() != token.CLOSE_PAREN {
			_, err = p.expected(token.COMMA)
			if err != nil {
				return nil, err
			}
		}
	}
	_, err = p.expected(token.CLOSE_PAREN)
	if err != nil {
		return nil, err
	}

	var funcType ast.Type
	if p.currentTokenKind() == token.COLON {
		p.advance()
		funcType, err = p.parseType(Default)
		if err != nil {
			return nil, err
		}
	}

	// fun () = expr
	var bodyExpr ast.Expr
	if p.currentTokenKind() == token.ASSIGN {
		p.advance()
		bodyExpr, err = p.parseExpr(Default)
		if err != nil {
			return nil, err
		}
	}

	// fun () {}
	var block []ast.Stmt
	if p.currentTokenKind() == token.OPEN_BRACE {
		p.advance()
		for p.hasTokens() && p.currentTokenKind() != token.CLOSE_BRACE {
			stmt, stmtErr := p.parseStmt()
			if stmtErr != nil {
				return nil, stmtErr
			}

			block = append(block, stmt)
			if p.currentTokenKind() != token.CLOSE_BRACE {
				_, err = p.expected(token.COMMA)
				if err != nil {
					return nil, err
				}
			}
		}
	}

	return &ast.FunctionLiteral{
		Type:       funcType,
		Parameters: funcParameters,
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
