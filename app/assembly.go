package main

import (
	"strconv"
    "github.com/alecthomas/participle/v2"
	"github.com/alecthomas/participle/v2/lexer"
)

// X86 Language
// the parsing is not really necessary and this would be slightly better without it because I could use arrays with specific lengths
// which the parsing library does not allow
type X86Program struct {
	X86Directives []*X86Directive `@@*`
    X86Instrs []*X86Instr         `@@*`
}

type X86Directive struct {
	Name *string `( "."@Ident`
	Arg *string  `  @Ident ) EOL`
}

type X86Instr struct {
	Label *string  `( @Ident":"`
    Addq []*X86Arg `  | "addq" @@ "," @@`
	Movq []*X86Arg `  | "movq" @@ "," @@`
	Retq string    `  | @"retq" ) EOL`
}

type X86Arg struct {
    X86Int *int               `"$"@Int`
    X86Reg *Register          `| @@`
	X86Offset *int            `| @Int`
	X86OffsetReg *Register    `"("@@")"`
}

type Register struct {
	// argument passing: rdi rsi rdx rcx r8 r9
	// caller saved: rax rcx rdx rsi rdi r8 r9 r10 r11 -> the caller needs to save these, the callee can use freely
	// callee saved: rsp rbp rbx r12 r13 r14 r15 -> callee can use these, but must restore them, caller can use freely
	Name string `"%"@Ident`
}

func X86ProgramToString(prog *X86Program) string {
	s := ""
	for _, dir := range prog.X86Directives {
		s += X86DirectiveToString(dir) + "\n"
	}
	for _, instr := range prog.X86Instrs {
		s += X86InstrToString(instr) + "\n"
	}
	return s
}

func X86DirectiveToString(directive *X86Directive) string {
	return *directive.Name + " " + *directive.Arg
}

func X86InstrToString(instr *X86Instr) string {
	switch {
	case instr.Label != nil:
		return *instr.Label + ":"
	case instr.Addq != nil:
		return "\taddq " + X86ArgToString(instr.Addq[0]) + " " + X86ArgToString(instr.Addq[1])
	case instr.Movq != nil:
		return "\tmovq " + X86ArgToString(instr.Movq[0]) + " " + X86ArgToString(instr.Movq[1])
	default:
		return "Haven't implemented print for this one yet"
	}
}

func X86ArgToString(arg *X86Arg) string {
	switch {
	case arg.X86Int != nil:
		return strconv.Itoa(*arg.X86Int)
	case arg.X86Reg != nil:
		return "%" + arg.X86Reg.Name
	default:
		// must be a stack location
		return strconv.Itoa(*arg.X86Offset) + "(%" + arg.X86OffsetReg.Name + ")"
	}
}

func GetLocation(i int) *X86Arg {
	assignableRegisters := []Register{
		Register{Name: "rcx"},
		Register{Name: "rdx"},
		Register{Name: "rsi"},
		Register{Name: "rdi"},
		Register{Name: "r8"},
		Register{Name: "r9"},
		Register{Name: "r10"},
		Register{Name: "rbx"},
		Register{Name: "r12"},
		Register{Name: "r13"},
		Register{Name: "r14"},
	}

	if i < len(assignableRegisters) {
		return &X86Arg{X86Reg: &assignableRegisters[i]}
	} else {
		offset := (i - len(assignableRegisters) + 1) * -8
		return &X86Arg{
			X86Offset: &offset,
			X86OffsetReg: &Register{Name: "rbp"}, // base pointer
		}
	}
}



func GetArgumentRegisters() []*Register {
	return []*Register{
		&Register{Name: "rdi"},
		&Register{Name: "rsi"},
		&Register{Name: "rdx"},
		&Register{Name: "rcx"},
		&Register{Name: "r8"},
		&Register{Name: "r9"},
	}
}

func GetCalleeClobbered() []*Register {
	return []*Register{
		&Register{Name: "rax"},
		&Register{Name: "rcx"},
		&Register{Name: "rdx"},
		&Register{Name: "rsi"},
		&Register{Name: "rdi"},
		&Register{Name: "r8"},
		&Register{Name: "r9"},
		&Register{Name: "r10"},
		&Register{Name: "r11"},
	}
}

func GetCalleeSaved() []*Register {
	return []*Register{
		&Register{Name: "rsp"},
		&Register{Name: "rbp"},
		&Register{Name: "rbx"},
		&Register{Name: "r12"},
		&Register{Name: "r13"},
		&Register{Name: "r14"},
		&Register{Name: "r15"},
	}
}

func GetX86Parser() *participle.Parser[X86Program] {
	basicLexer := lexer.MustSimple([]lexer.SimpleRule{
		{"Comment", `(?i)rem[^\n]*`},
		{"String", `"(\\"|[^"])*"`},
		{"Ident", `[a-zA-Z_]\w*`},
		{"Int", `[-]?(\d+)`},
		{"Punct", `[-[!@#$%^&*()+_={}\|:;"'<,>.?/]|]`},
		{"EOL", `[\n\r]+`},
		{"whitespace", `[ \t]+`},
	})

	parser := participle.MustBuild[X86Program](
		participle.Lexer(basicLexer),
	)

    return parser
}