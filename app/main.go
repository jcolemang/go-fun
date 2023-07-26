package main

import (
	"os"

    "github.com/alecthomas/participle/v2"
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
    X86Instrs []*X86Instrs ` "globl" "main" EOL "main:" EOL @@*`
}

type X86Instrs struct {
    Addq []*X86Arg `"addq" @@ "," @@`
}

type X86Arg struct {
    X86Int *int         `"$"@Int`
    X86Reg *X86Register `| @@`
}

type X86Register struct {
	Name *string `"%"@("rax")`
}

func GetX86Parser() *participle.Parser[X86Program] {
    parser := participle.MustBuild[X86Program]()
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
