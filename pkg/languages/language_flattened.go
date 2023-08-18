package languages

import (
	"fmt"
)

// Main Language
type FlatProgram struct {
	Statements []*FlatStatement
}

type FlatStatement struct {
	Expr *FlatExpr
	Assignment *FlatAssignment
    Return *FlatExpr // there can be only one
}

type FlatExpr struct {
    Bool *Bool
	Num *Num
    Var *Var
    App []*FlatExpr
    IfExpr *FlatIfExpr
}

type FlatIfExpr struct {
    IfCond *FlatExpr
    IfTrue []*FlatStatement
    IfTrueExpr *FlatExpr
    IfFalse []*FlatStatement
    IfFalseExpr *FlatExpr
}

type FlatAssignment struct {
	Ref *Var
	Expr *FlatExpr
}

func FlatProgramToString(prog *FlatProgram) string {
	var str string
	for _, s := range(prog.Statements) {
		str = str + FlatStatementToString(s) + "\n"
	}
	return str
}

func FlatStatementToString(statement *FlatStatement) string {
	switch {
	case statement.Expr != nil:
		return "\t" + FlatExprToString(statement.Expr)
	case statement.Assignment != nil:
		return "\t" + VarToString(statement.Assignment.Ref) + " = " + FlatExprToString(statement.Assignment.Expr)
	case statement.Return != nil:
		return "\treturn " + FlatExprToString(statement.Return)
	default:
		return "Got a nonsense statement and I don't want to deal with the error"
	}
}

func FlatExprToString(expr *FlatExpr) string {
	switch {
	case expr.Num != nil:
		return fmt.Sprint(*expr.Num.Int)
	case expr.Var != nil:
		return VarToString(expr.Var)
	case expr.App != nil:
		s := "( "
		for _, e := range(expr.App) {
			s = s + FlatExprToString(e) + " "
		}
		return s + " )"
	default:
		return "Got a nonsense expression and I don't want to deal with the error"
	}
}
