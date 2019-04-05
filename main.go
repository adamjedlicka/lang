package main

import (
	"github.com/adamjedlicka/lang/src/code"
	"github.com/adamjedlicka/lang/src/val"
	"github.com/adamjedlicka/lang/src/vm"
)

func main() {
	chunk := code.NewChunk()

	x := chunk.AddConstant(val.NewNumber(1.5))
	y := chunk.AddConstant(val.NewNumber(3.75))

	chunk.Write(code.OpConstant, 123)
	chunk.WriteRaw(y, 123)
	chunk.Write(code.OpConstant, 123)
	chunk.WriteRaw(x, 123)

	chunk.Write(code.OpNegate, 123)

	chunk.Write(code.OpAdd, 123)

	chunk.Write(code.OpReturn, 123)

	vm := vm.NewVM()
	vm.Interpret(chunk)
}
