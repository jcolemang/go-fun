package main

import (
	"errors"
    "language/pkg/languages"
)

func SelectInstructions(prog *languages.SimpleProgram, getVar func() *languages.Var) (*languages.VarAssemblyProgram, error) {
	var instrs []*languages.VarAssemblyInstr
	for _, s := range(prog.Statements) {
		stmtInstrs, err := SelectInstructionsStmt(s, getVar)
		if err != nil {
			return nil, err
		}
		instrs = append(instrs, stmtInstrs...)
	}
	return &languages.VarAssemblyProgram{Instrs: instrs}, nil
}

func SelectInstructionsStmt(stmt *languages.SimpleStatement, getVar func() *languages.Var) ([]*languages.VarAssemblyInstr, error) {
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
		targetVar := &languages.VarAssemblyVar{
			Generated: stmt.Assignment.Ref.Generated,
		}
		instrs, err := SelectInstructionsExpr(stmt.Assignment.Expr, targetVar)
		if err != nil {
			return nil, err
		}

		return instrs, nil
	case stmt.Return != nil:
		targetVar := &languages.VarAssemblyVar{
			Generated: getVar().Generated,
		}
		instrs, err := SelectInstructionsExpr(stmt.Return, targetVar)
		if err != nil {
			return nil, err
		}
        finalInstrs := []*languages.VarAssemblyInstr{
            &languages.VarAssemblyInstr{
                Mov: &[2]*languages.VarAssemblyImmediate{
                    &languages.VarAssemblyImmediate{
                        Register: languages.ReturnReg(),
                    },
                    &languages.VarAssemblyImmediate{
                        Var: targetVar,
                    },
                },
            },
            &languages.VarAssemblyInstr{
                Ret: &languages.Ret{},
            },
        }
		return append(instrs, finalInstrs...), nil
	default:
		return nil, errors.New("Unrecognized SimpleStatement")
	}
}

// without passing the variable through I think this would need to be able to generate a new
// variable to hold the value of the expression but that would I think just add a lot of extra
// unnecessary variables
func SelectInstructionsExpr(expr *languages.SimpleExpr, target *languages.VarAssemblyVar) ([]*languages.VarAssemblyInstr, error) {
	switch {
	case expr.Primitive != nil:
		switch {
		case expr.Primitive.Num != nil:
			if target != nil {
				var val *languages.VarAssemblyImmediate
				if expr.Primitive.Num.Int != nil {
					val = &languages.VarAssemblyImmediate{
						Int: expr.Primitive.Num.Int,
					}
				} else {
					return nil, errors.New("Unrecognized number type")
				}
				return []*languages.VarAssemblyInstr{
					&languages.VarAssemblyInstr{
						Mov: &[2]*languages.VarAssemblyImmediate{
							&languages.VarAssemblyImmediate{Var: target},
							val,
						},
					},
				}, nil
			} else {
				return []*languages.VarAssemblyInstr{}, nil
			}
		case expr.Primitive.Var != nil:
			return []*languages.VarAssemblyInstr{
				&languages.VarAssemblyInstr{
					Mov: &[2]*languages.VarAssemblyImmediate{
						&languages.VarAssemblyImmediate{
							Var: &languages.VarAssemblyVar{
								Generated: expr.Primitive.Var.Generated,
							},
						},
						&languages.VarAssemblyImmediate{Var: target},
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

func PrimitiveToImmediate(primitive *languages.SimplePrimitive) (*languages.VarAssemblyImmediate, error) {
	switch {
	case primitive.Num != nil && primitive.Num.Int != nil:
		return &languages.VarAssemblyImmediate{Int: primitive.Num.Int}, nil
	case primitive.Var != nil:
		return &languages.VarAssemblyImmediate{Var: &languages.VarAssemblyVar{Generated: primitive.Var.Generated}}, nil
	default:
		return nil, errors.New("Unrecognized primitive")
	}
}

func HandlePrimitive(primitive string, operands []*languages.SimplePrimitive, target *languages.VarAssemblyVar) ([]*languages.VarAssemblyInstr, error) {
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

        return []*languages.VarAssemblyInstr{
            &languages.VarAssemblyInstr{
                Add: &[3]*languages.VarAssemblyImmediate{
                    &languages.VarAssemblyImmediate{Var: target},
                    firstImm,
                    secondImm,
                },
            },
        }, nil
	case "print":
		return nil, errors.New("Need to define print")
	case "read":
		return nil, errors.New("Need to define read")
	default:
		return nil, errors.New("Unrecognized primitive")
	}
}
