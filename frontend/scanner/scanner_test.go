package scanner

import (
	"strings"
	"testing"

	"gotlin/frontend/token"
)

func TestScanner_SingleScan(t *testing.T) {
	input :=
		`var a: String = "testing"
		val b: Int = 42
	`
	tests := []struct {
		expectedType token.Kind
		expectedLit  string
	}{
		{token.VAR, "var"},
		{token.IDENTIFIER, "a"},
		{token.COLON, ":"},
		{token.STRING, "String"},
		{token.ASSIGN, "="},
		{token.STRINGLIT, "testing"},
		{token.VAL, "val"},
		{token.IDENTIFIER, "b"},
		{token.COLON, ":"},
		{token.INT, "Int"},
		{token.ASSIGN, "="},
		{token.INTLIT, "42"},
		{token.NEWLINE, "<NL>"},
		{token.EOF, "EOF"},
	}
	reader := strings.NewReader(input)
	scanner := NewScanner(reader)
	toks := scanner.ScanTokens()

	for i, tt := range tests {
		tok := toks[i]
		if tok.Kind != tt.expectedType {
			t.Fatalf("tests[%d] - tokentype wrong. expected=%q, got=%q",
				i, tt.expectedType, tok.Kind)
		}
		if tok.Spelling != tt.expectedLit {
			t.Fatalf("tests[%d] - literal wrong. expected=%q, got=%q",
				i, tt.expectedLit, tok.Spelling)
		}
	}
}
