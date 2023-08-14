package main

import (
    "errors"
    "language/pkg/languages"
)

func RemoveComplexOperands(prog *languages.IfStmtProgram, getVar func() *languages.Var) (*languages.SimpleProgram, error) {
    newStatements, err := RemoveComplexOperandsFromStatementList(prog.Statements, getVar)
    if err != nil {
        return nil, err
    }

	return &languages.SimpleProgram{
        Statements: newStatements,
    }, nil
}

func RemoveComplexOperandsFromStatementList(stmts []*languages.IfStmtStatement, getVar func() *languages.Var) ([]*languages.SimpleStatement, error) {
    var newStatements []*languages.SimpleStatement
	for _, s := range(stmts) {
		statements, newAssigns, err := RemoveComplexOperandsFromStatement(s, getVar)
		if err != nil {
			return nil, err
		}
		for _, a := range(newAssigns) {
			newStatements = append(newStatements, &languages.SimpleStatement{Assignment: a})
		}
		newStatements = append(newStatements, statements...)
	}

	return newStatements, nil
}

// looking at this now, I don't think the assignments needs to be here. It is always nil.
func RemoveComplexOperandsFromStatement(statement *languages.IfStmtStatement, getVar func() *languages.Var) ([]*languages.SimpleStatement, []*languages.Assignment[languages.SimpleExpr], error) {
	var newStatements []*languages.SimpleStatement
	switch {
	case statement.Expr != nil:
		newExpr, newAssigns, err := RemoveComplexOperandsFromExpr(statement.Expr, getVar)
		if err != nil {
			return nil, nil, err
		}
		for _, a := range(newAssigns) {
			newStatements = append(newStatements, &languages.SimpleStatement{Assignment: a})
		}
		newStatements = append(newStatements, &languages.SimpleStatement{Expr: newExpr})
		return newStatements, nil, nil
	case statement.Assignment != nil:
		newPrim, newAssigns, err := RemoveComplexOperandsFromExprAtomic(statement.Assignment.Expr, getVar)
		if err != nil {
			return nil, nil, err
		}
		for _, a := range(newAssigns) {
			newStatements = append(newStatements, &languages.SimpleStatement{Assignment: a})
		}
		newStatements = append(
            newStatements,
            &languages.SimpleStatement{
                Assignment: &languages.Assignment[languages.SimpleExpr]{
                    Ref: statement.Assignment.Ref,
                    Expr: &languages.SimpleExpr{
                        Primitive: newPrim,
                    },
                },
            },
        )

		return newStatements, nil, nil
    case statement.Return != nil:
		newExpr, newAssigns, err := RemoveComplexOperandsFromExpr(statement.Return, getVar)
		if err != nil {
			return nil, nil, err
		}
		for _, a := range(newAssigns) {
			newStatements = append(newStatements, &languages.SimpleStatement{Assignment: a})
		}
		newStatements = append(newStatements, &languages.SimpleStatement{Return: newExpr})
		return newStatements, nil, nil
    case statement.IfStmt != nil:
		newCondPrim, newCondAssigns, err := RemoveComplexOperandsFromExprAtomic(statement.IfStmt.IfCond, getVar)
		if err != nil {
			return nil, nil, err
		}

        newTrueStmts, err := RemoveComplexOperandsFromStatementList(statement.IfStmt.IfTrue, getVar)
		if err != nil {
			return nil, nil, err
		}

        newFalseStmts, err := RemoveComplexOperandsFromStatementList(statement.IfStmt.IfFalse, getVar)
		if err != nil {
			return nil, nil, err
		}

        newIfStmt := &languages.SimpleStatement{
            IfStmt: &languages.SimpleIfStmt{
                IfCond: newCondPrim,
                IfTrue: newTrueStmts,
                IfFalse: newFalseStmts,
            },
        }

        var newStatements []*languages.SimpleStatement
        for _, a := range newCondAssigns {
            newStatements = append(newStatements, &languages.SimpleStatement{
                Assignment: a,
            })
        }

		return append(newStatements, newIfStmt), newCondAssigns, nil
	default:
		return nil, nil, errors.New("Unrecognized statement in RemoveComplexOperandsFromStatement")
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
func RemoveComplexOperandsFromExpr(expr *languages.IfStmtExpr, getVar func() *languages.Var) (*languages.SimpleExpr, []*languages.Assignment[languages.SimpleExpr], error) {
    switch {
    case expr.Num != nil:
        return &languages.SimpleExpr{Primitive: &languages.SimplePrimitive{Num: expr.Num}}, []*languages.Assignment[languages.SimpleExpr]{}, nil
    case expr.Bool != nil:
        return &languages.SimpleExpr{Primitive: &languages.SimplePrimitive{Bool: expr.Bool}}, []*languages.Assignment[languages.SimpleExpr]{}, nil
    case expr.Var != nil:
        return &languages.SimpleExpr{Primitive: &languages.SimplePrimitive{Var: expr.Var}}, []*languages.Assignment[languages.SimpleExpr]{}, nil
    case expr.App != nil:
		var newExprs []*languages.SimplePrimitive
		var newAssignments []*languages.Assignment[languages.SimpleExpr]
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

		newApp := &languages.SimpleExpr{
			App: &languages.SimpleApplication{
				Operator: rator.Var,
				Operands: rands,
			},
		}
		return newApp, newAssignments, nil
    default:
        return nil, nil, errors.New("Unrecognized expression type in RemoveComplexOperandsFromExpr")
    }
}

func RemoveComplexOperandsFromExprAtomic(expr *languages.IfStmtExpr, getVar func() *languages.Var) (*languages.SimplePrimitive, []*languages.Assignment[languages.SimpleExpr], error) {
    switch {
    case expr.Num != nil:
        return &languages.SimplePrimitive{Num: expr.Num}, []*languages.Assignment[languages.SimpleExpr]{}, nil
    case expr.Bool != nil:
        return &languages.SimplePrimitive{Bool: expr.Bool}, []*languages.Assignment[languages.SimpleExpr]{}, nil
    case expr.Var != nil:
        return &languages.SimplePrimitive{Var: expr.Var}, []*languages.Assignment[languages.SimpleExpr]{}, nil
    case expr.App != nil:
        newApp, newAssignments, err := RemoveComplexOperandsFromExpr(expr, getVar)
        if err != nil {
            return nil, nil, err
        }

        newVar := getVar()
        newAssignment := &languages.Assignment[languages.SimpleExpr]{
            Ref: newVar,
            Expr: newApp,
        }
        newPrim := &languages.SimplePrimitive{
            Var: newVar,
        }

        return newPrim, append(newAssignments, newAssignment), nil
    default:
        return nil, nil, errors.New("Unrecognized expression type in RemoveComplexOperandsFromExpr")
    }
}

// As I code this I'm convincing myself this should maybe be its own pass, but I also think
// that ifs as expressions could themselves be considered "complex". In other words, removing
// the complex operands from `x = (if whatever something else)` really requires reordering the if.
// This is just as well an argument that it should be a separate pass, just a pass prior to this one,
// but this is _my_ program and I get to make it bad any way I please.
/*
func RemoveComplexOperandsFromIf(ifExpr *languages.FlatIfExpr, makeAtomic bool, getVar func() *languages.Var) (*languages.SimpleExpr, []*languages.SimpleAssignment, error) {
	    var newCondStatements []*languages.SimpleStatement
	    var newTrueStatements []*languages.SimpleStatement
	    var newFalseStatements []*languages.SimpleStatement
		newIfCond, newAssigns, err := RemoveComplexOperandsFromExpr(expr.IfExpr.IfCond, true, getVar)
		if err != nil {
			return nil, nil, err
		}
		for _, a := range(newAssigns) {
			newCondStatements = append(newCondStatements, &languages.SimpleStatement{Assignment: a})
		}

        branchVar := getVar()
		newIfTrueExpr, newTrueAssigns, err := RemoveComplexOperandsFromStatement(expr.IfExpr.IfTrue, false, getVar)
		if err != nil {
			return nil, nil, err
		}
		for _, a := range(newTrueAssigns) {
			newTrueStatements = append(newTrueStatements, &languages.SimpleStatement{Assignment: a})
		}
        newTrueStatements = append(
            newTrueStatements,
            &languages.SimpleStatement{
                Assignment: &languages.SimpleAssignment{
                    Ref: branchVar,
                    Expr: newIfTrueExpr,
                },
            },
        )

		newIfFalseExpr, newFalseAssigns, err := RemoveComplexOperandsFromExpr(statement.IfStmt.IfFalse, true, getVar)
		if err != nil {
			return nil, nil, err
		}
		for _, a := range(newFalseAssigns) {
			newFalseStatements = append(newFalseStatements, &languages.SimpleStatement{Assignment: a})
		}
        newFalseStatements = append(
            newFalseStatements,
            &languages.SimpleStatement{
                Assignment: &languages.SimpleAssignment{
                    Ref: branchVar,
                    Expr: newIfFalseExpr,
                },
            },
        )

		newStatements = append(
            newStatements,
            &languages.SimpleStatement{
                IfStmt: &languages.SimpleIfStmt{
                    IfCond: newIfCond,
                    IfTrue: newTrueStatements,
                    IfFalse: newFalseStatements,
                },
            },
        )
		return newStatements, nil, nil

}
*/

// There is a lot of overlap here with the above function but it felt better doing it this way to help limit the types
// Will probably see a cleaner way to do this in about 72 hours
func GenerateAtomicExpression(expr *languages.IfStmtExpr, getVar func() *languages.Var) (*languages.SimplePrimitive, []*languages.Assignment[languages.SimpleExpr], error) {
	switch {
	case expr.Num != nil:
		return &languages.SimplePrimitive{Num: expr.Num}, []*languages.Assignment[languages.SimpleExpr]{}, nil
	case expr.Var != nil:
		return &languages.SimplePrimitive{Var: expr.Var}, []*languages.Assignment[languages.SimpleExpr]{}, nil
	case expr.App != nil:
		var primitives []*languages.SimplePrimitive
		var newAssignments []*languages.Assignment[languages.SimpleExpr]
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
		newApp := &languages.SimpleExpr{
			App: &languages.SimpleApplication{
				Operator: rator.Var,
				Operands: rands,
			},
		}
		newVar := getVar()
		newAssignment := &languages.Assignment[languages.SimpleExpr]{
			Ref: newVar,
			Expr: newApp,
		}

		return &languages.SimplePrimitive{Var: newVar}, append(newAssignments, newAssignment), nil

	default:
		return nil, nil, errors.New("Unrecognized FlatExpr type")
	}
}

func IsExprSimple(expr *languages.FlatExpr) bool {
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

func IsExprPrimitive(expr *languages.FlatExpr) bool {
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
