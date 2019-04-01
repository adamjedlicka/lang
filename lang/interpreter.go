package lang

import (
	"fmt"
)

type Interpreter struct {
	expr Expr
}

func MakeInterpreter(expr Expr) Interpreter {
	i := Interpreter{}
	i.expr = expr

	return i
}

func (i *Interpreter) Interpret() (interface{}, error) {
	return i.evaluate(i.expr)
}

func (i *Interpreter) VisitLiteralExpr(expr Literal) (interface{}, error) {
	return expr.value, nil
}

func (i *Interpreter) VisitGroupingExpr(expr Grouping) (interface{}, error) {
	return i.evaluate(expr.expression)
}

func (i *Interpreter) VisitBinaryExpr(expr Binary) (interface{}, error) {
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

func (i *Interpreter) VisitUnaryExpr(expr Unary) (interface{}, error) {
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

func (i *Interpreter) evaluate(expr Expr) (interface{}, error) {
	return expr.Accept(i)
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
	if value == nil {
		return "null"
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
