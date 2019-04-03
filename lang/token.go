package lang

import (
	"fmt"
)

// Token represents one token in the language
type Token struct {
	tokenType TokenType
	lexeme    string
	literal   interface{}
	index     int
	line      int
	column    int
}

// MakeToken creates new token
func MakeToken(tokenType TokenType, lexeme string, literal interface{}, index, line, column int) Token {
	t := Token{}
	t.tokenType = tokenType
	t.lexeme = lexeme
	t.literal = literal
	t.index = index
	t.line = line
	t.column = column

	return t
}

func (t *Token) String() string {
	return fmt.Sprintf("%v %v %v", t.tokenType, t.lexeme, t.literal)
}
