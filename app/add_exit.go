package main

import (
    "errors"
)

// main reason to add here is just because it is the lowest level "language" with 
// basic assignments where naked expressions don't matter. Sure this would also 
// work elsewhere but this seems most reasonable at the moment
func AddExitVariable(prog *SimpleProgram, getVar func() *Var) (*SimpleExitProgram, error) {	
	statements := prog.Statements
	last, statements := statements[len(statements) - 1], statements[:len(statements) - 1]
	var exit *Var
	var newLast *SimpleStatement

	if last.Assignment != nil {
		exit = last.Assignment.Ref
		newLast = last
	} else if last.Expr != nil {
		exit = getVar()
		newLast = &SimpleStatement{
			Assignment: &SimpleAssignment{
				Ref: exit,
				Expr: last.Expr,
			},
		}
	} else {
		return nil, errors.New("Unrecognized SimpleStatement")
	}
	
	return &SimpleExitProgram{
		Statements: append(statements, newLast),
		Exit: exit,
	}, nil
}