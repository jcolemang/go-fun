package main

import (
    "github.com/alecthomas/participle/v2"
	"github.com/alecthomas/participle/v2/lexer"
)

// X86 Language
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
    X86Reg *X86Register       `| @@`
	X86Offset *int            `| @Int`
	X86OffsetReg *X86Register `  "("@@")"`
}

type X86Register struct {
	Name *string `"%"@Ident`
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