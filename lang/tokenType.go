package lang

// TokenType represents type of the Token
type TokenType string

// List of all tokens
const (
	LeftParen  TokenType = "LEFT_PAREN"
	RightParen TokenType = "RIGHT_PAREN"
	LeftBrace  TokenType = "RIGHT_BRACE"
	RightBrace TokenType = "RIGHT_BRACE"
	Comma      TokenType = "COMMA"
	Dot        TokenType = "DOT"
	Minus      TokenType = "MINUS"
	Plus       TokenType = "PLUS"
	Star       TokenType = "STAR"
	Slash      TokenType = "SLASH"

	Bang         TokenType = "BANG"
	BangEqual    TokenType = "BANG_EQUAL"
	Equal        TokenType = "EQUAL"
	EqualEqual   TokenType = "EQUAL_EQUAL"
	Greater      TokenType = "GREATER"
	GreaterEqual TokenType = "GREATER_EQUAL"
	Less         TokenType = "LESS"
	LessEqual    TokenType = "LESS_EQUAL"

	Identifier TokenType = "IDENTIFIER"
	String     TokenType = "STRING"
	Number     TokenType = "NUMBER"

	And    TokenType = "AND"
	Class  TokenType = "CLASS"
	Else   TokenType = "ELSE"
	False  TokenType = "FALSE"
	Func   TokenType = "FUNC"
	For    TokenType = "FOR"
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
