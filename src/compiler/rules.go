package compiler

type ParseFn func(*Parser)

type ParseRule struct {
	prefix     ParseFn
	infix      ParseFn
	precedence Precedence
}

var rules []ParseRule

func init() {
	rules = []ParseRule{
		{(*Parser).grouping, nil, PrecedenceCall},           // TokenLeftParen
		{nil, nil, PrecedenceNone},                          // TokenRightParen
		{nil, nil, PrecedenceNone},                          // TokenLeftBrace
		{nil, nil, PrecedenceNone},                          // TokenRightBrace
		{nil, nil, PrecedenceNone},                          // TokenComma
		{nil, nil, PrecedenceCall},                          // TokenDot
		{(*Parser).unary, (*Parser).binary, PrecedenceTerm}, // TokenMinus
		{nil, (*Parser).binary, PrecedenceTerm},             // TokenPlus
		{nil, nil, PrecedenceNone},                          // TokenSemicolon
		{nil, (*Parser).binary, PrecedenceFactor},           // TokenSlash
		{nil, (*Parser).binary, PrecedenceFactor},           // TokenStar
		{nil, nil, PrecedenceNone},                          // TokenBang
		{nil, nil, PrecedenceEquality},                      // TokenBangEqual
		{nil, nil, PrecedenceNone},                          // TokenEqual
		{nil, nil, PrecedenceEquality},                      // TokenEqualEqual
		{nil, nil, PrecedenceComparison},                    // TokenGreater
		{nil, nil, PrecedenceComparison},                    // TokenGreaterEqual
		{nil, nil, PrecedenceComparison},                    // TokenLess
		{nil, nil, PrecedenceComparison},                    // TokenLessEqual
		{nil, nil, PrecedenceNone},                          // TokenIdentifier
		{nil, nil, PrecedenceNone},                          // TokenString
		{(*Parser).number, nil, PrecedenceNone},             // TokenNumber
		{nil, nil, PrecedenceAnd},                           // TokenAnd
		{nil, nil, PrecedenceNone},                          // TokenClass
		{nil, nil, PrecedenceNone},                          // TokenElse
		{nil, nil, PrecedenceNone},                          // TokenFalse
		{nil, nil, PrecedenceNone},                          // TokenFor
		{nil, nil, PrecedenceNone},                          // TokenFn
		{nil, nil, PrecedenceNone},                          // TokenIf
		{nil, nil, PrecedenceNone},                          // TokenNull
		{nil, nil, PrecedenceOr},                            // TokenOr
		{nil, nil, PrecedenceNone},                          // TokenPrint
		{nil, nil, PrecedenceNone},                          // TokenReturn
		{nil, nil, PrecedenceNone},                          // TokenSuper
		{nil, nil, PrecedenceNone},                          // TokenThis
		{nil, nil, PrecedenceNone},                          // TokenTrue
		{nil, nil, PrecedenceNone},                          // TokenVar
		{nil, nil, PrecedenceNone},                          // TokenWhile
		{nil, nil, PrecedenceNone},                          // TokenError
		{nil, nil, PrecedenceNone},                          // TokenEOF
	}
}
