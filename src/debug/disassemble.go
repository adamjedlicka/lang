package debug

import (
	"fmt"

	"github.com/adamjedlicka/lang/src/code"
	"github.com/adamjedlicka/lang/src/val"
)

func DisassembleChunk(chunk *code.Chunk, name string) {
	fmt.Printf("== %s ==\n", name)

	for offset := 0; offset < chunk.Len(); {
		offset = DisassembleInstruction(chunk, offset)
	}
}

func DisassembleInstruction(chunk *code.Chunk, offset int) int {
	fmt.Printf("%04d ", offset)

	if offset > 0 && chunk.GetLine(offset) == chunk.GetLine(offset-1) {
		fmt.Print("   | ")
	} else {
		fmt.Printf("%4d ", chunk.GetLine(offset))
	}

	instruction := chunk.Get(offset)
	switch instruction {
	case code.OpConstant:
		return constantInstruction("OpConstant", chunk, offset)
	case code.OpAdd:
		return simpleInstruction("OpAdd", offset)
	case code.OpSubtract:
		return simpleInstruction("OpSubtract", offset)
	case code.OpMultiply:
		return simpleInstruction("OpMultiply", offset)
	case code.OpDivide:
		return simpleInstruction("OpDivide", offset)
	case code.OpNegate:
		return simpleInstruction("OpNegate", offset)
	case code.OpReturn:
		return simpleInstruction("OpReturn", offset)
	}

	fmt.Printf("Unknown opcode %d\n", instruction)
	return offset + 1
}

func constantInstruction(name string, chunk *code.Chunk, offset int) int {
	constant := chunk.GetRaw(offset + 1)
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
