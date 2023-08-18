package languages

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
    IfStmt *SimpleIfStmt
	Assignment *Assignment[SimpleExpr]
    Return *SimpleExpr
}

type SimpleExpr struct {
	Primitive *Primitive
    App *PrimitiveApplication
}

type SimpleIfStmt struct {
    IfCond *Primitive
    IfTrue []*SimpleStatement // cannot pull out statements for evaluation in each branch
    IfFalse []*SimpleStatement
}

type SimpleApplication struct {
	Operator *Var
	Operands []*Primitive
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
		return "\t" + SimpleExprToString(statement.Expr)
	case statement.Assignment != nil:
		return "\t" + VarToString(statement.Assignment.Ref) + " = " + SimpleExprToString(statement.Assignment.Expr)
	case statement.Return != nil:
		return "\treturn " + SimpleExprToString(statement.Return)
	default:
		return "Got a nonsense statement and I don't want to deal with the error"
	}
}

func SimplePrimitiveToString(prim *Primitive) string {
	if prim.Num != nil {
		return fmt.Sprint(*prim.Num.Int)
	} else if prim.Bool != nil {
        if prim.Bool.True != nil {
            return "#true"
        } else {
            return "#false"
        }
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
