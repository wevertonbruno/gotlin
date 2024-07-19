package parser

import (
	"fmt"

	"gotlin/frontend/ast"
	"gotlin/frontend/token"
)

func (p *Parser) parseStmt() (ast.Stmt, error) {
	kind := p.currentTokenKind()
	stmtHandler, exists := p.lookupTable.GetStmtHandlerIfExists(kind)
	if exists {
		return stmtHandler()
	}
	expr, err := p.parseExpr(Default)
	if err != nil {
		return nil, err
	}

	_, err = p.expected(token.SEMICOLON, token.NEWLINE)
	if err != nil {
		return nil, err
	}

	return &ast.ExprStmt{
		Expr: expr,
	}, nil
}

func (p *Parser) parseVariableDeclStmt() (ast.Stmt, error) {
	immutable := p.advance().Kind == token.VAL
	identifier, err := p.expected(token.IDENTIFIER)
	if err != nil {
		return nil, err
	}

	//Parse Type
	var explicitType ast.Type
	if p.currentTokenKind() == token.COLON {
		p.advance()
		explicitType, err = p.parseType(Default)
		if err != nil {
			return nil, err
		}
	}

	var assignedValue ast.Expr
	if p.currentTokenKind() != token.SEMICOLON && p.currentTokenKind() != token.NEWLINE {
		_, err = p.expected(token.ASSIGN)
		if err != nil {
			return nil, err
		}

		assignedValue, err = p.parseExpr(Assignment)
		if err != nil {
			return nil, err
		}
	} else if explicitType == nil {
		return nil, NewError("This variable must either have a type annotation or be initialized")
	}

	_, err = p.expected(token.SEMICOLON, token.NEWLINE)
	if err != nil {
		return nil, err
	}

	return &ast.VariableDecl{
		Name:      identifier,
		Type:      explicitType,
		Value:     assignedValue,
		Immutable: immutable,
	}, nil
}

func (p *Parser) parseAssignmentStmt() (ast.Stmt, error) {
	assigne, err := p.parseExpr(Default)
	if err != nil {
		return nil, err
	}

	if p.currentTokenKind() == token.ASSIGN {
		p.advance()

		if _, ok := assigne.(*ast.IdentifierExpr); !ok {
			return nil, NewError(fmt.Sprintf("Variable expected, got %v", assigne))
		}

		right, err2 := p.parseExpr(Default)
		if err2 != nil {
			return nil, err2
		}
		_, err2 = p.expected(token.SEMICOLON, token.NEWLINE)
		if err2 != nil {
			return nil, err2
		}

		return &ast.AssignStmt{Assigne: assigne, Value: right}, nil
	}
	_, err = p.expected(token.SEMICOLON, token.NEWLINE)
	if err != nil {
		return nil, err
	}

	return &ast.ExprStmt{
		Expr: assigne,
	}, nil
}

func (p *Parser) parseClassDeclStmt() (ast.Stmt, error) {
	_, err := p.expected(token.CLASS)
	if err != nil {
		return nil, err
	}

	className, err := p.expected(token.IDENTIFIER)
	if err != nil {
		return nil, err
	}

}
