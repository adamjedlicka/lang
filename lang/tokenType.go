package lang

// TokenType represents type of the Token
type TokenType string

// List of all tokens
const (
	// Single character tokens
	Comma      TokenType = "COMMA"
	Dot        TokenType = "DOT"
	LeftBrace  TokenType = "RIGHT_BRACE"
	LeftParen  TokenType = "LEFT_PAREN"
	Minus      TokenType = "MINUS"
	Plus       TokenType = "PLUS"
	RightBrace TokenType = "RIGHT_BRACE"
	RightParen TokenType = "RIGHT_PAREN"
	Semicolon  TokenType = "SEMICOLON"
	Slash      TokenType = "SLASH"
	Star       TokenType = "STAR"

	// One or two character tokens
	Bang         TokenType = "BANG"
	BangEqual    TokenType = "BANG_EQUAL"
	Equal        TokenType = "EQUAL"
	EqualEqual   TokenType = "EQUAL_EQUAL"
	Greater      TokenType = "GREATER"
	GreaterEqual TokenType = "GREATER_EQUAL"
	Less         TokenType = "LESS"
	LessEqual    TokenType = "LESS_EQUAL"

	// Literals
	Identifier TokenType = "IDENTIFIER"
	Number     TokenType = "NUMBER"
	String     TokenType = "STRING"

	// Keywords
	And    TokenType = "AND"
	Class  TokenType = "CLASS"
	Else   TokenType = "ELSE"
	False  TokenType = "FALSE"
	For    TokenType = "FOR"
	Func   TokenType = "FUNC"
	If     TokenType = "IF"
	Null   TokenType = "NULL"
	Or     TokenType = "OR"
	Print  TokenType = "PRINT"
	Return TokenType = "RETURN"
	Super  TokenType = "SUPER"
	This   TokenType = "THIS"
	True   TokenType = "TRUE"
	Var    TokenType = "VAR"
	While  TokenType = "WHILE"

	EOF TokenType = "EOF"
)
