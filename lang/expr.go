package lang

type Expr interface {
	Accept(ExprVisitor) (interface{}, error)
}

type ExprVisitor interface {
	VisitBinaryExpr(Binary) (interface{}, error)
	VisitGroupingExpr(Grouping) (interface{}, error)
	VisitLiteralExpr(Literal) (interface{}, error)
	VisitUnaryExpr(Unary) (interface{}, error)
	VisitVariableExpr(Variable) (interface{}, error)
}

type Binary struct {
	left     Expr
	operator Token
	right    Expr
}

func MakeBinary(left Expr, operator Token, right Expr) Binary {
	return Binary{
		left:     left,
		operator: operator,
		right:    right,
	}
}

func (b Binary) Accept(visitor ExprVisitor) (interface{}, error) {
	return visitor.VisitBinaryExpr(b)
}

type Grouping struct {
	expression Expr
}

func MakeGrouping(expression Expr) Grouping {
	return Grouping{
		expression: expression,
	}
}

func (g Grouping) Accept(visitor ExprVisitor) (interface{}, error) {
	return visitor.VisitGroupingExpr(g)
}

type Literal struct {
	value interface{}
}

func MakeLiteral(value interface{}) Literal {
	return Literal{
		value: value,
	}
}

func (l Literal) Accept(visitor ExprVisitor) (interface{}, error) {
	return visitor.VisitLiteralExpr(l)
}

type Unary struct {
	operator Token
	right    Expr
}

func MakeUnary(operator Token, right Expr) Unary {
	return Unary{
		operator: operator,
		right:    right,
	}
}

func (u Unary) Accept(visitor ExprVisitor) (interface{}, error) {
	return visitor.VisitUnaryExpr(u)
}

type Variable struct {
	name Token
}

func MakeVariable(name Token) Variable {
	return Variable{
		name: name,
	}
}

func (u Variable) Accept(visitor ExprVisitor) (interface{}, error) {
	return visitor.VisitVariableExpr(u)
}
