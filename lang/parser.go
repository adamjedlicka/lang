package lang

// Parser represents the language parser
type Parser struct {
	tokens  []Token
	current int
}

// MakeParser creates new parser
func MakeParser() Parser {
	return Parser{}
}

func (p *Parser) Parse(tokens []Token) ([]Stmnt, error) {
	p.tokens = tokens
	p.current = 0

	stmnts := make([]Stmnt, 0)

	for !p.isAtEnd() {
		stmnt, err := p.declaration()
		if err != nil {
			return nil, err
		}

		stmnts = append(stmnts, stmnt)
	}

	return stmnts, nil
}

// declaration → varDeclaration
//             | statement ;
func (p *Parser) declaration() (Stmnt, error) {
	if p.match(Var) {
		return p.varDeclaration()
	}

	return p.statement()

	// TODO : Synchronization
}

// varDeclaration → "var" IDENTIFIER ( "=" expression )? ";" ;
func (p *Parser) varDeclaration() (Stmnt, error) {
	name, err := p.consume(Identifier, "Expect variable name.")
	if err != nil {
		return nil, err
	}

	var initializer Expr
	if p.match(Equal) {
		initializer, err = p.expression()
		if err != nil {
			return nil, err
		}
	}

	_, err = p.consume(Semicolon, "Expect ';' after variable declaration.")
	if err != nil {
		return nil, err
	}

	return MakeVarStmnt(name, initializer), nil
}

// statement → expressionStatement
//           | ifStatement
//           | printStatement
//           | block ;
func (p *Parser) statement() (Stmnt, error) {
	if p.match(If) {
		return p.ifStatement()
	} else if p.match(Print) {
		return p.printStatement()
	} else if p.match(LeftBrace) {
		return p.block()
	}

	return p.expressionStatement()
}

// expressionStatement  → expression ";" ;
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

// ifStatement → "if" expression "{" block "}" ( "else" block )? ;
func (p *Parser) ifStatement() (Stmnt, error) {
	condition, err := p.expression()
	if err != nil {
		return nil, err
	}

	_, err = p.consume(LeftBrace, "Expect '{' after if condition.")
	if err != nil {
		return nil, err
	}

	thenBranch, err := p.block()
	if err != nil {
		return nil, err
	}

	var elseBranch Stmnt
	if p.match(Else) {
		_, err = p.consume(LeftBrace, "Expect '{' after 'else'.")
		if err != nil {
			return nil, err
		}

		elseBranch, err = p.block()
		if err != nil {
			return nil, err
		}
	}

	return MakeIfStmnt(condition, thenBranch, elseBranch), nil
}

// printStatement → "print" expression ";" ;
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

// block → "{" declaration* "}" ;
func (p *Parser) block() (Stmnt, error) {
	stmnts := make([]Stmnt, 0)

	for !p.check(RightBrace) && !p.isAtEnd() {
		stmnt, err := p.declaration()
		if err != nil {
			return nil, err
		}

		stmnts = append(stmnts, stmnt)
	}

	_, err := p.consume(RightBrace, "Expected '}' after block.")
	if err != nil {
		return nil, err
	}

	return MakeBlockStmnt(stmnts), nil
}

// expression → assignment ;
func (p *Parser) expression() (Expr, error) {
	return p.assignment()
}

// assignment → IDENTIFIER "=" assignment
//            | equality ;
func (p *Parser) assignment() (Expr, error) {
	expr, err := p.equality()
	if err != nil {
		return nil, err
	}

	if p.match(Equal) {
		equals := p.previous()
		value, err := p.assignment()
		if err != nil {
			return nil, err
		}

		if expr, ok := expr.(VariableExpr); ok {
			name := expr.name

			return MakeAssignExpr(name, value), nil
		}

		return nil, NewRuntimeError(equals.line, "Invalid assignment target.")
	}

	return expr, nil
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

		expr = MakeBinaryExpr(expr, operator, right)
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

		expr = MakeBinaryExpr(expr, operator, right)
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

		expr = MakeBinaryExpr(expr, operator, right)
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

		expr = MakeBinaryExpr(expr, operator, right)
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

		return MakeUnaryExpr(operator, right), nil
	}

	return p.primary()
}

// primary → "false" | "true" | "null"
//         | NUMBER | STRING
//         | "(" expression ")"
//         | IDENTIFIER ;
func (p *Parser) primary() (Expr, error) {
	if p.match(False) {
		return MakeLiteralExpr(false), nil
	} else if p.match(True) {
		return MakeLiteralExpr(true), nil
	} else if p.match(Null) {
		return MakeLiteralExpr(nil), nil
	} else if p.match(Number, String) {
		return MakeLiteralExpr(p.previous().literal), nil
	} else if p.match(LeftParen) {
		expr, err := p.expression()
		if err != nil {
			return nil, err
		}

		_, err = p.consume(RightParen, "Expected ')' after expression.")
		if err != nil {
			return nil, err
		}

		return MakeGroupingExpr(expr), nil
	} else if p.match(Identifier) {
		return MakeVariableExpr(p.previous()), nil
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
