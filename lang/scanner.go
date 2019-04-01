package lang

import (
	"strconv"
)

var keywords = map[string]TokenType{
	"and":    And,
	"class":  Class,
	"else":   Else,
	"false":  False,
	"for":    For,
	"func":   Func,
	"if":     If,
	"null":   Null,
	"or":     Or,
	"print":  Print,
	"return": Return,
	"super":  Super,
	"this":   This,
	"true":   True,
	"var":    Var,
	"while":  While,
}

// Scanner scans the source code and returns slice of Tokens
type Scanner struct {
	source string
	tokens []Token

	start   int
	current int
	line    int
}

// MakeScanner creates new scanner
func MakeScanner() Scanner {
	return Scanner{}
}

// ScanTokens scans the source code and returns list of tokens
func (s *Scanner) ScanTokens(source string) ([]Token, error) {
	s.source = source
	s.tokens = make([]Token, 0)

	s.start = 0
	s.current = 0
	s.line = 1

	for !s.isAtEnd() {
		s.start = s.current
		err := s.scanToken()
		if err != nil {
			return nil, err
		}
	}

	s.tokens = append(s.tokens, MakeToken(EOF, "", nil, s.line))

	return s.tokens, nil
}

func (s *Scanner) scanToken() error {
	c := s.advance()
	switch c {
	case '(':
		s.addToken(LeftParen, nil)
	case ')':
		s.addToken(RightParen, nil)
	case '{':
		s.addToken(LeftBrace, nil)
	case '}':
		s.addToken(RightBrace, nil)
	case ',':
		s.addToken(Comma, nil)
	case '.':
		s.addToken(Dot, nil)
	case '-':
		s.addToken(Minus, nil)
	case '+':
		s.addToken(Plus, nil)
	case ';':
		s.addToken(Semicolon, nil)
	case '*':
		s.addToken(Star, nil)
	case '!':
		if s.match('=') {
			s.addToken(BangEqual, nil)
		} else {
			s.addToken(Bang, nil)
		}
	case '=':
		if s.match('=') {
			s.addToken(EqualEqual, nil)
		} else {
			s.addToken(Equal, nil)
		}
	case '<':
		if s.match('=') {
			s.addToken(LessEqual, nil)
		} else {
			s.addToken(Less, nil)
		}
	case '>':
		if s.match('=') {
			s.addToken(GreaterEqual, nil)
		} else {
			s.addToken(Greater, nil)
		}
	case '/':
		if s.match('/') {
			for s.peek() != '\n' && !s.isAtEnd() {
				s.advance()
			}
		} else {
			s.addToken(Slash, nil)
		}
	case ' ':
	case '\r':
	case '\t':
	case '\n':
		s.line++
	case '"':
		err := s.string()
		if err != nil {
			return err
		}
	default:
		if s.isDigit(c) {
			s.number()
		} else if s.isAlpha(c) {
			s.identifier()
		} else {
			return NewScannerError(s.line, "Unexpected character.")
		}
	}

	return nil
}

func (s *Scanner) addToken(tokenType TokenType, literal interface{}) {
	text := s.source[s.start:s.current]
	s.tokens = append(s.tokens, MakeToken(tokenType, text, literal, s.line))
}

func (s *Scanner) advance() rune {
	s.current++

	return rune(s.source[s.current-1])
}

func (s *Scanner) match(expected rune) bool {
	if s.isAtEnd() {
		return false
	}

	if rune(s.source[s.current]) != expected {
		return false
	}

	s.current++

	return true
}

func (s *Scanner) peek() rune {
	if s.isAtEnd() {
		return 0
	}

	return rune(s.source[s.current])
}

func (s *Scanner) peekNext() rune {
	if s.current+1 >= len(s.source) {
		return 0
	}

	return rune(s.source[s.current+1])
}

func (s *Scanner) string() error {
	for s.peek() != '"' && !s.isAtEnd() {
		if s.peek() == '\n' {
			s.line++
		}

		s.advance()
	}

	if s.isAtEnd() {
		return NewScannerError(s.line, "Unterminated string.")
	}

	s.advance()

	s.addToken(String, s.source[s.start+1:s.current-1])

	return nil
}

func (s *Scanner) number() {
	for s.isDigit(s.peek()) {
		s.advance()
	}

	if s.peek() == '.' && s.isDigit(s.peekNext()) {
		s.advance()

		for s.isDigit(s.peek()) {
			s.advance()
		}
	}

	i, err := strconv.ParseFloat(s.source[s.start:s.current], 64)
	if err != nil {
		panic(err)
	}

	s.addToken(Number, i)
}

func (s *Scanner) identifier() {
	for s.isAlphaNumeric(s.peek()) {
		s.advance()
	}

	text := s.source[s.start:s.current]

	if keyType, ok := keywords[text]; ok {
		s.addToken(keyType, nil)
	} else {
		s.addToken(Identifier, nil)
	}

}

func (s *Scanner) isDigit(r rune) bool {
	return r >= '0' && r <= '9'
}

func (s *Scanner) isAlpha(r rune) bool {
	return (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r == '_')
}

func (s *Scanner) isAlphaNumeric(r rune) bool {
	return s.isDigit(r) || s.isAlpha(r)
}

func (s *Scanner) isAtEnd() bool {
	return s.current >= len(s.source)
}
