package main

import (
)

func PatchInstructions(prog *X86Program) *X86Program {
	newInstrs := make([]*X86Instr, 0)
	mainLabel := "main"
	newInstrs = append(newInstrs, &X86Instr{
		Label: &mainLabel,
	})

	for _, instrs := range prog.X86Instrs {
		newInstrs = append(newInstrs, PatchInstruction(instrs)...)
	}

	return &X86Program{
        X86Directives: append(prog.X86Directives, &X86Directive{Name: "globl", Arg: "main"}),
		X86Instrs: newInstrs,
	}
}

func PatchInstruction(instr *X86Instr) []*X86Instr {
	switch {
	// can just remove useless movqs
	// til how to use multiline statements!
	case instr.Movq != nil &&
			instr.Movq[0].X86Offset != nil &&
			instr.Movq[1].X86Offset != nil &&
			*instr.Movq[0].X86Offset == *instr.Movq[1].X86Offset &&
			instr.Movq[0].X86OffsetReg.Name == instr.Movq[1].X86OffsetReg.Name:
		return []*X86Instr{}
	// instruction cannot reference two memory locations
	case instr.Addq != nil && instr.Addq[0].X86Offset != nil && instr.Addq[1].X86Offset != nil:
		return []*X86Instr{
			&X86Instr{
				Movq: []*X86Arg{
					instr.Addq[0],
					&X86Arg{
						X86Reg: &Register{Name: "rax"},
					},
				},
			},
			&X86Instr{
				Addq: []*X86Arg{
					&X86Arg{
						X86Reg: &Register{Name: "rax"},
					},
					instr.Addq[1],
				},
			},
		}
	case instr.Movq != nil && instr.Movq[0].X86Offset != nil && instr.Movq[1].X86Offset != nil:
		return []*X86Instr{
			&X86Instr{
				Movq: []*X86Arg{
					instr.Addq[0],
					&X86Arg{
						X86Reg: &Register{Name: "rax"},
					},
				},
			},
			&X86Instr{
				Movq: []*X86Arg{
					&X86Arg{
						X86Reg: &Register{Name: "rax"},
					},
					instr.Addq[1],
				},
			},
		}

	default:
		return []*X86Instr{instr}
	}
}
