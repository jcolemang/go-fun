package languages

import (
)

// Main Language
type IfStmtProgram struct {
	Statements []*IfStmtStatement
}

// how on earth do I name this
type IfStmtStatement struct {
	Expr *IfStmtExpr
	Assignment *Assignment[IfStmtExpr]
    Return *IfStmtExpr
    IfStmt *IfStmtIfStmt
}

type IfStmtExpr struct {
    Bool *Bool
	Num *Num
    Var *Var
    App []*IfStmtExpr
}

// just for fun
type IfStmtIfStmt struct {
    IfCond *IfStmtExpr
    IfTrue []*IfStmtStatement
    IfFalse []*IfStmtStatement
}
