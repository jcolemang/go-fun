package main

import (
    "language/pkg/languages"
)

func PatchInstructions(prog *languages.ArmProgram) *languages.ArmProgram {
	newInstrs := make([]*languages.ArmInstr, 0)
	mainLabel := "_start"
	newInstrs = append(newInstrs, &languages.ArmInstr{
		Label: &mainLabel,
	})

	for _, instr := range prog.ArmInstrs {
		newInstrs = append(newInstrs, PatchInstruction(instr)...)
	}

	return &languages.ArmProgram{
        ArmDirectives: append(
            prog.ArmDirectives,
            []*languages.ArmDirective{
                &languages.ArmDirective{Name: "global", Arg: "_start"},
                &languages.ArmDirective{Name: "align", Arg: "4"},
            }...,
        ),
		ArmInstrs: newInstrs,
	}
}

func PatchInstruction(instr *languages.ArmInstr) []*languages.ArmInstr {
	switch {
	// can just remove useless Movs
	case instr.Mov != nil &&
			instr.Mov[0].ArmOffset != nil &&
			instr.Mov[1].ArmOffset != nil &&
			*instr.Mov[0].ArmOffset == *instr.Mov[1].ArmOffset &&
			instr.Mov[0].ArmOffsetReg.Name == instr.Mov[1].ArmOffsetReg.Name:
		return []*languages.ArmInstr{}
	case instr.Mov != nil &&
			instr.Mov[0].ArmReg != nil &&
			instr.Mov[1].ArmReg != nil &&
            *instr.Mov[0].ArmReg == *instr.Mov[1].ArmReg:
		return []*languages.ArmInstr{}
    // add doesn't like an int in the first argument
	case instr.Add != nil && !languages.IsRegister(*instr.Add[1]) && !languages.IsRegister(*instr.Add[2]):
		return []*languages.ArmInstr{
			&languages.ArmInstr{
				Mov: []*languages.ArmArg{
					&languages.ArmArg{
						ArmReg: languages.TempReg(),
					},
					instr.Add[1],
				},
			},
			&languages.ArmInstr{
				Add: []*languages.ArmArg{
                    instr.Add[0],
					&languages.ArmArg{
						ArmReg: languages.TempReg(),
					},
					instr.Add[2],
				},
			},
		}
    case instr.Add != nil && !languages.IsRegister(*instr.Add[1]):
		return []*languages.ArmInstr{
			&languages.ArmInstr{
				Add: []*languages.ArmArg{
                    instr.Add[0],
					instr.Add[2],
					instr.Add[1],
				},
			},
		}
	default:
		return []*languages.ArmInstr{instr}
	}
}
