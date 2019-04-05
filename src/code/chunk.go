package code

import (
	"github.com/adamjedlicka/lang/src/val"
)

type Chunk struct {
	data      []int
	lines     []int
	constants *val.ValueArray
}

func NewChunk() *Chunk {
	c := new(Chunk)
	c.data = make([]int, 0)
	c.lines = make([]int, 0)
	c.constants = val.NewValueArray()

	return c
}

func (c *Chunk) AddConstant(value val.Value) int {
	c.constants.Write(value)

	return c.constants.Len() - 1
}

func (c *Chunk) Write(instruction OpCode, line int) {
	c.WriteRaw(int(instruction), line)
}

func (c *Chunk) WriteRaw(data int, line int) {
	c.data = append(c.data, data)
	c.lines = append(c.lines, line)
}

func (c *Chunk) Get(offset int) int {
	return c.data[offset]
}

func (c *Chunk) GetConstant(offset int) val.Value {
	return c.constants.GetValue(offset)
}

func (c *Chunk) GetLine(offset int) int {
	return c.lines[offset]
}

func (c *Chunk) Len() int {
	return len(c.data)
}
