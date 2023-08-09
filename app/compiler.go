package main

import (
	"fmt"
	"os"
    "language/pkg/languages"

	"github.com/alecthomas/repr"
)


func CompileToFile(prog *languages.Program, location string, debug bool) error {
	Arm, err := Compile(prog, debug)
	if err != nil {
		return err
	}

	return os.WriteFile(location, []byte(languages.ArmProgramToString(Arm)), 0644)
}

func Compile(prog *languages.Program, debug bool) (*languages.ArmProgram, error) {
    if debug {
        fmt.Println("Initial program")
        repr.Println(prog)
    }

	getVar := languages.GetVarGenerator()

	// turning
	// (+ 1 (let ((x 2)) (let ((x 3)) x)))
	// into
	// (+ 1 (let ((x_1 2)) (let ((x_2 3)) x)))
	// In other words, takes care of the lexical scoping logic
	newProg, err := Uniquify(prog, getVar)
	if err != nil {
		return nil, err
	}

    if debug {
        fmt.Println("Program after Uniquify")
        repr.Println(newProg)
    }

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

    if debug {
        fmt.Println("Program after Flatten")
        fmt.Println(languages.FlatProgramToString(flatProg))
    }

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

    if debug {
        fmt.Println("Program after RemoveComplexOperands")
        fmt.Println(languages.SimpleProgramToString(simpleProg))
    }

	// Picks Arm instructions but keeps variables around
	varAssemblyProg, err := SelectInstructions(simpleProg, getVar)
	if err != nil {
		return nil, err
	}

    if debug {
        fmt.Println("Program after SelectInstructions")
        fmt.Println(languages.VarAssemblyProgramToString(varAssemblyProg))
    }

	// Assigns variables to registers
	assembly, err := AssignRegisters(varAssemblyProg, debug)
	if err != nil {
		return nil, err
	}

    if debug {
        fmt.Println("Program after AssignRegisters")
        fmt.Println(languages.ArmProgramToString(assembly))
    }

	// Removes some invalid and unnecessary instructions
	patchedAssembly := PatchInstructions(assembly)

    if debug {
        fmt.Println("Program after PatchInstructions")
        fmt.Println(languages.ArmProgramToString(assembly))
    }

	return patchedAssembly, nil
}
