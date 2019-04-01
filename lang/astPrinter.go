package lang

import (
	"fmt"
	"strings"
)

type AstPrinter struct {
	stmnts []Stmnt
}

func MakeAstPrinter(stmnts []Stmnt) AstPrinter {
	ap := AstPrinter{}
	ap.stmnts = stmnts

	return ap
}

func (ap AstPrinter) Print() {
	for _, stmnt := range ap.stmnts {
		text, err := stmnt.Accept(ap)
		if err != nil {
			fmt.Println(err)
			return
		}

		fmt.Println(text)
	}
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

func (ap AstPrinter) VisitVariableExpr(expr Variable) (interface{}, error) {
	return expr.name, nil
}

func (ap AstPrinter) VisitExpressionStmnt(stmnt ExpressionStmnt) (interface{}, error) {
	return stmnt.expr.Accept(ap)
}

func (ap AstPrinter) VisitPrintStmnt(stmnt PrintStmnt) (interface{}, error) {
	return ap.parenthesize("PRINT", stmnt.expr)
}

func (ap AstPrinter) VisitVarStmnt(stmnt VarStmnt) (interface{}, error) {
	return fmt.Sprintf("VAR %v = %v", stmnt.name, stmnt.initializer), nil
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
