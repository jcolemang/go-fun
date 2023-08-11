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
    var newStatements []*languages.IfStmtStatement
	for _, s := range(prog.Statements) {
		statements, err := UnexpressionFlatStatement(s, getVar)
		if err != nil {
			return nil, err
		}
		newStatements = append(newStatements, statements...)
	}

	return &languages.IfStmtProgram{
        Statements: newStatements,
    }, nil
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
            // var newCondStatements []*languages.SimpleStatement
            // var newTrueStatements []*languages.SimpleStatement
            // var newFalseStatements []*languages.SimpleStatement
            return nil, nil, errors.New("Haven't gotten here yet if")
        default:
            return nil, nil, errors.New("Unrecognized expression in UnexpressionFlatExpr")
    }
}
