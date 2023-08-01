package main

import (
)

// necessary? no. but I like it.
// this can almost certainly be removed, I don't think I need the exit variable here
type SimpleExitProgram struct {
	Statements []*SimpleStatement
	Exit *Var
}

func SimpleExitProgramToString(prog *SimpleExitProgram) string {
	var str string
	for _, s := range(prog.Statements) {
		str = str + SimpleStatementToString(s) + "\n"
	}
	return str
}