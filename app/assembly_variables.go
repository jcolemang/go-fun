package main

import (
)

type VarAssemblyProgram struct {
	Instrs []*VarAssemblyInstr
}

type VarAssemblyInstr struct {
	Addq [2]*VarAssemblyImmediate // val dest
	Movq [2]*VarAssemblyImmediate // val dest
}

type VarAssemblyImmediate struct {
	Var *VarAssemblyVar
	Int *int
	Register *Register
}

type VarAssemblyVar struct {
	Generated int
}