package token

import "fmt"

type Kind string

const (
	IDENTIFIER Kind = "<identifier>"
	NEWLINE    Kind = "<NL>"
	INTLIT     Kind = "<integer>"
	DOUBLELIT  Kind = "<double>"
	CHARLIT    Kind = "<char>"
	STRINGLIT  Kind = "<string>"
	BOOLEANLIT Kind = "<boolean>"

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

	PLUS      Kind = "+"
	DASH      Kind = "-"
	STAR      Kind = "*"
	SLASH     Kind = "/"
	EQ_EQ     Kind = "=="
	NOT_EQ    Kind = "!="
	LT        Kind = "<"
	LTE       Kind = "<="
	GT        Kind = ">"
	GTE       Kind = ">="
	QUEST_DOT Kind = "?."
	BANG_BANG Kind = "!!"
	NOT       Kind = "!"
	AND       Kind = "&&"
	OR        Kind = "||"

	ASSIGN        Kind = "="
	PLUS_ASSIGN   Kind = "+="
	MINUS_ASSIGN  Kind = "-="
	SEMICOLON     Kind = ";"
	COLON         Kind = ":"
	DOT           Kind = "."
	COMMA         Kind = ","
	QUESTION      Kind = "?"
	OPEN_PAREN    Kind = "("
	CLOSE_PAREN   Kind = ")"
	OPEN_BRACE    Kind = "{"
	CLOSE_BRACE   Kind = "}"
	OPEN_BRACKET  Kind = "["
	CLOSE_BRACKET Kind = "]"
	ELVIS         Kind = "?:"

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

		string(TRUE):  BOOLEANLIT,
		string(FALSE): BOOLEANLIT,
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
