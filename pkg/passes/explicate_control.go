package passes

import (
    "errors"
    "language/pkg/languages"
)

// I would like to:
// 1. form basic blocks. Each block starts with a label and ends with a goto.
// 2. Separate ifs from basic blocks so that each if is a comparison and two labels
func FormBlocks(prog *languages.SimpleProgram) (*languages.BlockProgram, error) {
    return nil, errors.New("Still working on it!")
}

func FormBlocksStmts(stmts []*languages.SimpleStatement, getLabel func() string) (error) {
    startLabel, endLabel := getLabel(), getLabel()
    var blockStmts []*BlockStatement
    for _, s := range stmts {

    }
    return errors.New("Still working on it")
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
