package main

import (
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