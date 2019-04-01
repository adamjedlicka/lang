package lang

import (
	"fmt"
	"strings"
)

type AstPrinter struct {
}

func MakeAstPrinter() AstPrinter {
	return AstPrinter{}
}

func (ap AstPrinter) Print(expr Expr) string {
	text, err := expr.Accept(ap)
	if err != nil {
		return fmt.Sprintf("%v", err)
	}

	return fmt.Sprintf("%v", text)
}

func (ap AstPrinter) VisitBinaryExpr(expr Binary) (interface{}, error) {
	return ap.parenthesize(expr.operator.lexeme, expr.left, expr.right)
}

func (ap AstPrinter) VisitGroupingExpr(expr Grouping) (interface{}, error) {
	return ap.parenthesize("group", expr.expression)
}

func (ap AstPrinter) VisitLiteralExpr(expr Literal) (interface{}, error) {
	if expr.value == nil {
		return "null", nil
	}

	return expr.value, nil
}

func (ap AstPrinter) VisitUnaryExpr(expr Unary) (interface{}, error) {
	return ap.parenthesize(expr.operator.lexeme, expr.right)
}

func (ap AstPrinter) parenthesize(name string, exprs ...Expr) (string, error) {
	sb := strings.Builder{}

	sb.WriteString("(")
	sb.WriteString(name)
	for _, expr := range exprs {
		sb.WriteString(" ")

		text, err := expr.Accept(ap)
		if err != nil {
			return "", err
		}

		sb.WriteString(fmt.Sprintf("%v", text))
	}
	sb.WriteString(")")

	return sb.String(), nil
}
