package compiler

type Token struct {
	tokenType TokenType
	lexeme    string
	line      int
	column    int
}

type TokenType uint8

// List of all possible token types.
const (
	// Single character tokens.
	TokenLeftParen TokenType = iota
	TokenRightParen
	TokenLeftBrace
	TokenRightBrace
	TokenComma
	TokenDot
	TokenMinus
	TokenPlus
	TokenSemicolon
	TokenSlash
	TokenStar

	// One or two character tokens.
	TokenBang
	TokenBangEqual
	TokenEqual
	TokenEqualEqual
	TokenGreater
	TokenGreaterEqual
	TokenLess
	TokenLessEqual

	// Literals.
	TokenIdentifier
	TokenString
	TokenNumber

	// Keywords.
	TokenAnd
	TokenClass
	TokenElse
	TokenFalse
	TokenFor
	TokenFn
	TokenIf
	TokenNull
	TokenOr
	TokenPrint
	TokenReturn
	TokenSuper
	TokenThis
	TokenTrue
	TokenVar
	TokenWhile

	TokenError
	TokenEOF
)
