package main

import (
	"errors"
	"strconv"
    "fmt"

    "language/pkg/languages"
    "language/pkg/graph"

    "github.com/alecthomas/repr"
	"golang.org/x/exp/slices"
)

func AssignRegisters(prog *languages.VarAssemblyProgram, debug bool) (*languages.ArmProgram, error) {
    liveAfterSets := UncoverLive(prog.Instrs)
	interferenceGraph := BuildInterferenceGraph(liveAfterSets)
	colorings := graph.ColorGraph(interferenceGraph)

    if debug {
        fmt.Println(graph.GraphToString(*interferenceGraph))
        repr.Println(colorings)
        fmt.Println(LiveAfterSetsToString(liveAfterSets))
    }

	newInstrs := make([]*languages.ArmInstr, len(prog.Instrs))
	for i, instr := range prog.Instrs {
		switch {
		case instr.Add != nil:
			first, second, third := instr.Add[0], instr.Add[1], instr.Add[2]
			firstArm, err := VarImmToArmArg(first, colorings)
			if err != nil {
				return nil, err
			}
			secondArm, err := VarImmToArmArg(second, colorings)
			if err != nil {
				return nil, err
			}
			thirdArm, err := VarImmToArmArg(third, colorings)
			if err != nil {
				return nil, err
			}
			newInstrs[i] = &languages.ArmInstr{
				Add: []*languages.ArmArg{
					firstArm,
					secondArm,
                    thirdArm,
				},
			}
		case instr.Mov != nil:
			// fair amount of repeated code with above case. Thought maybe fallthrough could help fix it but I'm not sure
			first, second := instr.Mov[0], instr.Mov[1]
			firstArm, err := VarImmToArmArg(first, colorings)
			if err != nil {
				return nil, err
			}
			secondArm, err := VarImmToArmArg(second, colorings)
			if err != nil {
				return nil, err
			}
			newInstrs[i] = &languages.ArmInstr{
				Mov: []*languages.ArmArg{
					firstArm,
					secondArm,
				},
			}
        case instr.Ret != nil:
			newInstrs[i] = &languages.ArmInstr{
				Ret: &languages.Ret{},
            }
		}

	}

	return &languages.ArmProgram{
		ArmDirectives: []*languages.ArmDirective{},
		ArmInstrs: newInstrs,
	}, nil
}

// Two locations interfere if in any instruction there is a write to one location while the other is live.
// Any live location subject to a write would be overwritten and so could not be read in the future of the
// program. The exception is in Mov (see below). Also note that no location can interfere with itself.
// Consider Mov s d and the live location v
// if s = v, then v does not interfere with d because v and d contain the same value.
// if d = v, then v does not interfere with d because v and d trivially contain the same value
// (this second if is the same in the other instructions)
func BuildInterferenceGraph(liveAfterSets []*LiveAfterInstr) *graph.Graph[Location] {
	newGraph := graph.NewGraph[Location]()
	for _, liveAfterSet := range liveAfterSets {
		read, written := LocationsReadWritten(liveAfterSet.Instr)
		switch {
		case liveAfterSet.Instr.Mov != nil:
            for writtenLoc, _ := range(written) {
                newGraph = graph.AddNode(*newGraph, writtenLoc)
                if len(read) == 0 {
                    for liveAfterLoc, _ := range(liveAfterSet.LiveAfter) {
                        if liveAfterLoc != writtenLoc {
                            newGraph = graph.AddEdge(*newGraph, writtenLoc, liveAfterLoc)
                        }
                    }
                } else {
                    for readLoc, _ := range(read) {
                        newGraph = graph.AddNode(*newGraph, readLoc)
                        for liveAfterLoc, _ := range(liveAfterSet.LiveAfter) {
                            if liveAfterLoc != writtenLoc && liveAfterLoc != readLoc {
                                newGraph = graph.AddEdge(*newGraph, writtenLoc, liveAfterLoc)
                            }
                        }
                    }
                }
			}
		default:
			for writtenLoc, _ := range(written) {
				newGraph = graph.AddNode(*newGraph, writtenLoc)
				for liveAfterLoc, _ := range(liveAfterSet.LiveAfter) {
					if writtenLoc != liveAfterLoc {
						newGraph = graph.AddEdge(*newGraph, writtenLoc, liveAfterLoc)
					}
				}
			}
		}
	}
	return newGraph
}

type LiveAfterInstr struct {
	Instr *languages.VarAssemblyInstr
	LiveAfter map[Location]bool
}

func LiveAfterSetsToString(instrs []*LiveAfterInstr) string {
    repr := ""
    for _, i := range(instrs) {
        instrStr := languages.VarAssemblyInstrToString(i.Instr)
        liveBeforeStr := ""
        for key, _ := range(i.LiveAfter) {
            locationStr := LocationToStr(key)
            liveBeforeStr = liveBeforeStr + locationStr + ", "
        }
        repr = repr + instrStr + " // " + liveBeforeStr + "\n"
    }
    return repr
}

func UncoverLive(instructions []*languages.VarAssemblyInstr) []*LiveAfterInstr {
	reversedInstrs := make([]*languages.VarAssemblyInstr, len(instructions))
	copy(reversedInstrs, instructions)
	slices.Reverse(reversedInstrs)

	processedInstrs := make([]*LiveAfterInstr, len(reversedInstrs))
	liveAfter := make(map[Location]bool)
	for i, instr := range(reversedInstrs) {
		lai := LiveAfterInstr{
			Instr: instr,
			LiveAfter: liveAfter,
		}
		liveBefore := LiveBefore(instr, liveAfter)
		liveAfter = liveBefore
		processedInstrs[i] = &lai
	}
    slices.Reverse(processedInstrs)
	return processedInstrs
}

type StackLocation struct {
	Offset int
	Register languages.Register
}

// converting everything to this location type makes things awkward when I need to look up the
// location coloring for an instruction as the instruction does not itself contain a location
type Location struct {
	// Register Register // why would registers need to get assigned to registers?
	Variable languages.VarAssemblyVar
	// Stack StackLocation // also, would would stack locations need to get assigned to registers?
}

// putting this here because of its dependence on this other silly function
func VarImmToArmArg(varImm *languages.VarAssemblyImmediate, colorings map[Location]int) (*languages.ArmArg, error) {
	switch {
	case varImm.Int != nil:
		return &languages.ArmArg{
			ArmInt: varImm.Int,
		}, nil
	case varImm.Var != nil: // this same type test is done twice here and in ImmToLoc which is silly
		loc, err := ImmToLoc(varImm)
		if err != nil {
			return nil, err
		}
		assignment, present := colorings[*loc]
		if !present {
			// return nil, errors.New("I made a mistake and there is an unassigned location")
			assignment = 0
		}
		return languages.GetLocation(assignment), nil
	case varImm.Register != nil:
		return &languages.ArmArg{
			ArmReg: varImm.Register,
		}, nil
	default: // must be a register, but I haven't handled that yet
		// return &ArmArg{
		// 	ArmOffset: nil,
		// 	ArmOffsetReg: nil,
		// }, nil
		return nil, errors.New("Unhandled VarAssemblyImmediate converting to ArmArg")
	}
}

func ImmToLoc(imm *languages.VarAssemblyImmediate) (*Location, error) {
	switch {
	case imm.Var != nil:
		return &Location{Variable: *imm.Var}, nil
	// case imm.Register != nil:
	// 	return &Location{Register: *imm.Register}, nil
	default:
		return nil, errors.New("Unhandled case in VarAssemblyImmediate conversion to Location because of my silly mistakes")
	}
}

func LocationToStr(l Location) string {
    switch {
    // case l.Register.Name != "":
    //     return l.Register.Name
    case l.Variable.Generated != 0:
        return "tmp" + strconv.Itoa(l.Variable.Generated)
    default:
        return "I'll get to this later"
    }
}

// Live before instruction k = live after k minus locations written by k, plus any locations read by k
// the logic here is that any locations written are OVERwritten and so are no longer live
// if a location is written to but never read, it will never be added to the set. the mere reference
// to a location does not imply that it will be live.
func LiveBefore(instruction *languages.VarAssemblyInstr, liveAfter map[Location]bool) map[Location]bool {
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

func LocationsReadWritten(instr *languages.VarAssemblyInstr) (map[Location]bool, map[Location]bool) {
	locationsRead := make(map[Location]bool)
    locationsWritten := make(map[Location]bool)
    switch {
    case instr.Add != nil:
        locationsRead = MergeMaps(locationsRead, LocationsReferenced(instr.Add[1]))
        locationsRead = MergeMaps(locationsRead, LocationsReferenced(instr.Add[2]))
        locationsWritten = MergeMaps(locationsWritten, LocationsReferenced(instr.Add[0]))
    case instr.Mov != nil:
        locationsRead = MergeMaps(locationsRead, LocationsReferenced(instr.Mov[1]))
        locationsWritten = MergeMaps(locationsWritten, LocationsReferenced(instr.Mov[0]))
    }
    return locationsRead, locationsWritten
}

func LocationsReferenced(arg *languages.VarAssemblyImmediate) map[Location]bool {
    locations := make(map[Location]bool)
    switch {
    case arg.Var != nil:
        locations[Location{Variable: *arg.Var}] = true
    // case arg.Register != nil:
    //     locations[Location{Register: *arg.Register}] = true
    }
    return locations
}
