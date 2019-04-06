package compiler

import "unicode"

type Scanner struct {
	source  []rune
	start   int
	current int
	line    int
}

func NewScanner(source []rune) *Scanner {
	return &Scanner{
		source:  source,
		start:   0,
		current: 0,
		line:    1,
	}
}

func (s *Scanner) scanToken() Token {
	s.skipWhitespace()

	s.start = s.current

	if s.isAtEnd() {
		return s.makeToken(TokenEOF)
	}

	r := s.advance()

	if s.isAlpha(r) {
		return s.identifier()
	}

	if s.isDigit(r) {
		return s.number()
	}

	switch r {
	case '(':
		return s.makeToken(TokenLeftParen)
	case ')':
		return s.makeToken(TokenRightParen)
	case '{':
		return s.makeToken(TokenLeftBrace)
	case '}':
		return s.makeToken(TokenRightBrace)
	case ';':
		return s.makeToken(TokenSemicolon)
	case ',':
		return s.makeToken(TokenComma)
	case '.':
		return s.makeToken(TokenDot)
	case '-':
		return s.makeToken(TokenMinus)
	case '+':
		return s.makeToken(TokenPlus)
	case '/':
		return s.makeToken(TokenSlash)
	case '*':
		return s.makeToken(TokenStar)
	case '!':
		if s.match('=') {
			return s.makeToken(TokenBangEqual)
		}
		return s.makeToken(TokenBang)
	case '=':
		if s.match('=') {
			return s.makeToken(TokenEqualEqual)
		}
		return s.makeToken(TokenEqual)
	case '<':
		if s.match('=') {
			return s.makeToken(TokenLessEqual)
		}
		return s.makeToken(TokenLess)
	case '>':
		if s.match('=') {
			return s.makeToken(TokenGreaterEqual)
		}
		return s.makeToken(TokenGreater)
	case '"':
		return s.string()
	}

	return s.errorToken("Unexpected character.")
}

func (s *Scanner) identifier() Token {
	for s.isAlphaNumeric(s.peek()) {
		s.advance()
	}

	return s.makeToken(s.identifierType())
}

func (s *Scanner) string() Token {
	for s.peek() != '"' && !s.isAtEnd() {
		if s.peek() == '\n' {
			s.line++
		}

		s.advance()
	}

	if s.isAtEnd() {
		return s.errorToken("Unterminated string.")
	}

	s.advance()

	return s.makeToken(TokenString)
}

func (s *Scanner) number() Token {
	for s.isDigit(s.peek()) {
		s.advance()
	}

	if s.peek() == '.' && s.isDigit(s.peekNext()) {
		s.advance()

		for s.isDigit(s.peek()) {
			s.advance()
		}
	}

	return s.makeToken(TokenNumber)
}

func (s *Scanner) makeToken(tokenType TokenType) Token {
	rowStart := 0

	for i := s.start - 1; i >= 0; i-- {
		if s.source[i] == '\n' {
			rowStart = i
		}
	}

	return Token{
		tokenType: tokenType,
		lexeme:    string(s.source[s.start:s.current]),
		line:      s.line,
		column:    s.start - rowStart,
	}
}

func (s *Scanner) errorToken(message string) Token {
	rowStart := 0

	for i := s.start - 1; i >= 0; i-- {
		if s.source[i] == '\n' {
			rowStart = i
		}
	}

	return Token{
		tokenType: TokenError,
		lexeme:    message,
		line:      s.line,
		column:    s.start - rowStart,
	}
}

func (s *Scanner) advance() rune {
	s.current++

	return s.source[s.current-1]
}

func (s *Scanner) match(r rune) bool {
	if s.isAtEnd() {
		return false
	}

	if r != s.source[s.current] {
		return false
	}

	s.current++

	return true
}

func (s *Scanner) peek() rune {
	return s.source[s.current]
}

func (s *Scanner) peekNext() rune {
	if s.isAtEnd() {
		return 0
	}

	return s.source[s.current+1]
}

func (s *Scanner) skipWhitespace() {
	if s.isAtEnd() {
		return
	}

	for {
		switch s.peek() {
		case ' ':
			fallthrough
		case '\r':
			fallthrough
		case '\t':
			s.advance()
		case '\n':
			s.line++
			s.advance()
		case '/':
			if s.peekNext() == '/' {
				for s.peek() != '\n' && !s.isAtEnd() {
					s.advance()
				}
			} else {
				return
			}
		default:
			return
		}
	}
}

func (s *Scanner) isAtEnd() bool {
	return s.current == len(s.source)-1
}

func (s *Scanner) isDigit(r rune) bool {
	return unicode.IsDigit(r)
}

func (s *Scanner) isAlpha(r rune) bool {
	return r == '_' || unicode.IsLetter(r)
}

func (s *Scanner) isAlphaNumeric(r rune) bool {
	return s.isDigit(r) || s.isAlpha(r)
}

func (s *Scanner) identifierType() TokenType {
	switch s.source[s.start] {
	case 'a':
		return s.checkKeyword(1, 2, "nd", TokenAnd)
	case 'c':
		return s.checkKeyword(1, 4, "lass", TokenClass)
	case 'e':
		return s.checkKeyword(1, 3, "lse", TokenElse)
	case 'f':
		if s.current-s.start > 1 {
			switch s.source[s.start+1] {
			case 'a':
				return s.checkKeyword(2, 3, "lse", TokenElse)
			case 'n':
				return s.checkKeyword(2, 1, "n", TokenFn)
			case 'o':
				return s.checkKeyword(2, 1, "r", TokenFor)
			}
		}
	case 'i':
		return s.checkKeyword(1, 1, "f", TokenIf)
	case 'n':
		return s.checkKeyword(1, 3, "ull", TokenNull)
	case 'o':
		return s.checkKeyword(1, 1, "r", TokenOr)
	case 'p':
		return s.checkKeyword(1, 4, "rint", TokenPrint)
	case 'r':
		return s.checkKeyword(1, 5, "eturn", TokenReturn)
	case 's':
		return s.checkKeyword(1, 4, "uper", TokenSuper)
	case 't':
		if s.current-s.start > 1 {
			switch s.source[s.start+1] {
			case 'h':
				return s.checkKeyword(2, 2, "is", TokenThis)
			case 'r':
				return s.checkKeyword(2, 2, "ue", TokenTrue)
			}
		}
	case 'v':
		return s.checkKeyword(1, 2, "ar", TokenVar)
	case 'w':
		return s.checkKeyword(1, 4, "hile", TokenWhile)
	}

	return TokenIdentifier
}

func (s *Scanner) checkKeyword(start, length int, rest string, tokenType TokenType) TokenType {
	if s.current-s.start == start+length && rest == string(s.source[s.start+start:s.start+start+length]) {
		return tokenType
	}

	return TokenIdentifier
}
