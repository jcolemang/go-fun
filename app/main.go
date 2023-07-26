package main

import (
	"os"

    "github.com/alecthomas/participle/v2"
	"github.com/alecthomas/participle/v2/lexer"
	"github.com/alecthomas/kong"
	"github.com/alecthomas/repr"
)

// Main Language
type Program struct {
	Expr *Expr `@@`
}

type Expr struct {
	NumVal *Num      `@@`
    VarExpr *Var     `| @@`
    AppExpr []*Expr  `| "(" @@ @@* ")" `
}

type Num struct {
	Value *float64 `@Float | @Int`
}

type Var struct {
    Value *string `@Ident`
}

func GetLanguageParser() *participle.Parser[Program] {
    parser := participle.MustBuild[Program]()
    return parser
}

// X86 Language
type X86Program struct {
	X86Directives []*X86Directive `@@*`
    X86Instrs []*X86Instr         `@@*`
}

type X86Directive struct {
	Name *string `( "."@Ident`
	Arg *string `   @Ident ) EOL`
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


// Code to actually run
var cli struct {
	Files []string `arg:"" optional:"" type:"existingfile" help:"GraphQL schema files to parse."`
}

func main() {
    parser := GetX86Parser()
	ctx := kong.Parse(&cli)
	for _, file := range cli.Files {
		r, err := os.Open(file)
		ctx.FatalIfErrorf(err)
		ast, err := parser.Parse(file, r)
		r.Close()
		repr.Println(ast)
		ctx.FatalIfErrorf(err)
	}
}
