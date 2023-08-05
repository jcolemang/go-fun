package main

import (
    "fmt"
)

type VarAssemblyProgram struct {
	Instrs []*VarAssemblyInstr
}

type VarAssemblyInstr struct {
	Add *[3]*VarAssemblyImmediate
	Mov *[2]*VarAssemblyImmediate
    Ret *VarAssemblyImmediate
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
    case instr.Add != nil:
        first, second := instr.Add[0], instr.Add[1]
        return "addq " + VarAssemblyImmediateToString(first) + " " + VarAssemblyImmediateToString(second)
    case instr.Mov != nil:
        first, second := instr.Mov[0], instr.Mov[1]
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
