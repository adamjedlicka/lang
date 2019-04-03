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

// declaration → fnDeclaration
//             | varDeclaration
//             | statement ;
func (p *Parser) declaration() (Stmnt, error) {
	if p.match(Func) {
		return p.function("function")
	} else if p.match(Var) {
		return p.varDeclaration()
	}

	return p.statement()

	// TODO : Synchronization
}

// function → IDENTIFIER "(" parameters? ")" block ;
func (p *Parser) function(kind string) (Stmnt, error) {
	name, err := p.consume(Identifier, "Expect "+kind+" name.")
	if err != nil {
		return nil, err
	}

	_, err = p.consume(LeftParen, "Expect '(' after "+kind+" name.")
	if err != nil {
		return nil, err
	}

	parameters := make([]Token, 0)
	if !p.check(RightParen) {
		for {
			parameter, err := p.consume(Identifier, "Exptect parameter name.")
			if err != nil {
				return nil, err
			}

			parameters = append(parameters, parameter)

			if !p.match(Comma) {
				break
			}
		}
	}

	_, err = p.consume(RightParen, "Expect ')' after parameters.")
	if err != nil {
		return nil, err
	}

	_, err = p.consume(LeftBrace, "Expect '{' before "+kind+" body.")
	if err != nil {
		return nil, err
	}

	block, err := p.block()
	if err != nil {
		return nil, err
	}

	body := (block.(BlockStmnt)).stmnts

	return MakeFnStmnt(name, parameters, body), nil
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
//           | forStatement
//           | returnStatement
//           | whileStatement
//           | printStatement
//           | block ;
func (p *Parser) statement() (Stmnt, error) {
	if p.match(If) {
		return p.ifStatement()
	} else if p.match(For) {
		return p.forStatement()
	} else if p.match(Return) {
		return p.returnStatement()
	} else if p.match(While) {
		return p.whileStatement()
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

// ifStatement → "if" expression block ( "else" block )? ;
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

// forStmt → "for" ( varDeclaration | expressionStatement | ";" )
//                 expression? ";"
//                 expression? ")" block ;
func (p *Parser) forStatement() (Stmnt, error) {
	var err error
	var initializer Stmnt
	var condition Expr
	var increment Expr

	// initializer
	if p.match(Semicolon) {
		initializer = nil
	} else if p.match(Var) {
		initializer, err = p.varDeclaration()
		if err != nil {
			return nil, err
		}
	} else {
		initializer, err = p.expressionStatement()
		if err != nil {
			return nil, err
		}
	}

	// condition
	if !p.check(Semicolon) {
		condition, err = p.expression()
		if err != nil {
			return nil, err
		}
	}
	_, err = p.consume(Semicolon, "Expect ';' after loop condition;")
	if err != nil {
		return nil, err
	}

	// increment
	if !p.check(LeftBrace) {
		increment, err = p.expression()
		if err != nil {
			return nil, err
		}
	}
	_, err = p.consume(LeftBrace, "Expect '{' after for clause.")
	if err != nil {
		return nil, err
	}

	block, err := p.block()
	if err != nil {
		return nil, err
	}

	stmnts := (block.(BlockStmnt)).stmnts

	if increment != nil {
		stmnts = append(stmnts, MakeExpressionStmnt(increment))
	}

	if condition == nil {
		condition = MakeLiteralExpr(true)
	}

	while := MakeWhileStmnt(condition, MakeBlockStmnt(stmnts))

	if initializer != nil {
		return MakeBlockStmnt([]Stmnt{initializer, while}), nil
	}

	return while, nil
}

// returnStatement → "return" expression? ";" ;
func (p *Parser) returnStatement() (Stmnt, error) {
	var err error
	var value Expr

	keyword := p.previous()

	if !p.check(Semicolon) {
		value, err = p.expression()
		if err != nil {
			return nil, err
		}
	}

	_, err = p.consume(Semicolon, "Expect ';' after return value.")
	if err != nil {
		return nil, err
	}

	return MakeReturnStmnt(keyword, value), nil
}

// whileStatement → "if" expression block ;
func (p *Parser) whileStatement() (Stmnt, error) {
	condition, err := p.expression()
	if err != nil {
		return nil, err
	}

	_, err = p.consume(LeftBrace, "Expect '{' after while condition.")
	if err != nil {
		return nil, err
	}

	body, err := p.block()
	if err != nil {
		return nil, err
	}

	return MakeWhileStmnt(condition, body), nil
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

// expression → "fn" lambda
//            | assignment ;
func (p *Parser) expression() (Expr, error) {
	if p.match(Func) {
		return p.lambda()
	}

	return p.assignment()
}

// lambda → "fn" "(" parameters? ")" block ;
func (p *Parser) lambda() (Expr, error) {
	_, err := p.consume(LeftParen, "Expect '(' after 'fn'.")
	if err != nil {
		return nil, err
	}

	parameters := make([]Token, 0)
	if !p.check(RightParen) {
		for {
			parameter, err := p.consume(Identifier, "Exptect parameter name.")
			if err != nil {
				return nil, err
			}

			parameters = append(parameters, parameter)

			if !p.match(Comma) {
				break
			}
		}
	}

	_, err = p.consume(RightParen, "Expect ')' after parameters.")
	if err != nil {
		return nil, err
	}

	_, err = p.consume(LeftBrace, "Expect '{' before lambda body.")
	if err != nil {
		return nil, err
	}

	block, err := p.block()
	if err != nil {
		return nil, err
	}

	body := (block.(BlockStmnt)).stmnts

	return MakeLambdaExpr(parameters, body), nil
}

// assignment → IDENTIFIER "=" assignment
//            | or ;
func (p *Parser) assignment() (Expr, error) {
	expr, err := p.or()
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

// logic_or → and ( "or" and )* ;
func (p *Parser) or() (Expr, error) {
	expr, err := p.and()
	if err != nil {
		return nil, err
	}

	if p.match(Or) {
		operator := p.previous()
		right, err := p.and()
		if err != nil {
			return nil, err
		}

		expr = MakeLogicalExpr(operator, expr, right)
	}

	return expr, nil
}

// logic_and → equality ( "and" equality )* ;
func (p *Parser) and() (Expr, error) {
	expr, err := p.equality()
	if err != nil {
		return nil, err
	}

	if p.match(And) {
		operator := p.previous()
		right, err := p.equality()
		if err != nil {
			return nil, err
		}

		expr = MakeLogicalExpr(operator, expr, right)
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
//       | call ;
func (p *Parser) unary() (Expr, error) {
	if p.match(Bang, Minus) {
		operator := p.previous()
		right, err := p.unary()
		if err != nil {
			return nil, err
		}

		return MakeUnaryExpr(operator, right), nil
	}

	return p.call()
}

// call → primary ( "(" arguments? ")" )* ;
func (p *Parser) call() (Expr, error) {
	expr, err := p.primary()
	if err != nil {
		return nil, err
	}

	for {
		if p.match(LeftParen) {
			expr, err = p.finishCall(expr)
			if err != nil {
				return nil, err
			}
		} else {
			break
		}
	}

	return expr, nil
}

func (p *Parser) finishCall(callee Expr) (Expr, error) {
	arguments := make([]Expr, 0)

	if !p.check(RightParen) {
		for {
			expr, err := p.expression()
			if err != nil {
				return nil, err
			}

			arguments = append(arguments, expr)

			if !p.match(Comma) {
				break
			}
		}
	}

	paren, err := p.consume(RightParen, "Expect ')' after arguments.")
	if err != nil {
		return nil, err
	}

	return MakeCallExpr(callee, paren, arguments), nil
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
