package parser

import (
	"fmt"

	"gotlin/frontend/ast"
	"gotlin/frontend/token"
)

type Scanner interface {
	ScanTokens() []token.Token
}

type Parser struct {
	scanner     Scanner
	lookupTable *LookupTable
	tokens      []token.Token
	cursor      int
}

func New(scanner Scanner) *Parser {
	p := &Parser{scanner: scanner, cursor: 0}
	p.lookupTable = NewLookupTable().
		//Literals
		AddNudHandler(token.IntLit, p.parsePrimaryExpr).
		AddNudHandler(token.StringLit, p.parsePrimaryExpr).
		AddNudHandler(token.BooleanLit, p.parsePrimaryExpr).
		AddNudHandler(token.IDENTIFIER, p.parsePrimaryExpr).

		//Logical
		AddLedHandler(token.AND, Logical, p.parseBinaryExpr).
		AddLedHandler(token.OR, Logical, p.parseBinaryExpr).
		AddLedHandler(token.ELVIS, Logical, p.parseBinaryExpr).
		AddLedHandler(token.OpNotNull, Call, p.parseNotNullExpr).

		//Relational
		AddLedHandler(token.OpLt, Relational, p.parseBinaryExpr).
		AddLedHandler(token.OpLte, Relational, p.parseBinaryExpr).
		AddLedHandler(token.OpGt, Relational, p.parseBinaryExpr).
		AddLedHandler(token.OpGte, Relational, p.parseBinaryExpr).
		AddLedHandler(token.OpEq, Relational, p.parseBinaryExpr).
		AddLedHandler(token.OpNotEq, Relational, p.parseBinaryExpr).

		//Additive
		AddLedHandler(token.OpPlus, Additive, p.parseBinaryExpr).
		AddLedHandler(token.OpMinus, Additive, p.parseBinaryExpr).
		AddLedHandler(token.OpDivide, Multiplicative, p.parseBinaryExpr).
		AddLedHandler(token.OpMulti, Multiplicative, p.parseBinaryExpr).

		//Unary
		AddNudHandler(token.OpMinus, p.parseUnaryExpr).
		AddNudHandler(token.OpPlus, p.parseUnaryExpr).
		AddNudHandler(token.NOT, p.parseUnaryExpr).
		AddNudHandler(token.LParen, p.parseGroupingExpr).

		//Statements
		AddStmtHandler(token.VAR, p.parseVariableDeclStmt).
		AddStmtHandler(token.VAL, p.parseVariableDeclStmt).
		AddStmtHandler(token.IDENTIFIER, p.parseAssignmentStmt).

		// Types
		AddTypeNudHandler(token.IDENTIFIER, p.parseUserType).
		AddTypeNudHandler(token.LBracket, p.parseArrayType). // TODO Check array syntax
		AddTypeLedHandler(token.QUESTION, Call, p.parseNullableType)

	return p
}

func (p *Parser) Parse() *ast.Program {
	p.tokens = p.scanner.ScanTokens()
	p.cursor = 0

	program := &ast.Program{
		Statements: []ast.Stmt{},
	}

	p.skipNewLines()
	for p.hasTokens() {
		stmt, err := p.parseStmt()
		if err != nil {
			panic(err) // TODO handle error synchronize parser
		}
		program.Statements = append(program.Statements, stmt)
		p.skipNewLines()
	}

	return program
}

func (p *Parser) skipNewLines() {
	for p.currentTokenKind() == token.NEWLINE {
		p.advance()
	}
}

func (p *Parser) currentToken() token.Token {
	return p.tokens[p.cursor]
}

func (p *Parser) currentTokenKind() token.Kind {
	return p.currentToken().Kind
}

func (p *Parser) advance() token.Token {
	tk := p.tokens[p.cursor]
	p.cursor++
	return tk
}

func (p *Parser) hasTokens() bool {
	return p.cursor < len(p.tokens) && p.currentTokenKind() != token.EOF
}

func (p *Parser) expected(kinds ...token.Kind) (token.Token, error) {
	for _, k := range kinds {
		if p.currentTokenKind() == k {
			return p.advance(), nil
		}
	}

	// TODO Improve error message
	return token.Token{}, NewError(fmt.Sprintf("expected one of %v; got %s", kinds, p.currentTokenKind()))
}
