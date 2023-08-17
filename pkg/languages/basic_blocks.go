package languages

import (
)

type BlockProgram struct {
    Blocks []IBlock
}

type IBlock interface {
    BlockLabel() string
}
func (b BasicBlock) BlockLabel() string {return b.Label}
func (b IfBlock) BlockLabel() string {return b.Label}

type BasicBlock struct {
    Label string
    Statements []IBlockStatement
    Terminator IBlockTerminator
}

type IBlockTerminator interface {
    IsBlockTerminator()
}

func (j BlockJump) IsBlockTerminator() {}
func (r BlockReturn) IsBlockTerminator() {}

type BlockJump struct {
    Label string
}

type BlockReturn struct {
    Val BlockExpr
}

type IfBlock struct {
    Label string
    IfCond Primitive
    IfTrue string
    IfFalse string
}

type IBlockStatement interface {
    IsBlockStatement()
}

func (e Assignment[BlockExpr]) IsBlockStatement() {}
func (e BlockExpr) IsBlockStatement() {}
func (e Goto) IsBlockStatement() {}

type Goto struct {
    Label string
}

type IBlockExpr interface {
    IsBlockExpr()
}
func (p Primitive) IsBlockExpr() {}
func (p PrimitiveApplication) IsBlockExpr() {}

type BlockExpr struct {
    Expr IBlockExpr
}

