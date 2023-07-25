package main

import (
	"fmt"
	"os"

	"github.com/alecthomas/kong"
	"github.com/alecthomas/repr"

	"github.com/alecthomas/participle/v2"
	// "github.com/alecthomas/participle/v2/lexer"
)

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

var cli struct {
	EBNF  bool     `help"Dump EBNF."`
	Files []string `arg:"" optional:"" type:"existingfile" help:"GraphQL schema files to parse."`
}

func main() {
    parser := participle.MustBuild[Program]()
	ctx := kong.Parse(&cli)
	if cli.EBNF {
		fmt.Println(parser.String())
		ctx.Exit(0)
	}
	for _, file := range cli.Files {
		r, err := os.Open(file)
		ctx.FatalIfErrorf(err)
		ast, err := parser.Parse(file, r)
		r.Close()
		repr.Println(ast)
		ctx.FatalIfErrorf(err)
	}
}