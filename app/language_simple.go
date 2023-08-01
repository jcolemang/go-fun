package main

import (
	"fmt"
)

// Simplifies language with significant constraints around
// application and assignment

type SimpleProgram struct {
	Statements []*SimpleStatement
}

type SimpleStatement struct {
	Expr *SimpleExpr
	Assignment *SimpleAssignment
}

type SimpleExpr struct {
	Primitive *SimplePrimitive
    App *SimpleApplication
}

type SimplePrimitive struct {
	Num *Num
	Var *Var
}

type SimpleApplication struct {
	Operator *Var
	Operands []*SimplePrimitive
}

type SimpleAssignment struct {
	Ref *Var
	Expr *SimpleExpr
}

func SimpleProgramToString(prog *SimpleProgram) string {
	var str string
	for _, s := range(prog.Statements) {
		str = str + SimpleStatementToString(s) + "\n"
	}
	return str
}

func SimpleStatementToString(statement *SimpleStatement) string {
	switch {
	case statement.Expr != nil:
		return SimpleExprToString(statement.Expr)
	case statement.Assignment != nil:
		return VarToString(statement.Assignment.Ref) + " = " + SimpleExprToString(statement.Assignment.Expr)
	default:
		return "Got a nonsense statement and I don't want to deal with the error"
	}
}

func SimplePrimitiveToString(prim *SimplePrimitive) string {
	if prim.Num != nil {
		return fmt.Sprint(*prim.Num.Int)
	} else {
		return VarToString(prim.Var)
	}
}

func SimpleExprToString(expr *SimpleExpr) string {
	switch {
	case expr.Primitive != nil:
		return SimplePrimitiveToString(expr.Primitive)
	case expr.App != nil:
		s := "( "
		s = s + VarToString(expr.App.Operator) + " "
		for _, p := range(expr.App.Operands) {
			s = s + SimplePrimitiveToString(p) + " "
		}
		return s + " )"
	default:
		return "Got a nonsense expression and I don't want to deal with the error"
	}
}