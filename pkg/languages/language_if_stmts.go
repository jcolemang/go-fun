package languages

import (
)

// Main Language
type IfStmtProgram struct {
	Statements []*IfStmtStatement
}

// how on earth do I name this
type IfStmtStatement struct {
	Expr *FlatExpr
	Assignment *Assignment[IfStmtExpr]
    Return *FlatExpr
    IfStmt *FlatIfStmt
}

type IfStmtExpr struct {
    Bool *Bool
	Num *Num
    Var *Var
    App []*FlatExpr
}

type FlatIfStmt struct {
    IfCond *IfStmtExpr
    IfTrue []*IfStmtStatement
    IfFalse []*IfStmtStatement
}
