package token

import "fmt"

type Kind string

const (
	IDENTIFIER Kind = "<identifier>"
	NEWLINE    Kind = "<NL>"
	IntLit     Kind = "<integer>"
	DoubleLit  Kind = "<double>"
	CharLit    Kind = "<char>"
	StringLit  Kind = "<string>"
	BooleanLit Kind = "<boolean>"

	IF       Kind = "if"
	ELSE     Kind = "else"
	FUNCTION Kind = "fun"
	CLASS    Kind = "class"
	WHILE    Kind = "while"
	VAR      Kind = "var"
	VAL      Kind = "val"
	PRINT    Kind = "print"
	RETURN   Kind = "return"
	INT      Kind = "Int"
	TRUE     Kind = "true"
	FALSE    Kind = "false"

	OpPlus     Kind = "+"
	OpMinus    Kind = "-"
	OpMulti    Kind = "*"
	OpDivide   Kind = "/"
	OpEq       Kind = "=="
	OpNotEq    Kind = "!="
	OpLt       Kind = "<"
	OpLte      Kind = "<="
	OpGt       Kind = ">"
	OpGte      Kind = ">="
	OpSafeCall Kind = "?."
	OpNotNull  Kind = "!!"
	NOT        Kind = "!"
	AND        Kind = "&&"
	OR         Kind = "||"

	ASSIGN       Kind = "="
	PLUS_ASSIGN  Kind = "+="
	MINUS_ASSIGN Kind = "-="
	SEMICOLON    Kind = ";"
	COLON        Kind = ":"
	DOT          Kind = "."
	QUESTION     Kind = "?"
	LParen       Kind = "("
	RParen       Kind = ")"
	LBrace       Kind = "{"
	RBrace       Kind = "}"
	LBracket     Kind = "["
	RBracket     Kind = "]"
	ELVIS        Kind = "?:"

	EOF   Kind = "EOF"
	ERROR Kind = "ERROR"
)

var (
	reservedKeywords = map[string]Kind{
		string(IF):       IF,
		string(ELSE):     ELSE,
		string(FUNCTION): FUNCTION,
		string(WHILE):    WHILE,
		string(VAR):      VAR,
		string(VAL):      VAL,
		string(PRINT):    PRINT,
		string(RETURN):   RETURN,
		string(CLASS):    CLASS,

		string(TRUE):  BooleanLit,
		string(FALSE): BooleanLit,
	}
)

type Pos struct {
	Line uint
	Col  uint
}

type Token struct {
	Kind     Kind
	Spelling string
	Position Pos
}

func NewToken(kind Kind, line uint, col uint) Token {
	return Token{
		Kind:     kind,
		Spelling: string(kind),
		Position: Pos{Line: line, Col: col},
	}
}

func NewTokenLiteral(kind Kind, spell string, line uint, col uint) Token {
	if kind == IDENTIFIER {
		if val, exists := reservedKeywords[spell]; exists {
			kind = val
		}
	}
	return Token{
		Kind:     kind,
		Spelling: spell,
		Position: Pos{Line: line, Col: col},
	}
}

func (t Token) String() string {
	return t.Spelling
}

func (k Kind) String() string {
	return string(k)
}

func (k Pos) String() string {
	return fmt.Sprintf("[%d, %d]", k.Line, k.Col)
}
