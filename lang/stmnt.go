package lang

type Stmnt interface {
	Accept(StmntVisitor) (interface{}, error)
}

type StmntVisitor interface {
	VisitBlockStmnt(BlockStmnt) (interface{}, error)
	VisitExpressionStmnt(ExpressionStmnt) (interface{}, error)
	VisitPrintStmnt(PrintStmnt) (interface{}, error)
	VisitVarStmnt(VarStmnt) (interface{}, error)
}

type BlockStmnt struct {
	stmnts []Stmnt
}

func MakeBlockStmnt(stmnts []Stmnt) BlockStmnt {
	return BlockStmnt{
		stmnts: stmnts,
	}
}

func (s BlockStmnt) Accept(visitor StmntVisitor) (interface{}, error) {
	return visitor.VisitBlockStmnt(s)
}

type ExpressionStmnt struct {
	expr Expr
}

func MakeExpressionStmnt(expr Expr) ExpressionStmnt {
	return ExpressionStmnt{
		expr: expr,
	}
}

func (s ExpressionStmnt) Accept(visitor StmntVisitor) (interface{}, error) {
	return visitor.VisitExpressionStmnt(s)
}

type PrintStmnt struct {
	expr Expr
}

func MakePrintStmnt(expr Expr) PrintStmnt {
	return PrintStmnt{
		expr: expr,
	}
}

func (s PrintStmnt) Accept(visitor StmntVisitor) (interface{}, error) {
	return visitor.VisitPrintStmnt(s)
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

func (s VarStmnt) Accept(visitor StmntVisitor) (interface{}, error) {
	return visitor.VisitVarStmnt(s)
}
