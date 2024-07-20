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
		AddNudHandler(token.INTLIT, p.parsePrimaryExpr).
		AddNudHandler(token.STRINGLIT, p.parsePrimaryExpr).
		AddNudHandler(token.BOOLEANLIT, p.parsePrimaryExpr).
		AddNudHandler(token.IDENTIFIER, p.parsePrimaryExpr).
		AddNudHandler(token.FUNCTION, p.parseFunctionLiteral).

		//Logical
		AddLedHandler(token.AND, Logical, p.parseBinaryExpr).
		AddLedHandler(token.OR, Logical, p.parseBinaryExpr).
		AddLedHandler(token.ELVIS, Logical, p.parseBinaryExpr).
		AddLedHandler(token.BANG_BANG, Call, p.parseNotNullExpr).

		//Relational
		AddLedHandler(token.LT, Relational, p.parseBinaryExpr).
		AddLedHandler(token.LTE, Relational, p.parseBinaryExpr).
		AddLedHandler(token.GT, Relational, p.parseBinaryExpr).
		AddLedHandler(token.GTE, Relational, p.parseBinaryExpr).
		AddLedHandler(token.EQ_EQ, Relational, p.parseBinaryExpr).
		AddLedHandler(token.NOT_EQ, Relational, p.parseBinaryExpr).

		//Additive
		AddLedHandler(token.PLUS, Additive, p.parseBinaryExpr).
		AddLedHandler(token.DASH, Additive, p.parseBinaryExpr).
		AddLedHandler(token.SLASH, Multiplicative, p.parseBinaryExpr).
		AddLedHandler(token.STAR, Multiplicative, p.parseBinaryExpr).

		// Call
		AddLedHandler(token.OPEN_PAREN, Call, p.parseCallExpr).

		//Unary
		AddNudHandler(token.DASH, p.parseUnaryExpr).
		AddNudHandler(token.PLUS, p.parseUnaryExpr).
		AddNudHandler(token.NOT, p.parseUnaryExpr).
		AddNudHandler(token.OPEN_PAREN, p.parseGroupingExpr).

		//Statements
		AddStmtHandler(token.VAR, p.parseVariableDeclStmt).
		AddStmtHandler(token.VAL, p.parseVariableDeclStmt).
		AddStmtHandler(token.IDENTIFIER, p.parseAssignmentStmt).
		AddStmtHandler(token.CLASS, p.parseClassDeclStmt).

		// Types
		AddTypeNudHandler(token.IDENTIFIER, p.parseUserType).
		AddTypeNudHandler(token.OPEN_BRACKET, p.parseArrayType). // TODO Check array syntax
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
