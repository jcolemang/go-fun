package passes

import (
    "errors"
    "language/pkg/languages"
)

func IfExprsToStmts(prog *languages.FlatProgram, getVar func() *languages.Var) (*languages.IfStmtProgram, error) {
    var newStatements []*languages.IfStmtStatement
	for _, s := range(prog.Statements) {
		statements, newAssigns, err := RemoveComplexOperandsFromStatement(s, getVar)
		if err != nil {
			return nil, err
		}
		for _, a := range(newAssigns) {
			newStatements = append(newStatements, &languages.SimpleStatement{Assignment: a})
		}
		newStatements = append(newStatements, statements...)
	}

	return &languages.SimpleProgram{
        Statements: newStatements,
    }, nil
}

func 

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
