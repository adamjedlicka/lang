package lang

type Stmnt interface {
	Accept(StmntVisitor) (interface{}, error)
}

type StmntVisitor interface {
	VisitExpressionStmnt(ExpressionStmnt) (interface{}, error)
	VisitPrintStmnt(PrintStmnt) (interface{}, error)
	VisitVarStmnt(VarStmnt) (interface{}, error)
}

type ExpressionStmnt struct {
	expr Expr
}

func MakeExpressionStmnt(expr Expr) ExpressionStmnt {
	return ExpressionStmnt{
		expr: expr,
	}
}

func (e ExpressionStmnt) Accept(visitor StmntVisitor) (interface{}, error) {
	return visitor.VisitExpressionStmnt(e)
}

type PrintStmnt struct {
	expr Expr
}

func MakePrintStmnt(expr Expr) PrintStmnt {
	return PrintStmnt{
		expr: expr,
	}
}

func (e PrintStmnt) Accept(visitor StmntVisitor) (interface{}, error) {
	return visitor.VisitPrintStmnt(e)
}

type VarStmnt struct {
	name        Token
	initializer Expr
}

func MakeVarStmnt(name Token, initializer Expr) VarStmnt {
	return VarStmnt{
		name:        name,
		initializer: initializer,
	}
}

func (e VarStmnt) Accept(visitor StmntVisitor) (interface{}, error) {
	return visitor.VisitVarStmnt(e)
}
