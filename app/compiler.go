package main

import (
	"fmt"
	"os"

	"github.com/alecthomas/repr"
)


func CompileToFile(prog *Program, location string) error {
	x86, err := Compile(prog)
	if err != nil {
		return err
	}

	return os.WriteFile(location, []byte(X86ProgramToString(x86)), 0644)
}

func Compile(prog *Program) (*X86Program, error) {
	fmt.Println("Initial program")
	repr.Println(prog)

	getVar := GetVarGenerator()

	// turning
	// (+ 1 (let ((x 2)) (let ((x 3)) x)))
	// into
	// (+ 1 (let ((x_1 2)) (let ((x_2 3)) x)))
	// In other words, takes care of the lexical scoping logic
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
	// could more appropriately be called "remove lets"
	flatProg, err := Flatten(newProg)
	if err != nil {
		return nil, err
	}

	fmt.Println("Program after Flatten")
	fmt.Println(FlatProgramToString(flatProg))

	// turning
	// (+ 1 (+ 2 (+ 3 4)))
	// into
	// tmp1 = (+ 3 4)
	// tmp2 = (+ 2 tmp1)
	// (+ 1 tmp2)
	// In other works, squashes out subexpressions
	simpleProg, err := RemoveComplexOperands(flatProg, getVar)
	if err != nil {
		return nil, err
	}

	fmt.Println("Program after RemoveComplexOperands")
	fmt.Println(SimpleProgramToString(simpleProg))

	// Add a final variable to later be used as an exit code or something
	// honestly I don't think this is strictly necessary
	simpleExitProg, err := AddExitVariable(simpleProg, getVar)
	if err != nil {
		return nil, err
	}

	fmt.Println("Program after AddExitVariable")
	fmt.Println(SimpleExitProgramToString(simpleExitProg))

	// Picks X86 instructions but keeps variables around
	varAssemblyProg, err := SelectInstructions(simpleExitProg)
	if err != nil {
		return nil, err
	}

	fmt.Println("Program after SelectInstructions")
	fmt.Println(VarAssemblyProgramToString(varAssemblyProg))

	// Assigns variables to registers
	assembly, err := AssignRegisters(varAssemblyProg)
	if err != nil {
		return nil, err
	}

	// Removes some invalid and unnecessary instructions
	patchedAssembly := PatchInstructions(assembly)

	return patchedAssembly, nil
}