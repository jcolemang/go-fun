package main

import (
	"errors"
    "language/pkg/languages"
)

func SelectInstructions(prog *languages.BlockProgram, getVar func() *languages.Var) (*languages.VarAssemblyProgram, error) {
	var instrs []*languages.VarAssemblyInstr
	for _, b := range(prog.Blocks) {
        blockInstrs, err := SelectInstructionsBlock(b, getVar)
        if err != nil {
            return nil, err
        }
        instrs = append(instrs, blockInstrs...)
	}
	return &languages.VarAssemblyProgram{Instrs: instrs}, nil
}

func SelectInstructionsBlock(block languages.IBlock, getVar func() *languages.Var) ([]*languages.VarAssemblyInstr, error) {
    switch b := block.(type) {
    case languages.BasicBlock:
        var instrs []*languages.VarAssemblyInstr
        for _, s := range(b.Statements) {
            stmtInstrs, err := SelectInstructionsStmt(s, getVar)
            if err != nil {
                return nil, err
            }
            instrs = append(instrs, stmtInstrs...)
        }
        terminatorInstrs, err := SelectInstructionsTerminator(b.Terminator, getVar)
        if err != nil {
            return nil, err
        }
        return append(instrs, terminatorInstrs...), nil
    }
    return nil, errors.New("Unrecognized block type in SelectInstructionsBlock")
}

func SelectInstructionsStmt(blockStmt languages.IBlockStatement, getVar func() *languages.Var) ([]*languages.VarAssemblyInstr, error) {
    switch stmt := blockStmt.(type) {
	case languages.BlockExpr:
		instrs, err := SelectInstructionsExpr(stmt.Expr, nil)
		if err != nil {
			return nil, err
		}
		return instrs, nil
	case languages.Assignment[languages.BlockExpr]:
		if stmt.Ref.Generated == 0 {
			return nil, errors.New("I've made a mistake in variable generation")
		}
		targetVar := &languages.VarAssemblyVar{
			Generated: stmt.Ref.Generated,
		}
		instrs, err := SelectInstructionsExpr(stmt.Expr.Expr, targetVar)
		if err != nil {
			return nil, err
		}

		return instrs, nil
	default:
		return nil, errors.New("Unrecognized BlockStatement in SelectInstructionsStmt")
	}
}

func SelectInstructionsTerminator(blockTerm languages.IBlockTerminator, getVar func() *languages.Var) ([]*languages.VarAssemblyInstr, error) {
    switch b := blockTerm.(type) {
	case languages.BlockReturn:
		targetVar := &languages.VarAssemblyVar{
			Generated: getVar().Generated,
		}
		instrs, err := SelectInstructionsExpr(b.Val.Expr, targetVar)
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
        return nil, errors.New("Unrecognized block terminator in SelectInstructionsTerminator")
    }
}

// without passing the variable through I think this would need to be able to generate a new
// variable to hold the value of the expression but that would I think just add a lot of extra
// unnecessary variables
func SelectInstructionsExpr(expr languages.IBlockExpr, target *languages.VarAssemblyVar) ([]*languages.VarAssemblyInstr, error) {
    switch b := expr.(type) {
	case languages.Primitive:
        if target == nil {
	        return []*languages.VarAssemblyInstr{}, nil
        }
		switch {
		case b.Num != nil:
            var val *languages.VarAssemblyImmediate
            if b.Num.Int != nil {
                val = &languages.VarAssemblyImmediate{
                    Int: b.Num.Int,
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
		case b.Bool != nil:
            var boolVal int
            if b.Bool.True != nil {
                boolVal = 1
            } else {
                boolVal = 0
            }
            val := &languages.VarAssemblyImmediate{
                Int: &boolVal,
            }
            return []*languages.VarAssemblyInstr{
                &languages.VarAssemblyInstr{
                    Mov: &[2]*languages.VarAssemblyImmediate{
                        &languages.VarAssemblyImmediate{Var: target},
                        val,
                    },
                },
            }, nil
		case b.Var != nil:
			return []*languages.VarAssemblyInstr{
				&languages.VarAssemblyInstr{
					Mov: &[2]*languages.VarAssemblyImmediate{
						&languages.VarAssemblyImmediate{
							Var: &languages.VarAssemblyVar{
								Generated: b.Var.Generated,
							},
						},
						&languages.VarAssemblyImmediate{Var: target},
					},
				},
			}, nil
		default:
			return nil, errors.New("Unrecognized primitive type in SelectInstructionsExpr")
		}
	case languages.PrimitiveApplication:
		switch {
		case b.Operator.Name != "":
			return nil, errors.New("An error was made and an unprocessed variable made it through")
		case b.Operator.Generated != 0:
			return nil, errors.New("User defined functions will go here")
		case b.Operator.Primitive != "":
			instrs, err := HandlePrimitive(b.Operator.Primitive, b.Operands, target)
			if err != nil {
				return nil, err
			}
			return instrs, nil
		default:
			return nil, errors.New("Unrecognized variable type")
		}
	default:
		return nil, errors.New("Unrecognized SimpleExpr type in SelectInstructionsExpr")
	}
}

func PrimitiveToImmediate(primitive *languages.Primitive) (*languages.VarAssemblyImmediate, error) {
	switch {
	case primitive.Num != nil && primitive.Num.Int != nil:
		return &languages.VarAssemblyImmediate{Int: primitive.Num.Int}, nil
	case primitive.Var != nil:
		return &languages.VarAssemblyImmediate{Var: &languages.VarAssemblyVar{Generated: primitive.Var.Generated}}, nil
	default:
		return nil, errors.New("Unrecognized primitive")
	}
}

func HandlePrimitive(primitive string, operands []*languages.Primitive, target *languages.VarAssemblyVar) ([]*languages.VarAssemblyInstr, error) {
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
