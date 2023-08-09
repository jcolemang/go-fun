package main

import (
	"errors"
    "language/pkg/languages"
)

func Flatten(progOrig *languages.Program) (*languages.FlatProgram, error) {
	flatExpr, assigns, err := FlattenExpr(progOrig.Expr)
	if err != nil {
		return nil, err
	}

	var newStatements []*languages.FlatStatement
	for _, a := range(assigns) {
		newStatements = append(newStatements, &languages.FlatStatement{Assignment: a})
	}
	newStatements = append(newStatements, &languages.FlatStatement{Return: flatExpr})

	return &languages.FlatProgram{
		Statements: newStatements,
	}, nil
}

func FlattenExpr(expr *languages.Expr) (*languages.FlatExpr, []*languages.FlatAssignment, error) {
	var assignments []*languages.FlatAssignment
	switch {
	case expr.Num != nil:
		return &languages.FlatExpr{Num: expr.Num}, assignments, nil
	case expr.Var != nil:
		return &languages.FlatExpr{Var: expr.Var}, assignments, nil
	case expr.Let != nil:
		var letAssignAssigns []*languages.FlatAssignment
		for _, a := range(expr.Let.LetAssignments) {
			flattenedExpr, subAssigns, err := FlattenExpr(a.Expr)
			if err != nil {
				return nil, nil, err
			}
			letAssignAssigns = append(letAssignAssigns, &languages.FlatAssignment{Ref: a.Ref, Expr: flattenedExpr})
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
		var newExprs []*languages.FlatExpr
		for _, e := range(expr.App) {
			subExpr, subAssigns, err := FlattenExpr(e)
			if err != nil {
				return nil, nil, err
			}
			newExprs = append(newExprs, subExpr)
			assignments = append(assignments, subAssigns...)
		}
		return &languages.FlatExpr{App: newExprs}, assignments, nil
	default:
		return nil, nil, errors.New("Unrecognized expression")
	}
}
