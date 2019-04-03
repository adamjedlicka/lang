package lang

import (
	"fmt"
)

// Token represents one token in the language
type Token struct {
	tokenType TokenType
	lexeme    string
	literal   interface{}
	line      int
	column    int
}

// MakeToken creates new token
func MakeToken(tokenType TokenType, lexeme string, literal interface{}, line, column int) Token {
	t := Token{}
	t.tokenType = tokenType
	t.lexeme = lexeme
	t.literal = literal
	t.line = line
	t.column = column

	return t
}

func (t *Token) String() string {
	return fmt.Sprintf("%v %v %v", t.tokenType, t.lexeme, t.literal)
}
