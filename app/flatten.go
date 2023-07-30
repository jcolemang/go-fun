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

func FlattenExpr(expr *Expr) (*FlatExpr, []*FlatAssignment, error) {
	var assignments []*FlatAssignment
	switch {
	case expr.Num != nil:
		return &FlatExpr{Num: expr.Num}, assignments, nil
	case expr.Var != nil:
		return &FlatExpr{Var: expr.Var}, assignments, nil
	case expr.Let != nil:
		// (let ((x (let ((y 1)) y)) x)
		// y = 1
		// x = y
		// x
		var letAssignAssigns []*FlatAssignment
		for _, a := range(expr.Let.LetAssignments) {
			flattenedExpr, subAssigns, err := FlattenExpr(a.Expr)
			if err != nil {
				return nil, nil, err
			}
			letAssignAssigns = append(letAssignAssigns, &FlatAssignment{Ref: a.Ref, Expr: flattenedExpr})
			assignments = append(assignments, subAssigns...)
		}

		assignments = append(assignments, letAssignAssigns...)
		flatBodyExpr, bodyAssigns, err := FlattenExpr(expr.Let.LetBody)
		if err != nil {
			return nil, nil, err
		}
		assignments = append(assignments, bodyAssigns...)

		return flatBodyExpr, assignments, nil
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