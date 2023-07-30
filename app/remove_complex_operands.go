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
		statements, newAssigns, err := RemoveComplexOperandsFromStatement(s, getVar)
		if err != nil {
			return nil, err
		}
		for _, a := range(newAssigns) {
			newStatements = append(newStatements, &FlatStatement{Assignment: a})
		}
		newStatements = append(newStatements, statements...)
	}

	return &FlatProgram{
        Statements: newStatements,
    }, nil
}

// Program after Flatten
// tmp2 = 3
// tmp1 = tmp2
// ( + 1 2 ( + tmp1 4  )  )

// Program after RemoveComplexOperands
// 3
// tmp2
// tmp3 = ( + tmp1 4  )
// ( + 1 2 tmp3  )

// PS C:\Users\James\Documents\Code\language> cat .\test-files\prog1
// (+ 1
//    2
//    (let ((x (let ((y 3)) y)))
//      (+ x 4)))
func RemoveComplexOperandsFromStatement(statement *FlatStatement, getVar func() *Var) ([]*FlatStatement, []*FlatAssignment, error) {
	var newStatements []*FlatStatement
	switch {
	case statement.Expr != nil:
		newExpr, newAssigns, err := RemoveComplexOperandsFromExpr(statement.Expr, false, getVar)
		if err != nil {
			return nil, nil, err
		}
		for _, a := range(newAssigns) {
			newStatements = append(newStatements, &FlatStatement{Assignment: a})
		}
		newStatements = append(newStatements, &FlatStatement{Expr: newExpr})
		return newStatements, nil, nil
	case statement.Assignment != nil:
		newExpr, newAssigns, err := RemoveComplexOperandsFromExpr(statement.Assignment.Expr, true, getVar)
		if err != nil {
			return nil, nil, err
		}
		for _, a := range(newAssigns) {
			newStatements = append(newStatements, &FlatStatement{Assignment: a})
		}
		newStatements = append(newStatements, &FlatStatement{Assignment: &FlatAssignment{Ref: statement.Assignment.Ref, Expr: newExpr}})

		return newStatements, nil, nil
	default:
		return nil, nil, errors.New("Unrecognized statement")
	}
}

// bool is necessary because we need to be able to distinguish between the (+ 1 2) in
// x = 1 + 2 and the (+ 1 2) in x = (+ x (+ 1 2)). The former only uses two "addresses"
// and the latter uses three. 
// 
func RemoveComplexOperandsFromExpr(expr *FlatExpr, makeAtomic bool, getVar func() *Var) (*FlatExpr, []*FlatAssignment, error) {
	var newExpr *FlatExpr
	var newAssignments []*FlatAssignment
    switch {
    case expr.Num != nil:
        return expr, []*FlatAssignment{}, nil
    case expr.Var != nil:
        return expr, []*FlatAssignment{}, nil
    case expr.App != nil:
		if IsExprSimple(expr) {
			if makeAtomic {
				newVar := getVar()
				newAssignment := &FlatAssignment{
					Ref: newVar,
					Expr: expr,
				}	
				newExpr = &FlatExpr{Var: newVar}
				newAssignments = append(newAssignments, newAssignment)
				return newExpr, newAssignments, nil
			} else {
				return expr, newAssignments, nil
			}
		} else {
			// Program after Flatten
			// ( + 1 ( + 2 ( + 3 4  )  ) ( + 5 6  )  )

			// Program after RemoveComplexOperands
			// tmp1 = ( + 3 4  )
			// tmp2 = ( tmp1  )
			// tmp3 = ( + 5 6  )
			// ( tmp3  )
			newExpr = &FlatExpr{App: []*FlatExpr{}}
			for _, e := range(expr.App) {
				subExpr, assigns, err := RemoveComplexOperandsFromExpr(e, true, getVar)
				if err != nil {
					return nil, nil, err
				}
				
				newExpr.App = append(newExpr.App, subExpr)
				newAssignments = append(newAssignments, assigns...)
			}
			// expression is now simple, don't need to duplicate the IsExprSimple branch logic
			finalExpr, finalAssigns, err := RemoveComplexOperandsFromExpr(newExpr, makeAtomic, getVar)
			if err != nil {
				return nil, nil, err
			}
			return finalExpr, append(newAssignments, finalAssigns...), nil
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