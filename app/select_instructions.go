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
	case expr.Primitive != nil:
		switch {
		case expr.Primitive.Num != nil:
			return &VarAssemblyImmediate{
				Int: expr.Primitive.Num.Value,
			}, []*VarAssemblyInstr{}, nil
		case expr.Primitive.Var != nil:
			return &VarAssemblyImmediate{
				Var: &VarAssemblyVar{
					Generated: expr.Primitive.Var.Generated,
				},
			}, []*VarAssemblyInstr{}, nil
		default:
			return nil, nil, errors.New("Unrecognized primitive type")
		}
	case expr.App != nil:
		switch {
		case expr.App.Operator.Name != "":
			return nil, nil, errors.New("An error was made and an unprocessed variable made it through")
		case expr.App.Operator.Generated != 0:
			return nil, nil, errors.New("User defined functions will go here")
		case expr.App.Operator.Primitive != "":
			switch expr.App.Operator.Primitive {
			default:
				return nil, nil, errors.New("Unrecognized primitive")
			}
		default:
			return nil, nil, errors.New("Unrecognized variable type")
		}

	default:
		return nil, nil, errors.New("Unrecognized SimpleExpr type")
	}	
}	

func HandlePrimitive(primitive string, operands []*SimplePrimitive) ([]*VarAssemblyInstr, error) {
	return nil, errors.New("Getting there")
}