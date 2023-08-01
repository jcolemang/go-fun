package main

import (
	"errors"
	"golang.org/x/exp/slices"
)

func AssignRegisters(prog *VarAssemblyProgram) (*X86Program, error) {
	return nil, errors.New("Haven't gotten here yet")
}

type LiveAfterInstr struct {
	Instr *VarAssemblyInstr
	LiveBefore map[Location]bool
}

func UncoverLive(instructions []*VarAssemblyInstr) []*LiveAfterInstr {
	reversedInstrs := make([]*VarAssemblyInstr, len(instructions))
	copy(reversedInstrs, instructions)
	slices.Reverse(reversedInstrs)

	processedInstrs := make([]*LiveAfterInstr, len(reversedInstrs))
	liveAfter := make(map[Location]bool)
	for i, instr := range(reversedInstrs) {
		liveBefore := LiveBefore(instr, liveAfter)
		lai := LiveAfterInstr{
			Instr: instr,
			LiveBefore: liveBefore,
		}
		liveAfter = liveBefore
		processedInstrs[i] = &lai
	}
	return processedInstrs
}

type StackLocation struct {
	Offset int
	Register *Register
}

type Location struct {
	Register *Register
	Variable *Var
	Stack *StackLocation
}

// Live before instruction k = live after k - locations written by k, plus any locations read by k
// the logic here is that any locations written are OVERwritten and so are no longer live
// if a location is written to but never read, it will never be added to the set. the mere reference
// to a location does not imply that it will be live.
func LiveBefore(instruction *VarAssemblyInstr, liveAfter map[Location]bool) map[Location]bool {
	locationsRead := make(map[Location]bool)
	locationsWritten := make(map[Location]bool)
	return nil
}