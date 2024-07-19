package virtualmachine

import (
	"bufio"
	"fmt"

	"gotlin/backend/virtualmachine/chunk"
	"gotlin/backend/virtualmachine/chunk/instruction"
)

type Result byte

const (
	ResultOk Result = iota
	ResultCompileError
	ResultRuntimeError

	StackMax = 256
)

type Compiler interface {
	Compile(reader *bufio.Reader)
}

type VM struct {
	chunk     *chunk.Chunk
	stack     stack
	ip        int
	debugMode bool
	compiler  Compiler
}

func New(compiler Compiler) *VM {
	return &VM{
		compiler: compiler,
	}
}

func (vm *VM) Interpret(source *bufio.Reader) Result {
	vm.compiler.Compile(source)
	return ResultOk
}

func (vm *VM) run() Result {
	for {
		// TODO - Debug
		vm.chunk.DisassembleInstruction(vm.ip)

		instr := vm.readByte()
		switch instr {
		case instruction.OpReturn:
			fmt.Printf("%g\n", vm.stack.pop())
			return ResultOk
		case instruction.OpConstant:
			constant := vm.readConstant()
			vm.stack.push(constant)
			break
		case instruction.OpConstantLong:
			constant := vm.readConstantLong()
			vm.stack.push(constant)
			break
		case instruction.OpNegate:
			vm.stack.push(-vm.stack.pop())
			break
		case instruction.OpAdd:
			r, l := vm.stack.pop(), vm.stack.pop()
			vm.stack.push(l + r)
			break
		case instruction.OpSubtract:
			r, l := vm.stack.pop(), vm.stack.pop()
			vm.stack.push(l - r)
			break
		case instruction.OpMultiply:
			r, l := vm.stack.pop(), vm.stack.pop()
			vm.stack.push(l * r)
			break
		case instruction.OpDivide:
			r, l := vm.stack.pop(), vm.stack.pop()
			vm.stack.push(l / r)
			break
		}
	}
}

func (vm *VM) readConstantLong() chunk.Value {
	index := int(vm.readByte())
	index |= int(vm.readByte()) << 8
	index |= int(vm.readByte()) << 16
	return vm.chunk.Constants[index]
}

func (vm *VM) readConstant() chunk.Value {
	return vm.chunk.Constants[vm.readByte()]
}

func (vm *VM) readByte() uint8 {
	b := vm.chunk.Instructions[vm.ip].Op
	vm.ip++
	return b
}
