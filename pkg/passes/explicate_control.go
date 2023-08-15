package passes

import (
    "fmt"
    "errors"
    "language/pkg/languages"
)

// I would like to:
// 1. form basic blocks. Each block starts with a label and ends with a goto.
// 2. Separate ifs from basic blocks so that each if is a comparison and two labels
func FormBlocks(prog *languages.SimpleProgram) (*languages.BlockProgram, error) {
    return nil, errors.New("Still working on it!")
}

func FormNextBlock(stmts []*languages.SimpleStatement, getLabel func() string) (languages.IBlock, error) {
    startLabel, endLabel := getLabel(), getLabel()
    var blockStmts []languages.IBlockStatement
    var terminator languages.IBlockTerminator
    blockIdx := 0

    if len(stmts) == 0 {
        return make([]languages.IBlock, 0), nil
    }

    for i, s := range stmts {
        blockIdx = i
        switch {
        case s.Expr != nil:
            bExpr, err := SimpleExprToBlockExpr(s.Expr)
            if err != nil {
                return nil, err
            }
            blockStmts = append(blockStmts, bExpr)
        case s.Assignment != nil:
            bExpr, err := SimpleExprToBlockExpr(s.Assignment.Expr)
            if err != nil {
                return nil, err
            }
            blockStmts = append(blockStmts, &languages.Assignment[languages.BlockExpr]{
                Ref: s.Assignment.Ref,
                Expr: bExpr,
            })
        case s.Return != nil:
            bExpr, err := SimpleExprToBlockExpr(s.Return)
            if err != nil {
                return nil, err
            }
            terminator = &languages.BlockReturn{
                Val: *bExpr,
            }
        }
    }

    // newBlock := &languages.
    return nil, errors.New("Still working on it")
}

func SimpleExprToBlockExpr(e *languages.SimpleExpr) (*languages.BlockExpr, error) {
    switch {
    case e.Primitive != nil:
        return &languages.BlockExpr{Expr: e.Primitive}, nil
    case e.App != nil:
        return &languages.BlockExpr{Expr: e.App}, nil
    default:
        return nil, errors.New("Could not convert SimpleExpr to BlockExpr")
    }
}

func ExplicateControl() error {
    return errors.New("Still working on it")
}

func NewLabelGenerator() func() string {
    labelNum := 0
    return func() string {
        labelNum = labelNum + 1
		return "label" + fmt.Sprint(labelNum)
    }
}
