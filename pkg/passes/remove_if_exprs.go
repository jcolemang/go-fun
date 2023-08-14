package passes

import (
    "errors"
    "language/pkg/languages"
)

// turning (+ 1 (if #t 2 3)) into
// (if #t
//   x = 2
//   x = 3)
// (+ 1 x)
func UnexpressionFlatProgram(prog *languages.FlatProgram, getVar func() *languages.Var) (*languages.IfStmtProgram, error) {
    newStatements, err := UnexpressionStatementList(prog.Statements, getVar)
    if err != nil {
        return nil, err
    }

	return &languages.IfStmtProgram{
        Statements: newStatements,
    }, nil
}

func UnexpressionStatementList(stmts []*languages.FlatStatement, getVar func() *languages.Var) ([]*languages.IfStmtStatement, error) {
    var newStatements []*languages.IfStmtStatement
	for _, s := range(stmts) {
		statements, err := UnexpressionFlatStatement(s, getVar)
		if err != nil {
			return nil, err
		}
		newStatements = append(newStatements, statements...)
	}

	return newStatements, nil
}

func UnexpressionFlatStatement(stmt *languages.FlatStatement, getVar func() *languages.Var) ([]*languages.IfStmtStatement, error) {
    switch {
        case stmt.Expr != nil:
            newExpr, newStmts, err := UnexpressionFlatExpr(stmt.Expr, getVar)
            if err != nil {
                return nil, err
            }
            return append(
                newStmts,
                &languages.IfStmtStatement{
                    Expr: newExpr,
                },
            ), nil
        case stmt.Assignment != nil:
            newExpr, newStmts, err := UnexpressionFlatExpr(stmt.Assignment.Expr, getVar)
            if err != nil {
                return nil, err
            }
            newStmt := &languages.IfStmtStatement{
                Assignment: &languages.Assignment[languages.IfStmtExpr]{
                    Ref: stmt.Assignment.Ref,
                    Expr: newExpr,
                },
            }
            return append(newStmts, newStmt), nil
        case stmt.Return != nil:
            newExpr, newStmts, err := UnexpressionFlatExpr(stmt.Return, getVar)
            if err != nil {
                return nil, err
            }
            newStmt := &languages.IfStmtStatement{
                Return: newExpr,
            }

            return append(newStmts, newStmt), nil
        default:
            return nil, errors.New("Unrecognized statement in UnexpressionFlatStatement")
    }
}


func UnexpressionFlatExpr(expr *languages.FlatExpr, getVar func() *languages.Var) (*languages.IfStmtExpr, []*languages.IfStmtStatement, error) {
    switch {
        case expr.Bool != nil:
            return &languages.IfStmtExpr{
                Bool: expr.Bool,
            }, make([]*languages.IfStmtStatement, 0), nil
        case expr.Num != nil:
            return &languages.IfStmtExpr{
                Num: expr.Num,
            }, make([]*languages.IfStmtStatement, 0), nil
        case expr.Var != nil:
            return &languages.IfStmtExpr{
                Var: expr.Var,
            }, make([]*languages.IfStmtStatement, 0), nil
        case expr.App != nil:
            var newStmts []*languages.IfStmtStatement
            var newAppExprs []*languages.IfStmtExpr
            for _, e := range expr.App {
                newExpr, opStatements, err := UnexpressionFlatExpr(e, getVar)
                if err != nil {
                    return nil, nil, err
                }
                newStmts = append(newStmts, opStatements...)
                newAppExprs = append(newAppExprs, newExpr)
            }
            return &languages.IfStmtExpr{
                App: newAppExprs,
            }, newStmts, nil
        case expr.IfExpr != nil:
            branchVar := getVar()
            ifExpr, ifStmt, newStmts, err := UnexpressionIfExpr(expr.IfExpr, branchVar, getVar)
            if err != nil {
                return nil, nil, err
            }
            newStmts = append(newStmts, &languages.IfStmtStatement{IfStmt: ifStmt})
            newStmts = append(newStmts, &languages.IfStmtStatement{
                Expr: ifExpr,
            })
            return ifExpr, newStmts, nil
        default:
            return nil, nil, errors.New("Unrecognized expression in UnexpressionFlatExpr")
    }
}

func UnexpressionIfExpr(ifExpr *languages.FlatIfExpr, branchVar *languages.Var, getVar func() *languages.Var) (*languages.IfStmtExpr, *languages.IfStmtIfStmt, []*languages.IfStmtStatement, error) {
    trueStmts, err := UnexpressionStatementList(ifExpr.IfTrue, getVar)
    if err != nil {
        return nil, nil, nil, err
    }

    trueExpr, trueExprStmts, err := UnexpressionFlatExpr(ifExpr.IfTrueExpr, getVar)
    if err != nil {
        return nil, nil, nil, err
    }
    trueStmts = append(trueStmts, trueExprStmts...)
    trueStmts = append(trueStmts, &languages.IfStmtStatement{
        Assignment: &languages.Assignment[languages.IfStmtExpr]{
            Ref: branchVar,
            Expr: trueExpr,
        },
    })

    falseStmts, err := UnexpressionStatementList(ifExpr.IfFalse, getVar)
    if err != nil {
        return nil, nil, nil, err
    }

    falseExpr, falseExprStmts, err := UnexpressionFlatExpr(ifExpr.IfFalseExpr, getVar)
    if err != nil {
        return nil, nil, nil, err
    }
    falseStmts = append(falseStmts, falseExprStmts...)
    falseStmts = append(falseStmts, &languages.IfStmtStatement{
        Assignment: &languages.Assignment[languages.IfStmtExpr]{
            Ref: branchVar,
            Expr: falseExpr,
        },
    })

    condExpr, condStmts, err := UnexpressionFlatExpr(ifExpr.IfCond, getVar)
    if err != nil {
        return nil, nil, nil, err
    }

    branchVarExpr := &languages.IfStmtExpr{Var: branchVar}
    ifStmt := &languages.IfStmtIfStmt{
        IfCond: condExpr,
        IfTrue: trueStmts,
        IfFalse: falseStmts,
    }

    return branchVarExpr, ifStmt, condStmts, nil
}
