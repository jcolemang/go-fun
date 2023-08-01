package main

import (
	"errors"
    "fmt"
	"golang.org/x/exp/slices"
)

func AssignRegisters(prog *VarAssemblyProgram) (*X86Program, error) {
    liveBeforeSets := UncoverLive(prog.Instrs)
    fmt.Println("Live before sets")
	fmt.Println(LiveBeforeSetsToString(liveBeforeSets))
	return nil, errors.New("Haven't gotten here yet")
}

func LiveBeforeSetsToString(instrs []*LiveBeforeInstr) string {
    repr := ""
    for _, i := range(instrs) {
        instrStr := VarAssemblyInstrToString(i.Instr)
        liveBeforeStr := ""
        for key, _ := range(i.LiveBefore) {
            locationStr := LocationToStr(key)
            liveBeforeStr = liveBeforeStr + locationStr + ", "
        }
        repr = repr + instrStr + " // " + liveBeforeStr + "\n"
    }
    return repr
}

type LiveBeforeInstr struct {
	Instr *VarAssemblyInstr
	LiveBefore map[Location]bool
}

func UncoverLive(instructions []*VarAssemblyInstr) []*LiveBeforeInstr {
	reversedInstrs := make([]*VarAssemblyInstr, len(instructions))
	copy(reversedInstrs, instructions)
	slices.Reverse(reversedInstrs)

	processedInstrs := make([]*LiveBeforeInstr, len(reversedInstrs))
	liveAfter := make(map[Location]bool)
	for i, instr := range(reversedInstrs) {
		liveBefore := LiveBefore(instr, liveAfter)
		lai := LiveBeforeInstr{
			Instr: instr,
			LiveBefore: liveBefore,
		}
		liveAfter = liveBefore
		processedInstrs[i] = &lai
	}
    slices.Reverse(processedInstrs)
	return processedInstrs
}

type StackLocation struct {
	Offset int
	Register Register
}

type Location struct {
	Register Register
	Variable VarAssemblyVar
	Stack StackLocation
}

func LocationToStr(l Location) string {
    switch {
    case l.Register.Name != "":
        return l.Register.Name
    case l.Variable.Generated != 0:
        return "tmp" + fmt.Sprint(l.Variable.Generated)
    default:
        return "I'll get to this later"
    }
}

// Live before instruction k = live after k - locations written by k, plus any locations read by k
// the logic here is that any locations written are OVERwritten and so are no longer live
// if a location is written to but never read, it will never be added to the set. the mere reference
// to a location does not imply that it will be live.
func LiveBefore(instruction *VarAssemblyInstr, liveAfter map[Location]bool) map[Location]bool {
    read, written := LocationsReadWritten(instruction)
    return MergeMaps(MapDifference(liveAfter, written), read)
}

func MergeMaps(m1 map[Location]bool, m2 map[Location]bool) map[Location]bool {
    newMap := make(map[Location]bool)
    for key, _ := range(m1) {
        newMap[key] = true
    }
    for key, _ := range(m2) {
        newMap[key] = true
    }
    return newMap
}

func MapDifference(m1 map[Location]bool, m2 map[Location]bool) map[Location]bool {
    newMap := make(map[Location]bool)
    for key, _ := range(m1) {
        _, present := m2[key]
        if !present {
            newMap[key] = true
        }
    }
    return newMap
}

func LocationsReadWritten(instr *VarAssemblyInstr) (map[Location]bool, map[Location]bool) {
	locationsRead := make(map[Location]bool)
    locationsWritten := make(map[Location]bool)
    switch {
    case instr.Addq != nil: 
        locationsRead = MergeMaps(locationsRead, LocationsReferenced(instr.Addq[0]))
        locationsRead = MergeMaps(locationsRead, LocationsReferenced(instr.Addq[1]))
        locationsWritten = MergeMaps(locationsWritten, LocationsReferenced(instr.Addq[1]))
    case instr.Movq != nil: 
        locationsRead = MergeMaps(locationsRead, LocationsReferenced(instr.Movq[0]))
        locationsWritten = MergeMaps(locationsWritten, LocationsReferenced(instr.Movq[1]))
    }
    return locationsRead, locationsWritten
}

func LocationsReferenced(arg *VarAssemblyImmediate) map[Location]bool {
    locations := make(map[Location]bool)
    switch {
    case arg.Var != nil:
        locations[Location{Variable: *arg.Var}] = true
    case arg.Register != nil:
        locations[Location{Register: *arg.Register}] = true
    }
    return locations
}
