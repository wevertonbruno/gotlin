package compiler

import (
	"bufio"
	"fmt"

	"gotlin/frontend/token"
)

type Scanner interface {
	ScanToken() token.Token
}

type Compiler struct {
	scanner Scanner
}

func New(scanner Scanner) *Compiler {
	return &Compiler{
		scanner: scanner,
	}
}

func (c *Compiler) Compile(reader *bufio.Reader) {
	var line uint = 0
	for {
		t := c.scanner.ScanToken()
		if t.Position[0] != line {
			fmt.Printf("%4d ", t.line)
			line = t.line
		} else {
			fmt.Printf("   | ")
		}
		fmt.Printf("%2d '%.*s'\n", t.type, t.length, t.start);
		if t.IsEOF() {
			break
		}
	}
}
