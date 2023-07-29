package main

import (
)

// Main Language
type FlatProgram struct {
	Expr *Expr
}

type FlatExpr struct {
	Num *Num
    Var *Var
	Let *LetExpr
    App []*Expr
}

type FlatLetExpr struct {
	LetAssignments []*Assignment
	LetBody *Expr
}

type FlatAssignment struct {
	Ref *Var
	Expr *Expr
}

type FlatNum struct {
	Value int
}

type FlatVar struct {
    Name string
    Temp int
}

