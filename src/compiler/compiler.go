package compiler

import (
	"fmt"

	"github.com/adamjedlicka/lang/src/code"
)

type Compiler struct {
	scanner *Scanner
}

func NewCompiler(source []rune) *Compiler {
	return &Compiler{
		scanner: NewScanner(source),
	}
}

func (c *Compiler) Compile() *code.Chunk {
	line := -1

	for {
		token := c.scanner.scanToken()

		if token.line != line {
			fmt.Printf("%4d ", token.line)
			line = token.line
		} else {
			fmt.Printf("   | ")
		}

		fmt.Printf("%2d '%s'\n", token.tokenType, token.lexeme)

		if token.tokenType == TokenEOF {
			break
		}
	}

	return nil
}
