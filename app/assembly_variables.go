package main

import (
)

type VarAssemblyProgram struct {
	Instrs []*VarAssemblyInstr
}

type VarAssemblyInstr struct {
	Addq [2]*VarAssemblyImmediate
	Movq [2]*VarAssemblyImmediate // dest val
}

type VarAssemblyImmediate struct {
	Var *VarAssemblyVar
	Int int
}

type VarAssemblyInt struct {
	Value int
}

type VarAssemblyVar struct {
	Generated int
}