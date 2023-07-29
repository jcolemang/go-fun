package main

import (
    "errors"
	"fmt"
)

// Need to turn
// (let ((a 5) (b 10)) (+ a (+ b 5)))
// into
// (let ((a 5) (b 10)) (let ((temp0 (+ b 5)) (+ a temp0))))
// in other words, this is making explicit the order of operations.
// in other other words, generating three address code
func RemoveComplexOperands(prog *Program) (*Program, error) {
	fmt.Println("Removing complex operands")
    GetNewVar := GetVarGenerator()

    newExpr, assignments, err := RemoveComplexOperandsFromExpr(prog.Expr, GetNewVar)
    if err != nil {
        return nil, err
    }
	if len(assignments) != 0 {
		return nil, errors.New("Left over assignments")
	}

	return &Program{
        Expr: newExpr,
    }, nil
}

// We recommend implementing this pass with an auxiliary method named rco_exp
// with two parameters: an LVar expression and a Boolean that specifies whether the
// expression needs to become atomic or not. The rco_exp method should return a
// pair consisting of the new expression and a list of pairs, associating new temporary
// variables with their initializing expressions.
// part of why this is annoying is uniquify should come first
func RemoveComplexOperandsFromExpr(expr *Expr, GetNewVar func() *Var) (*Expr, []*Assignment, error) {
    switch {
    case expr.Num != nil:
        return expr, []*Assignment{}, nil
    case expr.Var != nil:
        return expr, []*Assignment{}, nil
    case expr.Let != nil:
        return nil, []*Assignment{}, nil
    case expr.App != nil:
        // (+ (+ 1 (+ 2 3)) (+ 4 5)) ->
		// (let ((tmp1 (+ 2 3))
		//       (tmp2 (+ 1 tmp1))
		//       (tmp3 (+ 2 3))
	    //   (+ tmp1 tmp2 tmp3))
		// (+ 1 (let ((x 2)) x)) ->
		// 
		var subExprs []*Expr
		var assignments []*Assignment

		for _, e := range(expr.App) {
			// (+ (+ 1 2) 3)
			// get tmp1 = (+ 1 2)
			// add let to assign to returned subexpression
			if IsOperandPrimitive(e) {
				subExprs = append(subExprs, e)
			} else {
				// e here is (+ 1 2)
				// no assignments
				// should create a new variable and assign to e
				// then replace the current subexpression with that new variable
				e, assgns, err := RemoveComplexOperandsFromExpr(e, GetNewVar)
				if err != nil {
					return nil, nil, err
				}

				newVar := GetNewVar()
				newAssignment := Assignment{Ref: newVar, Expr: e}

				subExprs = append(subExprs, &Expr{Var: newVar})
				assignments = append(append(assignments, assgns...), &newAssignment)
			}
		}

		var newExpr Expr
		if len(assignments) == 0 {
			newExpr = Expr{
				App: subExprs,
			}
		} else {
			newExpr = Expr{
				Let: &LetExpr{
					LetAssignments: assignments,
					LetBody: &Expr{
						App: subExprs,
					},
				},
			}
		}


		// no new assignments to pass up
        return &newExpr, []*Assignment{}, nil
    default:
        return nil, nil, errors.New("Unrecognized expression type")
    }
}

func IsOperandPrimitive(expr *Expr) bool {
	switch {
    case expr.Num != nil:
        return true
    case expr.Var != nil:
        return true
    case expr.Let != nil:
        return false
    case expr.App != nil:
        return false
    default:
        // errors should be checked elsewhere
        return false
    }
}


