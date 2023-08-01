package main

import (
    "errors"
)

type Env struct {
	Parent *Env
	Vars map[Var]*Var
}

func Lookup(variable *Var, env *Env) (*Var, bool) {
	v, ok := env.Vars[*variable]
	if ok {
		return v, true
	}
	if env.Parent == nil {
		return nil, false
	}

	return Lookup(variable, env.Parent)
}

func Uniquify(prog *Program, getVar func() *Var) (*Program, error) {
	env := &Env{
		Vars: make(map[Var]*Var),
	}
	for _, v := range(GetBuiltIns()) {
		env.Vars[*v] = v
	}

	uniqExpr, err := UniquifyExpr(prog.Expr, env, getVar)
	if err != nil {
		return nil, err
	}
	return &Program{
		Expr: uniqExpr,
	}, nil
}

func UniquifyExpr(expr *Expr, env *Env, getVar func() *Var) (*Expr, error) {
	switch {
	case expr.Num != nil:
		return expr, nil
	case expr.Var != nil:
		v, ok := Lookup(expr.Var, env)
		if ok {
			return &Expr{Var: v}, nil
		} else {
			return nil, errors.New("Unbound variable: " + expr.Var.Name)
		}
	case expr.Let != nil:	
		boundVars := make(map[Var]*Var)
		var newAssignments []*Assignment
		for _, assignment := range(expr.Let.LetAssignments) {
			newVar := getVar()
			boundVars[*assignment.Ref] = newVar

			newAssignmentExpr, err := UniquifyExpr(assignment.Expr, env, getVar)
			if err != nil {
				return nil, err
			}
			newAssignments = append(newAssignments, &Assignment{
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

		return &Expr{
			Let: &LetExpr {
				LetAssignments: newAssignments,
				LetBody: newBody,
			},
		}, nil
	case expr.App != nil:
		var newExprs []*Expr
		for _, e := range(expr.App) {
			newExpr, err := UniquifyExpr(e, env, getVar)
			if err != nil {
				return nil, err
			}
			newExprs = append(newExprs, newExpr)
		}
		return &Expr{
			App: newExprs,
		}, nil
	default:
		return nil, errors.New("Unrecognized expression type")
	}
}