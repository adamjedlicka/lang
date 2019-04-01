package lang

// Parser represents the language parser
type Parser struct {
	tokens  []Token
	current int
}

// MakeParser creates new parser
func MakeParser(tokens []Token) Parser {
	p := Parser{}
	p.tokens = tokens
	p.current = 0

	return p
}

func (p *Parser) Parse() ([]Stmnt, error) {
	stmnts := make([]Stmnt, 0)

	for !p.isAtEnd() {
		stmnt, err := p.statement()
		if err != nil {
			return nil, err
		}

		stmnts = append(stmnts, stmnt)
	}

	return stmnts, nil
}

func (p *Parser) statement() (Stmnt, error) {
	if p.match(Print) {
		return p.printStatement()
	}

	return p.expressionStatement()
}

func (p *Parser) expressionStatement() (Stmnt, error) {
	expr, err := p.expression()
	if err != nil {
		return nil, err
	}

	_, err = p.consume(Semicolon, "Expected ; after expression.")
	if err != nil {
		return nil, err
	}

	return MakeExpressionStmnt(expr), nil
}

func (p *Parser) printStatement() (Stmnt, error) {
	expr, err := p.expression()
	if err != nil {
		return nil, err
	}

	_, err = p.consume(Semicolon, "Expected ; after value.")
	if err != nil {
		return nil, err
	}

	return MakePrintStmnt(expr), nil
}

func (p *Parser) expression() (Expr, error) {
	return p.equality()
}

// equality → comparison ( ( "!=" | "==" ) comparison )* ;
func (p *Parser) equality() (Expr, error) {
	expr, err := p.comparison()
	if err != nil {
		return nil, err
	}

	for p.match(BangEqual, EqualEqual) {
		operator := p.previous()
		right, err := p.comparison()
		if err != nil {
			return nil, err
		}

		expr = MakeBinary(expr, operator, right)
	}

	return expr, nil
}

// comparison → addition ( ( ">" | ">=" | "<" | "<=" ) addition )* ;
func (p *Parser) comparison() (Expr, error) {
	expr, err := p.addition()
	if err != nil {
		return nil, err
	}

	for p.match(Greater, GreaterEqual, Less, LessEqual) {
		operator := p.previous()
		right, err := p.addition()
		if err != nil {
			return nil, err
		}

		expr = MakeBinary(expr, operator, right)
	}

	return expr, nil
}

// addition → multiplication ( ( "-" | "+" ) multiplication )* ;
func (p *Parser) addition() (Expr, error) {
	expr, err := p.multiplication()
	if err != nil {
		return nil, err
	}

	for p.match(Minus, Plus) {
		operator := p.previous()
		right, err := p.multiplication()
		if err != nil {
			return nil, err
		}

		expr = MakeBinary(expr, operator, right)
	}

	return expr, nil
}

// multiplication → unary ( ( "/" | "*" ) unary )* ;
func (p *Parser) multiplication() (Expr, error) {
	expr, err := p.unary()
	if err != nil {
		return nil, err
	}

	for p.match(Slash, Star) {
		operator := p.previous()
		right, err := p.unary()
		if err != nil {
			return nil, err
		}

		expr = MakeBinary(expr, operator, right)
	}

	return expr, nil
}

// unary → ( "!" | "-" ) unary
//       | primary ;
func (p *Parser) unary() (Expr, error) {
	if p.match(Bang, Minus) {
		operator := p.previous()
		right, err := p.unary()
		if err != nil {
			return nil, err
		}

		return MakeUnary(operator, right), nil
	}

	return p.primary()
}

// primary → NUMBER | STRING | "false" | "true" | "null"
//         | "(" expression ")" ;
func (p *Parser) primary() (Expr, error) {
	if p.match(False) {
		return MakeLiteral(false), nil
	} else if p.match(True) {
		return MakeLiteral(true), nil
	} else if p.match(Null) {
		return MakeLiteral(nil), nil
	} else if p.match(Number, String) {
		return MakeLiteral(p.previous().literal), nil
	} else if p.match(LeftParen) {
		expr, err := p.expression()
		if err != nil {
			return nil, err
		}

		_, err = p.consume(RightParen, "Expected ')' after expression.")
		if err != nil {
			return nil, err
		}

		return MakeGrouping(expr), nil
	}

	return nil, NewParserError(p.peek(), "Unexpected token.")
}

func (p *Parser) consume(tokenType TokenType, message string) (Token, error) {
	if p.check(tokenType) {
		return p.advance(), nil
	}

	return Token{}, NewParserError(p.peek(), message)
}

func (p *Parser) match(tokenTypes ...TokenType) bool {
	for _, tokenType := range tokenTypes {
		if p.check(tokenType) {
			p.advance()

			return true
		}
	}

	return false
}

func (p *Parser) check(tokenType TokenType) bool {
	if p.isAtEnd() {
		return false
	}

	return p.peek().tokenType == tokenType
}

func (p *Parser) advance() Token {
	if !p.isAtEnd() {
		p.current++
	}

	return p.previous()
}

func (p *Parser) isAtEnd() bool {
	return p.peek().tokenType == EOF
}

func (p *Parser) peek() Token {
	return p.tokens[p.current]
}

func (p *Parser) previous() Token {
	return p.tokens[p.current-1]
}
