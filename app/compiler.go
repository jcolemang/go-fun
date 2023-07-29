package main

// step 1: Uniquify all variable names
// (+ 1 (let ((x 2)) (+ x (let ((x 3)) x)))))))) ->
// (+ 1 (let ((x_1 2)) (+ x_1 (let ((x_2 3)) x_2)))
// This will be done by Uniquify
// step 2: With unique variable names the lexical scoping rules of lets are no longer really useful
//         and can be replaced with good ol' assignments. This also means that not everything 
//         will still need to be nested expressions
// (+ 1 (let ((x_1 2)) (+ x_1 (let ((x_2 3)) x_2))) ->
// x_1 := 2
// x_2 := 3
// (+ 1 (+ x_1 x_2))
// step 3: Next comes removing the complex expressions as specified in the textbook