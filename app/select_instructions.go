package main

import (
	"errors"
)

func SelectInstructions(prog *SimpleExitProgram) (*VarAssemblyProgram, error) {
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
		// truly nothing to do with just the naked immediate here
		instrs, err := SelectInstructionsExpr(stmt.Expr, nil)
		if err != nil {
			return nil, err
		}
		return instrs, nil
	case stmt.Assignment != nil:
		if stmt.Assignment.Ref.Generated == 0 {
			return nil, errors.New("I've made a mistake in variable generation")
		}
		targetVar := &VarAssemblyVar{
			Generated: stmt.Assignment.Ref.Generated,
		}
		instrs, err := SelectInstructionsExpr(stmt.Assignment.Expr, targetVar)
		if err != nil {
			return nil, err
		}

		return instrs, nil
	default:
		return nil, errors.New("Unrecognized SimpleStatement")
	}
}

// without passing the variable through I think this would need to be able to generate a new
// variable to hold the value of the expression but that would I think just add a lot of extra
// unnecessary variables
func SelectInstructionsExpr(expr *SimpleExpr, target *VarAssemblyVar) ([]*VarAssemblyInstr, error) {
	switch {
	case expr.Primitive != nil:
		switch {
		case expr.Primitive.Num != nil:
			if target != nil {
				var val *VarAssemblyImmediate
				if expr.Primitive.Num.Int != nil {
					val = &VarAssemblyImmediate{
						Int: expr.Primitive.Num.Int,
					}
				} else {
					return nil, errors.New("Unrecognized number type")
				}
				return []*VarAssemblyInstr{
					&VarAssemblyInstr{
						Movq: [2]*VarAssemblyImmediate{
							&VarAssemblyImmediate{Var: target},
							val,
						},
					},
				}, nil
			} else {
				return []*VarAssemblyInstr{}, nil
			}
		case expr.Primitive.Var != nil:
			return []*VarAssemblyInstr{
				&VarAssemblyInstr{
					Movq: [2]*VarAssemblyImmediate{
						&VarAssemblyImmediate{Var: target},
						&VarAssemblyImmediate{
							Var: &VarAssemblyVar{
								Generated: expr.Primitive.Var.Generated,
							},
						},
					},
				},
			}, nil
		default:
			return nil, errors.New("Unrecognized primitive type")
		}
	case expr.App != nil:
		switch {
		case expr.App.Operator.Name != "":
			return nil, errors.New("An error was made and an unprocessed variable made it through")
		case expr.App.Operator.Generated != 0:
			return nil, errors.New("User defined functions will go here")
		case expr.App.Operator.Primitive != "":
			instrs, err := HandlePrimitive(expr.App.Operator.Primitive, expr.App.Operands, target)
			if err != nil {
				return nil, err
			}
			return instrs, nil
		default:
			return nil, errors.New("Unrecognized variable type")
		}
	default:
		return nil, errors.New("Unrecognized SimpleExpr type")
	}
}

func PrimitiveToImmediate(primitive *SimplePrimitive) (*VarAssemblyImmediate, error) {
	switch {
	case primitive.Num != nil && primitive.Num.Int != nil:
		return &VarAssemblyImmediate{Int: primitive.Num.Int}, nil
	case primitive.Var != nil:
		return &VarAssemblyImmediate{Var: &VarAssemblyVar{Generated: primitive.Var.Generated}}, nil
	default:
		return nil, errors.New("Unrecognized primitive")
	}
}

func HandlePrimitive(primitive string, operands []*SimplePrimitive, target *VarAssemblyVar) ([]*VarAssemblyInstr, error) {
	switch primitive {
	case "+":
		if len(operands) != 2 {
			return nil, errors.New("Cannot currently handle arbitrary numbers of arguments to addition")
		}
		first, second := operands[0], operands[1]
		firstImm, err := PrimitiveToImmediate(first)
		if err != nil {
			return nil, err
		}
		secondImm, err := PrimitiveToImmediate(second)
		if err != nil {
			return nil, err
		}

		switch {
		case first.Var != nil && first.Var.Generated == target.Generated:
			return []*VarAssemblyInstr{
				&VarAssemblyInstr{
					Addq: [2]*VarAssemblyImmediate{
						secondImm,
						firstImm,
					},
				},
			}, nil
		case second.Var != nil && second.Var.Generated == target.Generated:
			return []*VarAssemblyInstr{
				&VarAssemblyInstr{
					Addq: [2]*VarAssemblyImmediate{
						firstImm,
						secondImm,
					},
				},
			}, nil
		default:
			return []*VarAssemblyInstr{
				&VarAssemblyInstr{
					Movq: [2]*VarAssemblyImmediate{
						firstImm,
						&VarAssemblyImmediate{Register: &Register{Name: "rax"}},
					},
				},
				&VarAssemblyInstr{
					Addq: [2]*VarAssemblyImmediate{
						secondImm,
						&VarAssemblyImmediate{Register: &Register{Name: "rax"}},
					},
				},
				&VarAssemblyInstr{
					Movq: [2]*VarAssemblyImmediate{
						&VarAssemblyImmediate{Register: &Register{Name: "rax"}},
						&VarAssemblyImmediate{Var: target},
					},
				},
			}, nil
		}
	case "print":
		return nil, errors.New("Need to define print")
	case "read":
		return nil, errors.New("Need to define read")
	default:
		return nil, errors.New("Unrecognized primitive")
	}
}