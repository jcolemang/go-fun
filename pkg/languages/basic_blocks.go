package languages

import (
)

type BlockProgram[T Block] struct {
    Blocks []*T
	Statements []*SimpleStatement
}

type Block interface {
    BasicBlock | IfBlock
}

type BasicBlock struct {
    Label string
    Statements []*BlockStatement
    Jump string
}

type IfBlock struct {
    IfCond *Primitive
    IfTrue string
    IfFalse string
}

type BlockStatement struct {
	Expr *BlockExpr
	Assignment *Assignment[BlockExpr]
    Return *BlockExpr
    Goto *string
}

type BlockExpr struct {
	Primitive *Primitive
    App *PrimitiveApplication
}
