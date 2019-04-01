package lang

type Expr interface {
	Accept(ExprVisitor) (interface{}, error)
}

type ExprVisitor interface {
	VisitBinaryExpr(BinaryExpr) (interface{}, error)
	VisitGroupingExpr(GroupingExpr) (interface{}, error)
	VisitLiteralExpr(LiteralExpr) (interface{}, error)
	VisitUnaryExpr(UnaryExpr) (interface{}, error)
	VisitVariableExpr(VariableExpr) (interface{}, error)
}

type BinaryExpr struct {
	left     Expr
	operator Token
	right    Expr
}

func MakeBinaryExpr(left Expr, operator Token, right Expr) BinaryExpr {
	return BinaryExpr{
		left:     left,
		operator: operator,
		right:    right,
	}
}

func (b BinaryExpr) Accept(visitor ExprVisitor) (interface{}, error) {
	return visitor.VisitBinaryExpr(b)
}

type GroupingExpr struct {
	expression Expr
}

func MakeGroupingExpr(expression Expr) GroupingExpr {
	return GroupingExpr{
		expression: expression,
	}
}

func (g GroupingExpr) Accept(visitor ExprVisitor) (interface{}, error) {
	return visitor.VisitGroupingExpr(g)
}

type LiteralExpr struct {
	value interface{}
}

func MakeLiteralExpr(value interface{}) LiteralExpr {
	return LiteralExpr{
		value: value,
	}
}

func (l LiteralExpr) Accept(visitor ExprVisitor) (interface{}, error) {
	return visitor.VisitLiteralExpr(l)
}

type UnaryExpr struct {
	operator Token
	right    Expr
}

func MakeUnaryExpr(operator Token, right Expr) UnaryExpr {
	return UnaryExpr{
		operator: operator,
		right:    right,
	}
}

func (u UnaryExpr) Accept(visitor ExprVisitor) (interface{}, error) {
	return visitor.VisitUnaryExpr(u)
}

type VariableExpr struct {
	name Token
}

func MakeVariableExpr(name Token) VariableExpr {
	return VariableExpr{
		name: name,
	}
}

func (u VariableExpr) Accept(visitor ExprVisitor) (interface{}, error) {
	return visitor.VisitVariableExpr(u)
}
