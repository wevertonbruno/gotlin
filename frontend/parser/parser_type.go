package parser

import (
	"fmt"

	"gotlin/frontend/ast"
	"gotlin/frontend/token"
)

func (p *Parser) parseUserType() (ast.Type, error) {
	id, err := p.expected(token.IDENTIFIER)
	if err != nil {
		return nil, err
	}
	return &ast.TypeName{
		Name: id.Spelling,
	}, nil
}

func (p *Parser) parseArrayType() (ast.Type, error) {
	p.advance()
	_, err := p.expected(token.RBracket)
	if err != nil {
		return nil, err
	}

	underlying, err := p.parseType(Default)
	if err != nil {
		return nil, err
	}

	return &ast.ArrayType{
		Underlying: underlying,
	}, nil
}

func (p *Parser) parseNullableType(left ast.Type, precedence BindingPower) (ast.Type, error) {
	_, err := p.expected(token.QUESTION)
	if err != nil {
		return nil, err
	}

	return &ast.NullableType{
		Type: left,
	}, nil
}

func (p *Parser) parseType(precedence BindingPower) (ast.Type, error) {
	currKind := p.currentTokenKind()
	nudHandler, exists := p.lookupTable.GetTypeNUDHandlerIfExists(currKind)
	if !exists {
		return nil, NewError(fmt.Sprintf("Unhandled token %s", currKind))
	}

	left, err := nudHandler()
	if err != nil {
		return nil, err
	}

	for p.lookupTable.GetTypeBpHandler(p.currentTokenKind()) > precedence {
		currKind = p.currentTokenKind()
		ledHandler, existsLed := p.lookupTable.GetTypeLedHandlerIfExists(currKind)
		if !existsLed {
			return nil, NewError(fmt.Sprintf("Unhandled token %s", currKind))
		}

		left, err = ledHandler(left, p.lookupTable.GetTypeBpHandler(currKind))
		if err != nil {
			return nil, err
		}
	}

	return left, nil
}
