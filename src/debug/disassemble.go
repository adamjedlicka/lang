package debug

import (
	"fmt"

	"github.com/adamjedlicka/lang/src/val"

	"github.com/adamjedlicka/lang/src/code"
)

func DisassembleChunk(chunk *code.Chunk, name string) {
	fmt.Printf("== %s ==\n", name)

	for offset := 0; offset < chunk.Len(); {
		offset = disassembleInstruction(chunk, offset)
	}
}

func disassembleInstruction(chunk *code.Chunk, offset int) int {
	fmt.Printf("%04d ", offset)

	if offset > 0 && chunk.GetLine(offset) == chunk.GetLine(offset-1) {
		fmt.Print("   | ")
	} else {
		fmt.Printf("%4d ", chunk.GetLine(offset))
	}

	instruction := code.OpCode(chunk.Get(offset))
	switch instruction {
	case code.OpConstant:
		return constantInstruction("OP_CONSTANT", chunk, offset)
	case code.OpReturn:
		return simpleInstruction("OP_RETURN", offset)
	}

	fmt.Printf("Unknown opcode %d\n", instruction)
	return offset + 1
}

func constantInstruction(name string, chunk *code.Chunk, offset int) int {
	constant := chunk.Get(offset + 1)
	fmt.Printf("%-20s %4d '", name, constant)
	printValue(chunk.GetConstant(constant))
	fmt.Printf("'\n")

	return offset + 2
}

func simpleInstruction(name string, offset int) int {
	fmt.Printf("%s\n", name)
	return offset + 1
}

func printValue(value val.Value) {
	fmt.Print(value.String())
}
