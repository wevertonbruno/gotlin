package parser

import (
	"strings"
	"testing"

	"gotlin/frontend/scanner"
)

func TestParser_ParseProgram(t *testing.T) {
	tests := []struct {
		input    string
		stmtSize int
	}{
		{"4 + 4", 1},
		{"4 + 4; 2 + 1", 2},
		{"val v1: Int = 1; var v2: String", 2},
		{"var v1: Int = 1; var v2: Int = 1; v1 = v2 = 2", 3},
	}

	for _, test := range tests {
		s := scanner.NewScanner(strings.NewReader(test.input))
		p := New(s)
		program := p.Parse()
		if len(program.Statements) != test.stmtSize {
			t.Errorf("program.Statements is %d, want %d", len(program.Statements), test.stmtSize)
		}
	}
}
