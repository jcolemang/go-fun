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
    labelCounter := 0
    getLabel := func() string {
        labelCounter += 1
        return "label" + fmt.Sprint(labelCounter)
    }
    firstBlock, restBlocks, err := HelpFormBlocks(prog.Statements, []languages.IBlockStatement{}, getLabel)
    if err != nil {
        return nil, err
    }
    blocks := append([]languages.IBlock{firstBlock}, restBlocks...)
    blockProg := languages.BlockProgram{
        Blocks: blocks,
    }
    return &blockProg, nil
}

func HelpFormBlocks(stmts []*languages.SimpleStatement, currentBlockStmts []languages.IBlockStatement, getLabel func() string) (languages.IBlock, []languages.IBlock, error) {
    if len(stmts) == 0 {
        return nil, nil, errors.New("Error forming blocks: found no block terminator in HelpFormBlocks")
    }

    stmt, rest := stmts[0], stmts[1:]

    switch {
    case stmt.Expr != nil:
        bExpr, err := SimpleExprToBlockExpr(stmt.Expr)
        if err != nil {
            return nil, nil, err
        }
        return HelpFormBlocks(rest, append(currentBlockStmts, bExpr), getLabel)
    case stmt.Assignment != nil:
        bExpr, err := SimpleExprToBlockExpr(stmt.Assignment.Expr)
        if err != nil {
            return nil, nil, err
        }
        assignment := languages.Assignment[languages.BlockExpr]{
            Ref: stmt.Assignment.Ref,
            Expr: bExpr,
        }
        return HelpFormBlocks(rest, append(currentBlockStmts, assignment), getLabel)
    case stmt.Return != nil:
        bExpr, err := SimpleExprToBlockExpr(stmt.Return)
        if err != nil {
            return nil, nil, err
        }
        if len(rest) != 0 {
            return nil, nil, errors.New("Found unreachable statements in HelpFormBlocks")
        }
        blockLabel := getLabel()
        blockTerminator := languages.BlockReturn{
            Val: *bExpr,
        }
        block := languages.BasicBlock{
            Label: blockLabel,
            Statements: currentBlockStmts,
            Terminator: blockTerminator,
        }
        return block, []languages.IBlock{}, nil
    case stmt.IfStmt != nil:
        _, _, err := IfStmtToBlocks(stmt.IfStmt, getLabel)
        if err != nil {
            return nil, nil, err
        }
        return nil, nil, errors.New("Working on it")
    }
    return nil, nil, errors.New("IDK what happened in HelpFormBlocks")
}

func IfStmtToBlocks(stmt *languages.SimpleIfStmt, getLabel func() string) (languages.IBlock, []languages.IBlock, error) {
    firstTrueBlock, restTrueBlocks, err := HelpFormBlocks(stmt.IfTrue, []languages.IBlockStatement{}, getLabel)
    if err != nil {
        return nil, nil, err
    }
    firstFalseBlock, restFalseBlocks, err := HelpFormBlocks(stmt.IfFalse, []languages.IBlockStatement{}, getLabel)
    if err != nil {
        return nil, nil, err
    }
    ifBlock := languages.IfBlock{
        IfCond: *stmt.IfCond,
        IfTrue: firstTrueBlock.BlockLabel(),
        IfFalse: firstFalseBlock.BlockLabel(),
    }
    finalBlocks := append([]languages.IBlock{firstTrueBlock}, restTrueBlocks...)
    finalBlocks = append(finalBlocks, firstFalseBlock)
    finalBlocks = append(finalBlocks, restFalseBlocks...)
    return ifBlock, finalBlocks, nil
}

func SimpleExprToBlockExpr(e *languages.SimpleExpr) (*languages.BlockExpr, error) {
    switch {
    case e.Primitive != nil:
        return &languages.BlockExpr{Expr: *e.Primitive}, nil
    case e.App != nil:
        return &languages.BlockExpr{Expr: *e.App}, nil
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
