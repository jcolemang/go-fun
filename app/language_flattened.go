package main

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
    Return *FlatExpr
}

type FlatExpr struct {
	Num *Num
    Var *Var
    App []*FlatExpr
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
		return FlatExprToString(statement.Expr)
	case statement.Assignment != nil:
		return VarToString(statement.Assignment.Ref) + " = " + FlatExprToString(statement.Assignment.Expr)
	case statement.Return != nil:
		return "return " + FlatExprToString(statement.Return)
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
