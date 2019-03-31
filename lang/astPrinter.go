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
	return expr.Accept(ap).(string)
}

func (ap AstPrinter) VisitBinaryExpr(expr Binary) interface{} {
	return ap.parenthesize(expr.operator.lexeme, expr.left, expr.right)
}

func (ap AstPrinter) VisitGroupingExpr(expr Grouping) interface{} {
	return ap.parenthesize("group", expr.expression)
}

func (ap AstPrinter) VisitLiteralExpr(expr Literal) interface{} {
	if expr.value == nil {
		return "null"
	}

	return expr.value
}

func (ap AstPrinter) VisitUnaryExpr(expr Unary) interface{} {
	return ap.parenthesize(expr.operator.lexeme, expr.right)
}

func (ap AstPrinter) parenthesize(name string, exprs ...Expr) string {
	sb := strings.Builder{}

	sb.WriteString("(")
	sb.WriteString(name)
	for _, expr := range exprs {
		sb.WriteString(" ")

		sb.WriteString(fmt.Sprintf("%v", expr.Accept(ap)))
	}
	sb.WriteString(")")

	return sb.String()
}
