package parser

import (
	"fmt"

	"gotlin/frontend/ast"
	"gotlin/frontend/token"
)

func (p *Parser) parseStmt() (ast.Stmt, error) {
	kind := p.currentTokenKind()
	stmtHandler, exists := p.lookupTable.GetStmtHandlerIfExists(kind)
	var stmt ast.Stmt
	if exists {
		var err error
		stmt, err = stmtHandler()
		if err != nil {
			return nil, err
		}
	} else {
		expr, err := p.parseExpr(Default)
		if err != nil {
			return nil, err
		}

		stmt = &ast.ExprStmt{
			Expr: expr,
		}
	}

	// TODO Remove this expectation
	_, err := p.expected(token.SEMICOLON, token.NEWLINE)
	if err != nil {
		return nil, err
	}

	return stmt, nil
}

func (p *Parser) parseVariableDeclStmt() (ast.Stmt, error) {
	readOnly := p.advance().Kind == token.VAL
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
		Name:     identifier,
		Type:     explicitType,
		Value:    assignedValue,
		ReadOnly: readOnly,
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

	// Parses primary constructor
	var primaryConstructor *ast.ClassPrimaryConstructor
	if p.currentTokenKind() == token.OPEN_PAREN {
		p.advance()
		var parameters []ast.ClassParam
		for p.hasTokens() && p.currentTokenKind() != token.CLOSE_PAREN {
			parameter, err2 := p.expected(token.IDENTIFIER)
			if err2 != nil {
				return nil, err2
			}

			_, err2 = p.expected(token.COLON)
			if err2 != nil {
				return nil, err2
			}

			parameterType, err2 := p.parseType(Default)
			if err2 != nil {
				return nil, err2
			}

			var defaultValue ast.Expr
			if p.currentTokenKind() == token.ASSIGN {
				p.advance()
				defaultValue, err = p.parseExpr(Default)
				if err != nil {
					return nil, err
				}
			}

			parameters = append(parameters, ast.ClassParam{
				Name:         parameter.Spelling,
				Type:         parameterType,
				ReadOnly:     true,
				DefaultValue: defaultValue,
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

		if len(parameters) != 0 {
			primaryConstructor = &ast.ClassPrimaryConstructor{
				Parameters: parameters,
			}
		}
	}

	return &ast.ClassDeclStmt{
		Name:               className,
		PrimaryConstructor: primaryConstructor,
	}, nil
}
