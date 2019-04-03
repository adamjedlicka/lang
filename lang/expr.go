package lang

type Expr interface {
	Accept(ExprVisitor) (interface{}, error)
}

type ExprVisitor interface {
	VisitAssignExpr(AssignExpr) (interface{}, error)
	VisitBinaryExpr(BinaryExpr) (interface{}, error)
	VisitCallExpr(CallExpr) (interface{}, error)
	VisitGroupingExpr(GroupingExpr) (interface{}, error)
	VisitLambdaExpr(LambdaExpr) (interface{}, error)
	VisitLiteralExpr(LiteralExpr) (interface{}, error)
	VisitLogicalExpr(LogicalExpr) (interface{}, error)
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

type CallExpr struct {
	callee    Expr
	paren     Token
	arguments []Expr
}

func MakeCallExpr(callee Expr, paren Token, arguments []Expr) CallExpr {
	return CallExpr{
		callee:    callee,
		paren:     paren,
		arguments: arguments,
	}
}

func (e CallExpr) Accept(visitor ExprVisitor) (interface{}, error) {
	return visitor.VisitCallExpr(e)
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

type LambdaExpr struct {
	params []Token
	body   []Stmnt
}

func MakeLambdaExpr(params []Token, body []Stmnt) LambdaExpr {
	return LambdaExpr{
		params: params,
		body:   body,
	}
}

func (e LambdaExpr) Accept(visitor ExprVisitor) (interface{}, error) {
	return visitor.VisitLambdaExpr(e)
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

type LogicalExpr struct {
	operator Token
	left     Expr
	right    Expr
}

func MakeLogicalExpr(operator Token, left, right Expr) LogicalExpr {
	return LogicalExpr{
		operator: operator,
		left:     left,
		right:    right,
	}
}

func (e LogicalExpr) Accept(visitor ExprVisitor) (interface{}, error) {
	return visitor.VisitLogicalExpr(e)
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
