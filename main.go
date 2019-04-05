package main

import (
	"github.com/adamjedlicka/lang/src/code"
	"github.com/adamjedlicka/lang/src/debug"
	"github.com/adamjedlicka/lang/src/val"
)

func main() {
	chunk := code.NewChunk()

	constant := chunk.AddConstant(val.NewNumber(1.2))

	chunk.Write(code.OpConstant, 123)
	chunk.WriteRaw(constant, 123)

	chunk.Write(code.OpReturn, 123)

	debug.DisassembleChunk(chunk, "test chunk")
}
