package lang

type Stmnt interface {
	Accept(StmntVisitor) error
}

type StmntVisitor interface {
	VisitBlockStmnt(BlockStmnt) error
	VisitClassStmnt(ClassStmnt) error
	VisitExpressionStmnt(ExpressionStmnt) error
	VisitFnStmnt(FnStmnt) error
	VisitIfStmnt(IfStmnt) error
	VisitPrintStmnt(PrintStmnt) error
	VisitVarStmnt(VarStmnt) error
	VisitReturnStmnt(ReturnStmnt) error
	VisitWhileStmnt(WhileStmnt) error
}

type BlockStmnt struct {
	stmnts []Stmnt
}

func MakeBlockStmnt(stmnts []Stmnt) BlockStmnt {
	return BlockStmnt{
		stmnts: stmnts,
	}
}

func (s BlockStmnt) Accept(visitor StmntVisitor) error {
	return visitor.VisitBlockStmnt(s)
}

type ClassStmnt struct {
	name         Token
	declarations []VarStmnt
	methods      []FnStmnt
}

func MakeClassStmnt(name Token, declarations []VarStmnt, methods []FnStmnt) ClassStmnt {
	return ClassStmnt{
		name:         name,
		declarations: declarations,
		methods:      methods,
	}
}

func (s ClassStmnt) Accept(visitor StmntVisitor) error {
	return visitor.VisitClassStmnt(s)
}

type ExpressionStmnt struct {
	expr Expr
}

func MakeExpressionStmnt(expr Expr) ExpressionStmnt {
	return ExpressionStmnt{
		expr: expr,
	}
}

func (s ExpressionStmnt) Accept(visitor StmntVisitor) error {
	return visitor.VisitExpressionStmnt(s)
}

type FnStmnt struct {
	name   Token
	params []Token
	body   []Stmnt
}

func MakeFnStmnt(name Token, params []Token, body []Stmnt) FnStmnt {
	return FnStmnt{
		name:   name,
		params: params,
		body:   body,
	}
}

func (s FnStmnt) Accept(visitor StmntVisitor) error {
	return visitor.VisitFnStmnt(s)
}

type IfStmnt struct {
	condition  Expr
	thenBranch Stmnt
	elseBranch Stmnt
}

func MakeIfStmnt(condition Expr, thenBranch, elseBranch Stmnt) IfStmnt {
	return IfStmnt{
		condition:  condition,
		thenBranch: thenBranch,
		elseBranch: elseBranch,
	}
}

func (s IfStmnt) Accept(visitor StmntVisitor) error {
	return visitor.VisitIfStmnt(s)
}

type PrintStmnt struct {
	expr Expr
}

func MakePrintStmnt(expr Expr) PrintStmnt {
	return PrintStmnt{
		expr: expr,
	}
}

func (s PrintStmnt) Accept(visitor StmntVisitor) error {
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

func (s VarStmnt) Accept(visitor StmntVisitor) error {
	return visitor.VisitVarStmnt(s)
}

type ReturnStmnt struct {
	keyword Token
	value   Expr
}

func MakeReturnStmnt(keyword Token, value Expr) ReturnStmnt {
	return ReturnStmnt{
		keyword: keyword,
		value:   value,
	}
}

func (s ReturnStmnt) Accept(visitor StmntVisitor) error {
	return visitor.VisitReturnStmnt(s)
}

type WhileStmnt struct {
	condition Expr
	body      Stmnt
}

func MakeWhileStmnt(condition Expr, body Stmnt) WhileStmnt {
	return WhileStmnt{
		condition: condition,
		body:      body,
	}
}

func (s WhileStmnt) Accept(visitor StmntVisitor) error {
	return visitor.VisitWhileStmnt(s)
}
