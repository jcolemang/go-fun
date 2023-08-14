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
	case expr.Bool != nil:
		return &languages.FlatExpr{Bool: expr.Bool}, assignments, nil
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
	case expr.IfExpr != nil:
        condExpr, condSubAssigns, err := FlattenExpr(expr.IfExpr.IfCond)
        if err != nil {
            return nil, nil, err
        }

        trueExpr, trueSubAssigns, err := FlattenExpr(expr.IfExpr.IfTrue)
        if err != nil {
            return nil, nil, err
        }
        var trueStmts []*languages.FlatStatement
        for _, a := range trueSubAssigns {
            trueStmts = append(trueStmts, &languages.FlatStatement{Assignment: a})
        }

        falseExpr, falseSubAssigns, err := FlattenExpr(expr.IfExpr.IfFalse)
        if err != nil {
            return nil, nil, err
        }
        var falseStmts []*languages.FlatStatement
        for _, a := range falseSubAssigns {
            falseStmts = append(falseStmts, &languages.FlatStatement{Assignment: a})
        }

        assignments = append(assignments, condSubAssigns...)
		return &languages.FlatExpr{
            IfExpr: &languages.FlatIfExpr{
                IfCond: condExpr,
                IfTrue: trueStmts,
                IfTrueExpr: trueExpr,
                IfFalse: falseStmts,
                IfFalseExpr: falseExpr,
            },
        }, assignments, nil
	default:
		return nil, nil, errors.New("Unrecognized expression in flatten")
	}
}
