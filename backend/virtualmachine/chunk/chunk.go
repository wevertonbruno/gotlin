package chunk

import (
	"gotlin/backend/virtualmachine/chunk/instruction"
)

type Value float64

type Instruction struct {
	Op   byte
	Line int
}

type Chunk struct {
	Instructions []Instruction
	Constants    []Value
}

func New() *Chunk {
	return &Chunk{
		Instructions: make([]Instruction, 0),
		Constants:    make([]Value, 0),
	}
}

func (c *Chunk) Size() int {
	return len(c.Instructions)
}

func (c *Chunk) Write(op uint8, line int) {
	c.Instructions = append(c.Instructions, Instruction{
		Op:   op,
		Line: line,
	})
}

func (c *Chunk) WriteConstant(value Value, line int) {
	index := c.AddConstant(value)
	if index < 256 {
		c.Write(instruction.OpConstant, line)
		c.Write(uint8(index), line)
	} else {
		c.Write(instruction.OpConstantLong, line)
		c.Write(byte(index&0xff), line)
		c.Write(byte((index>>8)&0xff), line)
		c.Write(byte((index>>16)&0xff), line)
	}
}

func (c *Chunk) AddConstant(value Value) int {
	c.Constants = append(c.Constants, value)
	return len(c.Constants) - 1
}

func (c *Chunk) Destroy() {
	c.Instructions = nil
	c.Constants = nil
}
