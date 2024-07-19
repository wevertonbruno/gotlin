package executor

import (
	"bufio"
	"os"

	"gotlin/backend/virtualmachine"
)

type Repl struct {
	vm     *virtualmachine.VM
	reader *bufio.Reader
}

func NewRepl(vm *virtualmachine.VM) *Repl {
	return &Repl{
		vm:     vm,
		reader: bufio.NewReader(os.Stdin),
	}
}

func (r *Repl) Execute() {
	for {
		print(">> ")
		r.vm.Interpret(r.reader)
	}
}
