package compiler

import (
	"fmt"
	"os"
	"strconv"

	"github.com/adamjedlicka/lang/src/code"
	"github.com/adamjedlicka/lang/src/val"
)

type Parser struct {
	scanner *Scanner

	compilingChunk *code.Chunk

	hadError  bool
	panicMode bool
	current   Token
	previous  Token
}

func NewParser(scanner *Scanner) *Parser {
	return &Parser{
		scanner: scanner,

		compilingChunk: code.NewChunk(),

		hadError:  false,
		panicMode: false,
	}
}

func (p *Parser) Parse() *code.Chunk {
	p.advance()
	p.expression()
	p.consume(TokenEOF, "Expect EOF.")
	p.endCompiler()

	if p.hadError {
		return nil
	}

	return p.currentChunk()
}

func (p *Parser) advance() {
	p.previous = p.current

	for {
		p.current = p.scanner.scanToken()

		if p.current.tokenType != TokenError {
			return
		}

		p.errorAtCurrent(p.current.lexeme)
	}
}

func (p *Parser) expression() {
	p.parsePrecedence(PrecedenceAssignment)
}

func (p *Parser) number() {
	value, err := strconv.ParseFloat(p.previous.lexeme, 64)
	if err != nil {
		panic(err)
	}

	p.emitConstant(val.NewNumber(value))
}

func (p *Parser) grouping() {
	p.expression()
	p.consume(TokenRightParen, "Expect ')' after expression.")
}

func (p *Parser) unary() {
	operatorType := p.previous.tokenType

	// Compile the operand.
	p.parsePrecedence(PrecedenceUnary)

	// Emit the operator instruction.
	switch operatorType {
	case TokenMinus:
		p.emitInstruction(code.OpNegate)
	}
}

func (p *Parser) binary() {
	// Remember the oprator.
	operatorType := p.previous.tokenType

	// Compile the right operand.
	rule := p.getRule(operatorType)
	p.parsePrecedence(rule.precedence + 1)

	switch operatorType {
	case TokenPlus:
		p.emitInstruction(code.OpAdd)
	case TokenMinus:
		p.emitInstruction(code.OpSubtract)
	case TokenStar:
		p.emitInstruction(code.OpMultiply)
	case TokenSlash:
		p.emitInstruction(code.OpDivide)
	}
}

func (p *Parser) parsePrecedence(precedence Precedence) {
	p.advance()

	prefixRule := p.getRule(p.previous.tokenType).prefix
	if prefixRule == nil {
		p.error("Expect expression.")
		return
	}

	prefixRule(p)

	for precedence <= p.getRule(p.current.tokenType).precedence {
		p.advance()

		infixRule := p.getRule(p.previous.tokenType).infix

		infixRule(p)
	}
}

func (p *Parser) getRule(tokenType TokenType) *ParseRule {
	return &rules[int(tokenType)]
}

func (p *Parser) consume(tokenType TokenType, message string) {
	if p.current.tokenType == tokenType {
		p.advance()
		return
	}

	p.errorAtCurrent(message)
}

func (p *Parser) endCompiler() {
	p.emitReturn()
}

func (p *Parser) emitReturn() {
	p.emitInstruction(code.OpReturn)
}

func (p *Parser) emitInstruction(instruction code.OpCode) {
	p.currentChunk().Write(instruction, p.previous.line)
}

func (p *Parser) emitByte(data uint8) {
	p.currentChunk().WriteRaw(data, p.previous.line)
}

func (p *Parser) emitBytes(byte1, byte2 uint8) {
	p.emitByte(byte1)
	p.emitByte(byte2)
}

func (p *Parser) emitConstant(value val.Value) {
	offset := p.currentChunk().AddConstant(value)

	// Check if offset is greater or equal maximum value of uint8
	if offset >= (^uint8(0)) {
		p.error("Too many constants in one chunk.")
	}

	p.emitBytes(uint8(code.OpConstant), offset)
}

func (p *Parser) currentChunk() *code.Chunk {
	return p.compilingChunk
}

func (p *Parser) error(message string) {
	p.errorAt(p.previous, message)
}

func (p *Parser) errorAtCurrent(message string) {
	p.errorAt(p.current, message)
}

func (p *Parser) errorAt(token Token, message string) {
	if p.panicMode {
		return
	}

	p.panicMode = true

	fmt.Fprintf(os.Stderr, "[line %d:%d] Error", token.line, token.column)

	if token.tokenType == TokenEOF {
		fmt.Fprintf(os.Stderr, " at end")
	} else if token.tokenType == TokenError {
		// Nothing.
	} else {
		fmt.Fprintf(os.Stderr, " at '%s'", token.lexeme)
	}

	fmt.Fprintf(os.Stderr, ": %s\n", message)

	p.hadError = true
}
