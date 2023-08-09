package main

import (
    "errors"
    "language/pkg/languages"
)

type Env struct {
	Parent *Env
	Vars map[languages.Var]*languages.Var
}

func Lookup(variable *languages.Var, env *Env) (*languages.Var, bool) {
	v, ok := env.Vars[*variable]
	if ok {
		return v, true
	}
	if env.Parent == nil {
		return nil, false
	}

	return Lookup(variable, env.Parent)
}

func Uniquify(prog *languages.Program, getVar func() *languages.Var) (*languages.Program, error) {
	env := &Env{
		Vars: make(map[languages.Var]*languages.Var),
	}
	for _, v := range(languages.GetBuiltIns()) {
		env.Vars[*v] = v
	}

	uniqExpr, err := UniquifyExpr(prog.Expr, env, getVar)
	if err != nil {
		return nil, err
	}
	return &languages.Program{
		Expr: uniqExpr,
	}, nil
}

func UniquifyExpr(expr *languages.Expr, env *Env, getVar func() *languages.Var) (*languages.Expr, error) {
	switch {
	case expr.Num != nil:
		return expr, nil
	case expr.Var != nil:
		v, ok := Lookup(expr.Var, env)
		if ok {
			return &languages.Expr{Var: v}, nil
		} else {
			return nil, errors.New("Unbound variable: " + expr.Var.Name)
		}
	case expr.Let != nil:
		boundVars := make(map[languages.Var]*languages.Var)
		var newAssignments []*languages.Assignment
		for _, assignment := range(expr.Let.LetAssignments) {
			newVar := getVar()
			boundVars[*assignment.Ref] = newVar

			newAssignmentExpr, err := UniquifyExpr(assignment.Expr, env, getVar)
			if err != nil {
				return nil, err
			}
			newAssignments = append(newAssignments, &languages.Assignment{
				Ref: newVar,
				Expr: newAssignmentExpr,
			})
		}

		newEnv := &Env{
			Parent: env,
			Vars: boundVars,
		}
		newBody, err := UniquifyExpr(expr.Let.LetBody, newEnv, getVar)
		if err != nil {
			return nil, err
		}

		return &languages.Expr{
			Let: &languages.LetExpr {
				LetAssignments: newAssignments,
				LetBody: newBody,
			},
		}, nil
	case expr.App != nil:
		var newExprs []*languages.Expr
		for _, e := range(expr.App) {
			newExpr, err := UniquifyExpr(e, env, getVar)
			if err != nil {
				return nil, err
			}
			newExprs = append(newExprs, newExpr)
		}
		return &languages.Expr{
			App: newExprs,
		}, nil
	default:
		return nil, errors.New("Unrecognized expression type")
	}
}
