package main

// Need to turn 
// (let ((a 5) (b 10)) (+ a (+ b 5)))
// into 
// (let ((a 5) (b 10)) (let ((temp0 (+ b 5)) (+ a temp0))))
// in other words, this is making explicit the order of operations.
func RemoveComplexOperands(prog *Program) *Program {

	return nil
}