package chunk

import (
	"fmt"

	"gotlin/backend/virtualmachine/chunk/instruction"
)

func (c *Chunk) Disassemble(name string) {
	fmt.Printf("==== %s ====\n", name)
	for offset := 0; offset < c.Size(); {
		offset = c.DisassembleInstruction(offset)
	}
}

func (c *Chunk) DisassembleInstruction(offset int) int {
	fmt.Printf("%04d ", offset)

	instr := c.Instructions[offset]
	if offset > 0 && instr.Line == c.Instructions[offset-1].Line {
		fmt.Printf("   | ")
	} else {
		fmt.Printf("%4d ", instr.Line)
	}

	switch instr.Op {
	case instruction.OpReturn:
		return simpleInstruction("OP_RETURN", offset)
	case instruction.OpConstant:
		return constantInstruction("OP_CONSTANT", c, offset)
	case instruction.OpConstantLong:
		return constantLongInstruction("OP_CONSTANT_LONG", c, offset)
	case instruction.OpNegate:
		return simpleInstruction("OP_NEGATE", offset)
	case instruction.OpAdd:
		return simpleInstruction("OP_ADD", offset)
	case instruction.OpSubtract:
		return simpleInstruction("OP_SUBTRACT", offset)
	case instruction.OpMultiply:
		return simpleInstruction("OP_MULTIPLY", offset)
	case instruction.OpDivide:
		return simpleInstruction("OP_DIVIDE", offset)
	default:
		fmt.Printf("Unknown opcode %v\n", instr.Op)
		return offset + 1
	}
}

func simpleInstruction(name string, offset int) int {
	println(name)
	return offset + 1
}

func constantInstruction(name string, chunk *Chunk, offset int) int {
	constant := chunk.Instructions[offset+1].Op
	fmt.Printf("%-16s %4d '", name, constant)
	fmt.Printf("%g", chunk.Constants[constant])
	fmt.Println("'")
	return offset + 2
}

func constantLongInstruction(name string, chunk *Chunk, offset int) int {
	constant := int(chunk.Instructions[offset+1].Op) | int(chunk.Instructions[offset+2].Op)<<8 | int(chunk.Instructions[offset+3].Op)<<16
	fmt.Printf("%-16s %4d '", name, constant)
	fmt.Printf("%g", chunk.Constants[constant])
	fmt.Println("'")
	return offset + 4
}
