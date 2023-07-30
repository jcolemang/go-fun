package main

import (
	"errors"
)

func Flatten(progOrig *Program) (*FlatProgram, error) {
	flatExpr, assigns, err := FlattenExpr(progOrig.Expr)
	if err != nil {
		return nil, err
	}

	var newExprs []*FlatExpr
	for _, a := range(assigns) {
		newExprs = append(newExprs, &FlatExpr{Assignment: a})
	}
	newExprs = append(newExprs, flatExpr)

	return &FlatProgram{
		Exprs: newExprs,
	}, nil
}

func FlattenExpr(expr *Expr) (*FlatExpr, []*Assignment, error) {
	var assignments []*Assignment
	switch {
	case expr.Num != nil:
		return &FlatExpr{Num: expr.Num}, assignments, nil
	case expr.Var != nil:
		return &FlatExpr{Var: expr.Var}, assignments, nil
	case expr.Let != nil:
		return nil, nil, errors.New("I'll come back to this")
	case expr.App != nil:
		var newExprs []*FlatExpr
		for _, e := range(expr.App) {
			subExpr, subAssigns, err := FlattenExpr(e)
			if err != nil {
				return nil, nil, err
			}
			newExprs = append(newExprs, subExpr)
			assignments = append(assignments, subAssigns...)
		}
		return &FlatExpr{App: newExprs}, assignments, nil
	default:
		return nil, nil, errors.New("Unrecognized expression")
	}
}