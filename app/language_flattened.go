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
	Assignment *Assignment
    App []*FlatExpr
}

