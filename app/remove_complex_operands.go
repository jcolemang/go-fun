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
func RemoveComplexOperands(prog *Program) (*Program, error) {
    GetNewVar := GetRandomGenerator()

    newExpr, err := RemoveComplexOperandsFromExpr(prog.Expr, GetNewVar)
    if err != nil {
        return nil, err
    }

	return &Program{
        Expr: newExpr,
    }, nil
}

func RemoveComplexOperandsFromExpr(expr *Expr, GetNewVar *func() int) (*Expr, error) {
    switch {
    case expr.Num != nil:
        return expr, nil
    case expr.Var != nil:
        return expr, nil
    case expr.Let != nil:
        return nil, nil
    case expr.App != nil:
        // (+ (+ 1 2) (+ 2 3)) ->
        // (let ((tmp1 (+ 1 2)))
        //   (+ tmp1 (+ 2 3))) ->
        // (let ((tmp1 (+ 1 2)))
        //   (let ((tmp2 (+ 2 3))
        //     (+ tmp1 tmp2))))
        for _, e := range(expr.App) {
            
        }


        return nil, nil
    default:
        return nil, errors.New("Unrecognized expression type")
    }
}

func IsOperandComplex(expr *Expr) bool {
    switch {
    case expr.Num != nil:
        return true
    case expr.Var != nil:
        return true
    case expr.Let != nil:
        for _, la := range(expr.Let.LetAssignments) {
            if IsOperandComplex(la.Expr) {
                return true
            }
        }

        return IsOperandComplex(expr.Let.LetBody)
    case expr.App != nil:
        for _, e := range(expr.App) {
            if IsOperandComplex(e) {
                return true
            }
        }
        return false
    default:
        // errors should be checked elsewhere
        return false
    }
}

func GetRandomGenerator() *func() int {
    current := 0
    generator := func() int {
        current++
        return current
    }
    return &generator
}
