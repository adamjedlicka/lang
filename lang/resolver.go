package lang

type functionType uint8
type classType uint8

const (
	functionNone functionType = iota
	functionFunction
	functionMethod
	functionLambda
	functionInit
)

const (
	classNone classType = iota
	classClass
)

type Resolver struct {
	interpreter     *Interpreter
	scopes          []map[string]bool
	currentFunction functionType
	currentClass    classType
}

func MakeResolver(interpreter *Interpreter) Resolver {
	return Resolver{
		interpreter:     interpreter,
		scopes:          make([]map[string]bool, 0),
		currentFunction: functionNone,
		currentClass:    classNone,
	}
}

func (r *Resolver) Resolve(stmnts []Stmnt) error {
	return r.resolveStmnts(stmnts)
}

func (r *Resolver) VisitBlockStmnt(stmnt BlockStmnt) error {
	r.beginScope()

	err := r.resolveStmnts(stmnt.stmnts)
	if err != nil {
		return err
	}

	r.endScope()

	return nil
}

func (r *Resolver) VisitClassStmnt(stmnt ClassStmnt) error {
	enclosingClass := r.currentClass
	r.currentClass = classClass

	err := r.declare(stmnt.name)
	if err != nil {
		return err
	}

	r.define(stmnt.name)

	r.beginScope()

	r.scope()["this"] = true

	for _, declaration := range stmnt.declarations {
		err := r.resolveStmnt(declaration)
		if err != nil {
			return err
		}
	}

	for _, method := range stmnt.methods {
		declaration := functionMethod

		if method.name.lexeme == "init" {
			declaration = functionInit
		}

		err := r.resolveFunction(method, declaration)
		if err != nil {
			return err
		}
	}

	r.endScope()

	r.currentClass = enclosingClass

	return nil
}

func (r *Resolver) VisitExpressionStmnt(stmnt ExpressionStmnt) error {
	return r.resolveExpr(stmnt.expr)
}

func (r *Resolver) VisitFnStmnt(stmnt FnStmnt) error {
	err := r.declare(stmnt.name)
	if err != nil {
		return err
	}

	r.define(stmnt.name)

	return r.resolveFunction(stmnt, functionFunction)
}

func (r *Resolver) VisitIfStmnt(stmnt IfStmnt) error {
	err := r.resolveExpr(stmnt.condition)
	if err != nil {
		return err
	}

	err = r.resolveStmnt(stmnt.thenBranch)
	if err != nil {
		return err
	}

	if stmnt.elseBranch != nil {
		err = r.resolveStmnt(stmnt.elseBranch)
		if err != nil {
			return err
		}
	}

	return nil
}

func (r *Resolver) VisitPrintStmnt(stmnt PrintStmnt) error {
	return r.resolveExpr(stmnt.expr)
}

func (r *Resolver) VisitVarStmnt(stmnt VarStmnt) error {
	err := r.declare(stmnt.name)
	if err != nil {
		return err
	}

	if stmnt.initializer != nil {
		err := r.resolveExpr(stmnt.initializer)
		if err != nil {
			return err
		}
	}

	r.define(stmnt.name)

	return nil
}

func (r *Resolver) VisitReturnStmnt(stmnt ReturnStmnt) error {
	if r.currentFunction == functionNone {
		return NewResolverError(stmnt.keyword, "Cannot return from top-level code.")
	}

	if stmnt.value != nil && r.currentFunction == functionInit {
		return NewResolverError(stmnt.keyword, "Cannot return a value from an initializer.")
	}

	if stmnt.value != nil {
		err := r.resolveExpr(stmnt.value)
		if err != nil {
			return err
		}
	}

	return nil
}

func (r *Resolver) VisitWhileStmnt(stmnt WhileStmnt) error {
	err := r.resolveExpr(stmnt.condition)
	if err != nil {
		return err
	}

	err = r.resolveStmnt(stmnt.body)
	if err != nil {
		return err
	}

	return nil
}

func (r *Resolver) VisitAssignExpr(expr AssignExpr) (interface{}, error) {
	err := r.resolveExpr(expr.value)
	if err != nil {
		return nil, err
	}

	r.resolveLocal(expr, expr.name)

	return nil, nil
}

func (r *Resolver) VisitBinaryExpr(expr BinaryExpr) (interface{}, error) {
	err := r.resolveExpr(expr.left)
	if err != nil {
		return nil, err
	}

	err = r.resolveExpr(expr.right)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (r *Resolver) VisitCallExpr(expr CallExpr) (interface{}, error) {
	err := r.resolveExpr(expr.callee)
	if err != nil {
		return nil, err
	}

	for _, argument := range expr.arguments {
		err := r.resolveExpr(argument)
		if err != nil {
			return nil, err
		}
	}

	return nil, nil
}

func (r *Resolver) VisitGetExpr(expr GetExpr) (interface{}, error) {
	return nil, r.resolveExpr(expr.object)
}

func (r *Resolver) VisitGroupingExpr(expr GroupingExpr) (interface{}, error) {
	return nil, r.resolveExpr(expr.expression)
}

func (r *Resolver) VisitLambdaExpr(expr LambdaExpr) (interface{}, error) {
	return nil, r.resolveLambda(expr)
}

func (r *Resolver) VisitLiteralExpr(expr LiteralExpr) (interface{}, error) {
	return nil, nil
}

func (r *Resolver) VisitLogicalExpr(expr LogicalExpr) (interface{}, error) {
	err := r.resolveExpr(expr.left)
	if err != nil {
		return nil, err
	}

	err = r.resolveExpr(expr.right)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (r *Resolver) VisitSetExpr(expr SetExpr) (interface{}, error) {
	err := r.resolveExpr(expr.value)
	if err != nil {
		return nil, err
	}

	err = r.resolveExpr(expr.object)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (r *Resolver) VisitThisExpr(expr ThisExpr) (interface{}, error) {
	if r.currentClass == classNone {
		return nil, NewResolverError(expr.keword, "Cannot use 'this' outside of a class.")
	}

	r.resolveLocal(expr, expr.keword)

	return nil, nil
}

func (r *Resolver) VisitUnaryExpr(expr UnaryExpr) (interface{}, error) {
	return nil, r.resolveExpr(expr.right)
}

func (r *Resolver) VisitVariableExpr(expr VariableExpr) (interface{}, error) {
	if len(r.scopes) != 0 {
		if value, ok := r.scope()[expr.name.lexeme]; ok && !value {
			return nil, NewResolverError(expr.name, "Cannot read local variable in its own initializer.")
		}
	}

	r.resolveLocal(expr, expr.name)

	return nil, nil
}

func (r *Resolver) declare(token Token) error {
	if len(r.scopes) == 0 {
		return nil
	}

	if _, ok := r.scope()[token.lexeme]; ok {
		return NewResolverError(token, "Variable with this name already declared in this scope.")
	}

	r.scope()[token.lexeme] = false

	return nil
}

func (r *Resolver) define(token Token) {
	if len(r.scopes) == 0 {
		return
	}

	r.scope()[token.lexeme] = true
}

func (r *Resolver) beginScope() {
	r.scopes = append(r.scopes, make(map[string]bool))
}

func (r *Resolver) endScope() {
	r.scopes = r.scopes[:len(r.scopes)-1]
}

func (r *Resolver) resolveStmnts(stmnts []Stmnt) error {
	for _, stmnt := range stmnts {
		err := r.resolveStmnt(stmnt)
		if err != nil {
			return err
		}
	}

	return nil
}

func (r *Resolver) resolveStmnt(stmnt Stmnt) error {
	return stmnt.Accept(r)
}

func (r *Resolver) resolveExpr(expr Expr) error {
	_, err := expr.Accept(r)

	return err
}

func (r *Resolver) resolveLocal(expr Expr, name Token) {
	for i := len(r.scopes) - 1; i >= 0; i-- {
		if _, ok := r.scopes[i][name.lexeme]; ok {
			r.interpreter.resolve(expr, len(r.scopes)-1-i)
			return
		}
	}
}

func (r *Resolver) resolveFunction(function FnStmnt, fnType functionType) error {
	enclosingFunction := r.currentFunction
	r.currentFunction = fnType

	r.beginScope()

	for _, param := range function.params {
		err := r.declare(param)
		if err != nil {
			return err
		}

		r.define(param)
	}

	err := r.resolveStmnts(function.body)
	if err != nil {
		return err
	}

	r.endScope()

	r.currentFunction = enclosingFunction

	return nil
}

func (r *Resolver) resolveLambda(lambda LambdaExpr) error {
	enclosingFunction := r.currentFunction
	r.currentFunction = functionLambda

	r.beginScope()

	for _, param := range lambda.params {
		err := r.declare(param)
		if err != nil {
			return err
		}

		r.define(param)
	}

	err := r.resolveStmnts(lambda.body)
	if err != nil {
		return err
	}

	r.endScope()

	r.currentFunction = enclosingFunction

	return nil
}

func (r *Resolver) scope() map[string]bool {
	return r.scopes[len(r.scopes)-1]
}
