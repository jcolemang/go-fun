package main

import (
    "errors"
)

func RemoveComplexOperands(prog *FlatProgram, getVar func() *Var) (*SimpleProgram, error) {
	var newStatements []*SimpleStatement
	for _, s := range(prog.Statements) {
		statements, newAssigns, err := RemoveComplexOperandsFromStatement(s, getVar)
		if err != nil {
			return nil, err
		}
		for _, a := range(newAssigns) {
			newStatements = append(newStatements, &SimpleStatement{Assignment: a})
		}
		newStatements = append(newStatements, statements...)
	}

	return &SimpleProgram{
        Statements: newStatements,
    }, nil
}

func RemoveComplexOperandsFromStatement(statement *FlatStatement, getVar func() *Var) ([]*SimpleStatement, []*SimpleAssignment, error) {
	var newStatements []*SimpleStatement
	switch {
	case statement.Expr != nil:
		newExpr, newAssigns, err := RemoveComplexOperandsFromExpr(statement.Expr, false, getVar)
		// newExpr, newAssigns, err := RemoveComplexOperandsFromExpr(statement.Expr, getVar)
		if err != nil {
			return nil, nil, err
		}
		for _, a := range(newAssigns) {
			newStatements = append(newStatements, &SimpleStatement{Assignment: a})
		}
		newStatements = append(newStatements, &SimpleStatement{Expr: newExpr})
		return newStatements, nil, nil
	case statement.Assignment != nil:
		newExpr, newAssigns, err := RemoveComplexOperandsFromExpr(statement.Assignment.Expr, true, getVar)
		// newExpr, newAssigns, err := RemoveComplexOperandsFromExpr(statement.Assignment.Expr, getVar)
		if err != nil {
			return nil, nil, err
		}
		for _, a := range(newAssigns) {
			newStatements = append(newStatements, &SimpleStatement{Assignment: a})
		}
		newStatements = append(newStatements, &SimpleStatement{Assignment: &SimpleAssignment{Ref: statement.Assignment.Ref, Expr: newExpr}})

		return newStatements, nil, nil
	default:
		return nil, nil, errors.New("Unrecognized statement")
	}
}

// bool is necessary because we need to be able to distinguish between the (+ 1 2) in
// x = 1 + 2 and the (+ 1 2) in x = (+ x (+ 1 2)). The former only uses two "addresses"
// and the latter uses three as it is a subexpression in a larger expression
// I think that makeAtomic can be removed and assumed false but I also have a suspicion that it will matter for user defined functions.
// Currently the output of a function call can go directly into a variable with addq but I'm not sure if that same 
// logic works for user defined functions. Does the return value go into a defined register to be used? Would we
// then have to worry about multiple function calls? Would that even possibly come up if we're the ones generating 
// the instructions? Questions for a later day.
func RemoveComplexOperandsFromExpr(expr *FlatExpr, makeAtomic bool, getVar func() *Var) (*SimpleExpr, []*SimpleAssignment, error) {
    switch {
    case expr.Num != nil:
        return &SimpleExpr{Primitive: &SimplePrimitive{Num: expr.Num}}, []*SimpleAssignment{}, nil
    case expr.Var != nil:
        return &SimpleExpr{Primitive: &SimplePrimitive{Var: expr.Var}}, []*SimpleAssignment{}, nil
    case expr.App != nil:
		var newExprs []*SimplePrimitive
		var newAssignments []*SimpleAssignment
		for _, e := range(expr.App) {
			// arguments to a function must be atomic
			subExpr, subExprAssignments, err := GenerateAtomicExpression(e, getVar)
			if err != nil {
				return nil, nil, err
			}
			newExprs = append(newExprs, subExpr)
			newAssignments = append(newAssignments, subExprAssignments...)
		}

		rator, rands := newExprs[0], newExprs[1:]

		newApp := &SimpleExpr{
			App: &SimpleApplication{
				Operator: rator.Var,
				Operands: rands,
			},
		}
		if makeAtomic {
			newVar := getVar()
			newAssignment := &SimpleAssignment{
				Ref: newVar,
				Expr: newApp,
			}
			newExpr := &SimpleExpr{
				Primitive: &SimplePrimitive{
					Var: newVar,
				},
			}
			return newExpr, append(newAssignments, newAssignment), nil
		} else {
		return newApp, newAssignments, nil
		}
    default:
        return nil, nil, errors.New("Unrecognized expression type")
    }
}

// There is a lot of overlap here with the above function but it felt better doing it this way to help limit the types
// Will probably see a cleaner way to do this in about 72 hours
func GenerateAtomicExpression(expr *FlatExpr, getVar func() *Var) (*SimplePrimitive, []*SimpleAssignment, error) {
	switch {
	case expr.Num != nil:
		return &SimplePrimitive{Num: expr.Num}, []*SimpleAssignment{}, nil
	case expr.Var != nil:
		return &SimplePrimitive{Var: expr.Var}, []*SimpleAssignment{}, nil
	case expr.App != nil:
		var primitives []*SimplePrimitive
		var newAssignments []*SimpleAssignment
		for _, e := range(expr.App) {
			primitive, subExprAssigns, err := GenerateAtomicExpression(e, getVar)
			if err != nil {
				return nil, nil, err
			}
			primitives = append(primitives, primitive)
			newAssignments = append(newAssignments, subExprAssigns...)
		}

		rator, rands := primitives[0], primitives[1:]
		if rator.Var == nil {
			return nil, nil, errors.New("Attempt to apply something non-apply-able")
		}
		newApp := &SimpleExpr{
			App: &SimpleApplication{
				Operator: rator.Var,
				Operands: rands,
			},
		}
		newVar := getVar()
		newAssignment := &SimpleAssignment{
			Ref: newVar,
			Expr: newApp,
		}

		return &SimplePrimitive{Var: newVar}, append(newAssignments, newAssignment), nil

	default:
		return nil, nil, errors.New("Unrecognized FlatExpr type")
	}
}

func IsExprSimple(expr *FlatExpr) bool {
	switch {
	case expr.Num != nil:
		return true
	case expr.Var != nil:
		return true
	case expr.App != nil:
		for _, e := range(expr.App) {
			if !IsExprPrimitive(e) {
				return false
			}
		}
		return true
	default:
		return false
	}
}

func IsExprPrimitive(expr *FlatExpr) bool {
	switch {
	case expr.Num != nil:
		return true
	case expr.Var != nil:
		return true
	case expr.App != nil:
		return false
	default:
		return false
	}
}