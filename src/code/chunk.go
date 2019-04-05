package code

import (
	"github.com/adamjedlicka/lang/src/val"
)

type Chunk struct {
	code      []uint8
	lines     []int
	constants *val.ValueArray
}

func NewChunk() *Chunk {
	c := new(Chunk)
	c.code = make([]uint8, 0)
	c.lines = make([]int, 0)
	c.constants = val.NewValueArray()

	return c
}

func (c *Chunk) AddConstant(value val.Value) uint8 {
	c.constants.Write(value)

	return c.constants.Len() - 1
}

func (c *Chunk) Write(instruction OpCode, line int) {
	c.WriteRaw(uint8(instruction), line)
}

func (c *Chunk) WriteRaw(data uint8, line int) {
	c.code = append(c.code, data)
	c.lines = append(c.lines, line)
}

func (c *Chunk) Get(offset int) OpCode {
	return OpCode(c.GetRaw(offset))
}

func (c *Chunk) GetRaw(offset int) uint8 {
	return c.code[offset]
}

func (c *Chunk) GetConstant(offset uint8) val.Value {
	return c.constants.GetValue(offset)
}

func (c *Chunk) GetLine(offset int) int {
	return c.lines[offset]
}

func (c *Chunk) Len() int {
	return len(c.code)
}
