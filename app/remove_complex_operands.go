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
func RemoveComplexOperandsFromExpr(expr *FlatExpr, makeAtomic bool, getVar func() *Var) (*SimpleExpr, []*SimpleAssignment, error) {
    switch {
    case expr.Num != nil:
        return &SimpleExpr{Num: expr.Num}, []*SimpleAssignment{}, nil
    case expr.Var != nil:
        return &SimpleExpr{Var: expr.Var}, []*SimpleAssignment{}, nil
    case expr.App != nil:
		var newExprs []*SimpleExpr
		var newAssignments []*SimpleAssignment
		for _, e := range(expr.App) {
			// arguments to a function must be atomic
			subExpr, subExprAssignments, err := RemoveComplexOperandsFromExpr(e, true, getVar)
			if err != nil {
				return nil, nil, err
			}
			newExprs = append(newExprs, subExpr)
			newAssignments = append(newAssignments, subExprAssignments...)
		}

		rator, rands := newExprs[0], newExprs[1:]
		if rator.Var == nil {
			return nil, nil, errors.New("Operator must be atomic")
		}

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
				Var: newVar,
			}
			return newExpr, append(newAssignments, newAssignment), nil
		} else {
			return newApp, newAssignments, nil
		}
    default:
        return nil, nil, errors.New("Unrecognized expression type")
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