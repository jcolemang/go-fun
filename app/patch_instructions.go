package main

import (
)

func PatchInstructions(prog *ArmProgram) *ArmProgram {
	newInstrs := make([]*ArmInstr, 0)
	mainLabel := "_start"
	newInstrs = append(newInstrs, &ArmInstr{
		Label: &mainLabel,
	})

	for _, instrs := range prog.ArmInstrs {
		newInstrs = append(newInstrs, PatchInstruction(instrs)...)
	}

	return &ArmProgram{
        ArmDirectives: append(
            prog.ArmDirectives,
            []*ArmDirective{
                &ArmDirective{Name: "global", Arg: "_start"},
                &ArmDirective{Name: "align", Arg: "4"},
            }...,
        ),
		ArmInstrs: newInstrs,
	}
}

func PatchInstruction(instr *ArmInstr) []*ArmInstr {
	switch {
	// can just remove useless Movs
	// til how to use multiline statements!
	case instr.Mov != nil &&
			instr.Mov[0].ArmOffset != nil &&
			instr.Mov[1].ArmOffset != nil &&
			*instr.Mov[0].ArmOffset == *instr.Mov[1].ArmOffset &&
			instr.Mov[0].ArmOffsetReg.Name == instr.Mov[1].ArmOffsetReg.Name:
		return []*ArmInstr{}
	case instr.Add != nil && !IsRegister(*instr.Add[1]) && !IsRegister(*instr.Add[2]):
		return []*ArmInstr{
			&ArmInstr{
				Mov: []*ArmArg{
					&ArmArg{
						ArmReg: TempReg(),
					},
					instr.Add[1],
				},
			},
			&ArmInstr{
				Add: []*ArmArg{
                    instr.Add[0],
					&ArmArg{
						ArmReg: TempReg(),
					},
					instr.Add[2],
				},
			},
		}
    case instr.Add != nil && !IsRegister(*instr.Add[1]):
		return []*ArmInstr{
			&ArmInstr{
				Add: []*ArmArg{
                    instr.Add[0],
					instr.Add[2],
					instr.Add[1],
				},
			},
		}
	default:
		return []*ArmInstr{instr}
	}
}
