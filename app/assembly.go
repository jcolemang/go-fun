package main

import (
	"strconv"
)

// Arm Language
// the parsing is not really necessary and this would be slightly better without it because I could use arrays with specific lengths
// which the parsing library does not allow
type ArmProgram struct {
	ArmDirectives []*ArmDirective
    ArmInstrs []*ArmInstr
}

type ArmDirective struct {
	Name string
	Arg string
}

type ArmInstr struct {
	Label *string
    Add []*ArmArg
	Mov []*ArmArg
	Ret string
}

type ArmArg struct {
    ArmInt *int
    ArmReg *Register
	ArmOffset *int
	ArmOffsetReg *Register
}

func IsRegister(arg ArmArg) bool {
    return arg.ArmReg != nil
}

type Register struct {
	// argument passing: rdi rsi rdx rcx r8 r9
	// caller saved: rax rcx rdx rsi rdi r8 r9 r10 r11 -> the caller needs to save these, the callee can use freely
	// callee saved: rsp rbp rbx r12 r13 r14 r15 -> callee can use these, but must restore them, caller can use freely
	Name string
}

func TempReg() *Register {
    return &Register{
        Name: "x0",
    }
}

func ArmProgramToString(prog *ArmProgram) string {
	s := ""
	for _, dir := range prog.ArmDirectives {
		s += "\t" + ArmDirectiveToString(dir) + "\n"
	}
	for _, instr := range prog.ArmInstrs {
		s += ArmInstrToString(instr) + "\n"
	}
	return s
}

func ArmDirectiveToString(directive *ArmDirective) string {
	return "." + directive.Name + " " + directive.Arg
}

func ArmInstrToString(instr *ArmInstr) string {
	switch {
	case instr.Label != nil:
		return *instr.Label + ":"
	case instr.Add != nil:
		return "\tadd " + ArmArgToString(instr.Add[0]) + ", " + ArmArgToString(instr.Add[1]) + ", " + ArmArgToString(instr.Add[2])
	case instr.Mov != nil:
		return "\tmov " + ArmArgToString(instr.Mov[0]) + ", " + ArmArgToString(instr.Mov[1])
	default:
		return "Haven't implemented print for this one yet"
	}
}

func ArmArgToString(arg *ArmArg) string {
	switch {
	case arg.ArmInt != nil:
		return "#" + strconv.Itoa(*arg.ArmInt)
	case arg.ArmReg != nil:
		return arg.ArmReg.Name
	default:
		// must be a stack location
		return strconv.Itoa(*arg.ArmOffset) + "(%" + arg.ArmOffsetReg.Name + ")"
	}
}

func GetLocation(i int) *ArmArg {
	assignableRegisters := []Register{
		Register{Name: "x9"}, // temporary registers
		Register{Name: "x10"},
		Register{Name: "x11"},
		Register{Name: "x12"},
		Register{Name: "x13"},
		Register{Name: "x14"},
		Register{Name: "x15"},
		Register{Name: "x19"}, // callee saved
		Register{Name: "x20"},
		Register{Name: "x21"},
		Register{Name: "x22"},
		Register{Name: "x23"},
		Register{Name: "x24"},
		Register{Name: "x25"},
		Register{Name: "x26"},
		Register{Name: "x27"},
		Register{Name: "x28"},
	}

	if i < len(assignableRegisters) {
		return &ArmArg{ArmReg: &assignableRegisters[i]}
	} else {
		offset := (i - len(assignableRegisters) + 1) * -8
		return &ArmArg{
			ArmOffset: &offset,
			ArmOffsetReg: &Register{Name: "sp"}, // base pointer
		}
	}
}

func GetArgumentRegisters() []*Register {
    // 0 - 7
	return []*Register{
		&Register{Name: "rdi"},
		&Register{Name: "rsi"},
		&Register{Name: "rdx"},
		&Register{Name: "rcx"},
		&Register{Name: "r8"},
		&Register{Name: "r9"},
	}
}

func GetCalleeClobbered() []*Register {
    // 0 - 7 (inclusive)
	return []*Register{
		&Register{Name: "rax"},
		&Register{Name: "rcx"},
		&Register{Name: "rdx"},
		&Register{Name: "rsi"},
		&Register{Name: "rdi"},
		&Register{Name: "r8"},
		&Register{Name: "r9"},
		&Register{Name: "r10"},
		&Register{Name: "r11"},
	}
}

func GetCalleeSaved() []*Register {
    // 19 - 28 (inclusive)
	return []*Register{
		&Register{Name: "rsp"},
		&Register{Name: "rbp"},
		&Register{Name: "rbx"},
		&Register{Name: "r12"},
		&Register{Name: "r13"},
		&Register{Name: "r14"},
		&Register{Name: "r15"},
	}
}
