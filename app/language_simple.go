package main

import (
)

// Main Language
type SimpleProgram struct {
	Statements []*SimpleStatement
}

type SimpleStatement struct {
	Expr *SimpleExpr
	Assignment *SimpleAssignment
}

type SimpleExpr struct {
	Num *Num
    Var *Var
    App *SimpleApplication
}

type SimpleApplication struct {
	Operator *Var
	Operands []*SimpleExpr
}

type SimpleAssignment struct {
	Ref *Var
	Expr *SimpleExpr
}