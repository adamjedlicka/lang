package lang

import (
	"fmt"
)

type ParserError struct {
	token   Token
	message string
}

func NewParserError(token Token, message string) ParserError {
	e := ParserError{}
	e.token = token
	e.message = message

	return e
}

func (e ParserError) Error() string {
	if e.token.tokenType == EOF {
		return fmt.Sprintf("[line %v] ParserError at end: '%v", e.token.line, e.message)
	}

	return fmt.Sprintf("[line %v] ParserError at '%v': '%v", e.token.line, e.token.lexeme, e.message)
}
