package parser

import (
	"gotlin/frontend/ast"
	"gotlin/frontend/token"
)

type BindingPower byte

const (
	Default BindingPower = iota
	Comma
	Assignment
	Logical
	Relational
	Additive
	Multiplicative
	Unary
	Call
	Member
	Primary
)

type StmtHandler func() (ast.Stmt, error)
type NudHandler func() (ast.Expr, error)
type LedHandler func(left ast.Expr, precedence BindingPower) (ast.Expr, error)
type TypeNudHandler func() (ast.Type, error)
type TypeLedHandler func(left ast.Type, precedence BindingPower) (ast.Type, error)

type StmtLookup map[token.Kind]StmtHandler
type NudLookup map[token.Kind]NudHandler
type LedLookup map[token.Kind]LedHandler

type TypeNudLookup map[token.Kind]TypeNudHandler
type TypeLedLookup map[token.Kind]TypeLedHandler

type BpLookup map[token.Kind]BindingPower

type LookupTable struct {
	stmtTable    StmtLookup
	nudTable     NudLookup
	ledTable     LedLookup
	typeNudTable TypeNudLookup
	typeLedTable TypeLedLookup
	bpTable      BpLookup
	typeBpTable  BpLookup
}

func NewLookupTable() *LookupTable {
	return &LookupTable{
		stmtTable:    make(StmtLookup),
		nudTable:     make(NudLookup),
		ledTable:     make(LedLookup),
		typeNudTable: make(TypeNudLookup),
		typeLedTable: make(TypeLedLookup),
		bpTable:      make(BpLookup),
		typeBpTable:  make(BpLookup),
	}
}

func (r *LookupTable) AddStmtHandler(token token.Kind, handler StmtHandler) *LookupTable {
	r.bpTable[token] = Default
	r.stmtTable[token] = handler
	return r
}

func (r *LookupTable) AddLedHandler(token token.Kind, bp BindingPower, handler LedHandler) *LookupTable {
	r.ledTable[token] = handler
	r.bpTable[token] = bp
	return r
}

func (r *LookupTable) AddNudHandler(token token.Kind, handler NudHandler) *LookupTable {
	r.nudTable[token] = handler
	return r
}

func (r *LookupTable) AddTypeLedHandler(token token.Kind, bp BindingPower, handler TypeLedHandler) *LookupTable {
	r.typeLedTable[token] = handler
	r.typeBpTable[token] = bp
	return r
}

func (r *LookupTable) AddTypeNudHandler(token token.Kind, handler TypeNudHandler) *LookupTable {
	r.typeNudTable[token] = handler
	return r
}

func (r *LookupTable) GetNUDHandlerIfExists(kind token.Kind) (NudHandler, bool) {
	v, ok := r.nudTable[kind]
	return v, ok
}

func (r *LookupTable) GetLedHandlerIfExists(kind token.Kind) (LedHandler, bool) {
	v, ok := r.ledTable[kind]
	return v, ok
}

func (r *LookupTable) GetTypeNUDHandlerIfExists(kind token.Kind) (TypeNudHandler, bool) {
	v, ok := r.typeNudTable[kind]
	return v, ok
}

func (r *LookupTable) GetTypeLedHandlerIfExists(kind token.Kind) (TypeLedHandler, bool) {
	v, ok := r.typeLedTable[kind]
	return v, ok
}

func (r *LookupTable) GetBpHandler(kind token.Kind) BindingPower {
	return r.bpTable[kind]
}

func (r *LookupTable) GetTypeBpHandler(kind token.Kind) BindingPower {
	return r.typeBpTable[kind]
}

func (r *LookupTable) GetStmtHandlerIfExists(kind token.Kind) (StmtHandler, bool) {
	v, ok := r.stmtTable[kind]
	return v, ok
}
