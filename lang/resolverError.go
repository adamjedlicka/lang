package lang

import (
	"fmt"
)

type ResolverError struct {
	token   Token
	message string
}

func NewResolverError(token Token, message string) ResolverError {
	e := ResolverError{}
	e.token = token
	e.message = message

	return e
}

func (e ResolverError) Error() string {
	if e.token.tokenType == EOF {
		return fmt.Sprintf("[line %v] ResolverError at end: '%v", e.token.line, e.message)
	}

	return fmt.Sprintf("[line %v] ResolverError at '%v': '%v", e.token.line, e.token.lexeme, e.message)
}
