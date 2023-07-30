package main

import (
	"errors"
)

func SelectInstructions(prog *SimpleProgram) (*VarAssemblyProgram, error) {
	var instrs []*VarAssemblyInstr
	for _, s := range(prog.Statements) {
		stmtInstrs, err := SelectInstructionsStmt(s)
		if err != nil {
			return nil, err
		}
		instrs = append(instrs, stmtInstrs...)
	}
	return &VarAssemblyProgram{Instrs: instrs}, nil
}

func SelectInstructionsStmt(stmt *SimpleStatement) ([]*VarAssemblyInstr, error) {
	switch {
	case stmt.Expr != nil:	
		// truly nothing to do with a naked expression here
		_, instrs, err := SelectInstructionsExpr(stmt.Expr)
		if err != nil {
			return nil, err
		}
		return instrs, nil
	case stmt.Assignment != nil:
		imm, instrs, err := SelectInstructionsExpr(stmt.Assignment.Expr)
		if err != nil {
			return nil, err
		}
		if stmt.Assignment.Ref.Generated == 0 {
			return nil, errors.New("I've made a mistake in variable generation")
		}
		return append(instrs, &VarAssemblyInstr{
			Movq: [2]*VarAssemblyImmediate{
				&VarAssemblyImmediate{
					Var: &VarAssemblyVar{
						Generated: stmt.Assignment.Ref.Generated,
					},
				},
				imm,
			},
		}), nil
	default:
		return nil, errors.New("Unrecognized SimpleStatement")
	}
}

func SelectInstructionsExpr(expr *SimpleExpr) (*VarAssemblyImmediate, []*VarAssemblyInstr, error) {
	switch {
	case expr.Num != nil:
		return &VarAssemblyImmediate{
			Int: expr.Num.Value,
		}, []*VarAssemblyInstr{}, nil
	case expr.Var != nil:
		return &VarAssemblyImmediate{
			Var: &VarAssemblyVar{
				Generated: expr.Var.Generated,
			},
		}, []*VarAssemblyInstr{}, nil
	// case expr.App != nil:
	// 	rator, rands := expr.App[0], expr.App[1:]

	default:
		return nil, nil, errors.New("Haven't gotten to this yet")
	}	
}