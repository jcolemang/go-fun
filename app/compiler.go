package main

import (
	"fmt"

	"github.com/alecthomas/repr"
)

// step 1: Uniquify all variable names
// (+ 1 (let ((x 2)) (+ x (let ((x 3)) x)))))))) ->
// (+ 1 (let ((x_1 2)) (+ x_1 (let ((x_2 3)) x_2)))
// This will be done by Uniquify
// step 2: With unique variable names the lexical scoping rules of lets are no longer really useful
//         and can be replaced with good ol' assignments. This also means that not everything 
//         will still need to be nested expressions
// (+ 1 (let ((x_1 2)) (+ x_1 (let ((x_2 3)) x_2))) ->
// x_1 := 2
// x_2 := 3
// (+ 1 (+ x_1 x_2))
// step 3: Next comes removing the complex expressions as specified in the textbook, and continuing on from there.

// Other possible passes:
// expanding primitive mathematical expressions (+ 1 2 3) -> (+ 1 (+ 2 3)) to make arbitrary numbers of arguments possible

func Compile(prog *Program) (*VarAssemblyProgram, error) {	
	fmt.Println("Initial program")
	repr.Println(prog)

	getVar := GetVarGenerator()

	// turning 
	// (+ 1 (let ((x 2)) (let ((x 3)) x)))
	// into 
	// (+ 1 (let ((x_1 2)) (let ((x_2 3)) x)))
	newProg, err := Uniquify(prog, getVar)
	if err != nil {
		return nil, err
	}

	fmt.Println("Program after Uniquify")
	repr.Println(newProg)

	// turning
	// (+ 1 (let ((x_1 2)) (+ x_1 3)))
	// into 
	// x_1 = 2
	// tmp2 = (+ x_1 3)
	// (+ 1 tmp2)
	flatProg, err := Flatten(newProg)
	if err != nil {
		return nil, err
	}

	fmt.Println("Program after Flatten")
	repr.Println(flatProg)
	// fmt.Println(FlatProgramToString(flatProg))

	// turning
	// (+ 1 (+ 2 (+ 3 4)))
	// into
	// tmp1 = (+ 3 4)
	// tmp2 = (+ 2 tmp1)
	// (+ 1 tmp2)
	simpleProg, err := RemoveComplexOperands(flatProg, getVar)
	if err != nil {
		return nil, err
	}

	fmt.Println("Program after RemoveComplexOperands")
	repr.Println(simpleProg)
	// fmt.Println(FlatProgramToString(simpleProg))

	varAssemblyProg, err := SelectInstructions(simpleProg)
	if err != nil {
		return nil, err
	}

	fmt.Println("Program after SelectInstructions")
	repr.Println(varAssemblyProg)

	return varAssemblyProg, nil
}