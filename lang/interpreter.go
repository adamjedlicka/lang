package lang

import (
	"fmt"
	"strconv"
)

type Interpreter struct {
	globals *Env
	env     *Env
	stmnts  []Stmnt
}

func MakeInterpreter() Interpreter {
	env := MakeEnv(nil)

	env.values["time"] = Time{}

	return Interpreter{
		globals: env,
		env:     env,
	}
}

func (i *Interpreter) Interpret(stmnts []Stmnt) error {
	i.stmnts = stmnts

	for _, stmnt := range i.stmnts {
		err := stmnt.Accept(i)
		if err != nil {
			return err
		}
	}

	return nil
}

func (i *Interpreter) VisitLiteralExpr(expr LiteralExpr) (interface{}, error) {
	return expr.value, nil
}

func (i *Interpreter) VisitGroupingExpr(expr GroupingExpr) (interface{}, error) {
	return i.evaluate(expr.expression)
}

func (i *Interpreter) VisitBinaryExpr(expr BinaryExpr) (interface{}, error) {
	left, err := i.evaluate(expr.left)
	if err != nil {
		return nil, err
	}

	right, err := i.evaluate(expr.right)
	if err != nil {
		return nil, err
	}

	switch expr.operator.tokenType {
	case Minus:
		err := i.checkNumberOperands(expr.operator, left, right)
		if err != nil {
			return nil, err
		}

		return left.(float64) - right.(float64), nil
	case Slash:
		err := i.checkNumberOperands(expr.operator, left, right)
		if err != nil {
			return nil, err
		}

		return left.(float64) / right.(float64), nil
	case Star:
		err := i.checkNumberOperands(expr.operator, left, right)
		if err != nil {
			return nil, err
		}

		return left.(float64) * right.(float64), nil
	case Plus:
		if left, ok := left.(float64); ok {
			if right, ok := right.(float64); ok {
				return left + right, nil
			}
		}

		if left, ok := left.(string); ok {
			if right, ok := right.(string); ok {
				return left + right, nil
			}

			return fmt.Sprintf("%s%v", left, right), nil
		}

		return nil, NewRuntimeError(expr.operator.line, "Operands must be two numbers or two strings.")
	case Greater:
		err := i.checkNumberOperands(expr.operator, left, right)
		if err != nil {
			return nil, err
		}

		return left.(float64) > right.(float64), nil
	case GreaterEqual:
		err := i.checkNumberOperands(expr.operator, left, right)
		if err != nil {
			return nil, err
		}

		return left.(float64) >= right.(float64), nil
	case Less:
		err := i.checkNumberOperands(expr.operator, left, right)
		if err != nil {
			return nil, err
		}

		return left.(float64) < right.(float64), nil
	case LessEqual:
		err := i.checkNumberOperands(expr.operator, left, right)
		if err != nil {
			return nil, err
		}

		return left.(float64) <= right.(float64), nil
	case BangEqual:
		return !i.isEqual(left, right), nil
	case EqualEqual:
		return i.isEqual(left, right), nil
	}

	return nil, NewRuntimeError(expr.operator.line, "Error while evaluating binary operand.")
}

func (i *Interpreter) VisitCallExpr(expr CallExpr) (interface{}, error) {
	callee, err := i.evaluate(expr.callee)
	if err != nil {
		return nil, err
	}

	arguments := make([]interface{}, 0)
	for _, expr := range expr.arguments {
		argument, err := i.evaluate(expr)
		if err != nil {
			return nil, err
		}

		arguments = append(arguments, argument)
	}

	function, ok := callee.(Callable)
	if !ok {
		return nil, NewRuntimeError(expr.paren.line, "Can only call functions and classes.")
	}

	if function.Arity() != len(arguments) {
		return nil, NewRuntimeError(expr.paren.line,
			fmt.Sprintf("Expected %d arguments but got %d.", function.Arity(), len(arguments)))
	}

	return function.Call(i, arguments)
}

func (i *Interpreter) VisitLogicalExpr(expr LogicalExpr) (interface{}, error) {
	left, err := i.evaluate(expr.left)
	if err != nil {
		return nil, err
	}

	if expr.operator.tokenType == Or {
		if i.isTruthy(left) {
			return left, nil
		}
	} else {
		if !i.isTruthy(left) {
			return left, nil
		}
	}

	return i.evaluate(expr.right)
}

func (i *Interpreter) VisitUnaryExpr(expr UnaryExpr) (interface{}, error) {
	right, err := i.evaluate(expr.right)
	if err != nil {
		return nil, err
	}

	switch expr.operator.tokenType {
	case Minus:
		err := i.checkNumberOperand(expr.operator, right)
		if err != nil {
			return nil, err
		}

		return -right.(float64), nil
	case Bang:
		return !i.isTruthy(right), nil
	}

	return nil, NewRuntimeError(expr.operator.line, "Error while evaluating unary operand.")
}

func (i *Interpreter) VisitVariableExpr(expr VariableExpr) (interface{}, error) {
	return i.env.Get(expr.name)
}

func (i *Interpreter) VisitAssignExpr(expr AssignExpr) (interface{}, error) {
	value, err := i.evaluate(expr.value)
	if err != nil {
		return nil, err
	}

	err = i.env.Assign(expr.name, value)
	if err != nil {
		return nil, err
	}

	return value, nil
}

func (i *Interpreter) VisitLambdaExpr(expr LambdaExpr) (interface{}, error) {
	return MakeLambda(expr, i.env), nil
}

func (i *Interpreter) VisitBlockStmnt(stmnt BlockStmnt) error {
	return i.executeBlock(stmnt.stmnts, MakeEnv(i.env))
}

func (i *Interpreter) VisitExpressionStmnt(stmnt ExpressionStmnt) error {
	_, err := i.evaluate(stmnt.expr)

	return err
}

func (i *Interpreter) VisitFnStmnt(stmnt FnStmnt) error {
	function := MakeFunction(stmnt, i.env)

	return i.env.Define(stmnt.name, function)
}

func (i *Interpreter) VisitIfStmnt(stmnt IfStmnt) error {
	condition, err := i.evaluate(stmnt.condition)
	if err != nil {
		return err
	}

	if i.isTruthy(condition) {
		err = stmnt.thenBranch.Accept(i)
		if err != nil {
			return err
		}
	} else if stmnt.elseBranch != nil {
		err = stmnt.elseBranch.Accept(i)
		if err != nil {
			return err
		}
	}

	return nil
}

func (i *Interpreter) VisitPrintStmnt(stmnt PrintStmnt) error {
	value, err := i.evaluate(stmnt.expr)
	if err != nil {
		return err
	}

	fmt.Println(i.Stringify(value))

	return nil
}

func (i *Interpreter) VisitVarStmnt(stmnt VarStmnt) error {
	var err error
	var value interface{}

	if stmnt.initializer != nil {
		value, err = i.evaluate(stmnt.initializer)
		if err != nil {
			return err
		}
	}

	return i.env.Define(stmnt.name, value)
}

func (i *Interpreter) VisitReturnStmnt(stmnt ReturnStmnt) error {
	var err error
	var value interface{}

	if stmnt.value != nil {
		value, err = i.evaluate(stmnt.value)
		if err != nil {
			return err
		}
	}

	return MakeReturner(value)
}

func (i *Interpreter) VisitWhileStmnt(stmnt WhileStmnt) error {
	for {
		condition, err := i.evaluate(stmnt.condition)
		if err != nil {
			return err
		}

		if !i.isTruthy(condition) {
			return nil
		}

		err = stmnt.body.Accept(i)
		if err != nil {
			return err
		}
	}
}

func (i *Interpreter) evaluate(expr Expr) (interface{}, error) {
	return expr.Accept(i)
}

func (i *Interpreter) executeBlock(stmnts []Stmnt, env *Env) error {
	previous := i.env

	i.env = env

	for _, stmnt := range stmnts {
		err := stmnt.Accept(i)
		if err != nil {
			i.env = previous
			return err
		}
	}

	i.env = previous

	return nil
}

func (i *Interpreter) isTruthy(value interface{}) bool {
	switch value := value.(type) {
	case nil:
		return false
	case bool:
		return value
	}

	return true
}

func (i *Interpreter) isEqual(left, right interface{}) bool {
	if left == nil && right == nil {
		return true
	}

	return left == right
}

func (i *Interpreter) Stringify(value interface{}) string {
	switch value := value.(type) {
	case nil:
		return "null"
	case float64:
		return strconv.FormatFloat(value, 'f', -1, 64)
	}

	return fmt.Sprintf("%v", value)
}

func (i *Interpreter) checkNumberOperand(operator Token, operand interface{}) error {
	_, ok := operand.(float64)
	if ok {
		return nil
	}

	return NewRuntimeError(operator.line, "Operand must be a number.")
}

func (i *Interpreter) checkNumberOperands(operator Token, left, right interface{}) error {
	_, okLeft := left.(float64)
	_, okRight := right.(float64)
	if okLeft && okRight {
		return nil
	}

	return NewRuntimeError(operator.line, "Operands must be a numbers.")
}
