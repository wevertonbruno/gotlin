package main

import (
	"fmt"
	"os"
	"time"

	"github.com/sanity-io/litter"
	"gotlin/frontend/parser"
	"gotlin/frontend/scanner"
)

type Executor interface {
	Execute()
}

func main() {
	file, err := os.Open("example.gt")
	if err != nil {
		panic(err)
	}

	start := time.Now()
	s := scanner.NewScanner(file)
	p := parser.New(s)
	program := p.Parse()
	duration := time.Since(start)
	litter.Dump(program)
	fmt.Printf("\n\nExecution time: %s\n", duration)
}
