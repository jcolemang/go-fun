package main

import (
    "errors"
)

// Need to turn
// (let ((a 5) (b 10)) (+ a (+ b 5)))
// into
// (let ((a 5) (b 10)) (let ((temp0 (+ b 5)) (+ a temp0))))
// in other words, this is making explicit the order of operations.
// in other other words, generating three address code
func RemoveComplexOperands(prog *FlatProgram, getVar func() *Var) (*FlatProgram, error) {
	var newStatements []*FlatStatement
	for _, s := range(prog.Statements) {
		statement, newAssigns, err := RemoveComplexOperandsFromStatement(s, getVar)
		if err != nil {
			return nil, err
		}
		for _, a := range(newAssigns) {
			newStatements = append(newStatements, &FlatStatement{Assignment: a})
		}
		newStatements = append(newStatements, statement)
	}

	return &FlatProgram{
        Statements: newStatements,
    }, nil
}

func RemoveComplexOperandsFromStatement(statement *FlatStatement, getVar func() *Var) (*FlatStatement, []*FlatAssignment, error) {
	switch {
	case statement.Expr != nil:
		_, _, err := RemoveComplexOperandsFromExpr(statement.Expr, false, getVar)
		if err != nil {
			return nil, nil, err
		}
		return nil, nil, nil
	case statement.Assignment != nil:
		_, _, err := RemoveComplexOperandsFromExpr(statement.Assignment.Expr, true, getVar)
		if err != nil {
			return nil, nil, err
		}
		return nil, nil, nil
	default:
		return nil, nil, errors.New("Unrecognized expression")
	}
}

// We recommend implementing this pass with an auxiliary method named rco_exp
// with two parameters: an LVar expression and a Boolean that specifies whether the
// expression needs to become atomic or not. The rco_exp method should return a
// pair consisting of the new expression and a list of pairs, associating new temporary
// variables with their initializing expressions.
// part of why this is annoying is uniquify should come first
//
// bool is necessary because we need to be able to distinguish between the (+ 1 2) in
// x = 1 + 2 and the (+ 1 2) in x = (+ x (+ 1 2)). The former only uses two "addresses"
// and the latter uses three. 
// 
func RemoveComplexOperandsFromExpr(expr *FlatExpr, makeAtomic bool, GetNewVar func() *Var) (*FlatExpr, []*FlatAssignment, error) {
    switch {
    case expr.Num != nil:
        return expr, []*FlatAssignment{}, nil
    case expr.Var != nil:
        return expr, []*FlatAssignment{}, nil
    case expr.App != nil:
		if makeAtomic {

		} else {

		}
        return nil, nil, errors.New("Getting to this later")
    default:
        return nil, nil, errors.New("Unrecognized expression type")
    }
}

func IsExprAtomic(expr *FlatExpr) bool {
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