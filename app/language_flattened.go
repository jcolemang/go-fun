package main

import (
)

// Main Language
type FlatProgram struct {
	Statements []*FlatStatement
}

type FlatStatement struct {
	Expr *FlatExpr
	Assignment *FlatAssignment
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

