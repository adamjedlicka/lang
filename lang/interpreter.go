package lang

type Interpreter struct {
	expr Expr
}

func MakeInterpreter(expr Expr) Interpreter {
	i := Interpreter{}
	i.expr = expr

	return i
}

func (i *Interpreter) Interpret() interface{} {
	return i.evaluate(i.expr)
}

func (i *Interpreter) VisitLiteralExpr(expr Literal) interface{} {
	return expr.value
}

func (i *Interpreter) VisitGroupingExpr(expr Grouping) interface{} {
	return i.evaluate(expr.expression)
}

func (i *Interpreter) VisitBinaryExpr(expr Binary) interface{} {
	left := i.evaluate(expr.left)
	right := i.evaluate(expr.right)

	switch expr.operator.tokenType {
	case Minus:
		return left.(float64) - right.(float64)
	case Slash:
		return left.(float64) / right.(float64)
	case Star:
		return left.(float64) * right.(float64)
	case Plus:
		if left, ok := left.(float64); ok {
			if right, ok := right.(float64); ok {
				return left + right
			}
		}

		if left, ok := left.(string); ok {
			if right, ok := right.(string); ok {
				return left + right
			}
		}
	case Greater:
		return left.(float64) > right.(float64)
	case GreaterEqual:
		return left.(float64) >= right.(float64)
	case Less:
		return left.(float64) < right.(float64)
	case LessEqual:
		return left.(float64) <= right.(float64)
	case BangEqual:
		return !i.isEqual(left, right)
	case EqualEqual:
		return i.isEqual(left, right)
	}

	panic("UNREACHABLE")
}

func (i *Interpreter) VisitUnaryExpr(expr Unary) interface{} {
	right := i.evaluate(expr.right)

	switch expr.operator.tokenType {
	case Minus:
		return -right.(float64)
	case Bang:
		return !i.isTruthy(right)
	}

	panic("UNREACHABLE")
}

func (i *Interpreter) evaluate(expr Expr) interface{} {
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
