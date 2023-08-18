package languages

import (
    "fmt"
)

type VarAssemblyProgram struct {
    MainLabel string
	Instrs []*VarAssemblyInstr
}

type VarAssemblyInstr struct {
    Label *string
	Add *[3]*VarAssemblyImmediate
	Mov *[2]*VarAssemblyImmediate
    Ret *Ret
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
        first, second, third := instr.Add[0], instr.Add[1], instr.Add[2]
        return "\tadd " + VarAssemblyImmediateToString(first) + " " + VarAssemblyImmediateToString(second) + " " + VarAssemblyImmediateToString(third)
    case instr.Mov != nil:
        first, second := instr.Mov[0], instr.Mov[1]
        return "\tmov " + VarAssemblyImmediateToString(first) + " " + VarAssemblyImmediateToString(second)
    case instr.Ret != nil:
        return "\tret"
    default:
        return "Unrecognized thing and I don't wanna deal"
    }
}

func VarAssemblyImmediateToString(imm *VarAssemblyImmediate) string {
    switch {
        case imm.Var != nil:
            return "tmp" + fmt.Sprint(imm.Var.Generated)
        case imm.Int != nil:
            return "#" + fmt.Sprint(*imm.Int)
        case imm.Register != nil:
            return imm.Register.Name
    }
    return "Another unrecognized thing"
}
