package lang

type Expr interface {
	Accept(ExprVisitor) (interface{}, error)
}

type ExprVisitor interface {
	VisitAssignExpr(AssignExpr) (interface{}, error)
	VisitBinaryExpr(BinaryExpr) (interface{}, error)
	VisitGroupingExpr(GroupingExpr) (interface{}, error)
	VisitLiteralExpr(LiteralExpr) (interface{}, error)
	VisitUnaryExpr(UnaryExpr) (interface{}, error)
	VisitVariableExpr(VariableExpr) (interface{}, error)
}

type AssignExpr struct {
	name  Token
	value Expr
}

func MakeAssignExpr(name Token, value Expr) AssignExpr {
	return AssignExpr{
		name:  name,
		value: value,
	}
}

func (e AssignExpr) Accept(visitor ExprVisitor) (interface{}, error) {
	return visitor.VisitAssignExpr(e)
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

func (e BinaryExpr) Accept(visitor ExprVisitor) (interface{}, error) {
	return visitor.VisitBinaryExpr(e)
}

type GroupingExpr struct {
	expression Expr
}

func MakeGroupingExpr(expression Expr) GroupingExpr {
	return GroupingExpr{
		expression: expression,
	}
}

func (e GroupingExpr) Accept(visitor ExprVisitor) (interface{}, error) {
	return visitor.VisitGroupingExpr(e)
}

type LiteralExpr struct {
	value interface{}
}

func MakeLiteralExpr(value interface{}) LiteralExpr {
	return LiteralExpr{
		value: value,
	}
}

func (e LiteralExpr) Accept(visitor ExprVisitor) (interface{}, error) {
	return visitor.VisitLiteralExpr(e)
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

func (e UnaryExpr) Accept(visitor ExprVisitor) (interface{}, error) {
	return visitor.VisitUnaryExpr(e)
}

type VariableExpr struct {
	name Token
}

func MakeVariableExpr(name Token) VariableExpr {
	return VariableExpr{
		name: name,
	}
}

func (e VariableExpr) Accept(visitor ExprVisitor) (interface{}, error) {
	return visitor.VisitVariableExpr(e)
}
