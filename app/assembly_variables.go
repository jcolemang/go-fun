package main

import (
    "fmt"
)

type VarAssemblyProgram struct {
	Instrs []*VarAssemblyInstr
}

type VarAssemblyInstr struct {
	Addq *[2]*VarAssemblyImmediate // val dest
	Movq *[2]*VarAssemblyImmediate // val dest
}

type VarAssemblyImmediate struct {
	Var *VarAssemblyVar
	Int *int
	Register *Register
}

type VarAssemblyVar struct {
	Generated int
}

func VarAssemblyProgramToString(prog *VarAssemblyProgram) string {
    repr := ""
    for _, i := range(prog.Instrs) {
        repr = repr + VarAssemblyInstrToString(i) + "\n"
    }
    return repr
}

func VarAssemblyInstrToString(instr *VarAssemblyInstr) string {
    switch {
    case instr.Addq != nil:
        first, second := instr.Addq[0], instr.Addq[1]
        return "addq " + VarAssemblyImmediateToString(first) + " " + VarAssemblyImmediateToString(second)
    case instr.Movq != nil:
        first, second := instr.Movq[0], instr.Movq[1]
        return "movq " + VarAssemblyImmediateToString(first) + " " + VarAssemblyImmediateToString(second)
    default:
        return "Unrecognized thing and I don't wanna deal"
    }
}

func VarAssemblyImmediateToString(imm *VarAssemblyImmediate) string {
    switch {
        case imm.Var != nil:
            return "tmp" + fmt.Sprint(imm.Var.Generated)
        case imm.Int != nil:
            return "$" + fmt.Sprint(*imm.Int)
        case imm.Register != nil:
            return "%" + imm.Register.Name
    }
    return "Another unrecognized thing"
}
