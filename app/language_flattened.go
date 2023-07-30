package main

import (
)

// Main Language
type FlatProgram struct {
	Exprs []*FlatExpr
}

type FlatExpr struct {
	Num *Num
    Var *Var
	Assignment *FlatAssignment
    App []*FlatExpr
}

type FlatAssignment struct {
	Ref *Var
	Expr *FlatExpr
}

