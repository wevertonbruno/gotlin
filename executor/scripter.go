package executor

import (
	"bufio"
	"fmt"
	"os"

	"gotlin/backend/virtualmachine"
)

type FileReader interface {
	ReadFile(filename string) ([]byte, error)
}

type Scripter struct {
	vm        *virtualmachine.VM
	reader    *bufio.Reader
	closeFile func()
}

func NewScripter(vm *virtualmachine.VM, path string) *Scripter {
	file, err := os.Open(path)
	if err != nil {
		fmt.Printf("Error reading script file %s: %s", path, err)
		os.Exit(1)
	}

	return &Scripter{
		vm:     vm,
		reader: bufio.NewReader(file),
		closeFile: func() {
			_ = file.Close()
		},
	}
}

func (r *Scripter) Execute() {
	defer r.closeFile()
	result := r.vm.Interpret(r.reader)
	switch result {
	case virtualmachine.ResultCompileError:
		os.Exit(65)
	case virtualmachine.ResultRuntimeError:
		os.Exit(70)
	default:
		return
	}
}
